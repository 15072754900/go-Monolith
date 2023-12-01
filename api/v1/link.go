package v1

import (
	"gin-blog-hufeng/utils/r"
	"github.com/gin-gonic/gin"
)

type Link struct{}

func (*Link) GetList(c *gin.Context) {
	r.SuccessData(c, linkService.GetList())
}

func (*Link) SaveOrUpdate(c *gin.Context) {
	r.SendCode(c, linkService.SaveOrUpdate())
}

func (*Link) Delete(c *gin.Context) {
	r.SuccessData(c, linkService.Delete())
}
