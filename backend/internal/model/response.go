package model

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`    // 状态码，0 或 200 表示成功
	Data    interface{} `json:"data"`    // 数据内容
	Message string      `json:"message"` // 消息提示
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Data:    data,
		Message: "success",
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Data:    nil,
		Message: message,
	})
}

// InternalError 内部服务器错误
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}

// BadRequest 参数错误
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}
