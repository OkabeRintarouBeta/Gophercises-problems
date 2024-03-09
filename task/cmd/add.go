package cmd

import (
	"fmt"
	"strings"

	"github.com/okaberintaroubeta/task/db"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a task to your task list",
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")
		_, err := db.CreateTask(task)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Added task: %s\n", task)
		return
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
