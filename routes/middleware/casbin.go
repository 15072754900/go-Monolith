package middleware

import (
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
	"github.com/gin-gonic/gin"
)

// RBAC casbin 鉴权中间件
func RBAC() gin.HandlerFunc {
	utils.Casbin.LoadPolicy()
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		url, method := c.FullPath()[4:], c.Request.Method
		// 使用 casbin 自带的函数执行策略（规则）验证
		// fmt.Println("======>casbin 权限管理：", role, url, method)
		isPass, err := utils.Casbin.Enforcer().Enforce(role, url, method)
		// 权限验证未通过
		if err != nil || !isPass {
			r.SendCode(c, r.ERROR_PERMI_DENIED)
			c.Abort()
			// 这里的Abort执行，意味着后续的处理函数都不执行，现在在中间件里面，就是不去执行业务层而退出
			return
		} else {
			c.Next()
			// 权限通过就是继续套娃
		}
	}
}
