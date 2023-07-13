package v1

import (
	"fmt"
	"gin-blog-hufeng/model"
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
	"github.com/gin-gonic/gin"
)

type BlogInfo struct{}

// Report 用户登进后台时上报信息：如 IP 信息
func (*BlogInfo) Report(c *gin.Context) {
	r.SendCode(c, blogInfoService.Report(c))
}

// GetHomeInfo 获取博客首页信息
func (*BlogInfo) GetHomeInfo(c *gin.Context) {
	fmt.Println("1")
	r.SuccessData(c, blogInfoService.GetHomeInfo())
}

// GetBlogConfig 获取博客首页信息
func (*BlogInfo) GetBlogConfig(c *gin.Context) {
	r.SuccessData(c, blogInfoService.GetBlogConfig())
}

// UpdateBlogConfig 更新博客配置
func (*BlogInfo) UpdateBlogConfig(c *gin.Context) {
	r.SendCode(c, blogInfoService.UpdateBlogConfig(utils.BindJson[model.BlogConfigDetail](c)))
}

// GetAbout 获取关于
func (*BlogInfo) GetAbout(c *gin.Context) {
	r.SuccessData(c, blogInfoService.GetAbout())
}

// UpdateAbout 编辑关于
func (*BlogInfo) UpdateAbout(c *gin.Context) {
	r.SuccessData(c, blogInfoService.UpdateAbout(utils.BindJson[model.About](c)))
}
