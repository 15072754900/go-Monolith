package dao

import (
	"gin-blog-hufeng/model"
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/model/resp"
)

type Article struct{}

func (*Article) GetList(req req.GetArts) ([]resp.ArticleVo, int64) {
	var list = make([]resp.ArticleVo, 0)
	var total int64

	db := DB.Model(model.Article{})

	// 文章标题
	if req.Title != "" {
		db = db.Where("title LIKE ?", "%"+req.Title+"%")
	}

	// 查询内容多少在于req中是否有相关的请求标签，现在使用where逐项添加并检查

	// 文章标题，是否删除（？），状态，分类，类型。这些个检查条件
	if req.Title != "" {
		db = db.Where("title Like ?", "%"+req.Title+"%")
	}
	if req.IsDelete != nil {
		db = db.Where("is_delete", req.IsDelete)
	}
	if req.Status != 0 {
		db = db.Where("status", req.Status)
	}
	if req.CategoryId != 0 {
		db = db.Where("category_id", req.CategoryId)
	}
	if req.Type != 0 {
		db = db.Where("type", req.Type)
	}

	db.Preload("Category").Preload("Tags").Joins("LEFT JOIN article_tag ON article.id = article_tag.article_id").Group("id") // 去重
	// group的用法

	if req.TagId != 0 {
		db = db.Where("tag_id", req.TagId)
	}

	db.Count(&total).Limit(req.PageSize).Offset(req.PageSize * (req.PageNum - 1)).Order("is_top DESC, article.id DESC").Find(&list)

	return list, total
}
