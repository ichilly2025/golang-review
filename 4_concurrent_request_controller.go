package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// Task 定义任务接口
type Task func(ctx context.Context) (interface{}, error)

// Result 任务执行结果
type Result struct {
	Index int         // 任务索引
	Data  interface{} // 返回数据
	Err   error       // 错误信息
}

// ConcurrentController 并发请求控制器
type ConcurrentController struct {
	maxConcurrency int           // 最大并发数
	timeout        time.Duration // 超时时间
}

// NewConcurrentController 创建并发控制器
func NewConcurrentController(maxConcurrency int, timeout time.Duration) *ConcurrentController {
	return &ConcurrentController{
		maxConcurrency: maxConcurrency,
		timeout:        timeout,
	}
}

// Execute 执行任务列表
func (c *ConcurrentController) Execute(ctx context.Context, tasks []Task) []Result {
	taskCount := len(tasks)
	results := make([]Result, taskCount)
	taskChan := make(chan int, taskCount)
	resultChan := make(chan Result, taskCount)
	wg := sync.WaitGroup{}

	// 创建超时 context
	timeoutCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// 启动 worker pool
	for i := 0; i < c.maxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for taskIndex := range taskChan {
				select {
				case <-timeoutCtx.Done():
					// 超时或取消
					resultChan <- Result{
						Index: taskIndex,
						Data:  nil,
						Err:   timeoutCtx.Err(),
					}
					return
				default:
					// 执行任务
					data, err := tasks[taskIndex](timeoutCtx)
					resultChan <- Result{
						Index: taskIndex,
						Data:  data,
						Err:   err,
					}
				}
			}
		}()
	}

	// 分发任务
	go func() {
		for i := 0; i < taskCount; i++ {
			select {
			case <-timeoutCtx.Done():
				close(taskChan)
				return
			case taskChan <- i:
			}
		}
		close(taskChan)
	}()

	// 等待所有 worker 完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	for result := range resultChan {
		results[result.Index] = result
	}

	return results
}

// Post 文章数据结构
type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// httpRequest 发起 HTTP GET 请求
func httpRequest(ctx context.Context, url string) (interface{}, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("响应错误: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var post Post
	err = json.Unmarshal(body, &post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func main() {
	// 真实的 HTTP 请求任务
	urls := []string{
		"https://jsonplaceholder.typicode.com/posts/1",
		"https://jsonplaceholder.typicode.com/posts/2",
		"https://jsonplaceholder.typicode.com/posts/3",
		"https://jsonplaceholder.typicode.com/posts/4",
		"https://jsonplaceholder.typicode.com/posts/5",
	}

	tasks := make([]Task, len(urls))
	for i, url := range urls {
		currentURL := url // 避免闭包问题
		tasks[i] = func(ctx context.Context) (interface{}, error) {
			return httpRequest(ctx, currentURL)
		}
	}

	// 创建控制器：最大并发数 2，超时 10 秒
	controller := NewConcurrentController(2, 10*time.Second)

	// 创建可取消的 context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("开始执行任务...")
	startTime := time.Now()

	// 执行任务
	results := controller.Execute(ctx, tasks)

	fmt.Printf("任务执行完成，耗时: %v\n\n", time.Since(startTime))

	// 打印结果
	for _, result := range results {
		if result.Err != nil {
			fmt.Printf("Task %d 失败: %v\n", result.Index, result.Err)
		} else {
			post := result.Data.(Post)
			fmt.Printf("Task %d 成功: ID=%d, Title=%s\n", result.Index, post.ID, post.Title)
		}
	}
}
