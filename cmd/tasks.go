package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/BRO3886/gtasks/utils"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"google.golang.org/api/tasks/v1"
)

// tasksCmd represents the tasks command
var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "View, create, and delete tasks in a tasklist",
	// Long: `
	// View, create, list and delete tasks in a tasklist
	// for the currently signed in account.
	// `,
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

		list, err := utils.GetTaskLists(srv)
		if err != nil {
			log.Fatalf("Error %v", err)
		}

		fmt.Println("Choose a Tasklist:")
		var l []string
		for _, i := range list {
			l = append(l, i.Title)
		}

		prompt := promptui.Select{
			Label: "Select Tasklist",
			Items: l,
		}
		option, result, err := prompt.Run()
		fmt.Println("Creating task in " + result)

		tasks, err := utils.GetTasks(srv, list[option].Id)
		if err != nil {
			color.Red(err.Error())
			return
		}
		for index, i := range tasks {
			color.Green("[%d] %s\n", index+1, i.Title)
			fmt.Printf("%s: %s\n", color.YellowString("Description"), i.Notes)
			fmt.Printf("%s: %s\n", color.YellowString("Status"), i.Status)
			due, err := time.Parse(time.RFC3339, i.Due)
			if err != nil {
				fmt.Printf("No Due Date\n")
			} else {
				fmt.Printf("%s: %s\n", color.YellowString("Due"), due.Format("Mon Jan 2 2006 3:04PM"))
			}
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
		config := utils.ReadCredentials()
		client := getClient(config)

		srv, err := tasks.New(client)
		if err != nil {
			log.Fatalf("Unable to retrieve tasks Client %v", err)
		}

		list, err := utils.GetTaskLists(srv)
		if err != nil {
			log.Fatalf("Error %v", err)
		}

		fmt.Println("Choose a Tasklist:")
		var l []string
		for _, i := range list {
			l = append(l, i.Title)
		}

		prompt := promptui.Select{
			Label: "Select Tasklist",
			Items: l,
		}
		option, result, err := prompt.Run()
		fmt.Println("Creating task in " + result)

		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("Title: ")
		title := getInput(reader)
		fmt.Printf("Note: ")
		notes := getInput(reader)
		fmt.Printf("Due Date (dd/mm/yyyy): ")
		dateInput := getInput(reader)

		var dateString string

		if dateInput == "" {
			dateString = ""
		} else {
			arr := strings.Split(dateInput, "/")
			if len(arr) < 3 {
				color.Red("Date format incorrect")
				return
			}
			y, _ := strconv.Atoi(arr[2])
			if y < time.Now().Year() {
				color.Yellow("Please enter a valid year")
				return
			}
			d, _ := strconv.Atoi(arr[0])
			m, _ := strconv.Atoi(arr[1])

			t := time.Date(y, time.Month(m), d, 12, 0, 0, 0, time.UTC)
			dateString = t.Format(time.RFC3339)
		}
		task := &tasks.Task{Title: title, Notes: notes, Due: dateString}

		task, err = utils.CreateTask(srv, task, list[option].Id)
		if err != nil {
			color.Red("Unable to create task: %v", err)
			return
		}
		color.Green("Task created")
	},
}

func init() {
	tasksCmd.AddCommand(viewTasksCmd, createTaskCmd)
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
