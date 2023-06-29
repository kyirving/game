package main

import (
	"encoding/json"
	"fmt"
	_ "game/config"
	"game/utils"
	"strconv"
	"sync"
)

var wg = sync.WaitGroup{}

func main() {
	respChan := make(chan utils.Resp, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go check(strconv.Itoa(i), respChan)
	}

	go func() {
		for resp := range respChan {
			fmt.Println("准备发送报警")
			fmt.Println(resp.Code)
		}
	}()

	wg.Wait()
	fmt.Println("All goroutines completed")

	// 此方案也能实现，但需要所有线程执行完毕执行，可能出现通道内数据过多
	// for {
	// 	select {
	// 	case resp := <-respChan:
	// 		fmt.Println("准备发送报警")
	// 		fmt.Println(resp.Code)
	// 	case <-time.After(3 * time.Second):
	// 		fmt.Println("Timeout occurred")
	// 		return
	// 	default:

	// 		close(respChan)
	// 		return
	// 	}
	// }

	// 此代码会阻塞，可新起通道配合任务监测，关闭respChan通道退出程序，但并不优雅
	// for v := range respChan {
	// 	fmt.Println(v)
	// }
}

func check(status string, resp_chan chan<- utils.Resp) {
	defer wg.Done()

	params := make(map[string]string, 1)
	params["status"] = status

	api := "http://127.0.0.1:8000"

	result, err := utils.SendRequest(api, "POST", params)
	if err != nil {
		fmt.Println("utils.SendRequest fail :", err)
		return
	}

	resp := &utils.Resp{}
	err = json.Unmarshal(result, &resp)
	if err != nil {
		fmt.Println("json.Unmarshal fail :", err)
		return
	}

	if resp.Code != 200 {
		fmt.Println("发送数据")
		resp_chan <- *resp
	}
}
