package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

type DownloadInfo struct {
	Url      string
	FilePath string
}

func downloadFile(info DownloadInfo) error {
	// Download the file from the URL and save it to the file path
	response, err := http.Get(info.Url)
	if err != nil {
		return fmt.Errorf("Error downloading file err: %v url: %s , filePath: %s", err, info.Url, info.FilePath)
	}
	defer response.Body.Close()
	file, err := os.Create(info.FilePath)
	if err != nil {
		return fmt.Errorf("Error creating file err: %v filePath: %s", err, info.FilePath)
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("Error copying file err: %v filePath: %s", err, info.FilePath)
	}
	return nil
}

func worker(id int, jobs <-chan DownloadInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range jobs {
		fmt.Printf("Worker %d downloading file %s from %s\n", id, task.FilePath, task.Url)
		err := downloadFile(task)
		if err != nil {
			fmt.Printf("Worker %d error downloading file %s from %s err: %v\n", id, task.FilePath, task.Url, err)
		}
	}
}

func main() {
	tasks := []DownloadInfo{
		{Url: "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png", FilePath: "google.png"},
		{Url: "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png", FilePath: "google2.png"},
		{Url: "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png", FilePath: "google3.png"},
		{Url: "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png", FilePath: "google4.png"},
		{Url: "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png", FilePath: "google5.png"},
	}
	tasksChan := make(chan DownloadInfo, len(tasks))
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go worker(i, tasksChan, &wg)
	}
	for _, task := range tasks {
		tasksChan <- task
	}
	close(tasksChan)
	wg.Wait()
	fmt.Println("All files downloaded")
}
