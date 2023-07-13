package v1

import (
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
	"github.com/gin-gonic/gin"
)

type UserAuth struct{}

// Login 登录
func (*UserAuth) Login(c *gin.Context) {
	loginReq := utils.BindValidJson[req.Login](c)
	// BindValidJson属于是在定义时把一个数据类型绑定进行判断，但是在这里输入要验证的格式
	loginVo, code := userService.Login(c, loginReq.Username, loginReq.Password)

	r.SendData(c, code, loginVo)
}
