package main

import (
	"encoding/json"
	"fmt"
	_ "game/config"
	"game/utils"
	"strconv"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

func main() {
	chanMsg := make(chan utils.ChanMsg, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go check(strconv.Itoa(i), chanMsg)
	}

	go func() {
		for chanData := range chanMsg {
			fmt.Println("准备发送报警")
			utils.SendMessage(chanData.Stime, chanData.ServerId, chanData.Msg)
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

func check(serverid string, chan_msg chan<- utils.ChanMsg) {
	defer wg.Done()

	params := make(map[string]string, 1)
	params["status"] = serverid

	api := "http://127.0.0.1:8000"

	result, err := utils.SendRequest(api, "POST", params)
	if err != nil {
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
		chanMsg := &utils.ChanMsg{
			Stime:    time.Now().Format("2006-01-02 15:04:05"),
			ServerId: serverid,
			Msg:      resp.Msg,
		}
		chan_msg <- *chanMsg
	}
}
