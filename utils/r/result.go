package r

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 响应结构体
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// ReturnJson 返回 JSON 数据
func ReturnJson(c *gin.Context, httpCode, code int, msg string, data any) {
	// c.Header("","") // 根据需要在发送内容头部添加其他信息
	c.JSON(httpCode, Response{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

// 语法糖函数封装

// Send 自定义 httpCode, code, data ，然后写的理由是后面都需要这种，然后简写输入参数即可。
func Send(c *gin.Context, httpCode, code int, data any) {
	ReturnJson(c, httpCode, code, GetMsg(code), data)
}

// SendCode 自动根据 code 获取 message， 且 data == nil
func SendCode(c *gin.Context, code int) {
	Send(c, http.StatusOK, code, nil)
	// 成功结果，直接不发送问题
}

func SendData(c *gin.Context, code int, data any) {
	Send(c, http.StatusOK, code, data)
}

func SuccessData(c *gin.Context, data any) {
	Send(c, http.StatusOK, OK, data)
}
