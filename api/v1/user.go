package v1

import (
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
	"github.com/gin-gonic/gin"
)

type User struct{}

// 更新用户信息
func (*User) UpdateCurrent(c *gin.Context) {

}

// GetInfo 根据 Token 获取用户信息
func (*User) GetInfo(c *gin.Context) {
	r.SuccessData(c, userService.GetInfo(utils.GetFromContext[int](c, "user_info_id")))
}

func (*User) GetList(c *gin.Context) {
	r.SuccessData(c, userService.GetList)
}
