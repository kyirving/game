package utils

import (
	"encoding/json"
	"fmt"
	"game/config"
	"strconv"
)

func GenerateToken() {
	path := "/gettoken"
	api := fmt.Sprintf("%s%s", config.Config.Host, path)

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
		fmt.Println("发送消息")
	}

}

func SengMsg() {
	
}
