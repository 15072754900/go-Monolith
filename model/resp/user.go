package resp

import "time"

// 后台列表 VO
type UserVO struct {
}

// LoginVO 登录 VO
type LoginVO struct {
	ID         int `json:"id"`
	UserInfoId int `json:"user_info_id"`

	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Intro    string `json:"intro"`
	Website  string `json:"website"`

	IpAddress     string    `json:"ip_address"`
	IpSource      string    `json:"ip_source"`
	LastLoginTime time.Time `json:"last_login_time"`
	LoginType     int       `json:"login_type"`

	// 点赞 Set: 用于记录用户点赞过的文章，评论
	ArticleLikeSet []string `json:"article_like_set"`
	CommentLikeSet []string `json:"comment_like_set"`
	// TalkLikeSet

	Token string `json:"token"`
}

// 在线用户

// UserInfoVO 用户信息 VO
type UserInfoVO struct {
	ID       int    `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Intro    string `json:"intro"`
	Website  string `json:"website"`

	ArticleLikeSet []string `json:"article_like_set"`
	CommentLikeSet []string `json:"comment_like_set"`
}
