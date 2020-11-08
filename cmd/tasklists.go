/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/tasks/v1"
)

// tasklistsCmd represents the tasklists command
var tasklistsCmd = &cobra.Command{
	Use:   "tasklists",
	Short: "TODO: add",
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
		b, err := ioutil.ReadFile("credentials.json")
		if err != nil {
			log.Fatalf("Unable to read client secret file: %v", err)
		}
		config, err := google.ConfigFromJSON(b, tasks.TasksScope)
		if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
		}
		srv, err := tasks.New(getClient(config))
		if err != nil {
			log.Fatalf("Unable to retrieve tasks Client %v", err)
		}

		r, err := srv.Tasklists.List().MaxResults(10).Do()
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

func init() {
	tasklistsCmd.AddCommand(showlistsCmd)
	rootCmd.AddCommand(tasklistsCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tasklistsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
