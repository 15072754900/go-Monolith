package v1

import (
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
	"github.com/gin-gonic/gin"
)

type Tag struct{}

// GetList 条件标签列表（后台）
func (*Tag) GetList(c *gin.Context) {
	r.SuccessData(c, tagService.GetList(utils.BindPageQuery(c)))
}
