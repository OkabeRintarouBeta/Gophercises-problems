/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/okaberintaroubeta/task/db"
	"github.com/spf13/cobra"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task as complete",

	Run: func(cmd *cobra.Command, args []string) {
		var ids []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Fall to parse the argument:", arg)
			} else {
				ids = append(ids, id)
			}
		}

		tasks, err := db.ListAllTasks()
		if err != nil {
			fmt.Println("Something is wrong with listing the tasks...")
			return
		}
		for _, id := range ids {
			if id <= 0 || id > len(tasks) {
				fmt.Println("Invalid Task Number:", id)
				continue
			} else {
				task := tasks[id-1]
				err := db.DeleteTask(task.Key)
				if err != nil {
					fmt.Printf("Error when deleting the task %d: %v\n", id, err)
					continue
				}
				fmt.Printf("Finished task: %s\n", task.Value)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(doCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// doCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// doCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
