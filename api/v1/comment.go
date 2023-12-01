package v1

import (
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
	"github.com/gin-gonic/gin"
)

type Comment struct{}

func (*Comment) GetList(c *gin.Context) {
	var req = utils.BindValidQuery[req.GetComments](c)
	utils.CheckQueryPage(&req.PageSize, &req.PageNum)
	r.SuccessData(c, commentService.GetList(req))
}

// UpdateReview 必须知道需求是什么才能在函数中传递对应的参数，选取utils中的对应的提取内容与格式转换的函数
func (*Comment) UpdateReview(c *gin.Context) {
	r.SendCode(c, commentService.UpdateReview(utils.BindValidJson[req.UpdateReview](c)))
}

func (*Comment) Delete(c *gin.Context) {
	r.SendCode(c, commentService.Delete(utils.BindJson[[]int](c)))
}
