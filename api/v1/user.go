package v1

import (
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
	"github.com/gin-gonic/gin"
)

type User struct{}

// UpdateCurrent 更新用户信息，针对不同的端系统设置
func (*User) UpdateCurrent(c *gin.Context) {
	currentUser := utils.BindValidJson[req.UpdateCurrentUser](c) // 匹配获取当前用户信息
	currentUser.ID = utils.GetFromContent[int](c, "user_info_id")
	r.SendCode(c, userService.UpdateCurrent(currentUser))
}

// GetInfo 根据 Token 获取用户信息
func (*User) GetInfo(c *gin.Context) {
	r.SuccessData(c, userService.GetInfo(utils.GetFromContext[int](c, "user_info_id")))
}

func (*User) GetList(c *gin.Context) {
	r.SuccessData(c, userService.GetList)
}

// Update 更新当前用户信息，关联处理 [用户-角色] 关系
func (*User) Update(c *gin.Context) {
	r.SendCode(c, userService.Update(utils.BindJson[req.UpdateUser](c)))
}

// UpdateDisable 修改用户禁用状态
func (*User) UpdateDisable(c *gin.Context) {
	req := utils.BindValidJson[req.UpdateUserDisable](c) // 用户的信息在传输的上下文中提取出来，并进行修改或删除
	userService.UpdateDisable(req.ID, *req.IsDisable)
	r.Success(c)
}

// UpdatePassword 修改用户密码
func (*User) UpdatePassword(c *gin.Context) {
	r.SendCode(c, userService.UpdatePassword(utils.BindJson[req.UpdatePassword](c)))
}

func (*User) UpdateCurrentPassword(c *gin.Context) {
	// 这里需要用到携带的put的原始字段，并且需要和之前的比较，设计到dao和utils
	// 先获取一个参数一个类型，得到用户字段并执行类型转换，在函数里面设计
	r.SendCode(c, userService.UpdateCurrentPassword(
		utils.BindJson[req.UpdateAdminPassword](c),
		utils.GetFromContent[int](c, "user_info_id")))
}

func (*User) GetOnlineList(c *gin.Context) {
	r.SuccessData(c, userService.GetOnlineList(utils.BindPageQuery(c)))
}

func (*User) ForceOffline(c *gin.Context) {
	r.SendCode(c, userService.ForceOffline(utils.BindJson[req.ForceOfflineUser](c)))
}
