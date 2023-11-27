package service

import (
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/model/resp"
)

type Article struct{}

// 业务传来的数据和传输出去的数据，都在对应的地方放着，明确其设计思路

func (*Article) GetList(req req.GetArts) resp.PageResult[[]resp.ArticleVo] {

}
