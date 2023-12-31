package v1

import (
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
	"github.com/gin-gonic/gin"
)

type Resource struct{}

func (*Resource) GetTreeList(c *gin.Context) {
	r.SuccessData(c, resourceService.GetTreeList(utils.BindQuery[req.PageQuery](c)))
}

func (*Resource) SaveOrUpdate(c *gin.Context) {
	r.SendCode(c, resourceService.SaveOrUpdate(utils.BindJson[req.SaveOrUpdateResource](c)))
}

func (*Resource) Delete(c *gin.Context) {
	r.SendCode(c, resourceService.Delete(utils.GetIntParam(c, "id")))
}

// UpdateAnonymous 修改资源允许匿名访问
func (*Resource) UpdateAnonymous(c *gin.Context) {
	r.SendCode(c, resourceService.UpdateAnonymous(utils.BindValidJson[req.UpdateAnonymous](c)))
}

func (*Resource) GetOption(c *gin.Context) {
	r.SuccessData(c, resourceService.GetOptionList())
}
