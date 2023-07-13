package resp

// BlogHomeVO 后台首页 VO
type BlogHomeVO struct {
	ArticleCount int64 `json:"article_count"` // 文章数量
	UserCount    int64 `json:"user_count"`    // 用户数量
	MessageCount int64 `json:"message_count"` // 留言数量
	ViewCount    int   `json:"view_count"`    // 访问量 可能是因为数据量不会很大，直接使用小的类型？
}
