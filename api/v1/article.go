package v1

import (
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
	"github.com/gin-gonic/gin"
)

type Article struct{}

// 共性：都是面向一个对象的函数，都是接受一个参数的函数
// 处理相较于user，首先处理对象不同，面对的*gin.Context里的参数也不一样的需要，但是行为都是一致的，获取数据的获取数据，处理数据删除更新的执行其行为

func (*Article) GetList(c *gin.Context) {
	r.SuccessData(c, articleService.GetList(utils.BindQuery[req.GetArts](c)))
}

func (*Article) SaveOrUpdate(c *gin.Context) {

}

func (*Article) UpdateTop(c *gin.Context) {

}

func (*Article) GetInfo(c *gin.Context) {

}

func (*Article) SoftDelete(c *gin.Context) {

}

func (*Article) Delete(c *gin.Context) {

}

func (*Article) Export(c *gin.Context) {

}

func (*Article) Import(c *gin.Context) {

}
