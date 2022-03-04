package model

import "strconv"

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Tag  string      `json:"tag,omitempty"`
	Data interface{} `json:"data"`
}

func NewSuccessResponse(tag string, data interface{}) *Response {
	return &Response{Code: Success, Msg: "success", Tag: tag, Data: data}
}

func NewBindFailedResponse(tag string) *Response {
	return &Response{Code: WrongArgs, Msg: "wrong argument", Tag: tag}
}

func NewNoPermissionResponse() *Response {
	return &Response{Code: NoPermission, Msg: "No permission", Tag: strconv.Itoa(NoPermission), Data: nil}
}

func NewErrorResponse(tag string, code int, err error) *Response {
	return &Response{Code: code, Msg: err.Error(), Tag: tag, Data: nil}
}

const (
	Success           = 0  // 成功
	Failed            = -1 // 失败
	Error             = -2 // 一般性错误
	Unknown           = -3 // 未知错误
	WrongArgs         = -4 // 收到错误的参数
	UserOrPasswordErr = -5 // 用户名或密码错误
	NoPermission      = -6 // 没有权限
	UnsupportedAction = -7 // 不支持输入的操作

)
