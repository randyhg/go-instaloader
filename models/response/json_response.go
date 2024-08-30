package response

import "github.com/kataras/iris/v12"

const (
	ERROR   = 7
	SUCCESS = 0
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type JsonResult struct {
	IsEncrypted bool        `json:"isEncrypted"`
	ErrorCode   int         `json:"errorCode"`
	Message     string      `json:"message"`
	Success     bool        `json:"success"`
	Data        interface{} `json:"data"`
}

type Paging struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type JsonPageResult struct {
	IsEncrypted bool        `json:"isEncrypted"`
	ErrorCode   int         `json:"errorCode"`
	Message     string      `json:"message"`
	Success     bool        `json:"success"`
	Data        interface{} `json:"data"`
	Page        *Paging     `json:"page"`
}

func Result(code int, data interface{}, msg string, ctx iris.Context) {
	ctx.JSON(JsonResult{
		ErrorCode: code,
		Message:   msg,
		Success:   true,
		Data:      data,
	})
}

func PageResult(code int, data interface{}, msg string, paging *Paging, ctx iris.Context) {
	ctx.JSON(JsonPageResult{
		ErrorCode: code,
		Message:   msg,
		Success:   true,
		Data:      data,
		Page:      paging,
	})
}

func FailResult(code int, msg string, ctx iris.Context) {
	ctx.JSON(JsonResult{
		ErrorCode: code,
		Message:   msg,
		Success:   false,
		Data:      nil,
	})
}

func OkWithMessageV2(message string, data interface{}, ctx iris.Context) {
	Result(SUCCESS, data, message, ctx)
}

func OkWithPagination(message string, data interface{}, paging *Paging, ctx iris.Context) {
	PageResult(SUCCESS, data, message, paging, ctx)
}

func FailWithMessageV2(message string, ctx iris.Context) {
	FailResult(ERROR, message, ctx)
}

func JsonSuccess() *JsonResult {
	return &JsonResult{
		ErrorCode: 0,
		Data:      nil,
		Success:   true,
	}
}

func JsonData(data interface{}) *JsonResult {
	return &JsonResult{
		ErrorCode: 0,
		Message:   "",
		Data:      data,
		Success:   true,
	}
}

func JsonErrorMsg(message string) *JsonResult {
	return &JsonResult{
		ErrorCode: 0,
		Message:   message,
		Data:      nil,
		Success:   false,
	}
}

func JsonErrorCode(code int, message string) *JsonResult {
	return &JsonResult{
		ErrorCode: code,
		Message:   message,
		Data:      nil,
		Success:   false,
	}
}

func JsonErrorData(code int, message string, data interface{}) *JsonResult {
	return &JsonResult{
		ErrorCode: code,
		Message:   message,
		Data:      data,
		Success:   false,
	}
}
