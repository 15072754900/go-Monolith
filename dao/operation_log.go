package dao

import (
	"gin-blog-hufeng/model"
	"gin-blog-hufeng/model/req"
)

type OperationLog struct{}

func (*OperationLog) GetList(req req.PageQuery) ([]model.OperationLog, int64) {
	list := make([]model.OperationLog, 0)
	var total int64

	db := DB.Model(&model.OperationLog{})
	if req.KeyWord != "" {
		db = db.Where("opt_module LIKE ?", "%"+req.KeyWord+"%").Or("opt_desc LIKE ?", "%"+req.KeyWord+"%")
	}
	db.Count(&total).Order("id DESC").Limit(req.PageSize).Offset(req.PageSize * (req.PageNum - 1)).Find(&list)
	return list, total
}
