package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type User struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func main3() {
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		fmt.Println("请求失败：" + err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(fmt.Sprintf("响应错误: %d", resp.StatusCode))
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败：" + err.Error())
		return
	}

	user := User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Println("解析响应失败：" + err.Error())
		return
	}
	fmt.Printf("id: %d\ntitle: %v\nbody: %v\n", user.UserID, user.Title, user.Body)
}
