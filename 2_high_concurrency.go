package main

import (
	"fmt"
	"sync"
)

func main2() {
	// 带同步等待的并发工作池
	//workerPool3()
	// 分发与聚合的并发工作池
	//fanOutFanIn()
	// 多阶段任务的并发工作池
	pipeline()
}

// stage2: number * 2
// stage3: number + 1
func pipeline() {
	wg2 := sync.WaitGroup{}
	wg3 := sync.WaitGroup{}
	taskNumber := 5
	stage2Workers := 3 // number*2
	stage3Workers := 2 // number+1
	tasks := make(chan int, taskNumber)
	stage2Results := make(chan int, taskNumber)
	stage3Results := make(chan int, taskNumber)

	// stage2: number * 2
	for i := 0; i < stage2Workers; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			for task := range tasks {
				stage2Results <- task * 2
			}
		}()
	}

	// stage3: number + 1
	for i := 0; i < stage3Workers; i++ {
		wg3.Add(1)
		go func() {
			defer wg3.Done()
			for task := range stage2Results {
				stage3Results <- task + 1
			}
		}()
	}

	// product tasks
	for i := 0; i < taskNumber; i++ {
		tasks <- i
	}

	// close tasks
	close(tasks)

	// wait and close stage2 results
	wg2.Wait()
	close(stage2Results)

	// wait and close stage3 results
	wg3.Wait()
	close(stage3Results)

	// print stage3 results
	for result := range stage3Results {
		fmt.Println(result)
	}
}

// 分发与聚合
func fanOutFanIn() {
	wg := sync.WaitGroup{}
	taskNumber := 5
	workerNumber := 3
	tasks := make(chan int, taskNumber)
	results := make(chan int, taskNumber)

	// fan out
	for i := 0; i < workerNumber; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				results <- task
			}
		}()
	}

	for i := 0; i < taskNumber; i++ {
		tasks <- i
	}

	close(tasks)
	wg.Wait()
	close(results)

	// fan in
	total := 0
	for result := range results {
		total += result
	}
	fmt.Printf("total: %d\n", total)
}

// 用WaitGroup进行同步：Add, Done, Wait
// 先消费Task，后生产Task
func workerPool3() {
	wg := sync.WaitGroup{}
	taskNumber := 5
	workerNumber := 3
	tasks := make(chan int, taskNumber)

	for i := 0; i < workerNumber; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				fmt.Printf("%d is running.\n", task)
			}
		}()
	}

	for i := 0; i < taskNumber; i++ {
		tasks <- i
	}

	close(tasks)
	wg.Wait()
}
