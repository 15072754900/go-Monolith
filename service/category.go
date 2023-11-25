package service

import (
	"gin-blog-hufeng/dao"
	"gin-blog-hufeng/model"
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/model/resp"
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
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
	existByName := dao.GetOne(model.Category{}, "name = ?", req.Name)
	// 同名存在 && 存在的ID不等于当前要更新的ID -> 重复
	if !existByName.IsEmpty() && existByName.ID != req.ID {
		return r.ERROR_CATE_NAME_USED
	}
	category := utils.CopyProperties[model.Category](req)

	if req.ID != 0 {
		dao.Update(&category, "name")
	} else {
		dao.Create(&category)
	}
	return r.OK
}
