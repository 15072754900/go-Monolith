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
