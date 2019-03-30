package common

import "encoding/json"

type Job struct {
	JobName  string `json:"job_name"`
	ShellCmd string `json:"shell_cmd"`
	CronExpr string `json:"cron_expr"`
}

// HTTP接口应答
type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

// 应答方法
func BuildResponse(errno int, msg string, data interface{}) (resp []byte, err error) {
	// 1, 定义一个response
	var (
		response Response
	)

	response.Errno = errno
	response.Msg = msg
	response.Data = data

	// 2, 序列化json
	resp, err = json.Marshal(response)
	return
}
