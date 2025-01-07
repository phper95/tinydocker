package main

import (
	"github.com/spf13/cobra"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	var rootCmd = &cobra.Command{Use: "task"}

	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a new task",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Task added:", args[0])
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Listing tasks...")
		},
	}

	rootCmd.AddCommand(addCmd, listCmd)
	rootCmd.Execute()
}
