package middleware

import (
	"bytes"
	"fmt"
	"gin-blog-hufeng/dao"
	"gin-blog-hufeng/model"
	"gin-blog-hufeng/model/dto"
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

var optMap = map[string]string{
	"Article":      "文章",
	"BolgInfo":     "博客信息",
	"Category":     "分类",
	"Comment":      "评论",
	"FriendLink":   "友链",
	"Menu":         "菜单",
	"Message":      "留言",
	"OperationLog": "操作日志",
	"Resource":     "资源权限",
	"Role":         "角色",
	"Tag":          "标签",
	"User":         "用户",
	"Page":         "页面",
	// "Talk": "说说",
	// "Login": "登录",
	"POST":   "新增或修改",
	"PUT":    "修改",
	"DELETE": "删除",
} // 这只是定义一个数据实例，现在要输出其中的某些内容

func GetOptString(key string) string {
	return optMap[key]
}

// CustomResponseWriter 在 gin 中获取 Response Body 内容：对 gin 的ResponseWriter 进行包装，每次往请求方向响应数据时，将响应数据返回出去
type CustomResponseWriter struct {
	gin.ResponseWriter               // 数据接口
	body               *bytes.Buffer // 响应体缓存
}

// OperationLog 记录操作日志中间件
func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 优化记录文件上传
		// 不记录 GET 请求操作记录（比较多，但是就我而言，我可以尝试记录到一个专门的文件里面看一下） 和 文件上传操作记录（请求体太长）
		if c.Request.Method != "GET" && !strings.Contains(c.Request.RequestURI, "upload") {
			blw := &CustomResponseWriter{
				body:           bytes.NewBufferString(""),
				ResponseWriter: c.Writer,
			}
			c.Writer = blw // 就是加上一个响应体缓存，方便后续加上内容
			uuid := utils.GetFromContext[string](c, "uuid")

			// 从 session 中取用户的 SessionInfo
			val := sessions.Default(c).Get(KEY_USER + uuid)
			if val == nil {
				r.Send(c, http.StatusUnauthorized, r.ERROR_TOKEN_RUNTIME, nil)
				c.Abort()
				return
			}
			var sessionInfo dto.SessionInfo
			utils.Json.Unmarshal(val.(string), &sessionInfo) // 反序列化

			reqBody, _ := io.ReadAll(c.Request.Body)                // 记录请求信息
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody)) // 为什么要这样做？
			// 这里是一个将从 io 里获取的内容转换为 NopCloser 的一个类型的内容，禁止关闭，只使用读取
			// 在一些 如 http.Request.Body返回的就是一个实现了 io.ReadCloser 接口的对象，所以有必要进行类型转换
			// 同时注意这里的 buf.Buffer 的使用，进行了读取后面的reqBody内容，继而使用NopCloser转换
			OperationLog := model.OperationLog{
				OptModule:     GetOptString(getOptResource(c.HandlerName())), // FIXME:
				OptType:       GetOptString(c.Request.Method),
				OptUrl:        c.Request.RequestURI,
				OptMethod:     c.HandlerName(),
				OptDesc:       GetOptString(c.Request.Method) + GetOptString(getOptResource(c.HandlerName())), // TODO: 优化
				RequestParam:  string(reqBody),
				RequestMethod: c.Request.Method,
				UserId:        sessionInfo.UserInfoId,
				Nickname:      sessionInfo.Nickname,
				IpAddress:     sessionInfo.IpAddress,
				IpSource:      sessionInfo.IpSource,
			}
			c.Next()
			OperationLog.ResponseData = blw.body.String() // 从缓存中获取响应体内容
			// fmt.Println("操作日志记录： ", operationLog)
			dao.Create(&OperationLog)
		} else {
			fmt.Println("执行成功 在线上报")
			c.Next()
			// 不是就往下执行，直至判断成功
		}
	}
}

// example: "gin-blog-hufeng/api/v1.(*Resource).Delete-fm" => "Resource"
func getOptResource(handlerName string) string {
	s := strings.Split(handlerName, ".")[1]
	return s[2 : len(s)-1]
}
