package v1

import (
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
	"github.com/gin-gonic/gin"
)

type Category struct{}

// GetList 条件查询列表（后台）
func (*Category) GetList(c *gin.Context) {
	r.SuccessData(c, categoryService.GetList(utils.BindPageQuery(c)))
}

func (*Category) SaveOrUpdate(c *gin.Context) {
	r.SendCode(c, categoryService.SaveOrUpdate(utils.BindValidJson[req.AddOrEditCategory](c)))
}
