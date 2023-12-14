package dao

import (
	"gin-blog-hufeng/model"
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/model/resp"
)

type Role struct{}

// GetLabelsByUserInfoId 根绝 [userInfoId] 获取 [角色标签列表]
func (*Role) GetLabelsByUserInfoId(userInfoId int) (labels []string) {
	DB.Table("role r, user_role ur").
		Where("r.id = ur.role_id AND ur.user_id = ?", userInfoId).
		Pluck("label", &labels) // 将单列查询为切片
	return
}

func (*Role) GetList(req req.PageQuery) (list []resp.RoleVo, total int64) {
	list = make([]resp.RoleVo, 0)
	db := DB.Table("role")
	// 查询条件
	if req.KeyWord != "" {
		// 开始查询信息
		db = db.Where("name like ?", "%"+req.KeyWord+"%")
	}
	db.Count(&total).Limit(req.PageSize).Offset((req.PageNum - 1) * req.PageSize).Select("id, name, label, created_at, is_disable").Find(&list)
	return list, total
}

func (*Role) GetResourcesByRoleId(roleId int) (resourceIds []int) {
	DB.Model(&model.RoleResource{}).Where("role_id = ?", roleId).Pluck("resource_id", &resourceIds)
	return
}

func (*Role) GetMenusByRoleId(roleId int) (menuIds []int) {
	DB.Model(&model.RoleMenu{}).Where("role_id = ?", roleId).Pluck("menu_id", &menuIds)
	return
}

func (*Role) GetOption() []resp.OptionVo {
	list := make([]resp.OptionVo, 0)
	DB.Model(&model.Role{}).Select("id", "name").Find(&list)
	return list
}

func (*Role) GetLabelsByRoleIds(ids []int) (labels []string) {
	DB.Model(&model.Role{}).Where("role_id = ?", ids).Pluck("label", &labels)
	return
}
