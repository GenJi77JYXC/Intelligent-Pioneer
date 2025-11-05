package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 定义了所有 API 的标准响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// --- 响应码常量 ---
const (
	SuccessCode = 20000 // 与前端模板约定的成功码
	ErrorCode   = 50000 // 通用错误码
)

// Result 是一个辅助函数，用于快速构建和发送 JSON 响应
func Result(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// Success 快速发送一个成功的响应
func Success(c *gin.Context, data interface{}) {
	Result(c, SuccessCode, "success", data)
}

// Error 快速发送一个失败的响应
func Error(c *gin.Context, code int, msg string) {
	Result(c, code, msg, nil)
}

// ParamError 快速发送一个参数错误的响应
func ParamError(c *gin.Context, msg string) {
	Result(c, http.StatusBadRequest, "Invalid request parameters: "+msg, nil)
}
