package v1

import (
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
	"github.com/gin-gonic/gin"
)

type OperationLog struct{}

func (*OperationLog) GetList(c *gin.Context) {
	r.SuccessData(c, operationLogService.GetList(utils.BindPageQuery(c)))
}

func (*OperationLog) Delete(c *gin.Context) {
	r.SendCode(c, operationLogService.Delete(utils.BindJson[[]int](c)))
}
