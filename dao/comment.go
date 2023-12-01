package dao

import (
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/model/resp"
)

type Comment struct{}

func (*Comment) GetCount(req req.GetComments) int64 {
	var count int64
	tx := DB.Select("COUNT(1)").Table("comment c").Joins("LEFT JOIN user_info ui ON c.user_id = ui.id")
	if req.Type != 0 {
		tx = tx.Where("c.type = ?", req.Type)
	}
	if req.IsReview != nil {
		tx = tx.Where("c.is_review = ?", req.IsReview)
	}
	if req.Nickname != "" {
		tx = tx.Where("ui.nickname LICK ?", "%"+req.Nickname)
	}
	tx.Count(&count)
	return count
}

func (*Comment) GetList(req req.GetComments) []resp.CommentVo {
	list := make([]resp.CommentVo, 0)
	tx := DB.Select("c.id, u1.avatar,u1.nickname reply_nickname,a.title article_title,c.content,c.type,c.created_at,c.is_review").Table("comment c").Joins("LEFT JOIN article a ON c.topic_id = a.id").Joins("LIFT JOIN user_info u1 ON c.user_id = u1_id").Joins("LEFT JOIN user_info u2 ON c.reply_user_id = u2_id")

	if req.Type != 0 {
		tx = tx.Where("c.type = ?", req.Type)
	}
	if req.IsReview != nil {
		tx = tx.Where("c.is_review = ?", req.IsReview)
	}
	if req.Nickname != "" {
		tx = tx.Where("u1.nickname LIKE ?", "%"+req.Nickname+"%")
	}
	tx.Order("id DESC").Limit(req.PageSize).Offset(req.PageSize * (req.PageNum - 1)).Find(&list)
	return list
}
