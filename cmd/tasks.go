package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/BRO3886/google-tasks-cli/utils"
	"github.com/fatih/color"
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
		config := utils.ReadCredentials()
		client := getClient(config)

		srv, err := tasks.New(client)
		if err != nil {
			log.Fatalf("Unable to retrieve tasks Client %v", err)
		}

		list, err := getTaskLists(srv)
		if err != nil {
			log.Fatalf("Error %v", err)
		}

		fmt.Println("Choose a Tasklist:")
		for index, i := range list {
			fmt.Printf("[%d] %s\n", index+1, i.Title)
		}
		fmt.Printf("Choose an option: ")
		var option int
		if _, err := fmt.Scan(&option); err != nil {
			log.Fatalf("Unable to read option: %v", err)
		}
		fmt.Println("Tasks in '" + list[option-1].Title + "':\n")

		tasks, err := getTasks(srv, list[option-1].Id)
		if err != nil {
			log.Fatalf("Error %v", err)
		}
		for index, i := range tasks {
			color.Green("[%d] %s\n", index+1, i.Title)
			fmt.Printf("%s\n", i.Notes)
			due, err := time.Parse(time.RFC3339, i.Due)
			if err != nil {
				fmt.Printf("No Due Date\n")
			} else {
				color.Yellow("Due %s\n", due.Format("Mon Jan 2 2006 3:04PM"))
			}
		}

	},
}

var createTaskCmd = &cobra.Command{
	Use:   "create",
	Short: "Create task in a tasklist",
	Long: `
	Use this command to create tasks in a selected 
	tasklist for the currently signed in account
	`,
	Run: func(cmd *cobra.Command, args []string) {
		config := utils.ReadCredentials()
		client := getClient(config)

		srv, err := tasks.New(client)
		if err != nil {
			log.Fatalf("Unable to retrieve tasks Client %v", err)
		}

		list, err := getTaskLists(srv)
		if err != nil {
			log.Fatalf("Error %v", err)
		}

		fmt.Println("Choose a Tasklist:")
		for index, i := range list {
			fmt.Printf("[%d] %s\n", index+1, i.Title)
		}
		fmt.Printf("Choose an option: ")

		var option int
		if _, err := fmt.Scan(&option); err != nil {
			log.Fatalf("Unable to read option: %v", err)
		}

		fmt.Printf("Title: ")
		title := getInput()
		fmt.Printf("Note: ")
		notes := getInput()
		fmt.Printf("Due Date (dd/mm/yyyy): ")
		dateInput := getInput()

		var dateString string

		if dateInput == "" {
			dateString = ""
		} else {
			arr := strings.Split(dateInput, "/")
			if len(arr) < 3 {
				color.Red("Time format incorrect")
				return
			}
			y, _ := strconv.Atoi(arr[2])
			d, _ := strconv.Atoi(arr[0])
			m, _ := strconv.Atoi(arr[1])

			t := time.Date(y, time.Month(m), d, 12, 0, 0, 0, time.UTC)
			dateString = t.Format(time.RFC3339)
		}
		task := &tasks.Task{Title: title, Notes: notes, Due: dateString}

		task, err = createTask(srv, task, list[option-1].Id)
		if err != nil {
			log.Fatalf("Unable to create task: %v", err)
		}
		color.Green("Task created")
	},
}

func init() {
	tasksCmd.AddCommand(viewTasksCmd, createTaskCmd)
	rootCmd.AddCommand(tasksCmd)
}

func getTasks(srv *tasks.Service, id string) ([]*tasks.Task, error) {
	r, err := srv.Tasks.List(id).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve tasks. %v", err)
	}
	if len(r.Items) == 0 {
		return nil, errors.New("No Tasks found")
	}
	return r.Items, nil
}

func createTask(srv *tasks.Service, task *tasks.Task, tasklistID string) (*tasks.Task, error) {
	r, err := srv.Tasks.Insert(tasklistID, task).Do()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func getInput() string {
	reader := bufio.NewReader(os.Stdin)
	title, _ = reader.ReadString('\n')
	if runtime.GOOS == "windows" {
		title = strings.Replace(title, "\r\n", "", -1)
	} else {
		title = strings.Replace(title, "\n", "", -1)
	}
	return title
}
