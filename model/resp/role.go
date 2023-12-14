package resp

import "time"

// RoleVo 角色 + 资源id列表 + 菜单id列表 -> 关联性数据库
type RoleVo struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Label       string    `json:"label"`
	CreatedAt   time.Time `json:"created_at"`
	IsDisable   int       `json:"is_disable"`
	ResourceIds []int     `json:"resource_ids" gorm:"-"`
	MenuIds     []int     `json:"menu_ids" gorm:"-"`
}
