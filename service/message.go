package service

import (
	"gin-blog-hufeng/dao"
	"gin-blog-hufeng/model"
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/model/resp"
	"gin-blog-hufeng/utils/r"
)

type Message struct{}

func (*Message) GetList(req req.GetMessages) resp.PageResult[[]model.Message] {
	list, total := messageDao.GetList(req)
	return resp.PageResult[[]model.Message]{
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
		Total:    total,
		List:     list,
	}
}

func (*Message) Delete(ids []int) (code int) {
	// 什么时候需要，应该是需要使用到具体的结构体构建时在***Dao里面，如果是只是删除这些可以直接用dao（通用的）的
	dao.Delete(model.Message{}, "id IN ?", ids)
	return r.OK
}

func (*Message) UpdateReview(req req.UpdateReview) (code int) {
	maps := map[string]any{"is_review": req.IsReview}

	dao.UpdatesMap(&model.Message{}, maps, "id IN ?", req.Ids) // 错误全在dao里边设置完了
	return r.OK
}
