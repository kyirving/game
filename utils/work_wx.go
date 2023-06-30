package utils

import (
	"encoding/json"
	"fmt"
	"game/config"
	"strconv"
)

var message string = `**报警通知**
- 时间：%s
- 区服ID：%s

**详情**
> %s`

var AccessToken string

func generateToken() {
	path := "/gettoken"
	api := fmt.Sprintf("%s%s", config.Config.WorkWxConf.Host, path)

	params := make(map[string]string, 2)
	params["corpid"] = strconv.Itoa(config.Config.Corpid)
	params["corpsecret"] = config.Config.Corpsecret

	result, err := SendRequest(api, "GET", params)
	if err != nil {
		fmt.Println("gettoken fail : ", err)
		return
	}
	fmt.Println(result)

	var workWxresp = WorkWxResp{}
	err = json.Unmarshal(result, &workWxresp)
	if err != nil {
		fmt.Println("json.Unmarshal fail :", err)
		return
	}

	//获取成功
	if workWxresp.Errcode == 0 {
		AccessToken = workWxresp.AccessToken
	}
}

func SendMessage(stime, server_id, msg string) {
	data := make(map[string]interface{}, 3)
	content := make(map[string]string, 1)
	content["content"] = fmt.Sprintf(message, stime, server_id, msg)

	data["touser"] = config.Config.Touser
	data["msgtype"] = "markdown"
	data["agentid"] = strconv.Itoa(config.Config.Corpid)
	data["markdown"] = content

	api := fmt.Sprintf("%s/message/send?%s", config.Config.WorkWxConf.Host, AccessToken)
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Marshal fail :", err)
		return
	}
	result, err := SendRowRequest(api, b)
	if err != nil {
		fmt.Println("SendRowRequest fail :", err)
		return
	}

	var workWxresp = WorkWxResp{}
	err = json.Unmarshal(result, &workWxresp)
	if err != nil {
		fmt.Println("json.Unmarshal fail :", err)
		return
	}

	if workWxresp.Errcode != 0 {
		codes := make(map[int]string, 2)
		codes[40014] = "不合法的access_token"
		codes[41001] = "缺少access_token参数"
		codes[42001] = "access_token已过期"

		if _, ok := codes[workWxresp.Errcode]; ok {
			AccessToken = ""
			generateToken()
			SendMessage(stime, server_id, msg)
			//token
			fmt.Println("发送消息失败，正在重新发送")
		}
		fmt.Println("发送消息失败")
	}
	fmt.Println("发送通知成功")

}
