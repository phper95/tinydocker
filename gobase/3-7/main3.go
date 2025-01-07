package main

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Task struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var tasks []Task

func loadTasks() {
	file, err := os.ReadFile("tasks.json")
	if err == nil {
		json.Unmarshal(file, &tasks)
	}
}

func saveTasks() {
	data, _ := json.Marshal(tasks)
	os.WriteFile("tasks.json", data, 0644)
}

func addTask(title string) {
	tasks = append(tasks, Task{Title: title, Done: false})
	saveTasks()
	log.Println("Task added:", title)
}

func main() {
	var rootCmd = &cobra.Command{Use: "task"}

	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a new task",
		Run: func(cmd *cobra.Command, args []string) {
			addTask(args[0])
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Listing tasks...")
			loadTasks()
			for _, task := range tasks {
				log.Println(task)
			}
		},
	}

	rootCmd.AddCommand(addCmd, listCmd)
	rootCmd.Execute()
}
