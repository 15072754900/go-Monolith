package dao

import (
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/model/resp"
)

type Tag struct{}

func (*Tag) GetList(req req.PageQuery) ([]resp.TagVO, int64) {
	var datas = make([]resp.TagVO, 0)
	var total int64

	db := DB.Table("tag t").
		Select("t.id", "name", "COUNT(at.article_id) AS article_count", "t.created_at", "t.updated_at").
		Joins("LEFT JOIN article_tag at ON t.id = at.tag_id")
	if req.KeyWord != "" {
		db = db.Where("name LIKE ?", "%"+req.KeyWord+"%")
	}
	db.Group("t.id").Order("t.id DESC").
		Count(&total).
		Limit(req.PageSize).Offset(req.PageSize * (req.PageNum - 1)).
		Find(&datas)
	return datas, total
}