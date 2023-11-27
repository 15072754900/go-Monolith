package resp

import (
	"gin-blog-hufeng/model"
	"time"
)

// 记录文章相关信息，一些操作，包含前端根据这些内容设计界面展示
// 这就是具体的大体量内容的展示了

// ArticleVo 文章列表 前后端一起制定
type ArticleVo struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	CategoryId int            `json:"category_id"`
	Category   model.Category `gorm:"foreignkey:CategoryId" json:"category"`
	Tags       []model.Tag    `gorm:"many2many:article_tag;joinForeignKey:article_id" json:"tags"`

	UserId      int    `json:"user_id"`
	Title       string `json:"title"`
	Desc        string `json:"desc"`
	Content     string `json:"content"`
	Img         string `json:"img"`
	Type        int8   `json:"type"`
	Status      int8   `json:"status"`
	IsTop       *int8  `json:"is_top"`
	IsDelete    *int8  `json:"is_delete"`
	OriginalUrl string `json:"original_url"`

	LikeCount int `json:"like_count"`
	ViewCount int `json:"view_count"`
}
