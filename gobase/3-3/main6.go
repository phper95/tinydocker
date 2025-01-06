package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

// DownloadInfo 结构体保存每个下载任务的信息
type DownloadInfo struct {
	URL      string
	FilePath string
}

// downloadFile 下载单个文件
func downloadFile(url, filePath string) error {
	// 创建 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch URL %s: %w", url, err)
	}
	defer resp.Body.Close()

	// 创建目标文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	// 将 HTTP 响应写入文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}

	fmt.Printf("Downloaded: %s -> %s\n", url, filePath)
	return nil
}

// worker 处理单个下载任务
func worker(id int, tasks <-chan DownloadInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		fmt.Printf("[Worker %d] Downloading %s\n", id, task.URL)
		err := downloadFile(task.URL, task.FilePath)
		if err != nil {
			fmt.Printf("[Worker %d] Error downloading %s: %s\n", id, task.URL, err)
		}
	}
}

func main() {
	// 下载任务列表
	tasks := []DownloadInfo{
		{"https://mirrors.aliyun.com/pypi/packages/00/00/0188b746eefaea75d665b450c9165451a66aae541e5f73db4456eebc0289/loginhelper-0.0.5-py3-none-any.whl", "file1.whl"},
		{"https://mirrors.aliyun.com/pypi/packages/00/00/03734c91f740fb5f63a587110fa259d11a6e8f5b814ffe3069ef4837ecf2/aws_solutions_constructs.core-2.26.0-py3-none-any.whl?spm=a2c6h.25603864.0.0.5cbf5f26Gof6y0&file=aws_solutions_constructs.core-2.26.0-py3-none-any.whl", "file2.whl"},
		{"https://mirrors.aliyun.com/pypi/packages/00/00/041fb20ac0d321e42a953f67f3492227407d64fdba0f3f647eea7787b69a/file.io-cli-1.0.4.tar.gz?spm=a2c6h.25603864.0.0.70443731514iLB&file=file.io-cli-1.0.4.tar.gz", "file3.tag.gz"},
	}

	// 创建带缓冲的任务通道和 WaitGroup
	taskChan := make(chan DownloadInfo, len(tasks))
	var wg sync.WaitGroup

	// 限制 Goroutine 数量
	workerCount := 3
	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go worker(i, taskChan, &wg)
	}

	// 将任务发送到通道
	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	// 等待所有任务完成
	wg.Wait()
	fmt.Println("All downloads completed!")
}
