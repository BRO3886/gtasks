package cmd

import (
	"bufio"
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
	gtasks tasks -l "<task-list name>" view|add|rm|done
	
	[WITHOUT LIST FLAG]
	gtasks tasks view|add|rm|done
	* You would be prompted to select a tasklist
	`,
}

var viewTasksCmd = &cobra.Command{
	Use:   "view",
	Short: "View tasks in a tasklist",
	Long: `
	Use this command to view tasks in a selected 
	tasklist for the currently signed in account
	`,
	Run: func(cmd *cobra.Command, args []string) {
		srv := getService()
		tList := getTaskLists(srv)

		utils.Print("Tasks in %s:\n", tList.Title)

		tasks, err := api.GetTasks(srv, tList.Id, includeCompletedFlag)
		if err != nil {
			color.Red(err.Error())
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"No", "Title", "Description", "Status", "Due"})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
		// table.SetRowLine(true)
		// table.SetRowSeparator("-")

		for ind, task := range tasks {
			row := []string{fmt.Sprintf("%d", ind+1), task.Title, task.Notes}

			if task.Status == "needsAction" {
				row = append(row, "✖")
			} else if task.Status == "completed" {
				row = append(row, "✔")
			}

			due, err := time.Parse(time.RFC3339, task.Due)
			if err != nil {
				row = append(row, "No Due Date")
			} else {
				row = append(row, due.Local().Format("02 January 2006"))
			}

			table.Append(row)

		}
		table.Render()
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
		srv := getService()
		tList := getTaskLists(srv)
		utils.Warn("Creating task in %s\n", tList.Title)

		reader := bufio.NewReader(os.Stdin)

		utils.Print("Title: ")
		title := getInput(reader)
		utils.Print("Note: ")
		notes := getInput(reader)
		utils.Print("Due Date: ")
		dateInput := getInput(reader)

		var dateString string

		if dateInput != "" {
			// All possible examples: https://github.com/araddon/dateparse#extended-example
			t, err := dateparse.ParseAny(dateInput)
			if err != nil {
				utils.ErrorP("%s", "date Format incorrect")
			}

			dateString = t.Format(time.RFC3339)
		} else {
			dateString = ""
		}

		task := &tasks.Task{Title: title, Notes: notes, Due: dateString}

		_, err := api.CreateTask(srv, task, tList.Id)
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
		srv := getService()
		tList := getTaskLists(srv)
		tID := tList.Id

		tasks, err := api.GetTasks(srv, tID, false)
		if err != nil {
			color.Red(err.Error())
			return
		}

		var taskIndex int
		argProvided := false

		if len(args) == 1 {
			argProvided = true

			index, err := strconv.Atoi(args[0])
			if err != nil || index >= len(tasks) {
				utils.ErrorP("%s", "Incorrect task number\n")
			}

			taskIndex = index
		} else {
			utils.Print("Tasks in %s:\n", tList.Title)

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
				color.Red("Error: \n" + err.Error())
				return
			}

			taskIndex = option
		}

		if argProvided {
			taskIndex--
		}

		t := tasks[taskIndex]
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
		srv := getService()
		tList := getTaskLists(srv)
		utils.Print("Tasks in %s:\n", tList.Title)
		tID := tList.Id

		tasks, err := api.GetTasks(srv, tID, false)
		if err != nil {
			color.Red(err.Error())
			return
		}

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
			color.Red("Error: " + err.Error())
			return
		}

		t := tasks[option]
		t.Status = "completed"

		err = api.DeleteTask(srv, t.Id, tID)
		if err != nil {
			color.Red("Unable to delete task: %v", err)
			return
		}
		utils.Info("Deleted: %s\n", t.Title)
	},
}

var (
	includeCompletedFlag bool
	taskListFlag         string
)

func init() {
	viewTasksCmd.Flags().BoolVarP(&includeCompletedFlag, "include-completed", "i", false, "use this flag to include completed tasks")
	tasksCmd.PersistentFlags().StringVarP(&taskListFlag, "tasklist", "l", "", "use this flag to specify a tasklist")
	tasksCmd.AddCommand(viewTasksCmd, createTaskCmd, markCompletedCmd, deleteTaskCmd)
	rootCmd.AddCommand(tasksCmd)
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
