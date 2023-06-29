package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func SendRequest(api, method string, params map[string]string) (result []byte, err error) {

	values := url.Values{}
	for key, val := range params {
		values.Add(key, val)
	}
	// fmt.Println(values.)

	var req *http.Request

	switch method {
	case http.MethodPost:
		req, err = http.NewRequest(method, api, strings.NewReader(values.Encode()))
	case http.MethodGet:
		// 将参数编码到 URL 中
		urlWithParams := api + "?" + values.Encode()
		req, err = http.NewRequest(method, urlWithParams, nil)
	}
	if err != nil {
		fmt.Println("NewRequest fail : ", err)
		return
	}

	//todo 一定要设置请求头 否则服务端获取不到请求参数
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("client.Do fail : ", err)
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll fail : ", err)
		return
	}
	return b, nil
}
