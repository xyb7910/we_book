package web

type Result struct {
	// 业务代码错误
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
