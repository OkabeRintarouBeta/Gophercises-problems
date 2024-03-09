/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/okaberintaroubeta/task/db"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A list of all of your tasks",

	Run: func(cmd *cobra.Command, args []string) {
		tasks, _ := db.ListAllTasks()
		if len(tasks) == 0 {
			fmt.Println("No tasks now. Go and add some!")
		} else {
			for i, task := range tasks {
				fmt.Printf("%d:\t %s\n", i+1, task.Value)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
