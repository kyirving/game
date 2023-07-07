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

func SendMessage(message string) bool {
	data := make(map[string]interface{}, 3)
	content := make(map[string]string, 1)
	content["content"] = message

	data["msgtype"] = "markdown"
	data["markdown"] = content

	api := fmt.Sprintf("%s/webhook/send?key=%s", config.Config.WorkWxConf.Host, config.Config.WorkWxConf.WebhookKey)
	fmt.Println(api)
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Marshal fail :", err)
		return false
	}
	result, err := SendRowRequest(api, b)
	if err != nil {
		fmt.Println("SendRowRequest fail :", err)
		return false
	}

	var workWxresp = WorkWxResp{}
	err = json.Unmarshal(result, &workWxresp)
	if err != nil {
		fmt.Println("json.Unmarshal fail :", err)
		return false
	}

	if workWxresp.Errcode != 0 {
		fmt.Println("发送消息失败 :", workWxresp)
		return false
	}
	fmt.Println("发送通知成功")
	return true
}
