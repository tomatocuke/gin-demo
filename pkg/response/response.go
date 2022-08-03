package response

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Echo(ctx *gin.Context, data interface{}, err error) {
	switch x := err.(type) {
	// 成功
	case nil:
		Success(ctx, data)
	// 失败：可预知，前端执行动作
	case Action:
		Fail(ctx, x.Code, x.Msg)
	// 拒绝：可预知，前端显示信息
	case String:
		Fail(ctx, -1, x.Error())
	// 错误：不可预知。
	default:
		Error(ctx, err)
	}
}

func Success(ctx *gin.Context, data interface{}) {
	resp := Response{
		Code: 0,
		Msg:  "ok",
		Data: data,
	}
	done(ctx, resp)
}

func Fail(ctx *gin.Context, code int, msg string) {
	resp := Response{
		Code: code,
		Msg:  msg,
		Data: struct{}{},
	}
	done(ctx, resp)
}

var (
	errorHook func(errorMsg string)
)

func AddErrorHook(fn func(string)) {
	errorHook = fn
}

func Error(ctx *gin.Context, err error) {
	if errorHook != nil {
		errorHook(err.Error())
	}
	resp := Response{
		Code: -1,
		Msg:  err.Error(),
		Data: struct{}{},
	}
	done(ctx, resp)
}

func done(ctx *gin.Context, resp Response) {
	byteData, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	ctx.Abort()
	ctx.Data(http.StatusOK, "application/json", byteData)

	ctx.Set("response", byteData)
}