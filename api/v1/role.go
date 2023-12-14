package v1

import (
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
	"github.com/gin-gonic/gin"
)

type Role struct{}

func (*Role) GetTreeList(c *gin.Context) {
	r.SuccessData(c, roleService.GetTreeList(utils.BindPageQuery(c)))
}

func (*Role) SaveOrUpdate(c *gin.Context) {
	r.SendCode(c, roleService.SaveOrUpdate(utils.BindValidJson[req.SaveOrUpdateRole](c)))
}

func (*Role) Delete(c *gin.Context) {
	r.SendCode(c, roleService.Delete(utils.BindJson[[]int](c)))
}

func (*Role) GetOption(c *gin.Context) {
	r.SuccessData(c, roleService.GetOption())
}
