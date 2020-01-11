package common

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

func buildResponse(code int, msg string, data interface{}) (res []byte, err error) {
	var (
		response Response
	)
	response.Code = code
	response.Msg = msg
	response.Data = data

	res, err = json.Marshal(response)
	return
}

func WriteResponse(w http.ResponseWriter, code int, msg string, data interface{}) {
	res, err := buildResponse(code, msg, data)
	if err != nil {
		fmt.Println("对象转化成json失败：", err)
		return
	}
	w.Header().Add("Content-Type", "application/json;charset=utf-8")
	_, err = w.Write(res)
	if err != nil {
		fmt.Println("输出结果失败：", err)
	}
}
