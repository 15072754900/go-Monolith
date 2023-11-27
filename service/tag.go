package service

import (
	"gin-blog-hufeng/dao"
	"gin-blog-hufeng/model"
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/model/resp"
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
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

func (*Tag) SaveOrUpdate(req req.AddOrEditTag) int {
	existByName := dao.GetOne(model.Tag{}, "name = ?", req.Name)
	// 同名存在 && 存在的ID不等于当前要更新的ID -> 重复
	if !existByName.IsEmpty() && existByName.ID != req.ID {
		return r.ERROR_TAG_EXIST
	}
	tag := utils.CopyProperties[model.Tag](req)

	if req.ID != 0 {
		dao.Update(&tag, "name")
	} else {
		dao.Create(&tag)
	}
	return r.OK
}

func (*Tag) Delete(ids []int) (code int) {
	// 检查是否存在文章在该标签下面，存在则必须先删除文章，所以先报错
	count := dao.Count(model.Article{}, "tag_id = ?", ids)
	if count > 0 {
		return r.ERROR_TAG_ART_EXIST
	}
	dao.Delete(model.Article{}, "id in ?", ids)
	return r.OK
}

func (*Tag) GetOption() []resp.OptionVo {
	return tagDao.GetOption()
}
