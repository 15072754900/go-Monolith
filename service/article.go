package service

import (
	"gin-blog-hufeng/dao"
	"gin-blog-hufeng/model"
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/model/resp"
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type Article struct{}

// 业务传来的数据和传输出去的数据，都在对应的地方放着，明确其设计思路

func (*Article) GetList(req req.GetArts) resp.PageResult[[]resp.ArticleVo] {
	articleList, total := articleDao.GetList(req)

	// 相关数据还有书籍的点赞数量，但是这些都是在Redis中获取与保存的，所有有Redis的相关查询获取操作
	likeCountMap := utils.Redis.HGetAll(KEY_ARTICLE_LIKE_COUNT)               // 点赞数量 map
	viewCountZ := utils.Redis.ZRangeWithScores(KEY_ARTICLE_VIEW_COUNT, 0, -1) // -1表示所有
	viewCountMap := getViewCountMap(viewCountZ)                               //访问数量 map (这一步的意义和想法确定如何实现的)

	for i, article := range articleList {
		articleList[i].ViewCount = viewCountMap[article.ID]
		articleList[i].LikeCount, _ = strconv.Atoi(likeCountMap[strconv.Itoa(article.ID)]) // 进行格式转换的数据处理
	}

	return resp.PageResult[[]resp.ArticleVo]{
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
		Total:    total,
		List:     articleList}
}

// 获取点赞数量 map
func getViewCountMap(rz []redis.Z) map[int]int {
	m := make(map[int]int)
	for _, article := range rz {
		id, _ := strconv.Atoi(article.Member.(string))
		m[id] = int(article.Score)
	}
	return m
}

// SaveOrUpdate 这里删除或更新都是对名称数量更新但不是对内容更新
func (*Article) SaveOrUpdate(req req.SaveOrUpdateArt, userId int) (code int) {
	// 设置默认图片（blogConfig 中配置）
	if req.Img == "" {
		req.Img = blogInfoService.GetBlogConfig().ArticleCover // 获取默认配置的图片
	}

	article := utils.CopyProperties[model.Article](req)
	article.UserId = userId

	// 维护 [文章-分类] 关联
	category := saveArticleCategory(req)

	// 下面这代码不理解，直接添加id就行了
	if !category.IsEmpty() {
		article.CategoryId = category.ID
	}

	// 下面是更新的代码
	if article.ID == 0 {
		dao.Create(&article)
	} else {
		dao.Update(&article)
	}

	// 维护文章-标签关联
	saveArticleTag(req, article.ID)
	return r.OK
}

// 维护 文章-分类 关联
func saveArticleCategory(req req.SaveOrUpdateArt) model.Category {
	// 先查询是否存在，不存在则新建一个，然后返回，否则直接返回
	category := dao.GetOne(model.Category{}, "name = ?", req.CategoryName)
	if category.IsEmpty() && req.Status != model.DRAFT {
		category.Name = req.CategoryName
		dao.Create(&category)
	}
	return category
}

// 维护 文章-tag 关联
func saveArticleTag(req req.SaveOrUpdateArt, articleId int) {
	// 清除文章对应的标签关联
	if req.ID != 0 {
		dao.Delete(model.ArticleTag{}, "article_id = ?", req.ID)
	}
	// 遍历 req.TagNames 中传来的标签，不存在则新建
	var articleTags []model.ArticleTag
	for _, tagName := range req.TagNames {
		tag := dao.GetOne(model.Tag{}, "name = ?", tagName)
		// 标签不存在则新建
		if tag.IsEmpty() {
			tag.Name = tagName
			dao.Create(&tag)
		}
		articleTags = append(articleTags, model.ArticleTag{
			ArticleId: articleId,
			TagId:     tag.ID,
		})
	}
	dao.Create(&articleTags)
}

func (*Article) UpdateTop(req req.UpdateArtTop) (code int) {
	// 打包，格式转换，内容填充，数据库处理，返回值（那么错误处理呢）
	article := model.Article{
		Universal: model.Universal{ID: req.ID},
		IsTop:     req.IsTop,
	}
	dao.Update(&article)
	return r.OK
}
