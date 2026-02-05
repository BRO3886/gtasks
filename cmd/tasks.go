package cmd

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/BRO3886/gtasks/api"
	"github.com/BRO3886/gtasks/internal/utils"
	"github.com/araddon/dateparse"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"google.golang.org/api/tasks/v1"
)

// tasksCmd represents the tasks command
var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "View, create, and delete tasks in a tasklist",
	Long: `
	View, create, list and delete tasks in a tasklist
	for the currently signed in account.
	Usage:
	[WITH LIST FLAG]
	gtasks tasks -l "<task-list name>" view|add|rm|done|info

	[WITHOUT LIST FLAG]
	gtasks tasks view|add|rm|done|info
	* You would be prompted to select a tasklist
	`,
}

var viewTasksCmd = &cobra.Command{
	Use:   "view",
	Short: "View tasks in a tasklist",
	Long: `
	Use this command to view tasks in a selected 
	tasklist for the currently signed in account.
	You can control output with --format: table (default), json, csv.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		srv, err := api.GetService()
		if err != nil {
			utils.ErrorP("Failed to get service: %v\n", err)
			return
		}
		tList := getTaskLists(srv)

		taskItems, err := api.GetTasks(srv, tList.Id, viewTasksFlags.includeCompleted || viewTasksFlags.onlyCompleted, viewTasksFlags.max)
		if err != nil {
			color.Red(err.Error())
			return
		}

		utils.Sort(taskItems, viewTasksFlags.sort)

		var filteredTasks []*tasks.Task
		for _, task := range taskItems {
			if viewTasksFlags.onlyCompleted && task.Status == "needsAction" {
				continue
			}
			filteredTasks = append(filteredTasks, task)
		}

		switch viewTasksFlags.format {
		case "json":
			outputJSON(filteredTasks)
		case "csv":
			outputCSV(filteredTasks)
		default:
			outputTable(filteredTasks, tList.Title)
		}
	},
}

var createTaskCmd = &cobra.Command{
	Use:   "add",
	Short: "Add task in a tasklist",
	Long: `
	Use this command to add tasks in a selected 
	tasklist for the currently signed in account.
	
	Supports recurring tasks with --repeat flag:
	  gtasks tasks add -t "Standup" -d "2025-02-10" --repeat daily --repeat-count 5
	  gtasks tasks add -t "Weekly sync" -d "2025-02-10" --repeat weekly --repeat-until "2025-03-10"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		srv, err := api.GetService()
		if err != nil {
			utils.ErrorP("Failed to get service: %v\n", err)
			return
		}
		tList := getTaskLists(srv)
		utils.Warn("Creating task in %s\n", tList.Title)

		var title string
		var notes string
		var dateInput string

		if addTaskFlags.title == "" && (addTaskFlags.note != "" || addTaskFlags.due != "") {
			utils.ErrorP("Please specify a task title")
			return
		} else if addTaskFlags.title != "" {
			title = addTaskFlags.title
			notes = addTaskFlags.note
			dateInput = addTaskFlags.due

		} else {
			reader := bufio.NewReader(os.Stdin)

			utils.Print("Title: ")
			title = getInput(reader)
			utils.Print("Note: ")
			notes = getInput(reader)
			utils.Print("Due Date: ")
			dateInput = getInput(reader)
		}

		// Parse repeat pattern if specified
		repeatPattern, err := parseRepeatUnit(addTaskFlags.repeat)
		if err != nil {
			utils.ErrorP("%v\n", err)
			return
		}

		// If repeat is specified but no due date, require due date
		if repeatPattern != repeatNone && dateInput == "" {
			utils.ErrorP("Due date (--due) is required when using --repeat\n")
			return
		}

		var startDate time.Time
		if dateInput != "" {
			// All possible examples: https://github.com/araddon/dateparse#extended-example
			t, err := dateparse.ParseAny(dateInput)
			if err != nil {
				utils.ErrorP("Date format incorrect. Valid examples: https://github.com/araddon/dateparse#extended-example\n")
				return
			}
			startDate = t
		}

		// Parse repeat-until date if specified
		var untilDate *time.Time
		if addTaskFlags.repeatUntil != "" {
			t, err := dateparse.ParseAny(addTaskFlags.repeatUntil)
			if err != nil {
				utils.ErrorP("repeat-until date format incorrect. Valid examples: https://github.com/araddon/dateparse#extended-example\n")
				return
			}
			untilDate = &t
		}

		// Generate dates for recurring tasks
		var dates []time.Time
		if repeatPattern != repeatNone {
			dates = expandRepeatSchedule(startDate, repeatPattern, addTaskFlags.repeatCount, untilDate)
		} else if dateInput != "" {
			dates = []time.Time{startDate}
		} else {
			dates = []time.Time{} // No due date
		}

		// Create tasks
		if len(dates) == 0 {
			// No due date specified
			task := &tasks.Task{Title: title, Notes: notes}
			_, err = api.CreateTask(srv, task, tList.Id)
			if err != nil {
				utils.ErrorStyle.Printf("Unable to create task: %v", err)
				return
			}
			utils.Info("Task created\n")
		} else if len(dates) == 1 {
			// Single task with due date
			task := &tasks.Task{Title: title, Notes: notes, Due: dates[0].Format(time.RFC3339)}
			_, err = api.CreateTask(srv, task, tList.Id)
			if err != nil {
				utils.ErrorStyle.Printf("Unable to create task: %v", err)
				return
			}
			utils.Info("Task created\n")
		} else {
			// Multiple recurring tasks
			utils.Info("Creating %d recurring tasks...\n", len(dates))
			for i, d := range dates {
				task := &tasks.Task{Title: title, Notes: notes, Due: d.Format(time.RFC3339)}
				_, err = api.CreateTask(srv, task, tList.Id)
				if err != nil {
					utils.ErrorStyle.Printf("Unable to create task %d: %v\n", i+1, err)
					return
				}
			}
			utils.Info("Created %d tasks\n", len(dates))
		}
	},
}

var markCompletedCmd = &cobra.Command{
	Use:   "done",
	Short: "Mark tasks as done",
	Long: `
	Use this command to mark a task as completed
	in a selected tasklist for the currently signed in account
	`,
	Run: func(cmd *cobra.Command, args []string) {
		srv, err := api.GetService()
		if err != nil {
			utils.ErrorP("Failed to get service: %v\n", err)
			return
		}
		tList := getTaskLists(srv)
		tID := tList.Id

		tasks, err := api.GetTasks(srv, tID, false, 0)
		if err != nil {
			color.Red(err.Error())
			return
		}

		ind := getTaskIndex(args, tasks, tList.Title)
		t := tasks[ind]
		t.Status = "completed"

		_, err = api.UpdateTask(srv, t, tID)
		if err != nil {
			color.Red("Unable to mark task as completed: %v", err)
			return
		}
		utils.Info("Marked as complete: %s\n", t.Title)
	},
}

var undoTaskCmd = &cobra.Command{
	Use:   "undo",
	Short: "Mark a completed task as incomplete",
	Long: `
	Use this command to mark a completed task as incomplete
	in a selected tasklist for the currently signed in account
	`,
	Run: func(cmd *cobra.Command, args []string) {
		srv, err := api.GetService()
		if err != nil {
			utils.ErrorP("Failed to get service: %v\n", err)
			return
		}
		tList := getTaskLists(srv)
		tID := tList.Id

		// Get completed tasks only
		taskItems, err := api.GetTasks(srv, tID, true, 0)
		if err != nil {
			color.Red(err.Error())
			return
		}

		// Filter to only completed tasks
		var completedTasks []*tasks.Task
		for _, task := range taskItems {
			if task.Status == "completed" {
				completedTasks = append(completedTasks, task)
			}
		}

		if len(completedTasks) == 0 {
			utils.Info("No completed tasks to undo\n")
			return
		}

		ind := getTaskIndex(args, completedTasks, tList.Title)
		t := completedTasks[ind]
		t.Status = "needsAction"
		t.Completed = nil

		_, err = api.UpdateTask(srv, t, tID)
		if err != nil {
			color.Red("Unable to mark task as incomplete: %v", err)
			return
		}
		utils.Info("Marked as incomplete: %s\n", t.Title)
	},
}

var clearTasksCmd = &cobra.Command{
	Use:   "clear",
	Short: "Hide all completed tasks from the list",
	Long: `
	Use this command to hide all completed tasks from a tasklist.
	This marks completed tasks as hidden so they won't be returned
	by the API. Primarily affects tasks completed via the CLI.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		srv, err := api.GetService()
		if err != nil {
			utils.ErrorP("Failed to get service: %v\n", err)
			return
		}
		tList := getTaskLists(srv)

		// Confirmation prompt unless --force is set
		if !clearTasksFlags.force {
			prompt := promptui.Prompt{
				Label:     fmt.Sprintf("Clear all completed tasks from '%s'", tList.Title),
				IsConfirm: true,
			}
			_, err := prompt.Run()
			if err != nil {
				utils.Info("Cancelled\n")
				return
			}
		}

		err = api.ClearTasks(srv, tList.Id)
		if err != nil {
			color.Red("Unable to clear completed tasks: %v", err)
			return
		}
		utils.Info("Cleared completed tasks from %s\n", tList.Title)
	},
}

var deleteTaskCmd = &cobra.Command{
	Use:   "rm",
	Short: "Delete a task in a tasklist",
	Long: `
	Use this command to delete a task in a tasklist
	for the currently signed in account
	`,
	Run: func(cmd *cobra.Command, args []string) {
		srv, err := api.GetService()
		if err != nil {
			utils.ErrorP("Failed to get service: %v\n", err)
			return
		}
		tList := getTaskLists(srv)
		tID := tList.Id

		tasks, err := api.GetTasks(srv, tID, false, 0)
		if err != nil {
			color.Red(err.Error())
			return
		}

		ind := getTaskIndex(args, tasks, tList.Title)
		t := tasks[ind]

		err = api.DeleteTask(srv, t.Id, tID)
		if err != nil {
			color.Red("Unable to delete task: %v", err)
			return
		}
		utils.Info("Deleted: %s\n", t.Title)
	},
}

var infoTaskCmd = &cobra.Command{
	Use:   "info [task-number]",
	Short: "View detailed information about a task",
	Long: `
	Use this command to view detailed information about a task
	including links, notes, and other metadata.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		srv, err := api.GetService()
		if err != nil {
			utils.ErrorP("Failed to get service: %v\n", err)
			return
		}
		tList := getTaskLists(srv)
		tID := tList.Id

		tasks, err := api.GetTasks(srv, tID, infoTaskFlags.includeCompleted, 0)
		if err != nil {
			color.Red(err.Error())
			return
		}

		ind := getTaskIndex(args, tasks, tList.Title)
		t := tasks[ind]

		// Display detailed task information
		utils.Print("\n")
		utils.Print("Task: %s\n", t.Title)

		// Status
		status := "Needs action"
		if t.Status == "completed" {
			status = "Completed"
		}
		utils.Print("Status: %s\n", status)

		// Due date
		if t.Due != "" {
			due, err := time.Parse(time.RFC3339, t.Due)
			if err == nil {
				utils.Print("Due: %s\n", due.Local().Format("02 January 2006"))
			} else {
				utils.Print("Due: Not set\n")
			}
		} else {
			utils.Print("Due: Not set\n")
		}

		// Notes
		if t.Notes != "" {
			utils.Print("Notes: %s\n", t.Notes)
		} else {
			utils.Print("Notes: None\n")
		}

		// Links
		utils.Print("\n")
		if len(t.Links) > 0 {
			utils.Print("Links:\n")
			for _, link := range t.Links {
				utils.Print("  - %s\n", link.Link)
			}
		} else {
			utils.Print("Links: No links\n")
		}

		// WebViewLink
		if t.WebViewLink != "" {
			utils.Print("\nView in Google Tasks: %s\n", t.WebViewLink)
		}
		utils.Print("\n")
	},
}

var updateTaskCmd = &cobra.Command{
	Use:   "update [task-number]",
	Short: "Update an existing task",
	Long: `
	Use this command to update an existing task in a tasklist.
	
	Interactive mode (no flags): prompts for each field showing current values.
	Press Enter to keep the current value, or type a new value.
	
	Flag mode: only update fields that are explicitly provided.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		srv, err := api.GetService()
		if err != nil {
			utils.ErrorP("Failed to get service: %v\n", err)
			return
		}
		tList := getTaskLists(srv)
		tID := tList.Id

		taskItems, err := api.GetTasks(srv, tID, false, 0)
		if err != nil {
			color.Red(err.Error())
			return
		}

		ind := getTaskIndex(args, taskItems, tList.Title)
		t := taskItems[ind]

		utils.Info("Updating task: %s\n\n", t.Title)

		// Check if any flags were provided
		titleFlagSet := cmd.Flags().Changed("title")
		noteFlagSet := cmd.Flags().Changed("note")
		dueFlagSet := cmd.Flags().Changed("due")
		flagMode := titleFlagSet || noteFlagSet || dueFlagSet

		var newTitle, newNote, newDue string

		if flagMode {
			// Flag mode: only update fields that were explicitly set
			if titleFlagSet {
				newTitle = updateTaskFlags.title
			} else {
				newTitle = t.Title
			}
			if noteFlagSet {
				newNote = updateTaskFlags.note
			} else {
				newNote = t.Notes
			}
			if dueFlagSet {
				newDue = updateTaskFlags.due
			}
		} else {
			// Interactive mode: prompt for each field
			reader := bufio.NewReader(os.Stdin)

			// Title
			currentTitle := t.Title
			utils.Print("Title [%s]: ", currentTitle)
			newTitle = getInput(reader)
			if newTitle == "" {
				newTitle = currentTitle
			}

			// Note
			currentNote := t.Notes
			if currentNote == "" {
				utils.Print("Note []: ")
			} else {
				utils.Print("Note [%s]: ", currentNote)
			}
			newNote = getInput(reader)
			if newNote == "" {
				newNote = currentNote
			}

			// Due date
			currentDue := formatDueHuman(t.Due)
			if currentDue == "-" {
				utils.Print("Due []: ")
			} else {
				utils.Print("Due [%s]: ", currentDue)
			}
			newDue = getInput(reader)
		}

		// Apply changes
		t.Title = newTitle
		t.Notes = newNote

		// Parse and set due date if provided
		if newDue != "" {
			parsedDue, err := dateparse.ParseAny(newDue)
			if err != nil {
				utils.ErrorP("Date format incorrect. Valid examples: https://github.com/araddon/dateparse#extended-example\n")
				return
			}
			t.Due = parsedDue.Format(time.RFC3339)
		} else if !flagMode && newDue == "" {
			// Keep existing due date in interactive mode when user presses Enter
		} else if flagMode && dueFlagSet && newDue == "" {
			// Clear due date if --due="" was explicitly set
			t.Due = ""
		}

		_, err = api.UpdateTask(srv, t, tID)
		if err != nil {
			color.Red("Unable to update task: %v", err)
			return
		}
		utils.Info("\nUpdated: %s\n", t.Title)
	},
}

var (
	viewTasksFlags struct {
		includeCompleted bool
		onlyCompleted    bool
		sort             string
		format           string
		max              int
	}
	taskListFlag string
	addTaskFlags struct {
		title       string
		note        string
		due         string
		repeat      string
		repeatCount int
		repeatUntil string
	}
	clearTasksFlags struct {
		force bool
	}
	updateTaskFlags struct {
		title string
		note  string
		due   string
	}
	infoTaskFlags struct {
		includeCompleted bool
	}
)

func init() {
	createTaskCmd.Flags().StringVarP(&addTaskFlags.title, "title", "t", "", "use this flag to set a tasks title")
	createTaskCmd.Flags().StringVarP(&addTaskFlags.note, "note", "n", "", "use this flag to set a tasks note")
	createTaskCmd.Flags().StringVarP(&addTaskFlags.due, "due", "d", "", "due date (e.g., '2024-12-25', 'Dec 25', 'tomorrow')")
	createTaskCmd.Flags().StringVarP(&addTaskFlags.repeat, "repeat", "r", "", "repeat pattern: daily, weekly, monthly, yearly")
	createTaskCmd.Flags().IntVar(&addTaskFlags.repeatCount, "repeat-count", 0, "number of occurrences for repeating task")
	createTaskCmd.Flags().StringVar(&addTaskFlags.repeatUntil, "repeat-until", "", "end date for repeating task (e.g., '2025-03-01')")
	viewTasksCmd.Flags().BoolVarP(&viewTasksFlags.includeCompleted, "include-completed", "i", false, "use this flag to include completed tasks")
	viewTasksCmd.Flags().BoolVar(&viewTasksFlags.onlyCompleted, "completed", false, "use this flag to only show completed tasks")
	viewTasksCmd.Flags().StringVar(&viewTasksFlags.sort, "sort", "position", "use this flag to sort by [due,title,position]")
	viewTasksCmd.Flags().StringVar(&viewTasksFlags.format, "format", "table", "output format: table, json, csv")
	viewTasksCmd.Flags().IntVar(&viewTasksFlags.max, "max", 0, "maximum number of tasks to return (0 = all)")
	clearTasksCmd.Flags().BoolVarP(&clearTasksFlags.force, "force", "f", false, "skip confirmation prompt")
	updateTaskCmd.Flags().StringVarP(&updateTaskFlags.title, "title", "t", "", "new title for the task")
	updateTaskCmd.Flags().StringVarP(&updateTaskFlags.note, "note", "n", "", "new note for the task")
	updateTaskCmd.Flags().StringVarP(&updateTaskFlags.due, "due", "d", "", "new due date for the task")
	infoTaskCmd.Flags().BoolVarP(&infoTaskFlags.includeCompleted, "include-completed", "i", false, "include completed tasks when selecting by number")
	tasksCmd.PersistentFlags().StringVarP(&taskListFlag, "tasklist", "l", "", "use this flag to specify a tasklist")
	tasksCmd.AddCommand(viewTasksCmd, createTaskCmd, markCompletedCmd, undoTaskCmd, deleteTaskCmd, clearTasksCmd, infoTaskCmd, updateTaskCmd)
	rootCmd.AddCommand(tasksCmd)
}

type TaskOutput struct {
	Number      int    `json:"number"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status"`
	Due         string `json:"due,omitempty"`
}

func outputTable(tasks []*tasks.Task, listTitle string) {
	utils.Print("Tasks in %s:\n", listTitle)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"No", "Title", "Description", "Status", "Due"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetCenterSeparator("|")
	table.SetRowLine(false)
	table.SetRowSeparator("-")
	table.SetAutoWrapText(false)

	for ind, task := range tasks {
		row := []string{
			fmt.Sprintf("%d", ind+1),
			truncate(task.Title, 30),
			truncate(task.Notes, 40),
			statusLabel(task.Status),
			formatDueHuman(task.Due),
		}
		table.Append(row)
	}

	table.Render()
}

func truncate(s string, maxLen int) string {
	// Replace newlines with spaces for single-line display
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", "")
	
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func outputJSON(tasks []*tasks.Task) {
	var output []TaskOutput

	for ind, task := range tasks {
		output = append(output, TaskOutput{
			Number:      ind + 1,
			Title:       task.Title,
			Description: task.Notes,
			Status:      statusLabel(task.Status),
			Due:         formatDueISO(task.Due),
		})
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(output)
}

func outputCSV(tasks []*tasks.Task) {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	_ = writer.Write([]string{"No", "Title", "Description", "Status", "Due"})

	for ind, task := range tasks {
		_ = writer.Write([]string{
			fmt.Sprintf("%d", ind+1),
			task.Title,
			task.Notes,
			statusLabel(task.Status),
			formatDueHuman(task.Due),
		})
	}
}

func statusLabel(status string) string {
	switch status {
	case "completed":
		return "completed"
	case "needsAction":
		return "pending"
	default:
		return status
	}
}

func formatDueHuman(due string) string {
	if due == "" {
		return "-"
	}
	parsed, err := time.Parse(time.RFC3339, due)
	if err != nil {
		return "-"
	}
	return parsed.Local().Format("02 January 2006")
}

func formatDueISO(due string) string {
	if due == "" {
		return ""
	}
	parsed, err := time.Parse(time.RFC3339, due)
	if err != nil {
		return ""
	}
	return parsed.Local().Format("2006-01-02")
}

func getInput(reader *bufio.Reader) string {
	title, _ = reader.ReadString('\n')
	if runtime.GOOS == "windows" {
		title = strings.Replace(title, "\r\n", "", -1)
	} else {
		title = strings.Replace(title, "\n", "", -1)
	}
	return title
}

func getTaskIndex(args []string, tasks []*tasks.Task, title string) int {
	var taskIndex int
	argProvided := false
	if len(args) == 1 {
		argProvided = true

		index, err := strconv.Atoi(args[0])
		if err != nil || index > len(tasks) || index < 1 {
			utils.ErrorP("%s", "Incorrect task number\n")
		}

		taskIndex = index
	} else {
		utils.Print("Tasks in %s:\n", title)

		tString := []string{}
		for _, i := range tasks {
			tString = append(tString, i.Title)
		}

		prompt := promptui.Select{
			Label: "Select Task",
			Items: tString,
		}

		option, _, err := prompt.Run()
		if err != nil {
			utils.ErrorP("Error: %s\n", err.Error())
		}

		taskIndex = option
	}

	if argProvided {
		taskIndex--
	}

	return taskIndex
}

func getTaskLists(srv *tasks.Service) tasks.TaskList {
	list, err := api.GetTaskLists(srv)
	if err != nil {
		utils.ErrorP("Error %v", err)
	}

	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Title <= list[j].Title
	})

	index := -1

	if taskListFlag != "" {

		var titles []string
		for _, tasklist := range list {
			titles = append(titles, tasklist.Title)
		}

		index = sort.SearchStrings(titles, taskListFlag)

		if !(index >= 0 && index < len(list) && list[index].Title == taskListFlag) {
			utils.ErrorP("%s\n", "incorrect task-list name")
		}

	} else {
		utils.Print("Choose a Tasklist:")
		var l []string
		for _, i := range list {
			l = append(l, i.Title)
		}

		prompt := promptui.Select{
			Label: "Select Tasklist",
			Items: l,
		}
		option, _, err := prompt.Run()
		if err != nil {
			utils.ErrorP("Error: %s", err.Error())
		}

		index = option
	}

	return list[index]
}
