package v1

import (
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
	"github.com/gin-gonic/gin"
)

type Message struct{}

func (*Message) GetList(c *gin.Context) {
	var req = utils.BindValidQuery[req.GetMessages](c)
	utils.CheckQueryPage(&req.PageNum, &req.PageSize)
	r.SuccessData(c, messageService.GetList(req))
}

func (*Message) Delete(c *gin.Context) {
	r.SendCode(c, messageService.Delete(utils.BindJson[[]int](c)))
}

func (*Message) UpdateReview(c *gin.Context) {
	r.SendCode(c, messageService.UpdateReview(utils.BindValidJson[req.UpdateReview](c)))
}
