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

// SaveOrUpdate 保存文章需要知道作者信息，这里采用从用户token值获取信息
func (*Article) SaveOrUpdate(c *gin.Context) {
	r.SendCode(c, articleService.SaveOrUpdate(
		utils.BindValidJson[req.SaveOrUpdateArt](c),
		utils.GetFromContent[int](c, "user_info_id"),
	))
}

// UpdateTop 修改置顶信息
func (*Article) UpdateTop(c *gin.Context) {
	r.SendCode(c, articleService.UpdateTop(utils.BindValidJson[req.UpdateArtTop](c))) // 信息还是要验证是否准确
}

// GetInfo 获取文章详细信息
func (*Article) GetInfo(c *gin.Context) {
	r.SuccessData(c, articleService.GetInfo(utils.GetIntParam(c, "id")))
}

// SoftDelete 这个与delete的区别在于用put方法，以及...
func (*Article) SoftDelete(c *gin.Context) {
	req := utils.BindValidJson[req.SoftDelete](c)
	r.SendCode(c, articleService.SoftDelete(req.Ids, req.IsDelete))
}

func (*Article) Delete(c *gin.Context) {
	r.SendCode(c, articleService.Delete(utils.BindJson[[]int](c)))
}

func (*Article) Export(c *gin.Context) {
	r.SuccessData(c, articleService.export)
}

func (*Article) Import(c *gin.Context) {

}
