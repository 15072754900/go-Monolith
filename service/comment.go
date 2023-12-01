package service

import (
	"gin-blog-hufeng/dao"
	"gin-blog-hufeng/model"
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/model/resp"
	"gin-blog-hufeng/utils/r"
)

type Comment struct{}

func (*Comment) GetList(req req.GetComments) resp.PageResult[[]resp.CommentVo] {
	var list = make([]resp.CommentVo, 0) // 都是创建一个容器，后面添加，并在最后返回一个包含该变量内容的一个结构实例

	total := commentDao.GetCount(req)
	if total != 0 {
		list = commentDao.GetList(req)
	}

	return resp.PageResult[[]resp.CommentVo]{
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
		Total:    total,
		List:     list,
	}
}

// UpdateReview 修改评论审核，就是是否能够过审并显示
func (*Comment) UpdateReview(req req.UpdateReview) (code int) {
	maps := map[string]any{"is_review": req.IsReview}
	dao.UpdatesMap(&model.Comment{}, maps, "id IN ?", req.Ids)
	return r.OK
}

func (*Comment) Delete(ids []int) (code int) {
	dao.Delete(model.Comment{}, "id IN ?", ids)
	return r.OK
}
