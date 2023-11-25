package req

type Login struct {
	Username string `json:"username" validate:"required" label:"用户名"`
	Password string `json:"password" validate:"required,min=4,max=20" label="密码"`
	Code     string `json:"code" validate:"required" label:"邮箱验证码"`
}

// GetUsers 查询用户 （后台）
type GetUsers struct {
	PageQuery
	LoginType int8   `json:"login_type"`
	Nickname  string `json:"nickname"`
}

// UpdateUser 更新用户（后台）
type UpdateUser struct {
	UserInfoId int    `json:"id"`
	Nickname   string `json:"nickname"`
	RoleIds    []int  `json:"role_ids"`
}

// UpdateUserDisable 修改用户禁用状态
type UpdateUserDisable struct {
	ID        int  `json:"id"`
	IsDisable *int `json:"is_disable" validate:"required,min=0,max=1"`
}

type UpdatePassword struct {
	Username string `json:"username" validate:"required" label:"用户名"`
	Password string `json:"password" validate:"required" label:"密码"`
}

type UpdateAdminPassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type UpdateCurrentUser struct {
	ID       int    `json:"id"`
	Nickname string `json:"nickname" validate:"required"`
	Avatar   string `json:"avatar"`
	Intro    string `json:"intro"`
	Website  string `json:"website"`
	Email    string `json:"email"`
}

// ForceOfflineUser 强制下线需要用户信息计算其uuid
type ForceOfflineUser struct {
	UserIndoId int    `json:"user_info_id"`
	IpAddress  string `json:"ip_address"`
	Browser    string `json:"browser"`
	OS         string `json:"os"`
}
