package model

import (
	"reflect"
	"time"
)

// 权限控制 - 7 张表

// 角色

// 资源

// 角色-资源 关联

// UserRole 用户-角色 关联
type UserRole struct {
	UserId int `json:"User_id"`
	RoleId int `json:"role_id"`
}

// Menu 菜单
type Menu struct {
	Universal `mapstructure:"squash"`
	Name      string `gorm:"type:varchar(20);comment:菜单名" json:"name"`
	Path      string `gorm:"type:varchar(50);comment:菜单路径" json:"path"`
	Component string `gorm:"type:varchar(50);comment:组件" json:"component"`
	Icon      string `gorm:"type:varchar(50);comment:菜单图标" json:"icon"`
	ParentId  int    `gorm:"comment:父菜单id" json:"parent_id"`
	OrderNum  int8   `gorm:"type:tinyint;default:0;comment:显示排序" json:"order_num"`
	IsHidden  int8   `gorm:"type:tinyint(1);default:0;comment:是否隐藏(0-否 1-是)" json:"is_hidden"`
	KeepAlive int8   `gorm:"type:tinyint(1);default:1" json:"keep_alive"`
	Redirect  string `gorm:"type:varchar(50);comment:跳转路径" json:"redirect"`
}

// 角色-菜单 关联

// UserAuth 用户账户信息
type UserAuth struct {
	Universal     `mapstructure:"squash"` // 这些tag都是在编译期间完成的操作，包括那些gorm的自增，在编译期间进行数据库的语句书写
	UserInfoId    int                     `gorm:"comment:用户信息ID" json:"user_info_id"`
	Username      string                  `gorm:"type:varchar(50);comment:用户名" json:"username"`
	Password      string                  `gorm:"type:varchar(100);comment:密码" json:"password"`
	LoginType     int                     `gorm:"type:tinyint(1);comment:登录类型" json:"login_type"`
	IpAddress     string                  `gorm:"type:varchar(20);comment:登录IP地址" json:"ip_address"`
	IpSource      string                  `gorm:"type:varchar(50);comment:IP来源" json:"ip_source"`
	LastLoginTime time.Time               `gorm:"comment:上次登录时间" json:"last_login_time"`
}

func (u *UserAuth) IsEmpty() bool {
	return reflect.DeepEqual(u, &UserAuth{})
}
