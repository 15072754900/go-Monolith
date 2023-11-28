package service

import (
	"fmt"
	"gin-blog-hufeng/dao"
	"gin-blog-hufeng/model"
	"gin-blog-hufeng/model/resp"
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
	"github.com/gin-gonic/gin"
	"strings"
)

type BlogInfo struct{}

// Cache Aside Pattern: 经典的缓存 + 数据库的读写模式

// Report 用户上报信息：统计地域信息，访问数量
func (*BlogInfo) Report(c *gin.Context) (code int) {
	ipAddress := utils.IP.GetIpAddress(c)
	userAgent := utils.IP.GetUserAgent(c)
	browser := userAgent.Name + " " + userAgent.OSVersion.String()
	os := userAgent.OS + " " + userAgent.OSVersion.String()
	uuid := utils.Encryptor.MD5(ipAddress + browser + os)
	// 当前用户没有统计过访问人数（不在 用户set 中）
	if !utils.Redis.SIsMember(KEY_UNIQUE_VISITOR_SET, uuid) {
		// 统计地域信息
		ipSource := utils.IP.GetIpSource(ipAddress)
		if ipSource != "" { // 获取到具体的位置，提取出其中的 省份
			address := strings.Split(ipSource, "|")
			province := strings.ReplaceAll(address[2], "省", "")
			utils.Redis.HIncrBy(KEY_VISIT_AREA, province, 1) // 查询要检测的值，和要进行的操作，执行自增操作
		} else {
			utils.Redis.HIncrBy(KEY_VISIT_AREA, "未知", 1)
		}
		// 访问数量 + 1
		utils.Redis.Incr(KEY_VIEW_COUNT)
		// 将当前用户记录到 用户 set
		utils.Redis.SAdd(KEY_UNIQUE_VISITOR_SET, uuid)
	}
	return r.OK
}

// GetHomeInfo 博客后台首页信息 TODO：完善首页显示
func (b *BlogInfo) GetHomeInfo() resp.BlogHomeVO {
	articleCount := dao.Count(model.Article{}, "status = ? AND is_delete = ?", 1, 0)
	userCount := dao.Count(model.UserInfo{}, "")
	messageCount := dao.Count(model.Message{}, "")
	viewCount := utils.Redis.GetInt(KEY_VIEW_COUNT) // 查询观看人数
	fmt.Println("2")

	return resp.BlogHomeVO{
		ArticleCount: articleCount,
		UserCount:    userCount,
		MessageCount: messageCount,
		ViewCount:    viewCount,
	}
}

// GetBlogConfig 获取博客设置
func (*BlogInfo) GetBlogConfig() (respVO model.BlogConfigDetail) {
	// 尝试从 Redis 中取值
	blogConfig := utils.Redis.GetVal(KEY_BLOG_CONFIG)
	// Redis 中没有值，再查数据库，查到后设置到 Redis 中
	if blogConfig == "" {
		blogConfig = dao.GetOne(model.BlogConfig{}, "id", 1).Config
		utils.Redis.Set(KEY_BLOG_CONFIG, blogConfig, 0)
	}
	// 反序列化字符串为 golang 需要的文件名 对象
	utils.Json.Unmarshal(blogConfig, &respVO)
	return respVO
}

// UpdateBlogConfig 更新博客首页信息
func (*BlogInfo) UpdateBlogConfig(reqVO model.BlogConfigDetail) (code int) {
	blogConfig := model.BlogConfig{
		Universal: model.Universal{ID: 1},
		Config:    utils.Json.Marshal(reqVO), // 序列化 golang 对象
	}
	dao.Update(&blogConfig, "config")
	utils.Redis.Del(KEY_BLOG_CONFIG) // 从 Redis 中删除旧值
	return r.OK
}

// GetAbout 获取相关
func (*BlogInfo) GetAbout() string {
	return utils.Redis.GetVal(KEY_ABOUT)
}

// UpdateAbout 更新相关
func (*BlogInfo) UpdateAbout(data model.About) (code int) {
	utils.Redis.Set(KEY_ABOUT, data.Content, 0)
	return r.OK
}
