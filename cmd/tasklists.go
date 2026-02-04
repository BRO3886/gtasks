package cmd

import (
	"github.com/BRO3886/gtasks/api"
	"github.com/BRO3886/gtasks/internal/utils"
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
	gtasks tasklists add -t <TITLE>
	gtasks tasklists add --title <TITLE>

	Remove tasklist
	gtasks tasklists rm

	`,
}

var showlistsCmd = &cobra.Command{
	Use:   "view",
	Short: "view tasklists",
	Long:  `view task lists for the account currently signed in`,
	Run: func(cmd *cobra.Command, args []string) {
		srv, err := api.GetService()
		if err != nil {
			utils.ErrorP("Failed to get service: %v\n", err)
			return
		}
		list, err := api.GetTaskLists(srv)
		if err != nil {
			utils.ErrorP("Error: %v\n", err)
		}

		for index, i := range list {
			utils.Print("[%d] %s\n", index+1, i.Title)
		}
	},
}

var addListcmd = &cobra.Command{
	Use:   "add",
	Short: "add tasklist",
	Long:  `add tasklist for the currently signed in account`,
	Run: func(cmd *cobra.Command, args []string) {
		srv, err := api.GetService()
		if err != nil {
			utils.ErrorP("Failed to get service: %v\n", err)
			return
		}
		if title == "" {
			utils.Warn("%s\n", "Title should not be empty. Use -t for title.\nExamples:\ngtasks tasklists create -t <TITLE>\ngtasks tasklists create --title <TITLE>")
			return
		}
		t := &tasks.TaskList{Title: title}
		r, err := srv.Tasklists.Insert(t).Do()
		if err != nil {
			utils.ErrorP("Unable to create task list. %v", err)
		}
		title = ""
		utils.Info("task list created: %s", r.Title)
	},
}

var removeListCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove tasklist",
	Long:  `Remove a tasklist for the currently signed in account`,
	Run: func(cmd *cobra.Command, args []string) {
		srv, err := api.GetService()
		if err != nil {
			utils.ErrorP("Failed to get service: %v\n", err)
			return
		}
		list, err := api.GetTaskLists(srv)
		if err != nil {
			utils.ErrorP("Error %v", err)
		}

		utils.Print("Choose a Tasklist: ")
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
		utils.Print("%s: %s\n", utils.WarnStyle.Sprint("Deleting list..."), result)

		err = api.DeleteTaskList(srv, list[option].Id)
		if err != nil {
			utils.ErrorP("Error deleting tasklist: %s", err.Error())
			return
		}
		utils.Info("Tasklist deleted")
	},
}

var updateTitleCmd = &cobra.Command{
	Use:   "update",
	Short: "update tasklist title",
	Long:  `Update tasklist title for the currently signed in account`,
	Run: func(cmd *cobra.Command, args []string) {
		srv, err := api.GetService()
		if err != nil {
			utils.ErrorP("Failed to get service: %v\n", err)
			return
		}
		if title == "" {
			utils.Warn("Title should not be empty. Use -t for title.\nExamples:\ngtasks tasklists update -t <TITLE>\ngtasks tasklists update --title <TITLE>\n")
			return
		}

		list, err := api.GetTaskLists(srv)
		if err != nil {
			utils.ErrorP(utils.Error("Error %v", err))
		}

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
		t := list[option]
		t.Title = title

		_, err = api.UpdateTaskList(srv, &t)
		if err != nil {
			utils.ErrorP("Error updating tasklist: ", err.Error())
		}
		utils.Info("Tasklist title updated")
	},
}

var title string

func init() {
	addListcmd.Flags().StringVarP(&title, "title", "t", "", "title of task list (required)")
	updateTitleCmd.Flags().StringVarP(&title, "title", "t", "", "title of task list (required)")
	tasklistsCmd.AddCommand(showlistsCmd, addListcmd, removeListCmd, updateTitleCmd)
	rootCmd.AddCommand(tasklistsCmd)
}
