package dao

import (
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/model/resp"
)

type Category struct{}

func (*Category) GetList(req req.PageQuery) ([]resp.CategoryVo, int64) {
	var dates = make([]resp.CategoryVo, 0)
	var total int64

	db := DB.Table("category c").
		Select("c.id", "c.name", "COUNT(a.id) AS article_count", "c.created_at", "c.updated_at").
		Joins("LEFT JOIN article a ON c.id = a.category_id AND a.is_delete = 0 AND a.status = 1")

	// 条件查询
	if req.KeyWord != "" {
		db = db.Where("name LIKE ?", "%"+req.KeyWord+"%")
	}
	db.Group("c.id").
		Order("c.id DESC").
		Count(&total).
		Limit(req.PageSize).Offset(req.PageSize * (req.PageNum - 1)).
		Find(&dates)
	return dates, total
}
