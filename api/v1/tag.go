package v1

import (
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
	"github.com/gin-gonic/gin"
)

type Tag struct{}

// GetList 条件标签列表（后台）
func (*Tag) GetList(c *gin.Context) {
	r.SuccessData(c, tagService.GetList(utils.BindPageQuery(c)))
}

// 执行新增/编辑 标签
func (*Tag) SaveOrUpdate(c *gin.Context) {
	r.SendCode(c, tagService.SaveOrUpdate(utils.BindValidJson[req.AddOrEditTag](c)))
}

// 执行批量删除操作
func (*Tag) Delete(c *gin.Context) {
	r.SendCode(c, tagService.Delete(utils.BindJson[[]int](c))) // 每次请求与接受的c都是不同的，来自前端传来的数据
}

// 需要获取全局数据的时候就直接不穿参数

// GetOption 获取下拉框选项数据
func (*Tag) GetOption(c *gin.Context) {
	r.SuccessData(c, tagService.GetOption())
}

// GetFrontList 查询并获取标签列表（前台） 未使用到
//func (*Tag) GetFrontList(c *gin.Context) {
//	r.SuccessData(c, tagService.GetFrontList())
//}
