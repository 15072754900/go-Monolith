package service

import (
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
