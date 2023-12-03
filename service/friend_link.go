package service

import (
	"gin-blog-hufeng/dao"
	"gin-blog-hufeng/model"
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/model/resp"
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
)

type FriendLink struct{}

func (*FriendLink) GetList(req req.PageQuery) (data resp.PageResult[[]model.FriendLink]) {
	list, total := friendLinkDao.GetList(req)
	return resp.PageResult[[]model.FriendLink]{
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
		List:     list,
		Total:    total,
	}
}

func (*FriendLink) SaveOrUpdate(req req.SaveOrUpdateLink) (code int) {
	link := utils.CopyProperties[model.FriendLink](req) // vo->po
	// 将数据转换到另一种类型上
	if link.ID != 0 {
		dao.Update(&link)
	} else {
		dao.Create(&link)
	}
	return r.OK
}

func (*FriendLink) Delete(ids []int) (code int) {
	dao.Delete(ids, "id IN ?", ids)
	return r.OK
}
