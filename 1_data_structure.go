package main

import (
	"fmt"
	"sync"
)

func main1() {
	// bool
	var isTrue bool = true
	fmt.Printf("isTrue: %v\n", isTrue)

	// string
	var strDescription string = "desc"
	fmt.Printf("strDescription: %v\n", strDescription)

	// int
	var iNumber int64 = 64
	fmt.Printf("iNumber: %v\n", iNumber)

	// float
	var fPrice float64 = 3.14
	fmt.Printf("fPrice: %v\n", fPrice)

	// slice
	// 类型在左侧，花括号内传入初始值
	var arrSlice []string = []string{"str1", "str2", "str3"}
	fmt.Printf("arrSlice: %v\n", arrSlice)

	// map
	// 花括号内用冒号分隔key和value
	var mUsers map[int]string = map[int]string{
		1: "user1",
		2: "user2",
		3: "user3",
	}
	fmt.Printf("mUsers: %v\n", mUsers)

	// struct
	// 花括号内再套花括号来初始化struct
	type User struct {
		ID   int
		Name string
	}
	var users []User = []User{
		{ID: 1, Name: "user1"},
		{ID: 2, Name: "user2"},
		{ID: 3, Name: "user3"},
	}
	fmt.Printf("users: %v\n", users)

	var user interface{} = User{1, "user1"}
	fmt.Printf("user: %v\n", user)

	// channel: FIFO
	// defer: FILO
	// 缓冲区大小设置需要平衡内存和性能（空间换时间）
	var userChannel chan User = make(chan User, 3) // 缓冲区满3个数据，阻塞生产者，直到被取走
	var wg sync.WaitGroup

	// 启动 3 个并发消费者
	for i := 0; i < 3; i++ {
		wg.Add(1) // 计数+1，等待1个go routine
		go func() {
			defer wg.Done() // 计数-1，本go routine已经完成
			for u := range userChannel {
				fmt.Printf("userChannel: %v\n", u)
			}
		}()
	}

	// 发送数据
	for i := 0; i < 10; i++ {
		userChannel <- User{i + 1, fmt.Sprintf("user%d", i+1)}
	}
	close(userChannel)
	wg.Wait() // 阻塞，直到计数为0
}
