package middleware

import (
	"gin-blog-hufeng/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// Logger 日志中间件 // 中间件函数都是handlerFUNC 是因为handler是服务的名字，那么服务的函数就是加func
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next() // 从这里开始，后面的就是整个业务层跑起来直至完成的时间，再开始结束计时
		cost := time.Since(start)
		utils.Logger.Info(c.Request.URL.Path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			// zap.String("path", c.Request.URL.Path), // 为什么这里不给出来，还是说这里其实是动态中间件的范畴
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}
