package service

import (
	"gin-blog-hufeng/dao"
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/model/resp"
)

type Category struct{}

func (*Category) GetList(req req.PageQuery) resp.PageResult[[]resp.CategoryVo] {
	data, total := categoryDao.GetList(req)
	return resp.PageResult[[]resp.CategoryVo]{
		Total:    total,
		List:     data,
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
	}
}

func (*Category) SaveOrUpdate(req req.AddOrEditCategory) int {
	// 查询目标分类是否存在
	existByName := dao.GetOne(model)
}
