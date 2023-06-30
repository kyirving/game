package main

import (
	"encoding/json"
	"fmt"
	"game/config"
	"game/utils"
	"strconv"
	"sync"
	"time"
)

var (
	wg       sync.WaitGroup
	chanMsg  chan utils.ChanMsg
	chanTask chan string
	jobChan  chan string
	reqNum   int
)

func main() {
	reqNum = 1000
	chanMsg = make(chan utils.ChanMsg, reqNum)
	//任务统计
	chanTask = make(chan string, reqNum)

	//工作通道
	jobChan = make(chan string, reqNum)
	createJobChan(reqNum)
	createPool(config.Config.PoolNum)

	wg.Add(1)
	go CheckOK(chanTask, chanMsg, reqNum)

	//发送通知协程
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go sendMsgTask(chanMsg)
	}
	wg.Wait()
	fmt.Println("All goroutines finish")
}

//创建工作池
func createPool(num int) {
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for server_id := range jobChan {
				checkGame(server_id, chanMsg, chanTask)
			}
		}()
	}
}

func createJobChan(reqNum int) {
	for i := 1; i <= reqNum; i++ {
		jobChan <- strconv.Itoa(i)
	}
	defer close(jobChan)
}

func checkGame(serverid string, chan_msg chan<- utils.ChanMsg, chan_task chan<- string) {
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
		chanMsg := &utils.ChanMsg{
			Stime:    time.Now().Format("2006-01-02 15:04:05"),
			ServerId: serverid,
			Msg:      resp.Msg,
		}
		chan_msg <- *chanMsg
	}
	// 用于监控协程知道已经完成了几个任务
	chan_task <- serverid
}

// 任务统计协程
func CheckOK(chan_task <-chan string, chanMsg chan utils.ChanMsg, reqNum int) {
	defer wg.Done()
	var count int
	for {
		server_id := <-chan_task
		fmt.Printf("%s 完成了检查任务\n", server_id)
		count++
		if count == reqNum {
			fmt.Println("检查协助已执行完毕")
			close(chanMsg)
			break
		}
	}
}

func sendMsgTask(chanMsg chan utils.ChanMsg) {
	defer wg.Done()

	for chanData := range chanMsg {
		fmt.Println("准备发送报警: ", chanData)
		// utils.SendMessage(chanData.Stime, chanData.ServerId, chanData.Msg)
	}
	fmt.Println("准备退出")
}
