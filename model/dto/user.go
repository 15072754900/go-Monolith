package dto

import "gin-blog-hufeng/model/resp"

// SessionInfo Session 信息：记录用户详细信息 + 是否被强退
type SessionInfo struct {
	UserDetailDTO
	IsOffline int `json:"is_offline"`
}

// UserDetailDTO 用户详细信息：仅用于在后端内部进行传输
type UserDetailDTO struct {
	// 为什么仅在后端传输
	resp.LoginVO
	Password   string   `json:"password"`
	IsDisable  int8     `json:"is_disable"` // 我记得int8是rune的相同类型
	Browser    string   `json:"browser"`
	OS         string   `json:"os"`
	RoleLabels []string `json:"role_labels"`
}
