package utils

type Resp struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data []interface{} `json:"data"`
}

type ChanMsg struct {
	Stime    string
	ServerId string
	Msg      string
}

type WorkWxResp struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
