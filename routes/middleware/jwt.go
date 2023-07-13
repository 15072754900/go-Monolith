package middleware

import (
	"fmt"
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

// JWTAuth JWT 中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 约定 Token 放在 Header 的 Authorization 中，并使用 Bearer 开头
		token := c.Request.Header.Get("Authorization")
		// 若 token 为空
		if token == "" {
			r.SendCode(c, r.ERROR_TOKEN_NOT_EXIST)
			c.Abort()
			return
		}

		// token 的正确格式： `Bearer [tokenString]` 这种情况应该是加上一些其他的标志
		parts := strings.Split(token, " ")
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			r.SendCode(c, r.ERROR_TOKEN_TYPE_WRONG)
			c.Abort()
			return
		}

		// parts[1] 是获取到的 tokenString，使用 JWT 解析函数解析它
		claims, err := utils.GetJWT().ParseToken(parts[1])
		fmt.Println(claims)
		// token 解析失败
		if err != nil {
			r.SendData(c, r.ERROR_TOKEN_WRONG, err.Error())
			c.Abort()
			return
		}

		// 判断token 已过期
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			fmt.Print("时间出错了")
			r.SendCode(c, r.ERROR_TOKEN_RUNTIME)
			c.Abort()
			return
		}

		// 将当前请求的相关信息保存到请求的上下文 c 中
		// 后续的处理函数可以用 c.Get("xxx") 来获取当前请求的用户信息
		c.Set("user_info_id", claims.UserId)
		c.Set("role", claims.Role)
		c.Set("uuid", claims.UUID)
		fmt.Println("登录成功 token 正确")
		c.Next()
		// 后面的都是在所有前面的中间件的内部，直至完成后面的所有函数才能一个个中间件完成，当然，无需鉴权的独立于这里的中间件之外
	}
}
