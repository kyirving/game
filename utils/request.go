package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
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
		var url string
		if values.Encode() != "" {
			url = api + "?" + values.Encode()
		} else {
			url = api
		}
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		fmt.Println("NewRequest fail : ", err)
		return
	}

	//todo 一定要设置请求头 否则服务端获取不到请求参数
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 创建一个 HTTP 客户端
	client := &http.Client{
		Timeout: 20 * time.Second, // 设置超时时间为 10 秒
	}
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

func SendBatchRequest(api, method string, params []map[string]int) (result []byte, err error) {

	values := url.Values{}
	for key, val := range params {
		values.Set(fmt.Sprintf("param[%d][Ptid]", key), strconv.Itoa(val["Ptid"]))
		values.Set(fmt.Sprintf("param[%d][GameId]", key), strconv.Itoa(val["GameId"]))
		values.Set(fmt.Sprintf("param[%d][trueZoneId]", key), strconv.Itoa(val["trueZoneId"]))
	}
	var req *http.Request

	switch method {
	case http.MethodPost:
		req, err = http.NewRequest(method, api, strings.NewReader(values.Encode()))
	case http.MethodGet:
		// 将参数编码到 URL 中
		var url string
		if values.Encode() != "" {
			url = api + "?" + values.Encode()
		} else {
			url = api
		}
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		fmt.Println("NewRequest fail : ", err)
		return
	}

	//todo 一定要设置请求头 否则服务端获取不到请求参数
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// 创建一个 HTTP 客户端
	client := &http.Client{
		Timeout: 20 * time.Second, // 设置超时时间为 10 秒
	}
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

func SendRowRequest(api string, data []byte) (result []byte, err error) {

	req, err := http.NewRequest(http.MethodPost, api, bytes.NewReader(data))
	if err != nil {
		fmt.Println("NewRequest fail2 : ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// 创建一个 HTTP 客户端
	client := &http.Client{
		Timeout: 20 * time.Second, // 设置超时时间为 10 秒
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("client.Do fail2 : ", err)
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll fail2 : ", err)
		return
	}
	return b, nil
}
