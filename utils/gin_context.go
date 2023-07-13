package utils

import (
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/utils/r"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// Validate 参数合法性校验
func Validate(c *gin.Context, data any) {
	validMsg := Validator.Validate(data)
	if validMsg != "" {
		r.ReturnJson(c, http.StatusOK, r.ERROR_INVALID_PARAM, validMsg, nil)
		panic(nil)
	}
}

// BindJson Json 绑定
func BindJson[T any](c *gin.Context) (data T) {
	if err := c.ShouldBindJSON(&data); err != nil {
		Logger.Error("BindJson", zap.Error(err))
		panic(r.ERROR_REQUEST_PARAM)
	}
	return
}

func BindValidJson[T any](c *gin.Context) (data T) {
	// Json 绑定
	if err := c.ShouldBindJSON(&data); err != nil { // login 的数据绑定
		Logger.Error("BindJson", zap.Error(err))
		panic(r.ERROR_REQUEST_PARAM) // 根据之前写的，panic错误应该在recover里回复并输出到日志里
	}
	// 参数合法性校验
	//Validate(c, &data)
	return data
}

// GetFromContext 从 gin Context 上获取值，该值是 JWT middleware 解析 Token 后设置的
// 如果该值不存在，说明 Token 有问题
func GetFromContext[T any](c *gin.Context, key string) T {
	val, exist := c.Get(key)
	if !exist {
		panic(r.ERROR_TOKEN_RUNTIME)
	}
	return val.(T)
}

// BindPageQuery Params 分页绑定（处理了 PageSize 和 PageQuery）
func BindPageQuery(c *gin.Context) (data req.PageQuery) {
	if err := c.ShouldBindQuery(&data); err != nil {
		Logger.Error("BindQuery", zap.Error(err))
		panic(r.ERROR_REQUEST_PARAM)
	}
	// 检查分页参数
	CheckQueryPage(&data.PageSize, &data.PageNum)
	return
}

// CheckQueryPage 检查分页参数
func CheckQueryPage(pageSize, pageNum *int) {
	switch {
	case *pageSize >= 100:
		*pageSize = 100
	case *pageSize <= 0:
		*pageSize = 10
	}
	if *pageNum <= 0 {
		*pageNum = 1
	}
}

// BindQuery Param 绑定
func BindQuery[T any](c *gin.Context) (data T) {
	if err := c.ShouldBindQuery(&data); err != nil {
		Logger.Error("BindQuery", zap.Error(err))
		panic(r.ERROR_REQUEST_PARAM)
	}
	return
}

// GetFromContent 从 Gin Context 上获取值，该值是 JWT middleware 解析 Token 后设置的
// 如果该值不存在，说明 Token 有问题
func GetFromContent[T any](c *gin.Context, key string) T {
	val, exist := c.Get(key)
	if !exist {
		panic(r.ERROR_TOKEN_RUNTIME)
	}
	return val.(T)
}
