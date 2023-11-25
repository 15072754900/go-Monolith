package routes

import (
	"gin-blog-hufeng/config"
	"gin-blog-hufeng/routes/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

// BackRouter 后台管理页面的接口路由
// 开始写业务了，这里在gin里面gin.new就是一个http.handler是一个可以r.run或者其他方式启动的函数调用体
func BackRouter() http.Handler {
	gin.SetMode(config.Cfg.Server.AppMode)

	r := gin.New()
	r.SetTrustedProxies([]string{"*"}) // 设置请求白名单

	// 使用本地文件上传， 需要静态文件服务（暂时使用本地，后续升级为nas文件存储）
	if config.Cfg.Upload.OssType == "local" {
		r.Static("/public", "./public")
		r.StaticFS("/dir", http.Dir("./public")) // 将 public 目录内的文件列表展示
	}

	r.Use(gin.Logger())
	// 使用自定义的一些中间件
	r.Use(middleware.Logger())             // 要输出的内容的日志
	r.Use(middleware.ErrorRecovery(false)) // 自定义错误处理中间件
	r.Use(middleware.Cors())               // 跨域中间件

	// 基于 cookie 存储 session 这也用作一个中间件
	// 诶嘿，但是现在不写，后面再写 写了几天发现要跑起来还是要写session
	store := cookie.NewStore([]byte(config.Cfg.Session.Salt)) // 用于校验的authentication

	// session 存储时间跟 JWT 过期时间一致
	store.Options(sessions.Options{MaxAge: int(config.Cfg.JWT.Expire) * 3600})
	r.Use(sessions.Sessions(config.Cfg.Session.Name, store)) // Session 中间件

	// 无需鉴权的接口
	base := r.Group("/api")
	{
		// TODO: 用户注册 和 后台登录 应该记录到 日志
		base.POST("/login", userAuthAPI.Login)   // 后台登录
		base.POST("/report", blogInfoAPI.Report) // 上报信息
	}

	// 需要登录鉴权的接口
	auth := base.Group("") // "/admin"
	// !注意中间件的顺序
	// 洋葱结构的实现，类似的koa
	// 中间件：鉴权、权限、监听在线、记录日志
	auth.Use(middleware.JWTAuth())      // JWT 鉴权中间件
	auth.Use(middleware.RBAC())         // casbin 权限中间件
	auth.Use(middleware.ListenOnline()) // 监听在线用户
	auth.Use(middleware.OperationLog()) // 记录请求操作日志
	{
		auth.GET("/home", blogInfoAPI.GetHomeInfo) // 后台首页信息

		// 博客设置
		setting := auth.Group("/setting")
		{
			// 这就是从数据库里面拿到一个对应模型的数据完事，或者加上一个放在redis里面的
			setting.GET("/blog-config", blogInfoAPI.GetBlogConfig)    // 获取博客设置
			setting.GET("/about", blogInfoAPI.GetAbout)               // 获取关于我
			setting.PUT("/blog-config", blogInfoAPI.UpdateBlogConfig) // 编辑博客设置
			setting.PUT("/about", blogInfoAPI.UpdateAbout)            // 编辑关于我
		}

		// 用户模块 用户的上线下线修改密码获取列表
		user := auth.Group("/user")
		{
			user.GET("/list", userAPI.GetList)                           // 获取用户列表
			user.GET("/info", userAPI.GetInfo)                           // 获取当前用户信息
			user.PUT("", userAPI.Update)                                 // 更新用户信息
			user.PUT("/disable", userAPI.UpdateDisable)                  // 修改用户禁用状态
			user.PUT("/password", userAPI.UpdatePassword)                // 修改普通用户密码，不需要原密码
			user.PUT("/current/password", userAPI.UpdateCurrentPassword) // 修改管理员用户密码，不需要原密码
			user.PUT("/current", userAPI.UpdateCurrent)                  // 修改当前用户的密码
			user.GET("/online", userAPI.GetOnlineList)                   // 获取在线用户信息
			user.DELETE("/offline", userAPI.ForceOffline)                // 强制用户下线
		}

		// 分类模块
		category := auth.Group("/category")
		{
			category.GET("/list", categoryAPI.GetList)  // 分类列表
			category.POST("", categoryAPI.SaveOrUpdate) // 新增/编辑分类
		}

		// 标签模块
		tag := auth.Group("/tag")
		{
			tag.GET("/list", tagAPI.GetList)
		}

		// 文章模块

		// 评论模块

		// 留言模块

		// 友情链接

		// 资源模块

		// 菜单模块
		menu := auth.Group("/menu")
		{
			menu.GET("/list", menuAPI.GetTreeList)
			menu.GET("/user/list", menuAPI.GetUserMenu) // 获取当前用户
		}

		// 角色模块

		// 操作日志模块

		// 页面模块
	}

	return r
}
