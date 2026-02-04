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

		taskItems, err := api.GetTasks(srv, tList.Id, viewTasksFlags.includeCompleted || viewTasksFlags.onlyCompleted)
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
	tasklist for the currently signed in account
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

		var dateString string

		if dateInput != "" {
			// All possible examples: https://github.com/araddon/dateparse#extended-example
			t, err := dateparse.ParseAny(dateInput)
			if err != nil {
				utils.ErrorP("Date format incorrect. Valid examples: https://github.com/araddon/dateparse#extended-example\n")
				return
			}

			dateString = t.Format(time.RFC3339)
		} else {
			dateString = ""
		}

		task := &tasks.Task{Title: title, Notes: notes, Due: dateString}

		_, err = api.CreateTask(srv, task, tList.Id)
		if err != nil {
			utils.ErrorStyle.Printf("Unable to create task: %v", err)
			return
		}
		utils.Info("Task created\n")
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

		tasks, err := api.GetTasks(srv, tID, false)
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

		tasks, err := api.GetTasks(srv, tID, false)
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

		tasks, err := api.GetTasks(srv, tID, true)
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

var (
	viewTasksFlags struct {
		includeCompleted bool
		onlyCompleted    bool
		sort             string
		format           string
	}
	taskListFlag string
	addTaskFlags struct {
		title string
		note  string
		due   string
	}
)

func init() {
	createTaskCmd.Flags().StringVarP(&addTaskFlags.title, "title", "t", "", "use this flag to set a tasks title")
	createTaskCmd.Flags().StringVarP(&addTaskFlags.note, "note", "n", "", "use this flag to set a tasks note")
	createTaskCmd.Flags().StringVarP(&addTaskFlags.due, "due", "d", "", "due date (e.g., '2024-12-25', 'Dec 25', 'tomorrow')")
	viewTasksCmd.Flags().BoolVarP(&viewTasksFlags.includeCompleted, "include-completed", "i", false, "use this flag to include completed tasks")
	viewTasksCmd.Flags().BoolVar(&viewTasksFlags.onlyCompleted, "completed", false, "use this flag to only show completed tasks")
	viewTasksCmd.Flags().StringVar(&viewTasksFlags.sort, "sort", "position", "use this flag to sort by [due,title,position]")
	viewTasksCmd.Flags().StringVar(&viewTasksFlags.format, "format", "table", "output format: table, json, csv")
	tasksCmd.PersistentFlags().StringVarP(&taskListFlag, "tasklist", "l", "", "use this flag to specify a tasklist")
	tasksCmd.AddCommand(viewTasksCmd, createTaskCmd, markCompletedCmd, deleteTaskCmd, infoTaskCmd)
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

	for ind, task := range tasks {
		row := []string{
			fmt.Sprintf("%d", ind+1),
			task.Title,
			task.Notes,
			statusLabel(task.Status),
			formatDueHuman(task.Due),
		}
		table.Append(row)
	}

	table.Render()
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
