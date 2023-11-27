package service

import (
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/model/resp"
	"gin-blog-hufeng/utils"
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
