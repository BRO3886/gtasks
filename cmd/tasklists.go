package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/BRO3886/google-tasks-cli/utils"
	"github.com/spf13/cobra"
	"google.golang.org/api/tasks/v1"
)

// tasklistsCmd represents the tasklists command
var tasklistsCmd = &cobra.Command{
	Use:   "tasklists",
	Short: "View and create tasklists for currently signed-in account",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	Use:   "show",
	Short: "show tasklists",
	Long:  `Show task lists for the account currently signed in`,
	Run: func(cmd *cobra.Command, args []string) {
		config := utils.ReadCredentials()
		client := getClient(config)
		srv, err := tasks.New(client)
		if err != nil {
			log.Fatalf("Unable to retrieve tasks Client %v", err)
		}

		r, err := srv.Tasklists.List().Do()
		if err != nil {
			log.Fatalf("Unable to retrieve task lists. %v", err)
		}

		fmt.Println("Task Lists:")
		if len(r.Items) > 0 {
			for _, i := range r.Items {
				fmt.Printf("%s (%s)\n", i.Title, i.Id)
			}
		} else {
			fmt.Print("No task lists found.")
		}
	},
}

var createlistsCmd = &cobra.Command{
	Use:   "create",
	Short: "create tasklist",
	Long:  `TODO:add`,
	Run: func(cmd *cobra.Command, args []string) {
		config := utils.ReadCredentials()
		client := getClient(config)
		srv, err := tasks.New(client)
		if err != nil {
			log.Fatalf("Unable to retrieve tasks Client %v", err)
		}
		if title == "" {
			fmt.Println("Title should not be empty. Use -t for title.\nExamples:\ngtasks tasklists create -t <TITLE>\ngtasks tasklists create --title <TITLE>")
			return
		}
		t := &tasks.TaskList{Title: title}
		r, err := srv.Tasklists.Insert(t).Do()
		if err != nil {
			log.Fatalf("Unable to create task list. %v", err)
		}
		fmt.Println("Created: " + r.Title)
	},
}

var title string

func init() {
	createlistsCmd.Flags().StringVarP(&title, "title", "t", "", "title of task list (required)")
	tasklistsCmd.AddCommand(showlistsCmd, createlistsCmd)
	rootCmd.AddCommand(tasklistsCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tasklistsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
