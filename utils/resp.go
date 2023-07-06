package utils

// type Resp struct {
// 	Code int           `json:"code"`
// 	Msg  string        `json:"msg"`
// 	Data []interface{} `json:"data"`
// }

type Resp struct {
	RetCode int    `json:"retCode"`
	Msg     string `json:"message"`
}

type Server struct {
	ServerId   string `json:"server_id"`
	ServerName string `json:"server_name"`
	StratTime  string `json:"strat_time"`
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
