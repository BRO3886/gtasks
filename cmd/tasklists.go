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
	Long: `
	View and create tasklists for currently signed-in account
	
	View tasklists:
	gtasks tasklists show

	Create tasklist:
	gtasks tasklists create -t <TITLE>
	gtasks tasklists create --title <TITLE>

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
}

func getTaskLists(srv *tasks.Service) ([]*tasks.TaskList, error) {
	r, err := srv.Tasklists.List().Do()
	if err != nil {
		log.Fatalf("Unable to retrieve task lists. %v", err)
	}

	if len(r.Items) == 0 {
		return nil, errors.New("No Tasklist found")
	}
	return r.Items, nil
}
