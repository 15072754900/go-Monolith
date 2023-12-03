package dao

import (
	"gin-blog-hufeng/model"
	"gin-blog-hufeng/model/req"
)

type FriendLink struct{}

func (*FriendLink) GetList(req req.PageQuery) (list []model.FriendLink, total int64) {
	// 执行从数据库中获取内容的操作
	tx := DB.Model(&model.FriendLink{})
	if req.KeyWord != "" {
		tx = tx.Where("name LIKE ?", "%"+req.KeyWord+"%")
	}
	tx.Count(&total)
	err := tx.Limit(req.PageSize).Offset(req.PageSize * (req.PageNum - 1)).Find(&list).Error
	if err != nil {
		panic(err)
	}
	return
}
