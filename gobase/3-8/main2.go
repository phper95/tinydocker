package main

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

type Task struct {
	ID    int64
	Title string
	Done  bool
}

var tasks []Task

func loadTasks() {
	file, err := os.ReadFile("tasks.json")
	if err != nil {
		panic("Error loading tasks:" + err.Error())
	}
	err = json.Unmarshal(file, &tasks)
	if err != nil {
		log.Println("Error unmarshalling tasks:" + err.Error())
	}

}

func addTask(args ...string) {
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		log.Println("Error parsing ID:", err)
		return
	}
	// 添加任务逻辑
	log.Println("Adding task:", id, args[1])
	tasks = append(tasks, Task{ID: id, Title: args[1], Done: false})
	saveTask()
}

func saveTask() {
	data, err := json.Marshal(tasks)
	if err != nil {
		log.Println("Error marshalling tasks:", err)
		return
	}
	err = os.WriteFile("tasks.json", data, 0644)
	if err != nil {
		log.Println("Error saving tasks:", err)
	}
}

func main() {
	// 加载任务
	loadTasks()

	rootCmd := &cobra.Command{Use: "task"}

	// 创建 "add" 子命令，用于添加新任务
	addcmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new task",
		Run: func(cmd *cobra.Command, args []string) {
			// 在这里添加添加新任务逻辑
			log.Println("Adding a new task", args)
			if len(args) < 2 {
				log.Println("Error: invalid arguments")
				return
			}
			addTask(args...)
		},
	}

	// 创建 "list" 子命令，用于列出所有任务
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		Run: func(cmd *cobra.Command, args []string) {
			// 在这里添加列出所有任务逻辑
			log.Println("Listing all tasks", args)
			for _, task := range tasks {
				log.Println("ID:", task.ID, "Title:", task.Title, "Done:", task.Done)
			}
		},
	}

	// 将 "add" 和 "list" 子命令添加到根命令中
	rootCmd.AddCommand(addcmd, listCmd)

	// 执行根命令
	if err := rootCmd.Execute(); err != nil {
		log.Println("Error executing command:", err)
	}

}
