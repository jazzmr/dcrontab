package common

import "encoding/json"

type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func BuildResponse(code int, msg string, data interface{}) (res []byte, err error) {
	var (
		response Response
	)
	response.Code = code
	response.Msg = msg
	response.Data = data

	res, err = json.Marshal(response)
	return
}
