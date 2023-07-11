package game

import (
	"encoding/json"
	"fmt"
	"game/config"
	"game/utils"
	"math"
	"sync"
	"time"
)

var (
	wg         sync.WaitGroup
	chanMsg    chan utils.ChanMsg
	chanTask   chan int
	jobChan    chan []utils.Server
	reqNum     int
	msgMap     sync.Map
	batchCount int
	totalReq   int
)

func Run() {

	if isExit() {
		fmt.Println("维护中!!!")
		return
	}

	batchCount = 10
	server_list, err := getServer()
	if err != nil {
		fmt.Println("获取区服列表失败")
		return
	}

	//区服监测数
	reqNum = int(math.Ceil(float64(len(server_list)) / float64(batchCount)))
	// num := math.Ceil(float64(len(server_list) / batchCount))

	//区服监测响应信息
	chanMsg = make(chan utils.ChanMsg, reqNum)
	//任务统计
	chanTask = make(chan int, reqNum)
	//工作通道
	jobChan = make(chan []utils.Server, reqNum)
	createJobChan(server_list, batchCount)
	//线程池发送监测请求
	createPool(config.Config.GameConf.GameId, config.Config.GameConf.PtId, config.Config.PoolNum)

	wg.Add(1)
	go CheckOK(chanTask, chanMsg, reqNum)

	//发送通知协程
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go sendMsgTask(chanMsg)
	}

	wg.Wait()
	fmt.Println("All goroutines finish")
	fmt.Println("totalReq :", totalReq)
	// 遍历
	content := "**【盟重英雄冰雪单职业（244）】区服异常通知**\n"
	msgMap.Range(func(key, value interface{}) bool {
		if ChanMsg, ok := value.(utils.ChanMsg); ok {
			content += fmt.Sprintf("> server_id: %s \n> 异常信息 ： %s \n\n", ChanMsg.ServerId, ChanMsg.Msg)
		}
		return true
	})

	if content != "" {
		fmt.Println(content)
		fmt.Println("准备发送报警")
		is_send := utils.SendMessage(content)
		if !is_send {
			fmt.Println("企业微信报警失败！！！")
		}
	}
}

//创建工作池
func createPool(game_id, ptid, num int) {
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for servers := range jobChan {
				totalReq++
				start := time.Now()
				checkGame(game_id, ptid, servers, chanMsg, chanTask)
				elapsed := time.Since(start)
				fmt.Printf("request-runtime：%s\n", elapsed)
			}
		}()
	}
}

//往工作通道写任务供工作池使用
func createJobChan(server_list []utils.Server, batchCount int) {
	for i := 0; i < len(server_list); i += batchCount {
		end := i + batchCount
		if end > len(server_list) {
			end = len(server_list)
		}
		batch := server_list[i:end]
		jobChan <- batch

	}
	defer close(jobChan)
}

func checkGame(game_id, ptid int, servers []utils.Server, chan_msg chan<- utils.ChanMsg, chan_task chan<- int) {

	var batchParams []map[string]int
	for _, v := range servers {
		batchParams = append(batchParams, map[string]int{"Ptid": ptid, "GameId": game_id, "trueZoneId": v.ServerId})
	}

	result, err := utils.SendBatchRequest(config.Config.GameConf.ServerStatus, "POST", batchParams)
	if err != nil {
		// 用于监控协程知道已经完成了几个任务
		chan_task <- 1
		return
	}

	resp := []utils.Resp{}
	err = json.Unmarshal(result, &resp)
	if err != nil {
		chan_task <- 1
		fmt.Println("json.Unmarshal fail :", err)
		return
	}

	// fmt.Println(resp)
	for _, res := range resp {
		if res.RetCode != 0 {
			chanMsg := utils.ChanMsg{
				Stime:    time.Now().Format("2006-01-02 15:04:05"),
				ServerId: res.ServerId,
				Msg:      res.Msg,
			}
			fmt.Println("servers", servers)
			fmt.Println("chanMsg", chanMsg)
			chan_msg <- chanMsg
		}
	}
	chan_task <- 1
}

// 任务统计协程
func CheckOK(chan_task <-chan int, chanMsg chan utils.ChanMsg, reqNum int) {
	defer wg.Done()
	var count int
	for {
		<-chan_task
		// server_id := <-chan_task
		// fmt.Printf("%d 完成了检查任务\n", server_id)
		count++
		if count == reqNum {
			fmt.Println("监测协程已执行完毕")
			close(chanMsg)
			break
		}
	}
}

func sendMsgTask(chanMsg chan utils.ChanMsg) {
	defer wg.Done()

	for chanData := range chanMsg {
		msgMap.Store(chanData.ServerId, chanData)
	}
	fmt.Println("发送通知协程已退出")
}

//获取区服列表
func getServer() (server_list []utils.Server, err error) {
	server_list = []utils.Server{}
	resp, err := utils.SendRequest(config.Config.GameConf.ServerList, "GET", nil)
	if err != nil {
		fmt.Println("utils.SendRequest fail :", err)
		return
	}

	err = json.Unmarshal(resp, &server_list)
	if err != nil {
		fmt.Println("json.Unmarshal fail :", err)
	}
	return
}

func isExit() bool {
	now := time.Now()
	weekDay := int(now.Weekday())

	if weekDay == 2 {
		stime := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())
		etime := time.Date(now.Year(), now.Month(), now.Day(), 11, 30, 0, 0, now.Location())
		if now.After(stime) && now.Before(etime) {
			return true
		}
	}

	return false
}
