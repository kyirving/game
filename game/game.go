package game

import (
	"encoding/json"
	"fmt"
	"game/config"
	"game/utils"
	"sync"
	"time"
)

var (
	wg       sync.WaitGroup
	chanMsg  chan utils.ChanMsg
	chanTask chan string
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

	server_list = server_list[:10]

	//区服监测数
	reqNum = len(server_list)
	//区服监测响应信息
	chanMsg = make(chan utils.ChanMsg, reqNum)
	//任务统计
	chanTask = make(chan string, reqNum)
	//工作通道
	jobChan = make(chan utils.Server, reqNum)
	createJobChan(server_list)
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
	// 遍历
	content := ""
	msgMap.Range(func(key, value interface{}) bool {

		if ChanMsg, ok := value.(utils.ChanMsg); ok {
			
		}
		return true
	})

}

//创建工作池
func createPool(num int) {
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for server := range jobChan {
				checkGame(server, chanMsg, chanTask)
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

func checkGame(server utils.Server, chan_msg chan<- utils.ChanMsg, chan_task chan<- string) {
	params := make(map[string]string, 2)
	params["trueZoneId"] = server.ServerId

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

	if resp.RetCode != 1 {
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
		msgMap.Store(chanData.ServerId, chanData)
		// utils.SendMessage(chanData.Stime, chanData.ServerId, chanData.Msg)
	}
	fmt.Println("准备退出")
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
