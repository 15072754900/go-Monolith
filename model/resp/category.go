package resp

import "time"

// CategoryVo 前后台通用
type CategoryVo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Article   int       `json:"article_count"` // 文章数量
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"update_at"`
}
