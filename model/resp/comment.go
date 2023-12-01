package resp

import "time"

// CommentVo 后台评论 VO 评论有人物，有内容，有昵称图像等等
type CommentVo struct {
	ID            int       `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	Avatar        string    `json:"avatar"`
	Nickname      string    `json:"nickname"`
	ReplyNickname string    `json:"reply_nickname"`
	ArticleTitle  string    `json:"article_title"`
	Content       string    `json:"content"`
	Type          int       `json:"type"`
	IsReview      int       `json:"is_review"`
}
