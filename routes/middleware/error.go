package middleware

import (
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
)

// ErrorRecovery recovery 项目可能出现的 panic，并使用 zap 记录相关日志
func ErrorRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			// 后面的所有内容其实都是基于defer，因为要恢复的是整个程序的panic，所以用defer
			if err := recover(); err != nil { // recover 一定是在defer中使用的
				// check for a broken connection, as it is not only really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok { // err的类型断言
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							// 这个就是检查err类型断言之后的数据是否存在一些内容在里面
							brokenPipe = true // 反正就是断联了嘛
							// 一个broken pipe 指的是客户端向服务端传输信息，但是两者断开，一般发生在服务端关闭的时候
							// connection reset by peer 指由于服务器向客户端发送数据的时候，两者断开，但是客户端依旧在尝试读取数据
							// peer 指的是客户端与服务端之间的连接
						}
					}
				}

				// 已经发现错误就要解决错误：处理panic(xxx) 的操作
				// 看他是什么类型（通过类型断言成功），进行响应消息发送
				if code, ok := err.(int); ok { // panic(code) 根据错误代码获取 msg 基于当前err存在
					r.SendCode(c, code)
				} else if msg, ok := err.(string); ok { // panic(string) 返回 string
					r.ReturnJson(c, http.StatusOK, r.FAIL, msg, nil)
				} else if e, ok := err.(error); ok { // panic(error) 发送消息
					r.ReturnJson(c, http.StatusOK, r.FAIL, e.Error(), nil)
				} else { // 其他
					r.Send(c, http.StatusOK, r.FAIL, nil)
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false) // 获取的是客户端的请求的Request 并将其转化为字节切片
				if brokenPipe {
					utils.Logger.Error(c.Request.URL.Path,
						zap.Any("error", err),                      // 输出错误信息到日志
						zap.String("request", string(httpRequest)), // 输出HTTP请求作为一个字段写入日志
					)
					// If the connection is dead, we can't write a status to it.
					if e, ok := err.(error); ok {
						_ = c.Error(e)
					} else {
						c.Abort()
						return
					}
					c.Abort()
					return
				}
				if stack {
					// 这个stack是从外面传进来的
					utils.Logger.Error("[Recover from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					utils.Logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)))
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
