package game

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
	chanTask chan int
	jobChan  chan utils.Server
	reqNum   int
	msgMap   sync.Map
)

func Run() {
	server_list, err := getServer()
	if err != nil {
		fmt.Println("获取区服列表失败")
		return
	}

	//区服监测数
	reqNum = len(server_list)
	//区服监测响应信息
	chanMsg = make(chan utils.ChanMsg, reqNum)
	//任务统计
	chanTask = make(chan int, reqNum)
	//工作通道
	jobChan = make(chan utils.Server, reqNum)
	createJobChan(server_list)
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
	// 遍历
	// content := "**区服异常通知【测试】**\n"
	content := ""
	msgMap.Range(func(key, value interface{}) bool {
		if ChanMsg, ok := value.(utils.ChanMsg); ok {
			content += fmt.Sprintf("> server_id: %d \n> 异常信息 ： %s \n\n", ChanMsg.ServerId, ChanMsg.Msg)
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
			for server := range jobChan {
				checkGame(game_id, ptid, server, chanMsg, chanTask)
			}
		}()
	}
}

//往工作通道写任务供工作池使用
func createJobChan(server_list []utils.Server) {
	for _, server := range server_list {
		jobChan <- server
	}
	defer close(jobChan)
}

func checkGame(game_id, ptid int, server utils.Server, chan_msg chan<- utils.ChanMsg, chan_task chan<- int) {
	params := make(map[string]string, 2)
	params["trueZoneId"] = strconv.Itoa(server.ServerId)
	params["GameId"] = strconv.Itoa(game_id)
	params["PtId"] = strconv.Itoa(ptid)

	result, err := utils.SendRequest(config.Config.GameConf.ServerStatus, "GET", params)
	if err != nil {
		// 用于监控协程知道已经完成了几个任务
		chan_task <- server.ServerId
		return
	}

	resp := &utils.Resp{}
	err = json.Unmarshal(result, &resp)
	if err != nil {
		chan_task <- server.ServerId
		fmt.Println("json.Unmarshal fail :", err)
		return
	}

	if resp.RetCode != 0 {
		chanMsg := &utils.ChanMsg{
			Stime:    time.Now().Format("2006-01-02 15:04:05"),
			ServerId: server.ServerId,
			Msg:      resp.Msg,
		}
		chan_msg <- *chanMsg
	}
	chan_task <- server.ServerId
}

// 任务统计协程
func CheckOK(chan_task <-chan int, chanMsg chan utils.ChanMsg, reqNum int) {
	defer wg.Done()
	var count int
	for {
		server_id := <-chan_task
		fmt.Printf("%d 完成了检查任务\n", server_id)
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
