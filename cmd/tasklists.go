package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/BRO3886/gtasks/api"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"google.golang.org/api/tasks/v1"
)

// tasklistsCmd represents the tasklists command
var tasklistsCmd = &cobra.Command{
	Use:   "tasklists",
	Short: "View and create tasklists for currently signed-in account",
	Long: `
	View and create tasklists for currently signed-in account
	
	View tasklists:
	gtasks tasklists view

	Create tasklist:
	gtasks tasklists create -t <TITLE>
	gtasks tasklists create --title <TITLE>

	Remove tasklist
	gtasks tasklists rm

	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg. Use -h to show the list of available commands")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var showlistsCmd = &cobra.Command{
	Use:   "view",
	Short: "view tasklists",
	Long:  `view task lists for the account currently signed in`,
	Run: func(cmd *cobra.Command, args []string) {
		srv := getService()
		list, err := api.GetTaskLists(srv)
		if err != nil {
			log.Fatalf("Error %v", err)
		}

		for index, i := range list {
			fmt.Printf("[%d] %s\n", index+1, i.Title)
		}

	},
}

var createlistsCmd = &cobra.Command{
	Use:   "create",
	Short: "create tasklist",
	Long:  `Create tasklist for the currently signed in account`,
	Run: func(cmd *cobra.Command, args []string) {
		srv := getService()
		if title == "" {
			fmt.Println("Title should not be empty. Use -t for title.\nExamples:\ngtasks tasklists create -t <TITLE>\ngtasks tasklists create --title <TITLE>")
			return
		}
		t := &tasks.TaskList{Title: title}
		r, err := srv.Tasklists.Insert(t).Do()
		if err != nil {
			log.Fatalf("Unable to create task list. %v", err)
		}
		title = ""
		fmt.Println(color.GreenString("Created: ") + r.Title)
	},
}

var removeListCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove tasklist",
	Long:  `Remove a tasklist for the currently signed in account`,
	Run: func(cmd *cobra.Command, args []string) {
		srv := getService()
		list, err := api.GetTaskLists(srv)
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
		if err != nil {
			color.Red("Error: " + err.Error())
			return
		}
		fmt.Printf("%s: %s\n", color.YellowString("Deleting list"), result)

		err = api.DeleteTaskList(srv, list[option].Id)
		if err != nil {
			color.Red("Error deleting tasklist: " + err.Error())
			return
		}
		color.Green("Tasklist deleted")
	},
}

var updateTitleCmd = &cobra.Command{
	Use:   "update",
	Short: "update tasklist title",
	Long:  `Update tasklist title for the currently signed in account`,
	Run: func(cmd *cobra.Command, args []string) {
		srv := getService()
		if title == "" {
			fmt.Println("Title should not be empty. Use -t for title.\nExamples:\ngtasks tasklists update -t <TITLE>\ngtasks tasklists update --title <TITLE>")
			return
		}

		list, err := api.GetTaskLists(srv)
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
		option, _, err := prompt.Run()
		if err != nil {
			color.Red("Error: " + err.Error())
			return
		}
		t := list[option]
		t.Title = title

		_, err = api.UpdateTaskList(srv, &t)
		if err != nil {
			color.Red("Error updating tasklist: " + err.Error())
			return
		}
		color.Green("Tasklist title updated")
	},
}

var title string

func init() {
	createlistsCmd.Flags().StringVarP(&title, "title", "t", "", "title of task list (required)")
	updateTitleCmd.Flags().StringVarP(&title, "title", "t", "", "title of task list (required)")
	tasklistsCmd.AddCommand(showlistsCmd, createlistsCmd, removeListCmd, updateTitleCmd)
	rootCmd.AddCommand(tasklistsCmd)
}
