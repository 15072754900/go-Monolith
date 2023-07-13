package service

import (
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/model/resp"
)

type Tag struct{}

// GetList 分页查询标签
func (*Tag) GetList(req req.PageQuery) resp.PageResult[[]resp.TagVO] {
	data, total := tagDao.GetList(req)
	return resp.PageResult[[]resp.TagVO]{
		Total:    total,
		List:     data,
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
	}
}
