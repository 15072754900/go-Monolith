package routes

import (
	"gin-blog-hufeng/config"
	"gin-blog-hufeng/dao"
	"gin-blog-hufeng/utils"
	"log"
	"net/http"
	"time"
)

// InitGlobalVariable 初始化全局变量
func InitGlobalVariable() {
	// 初始化Viper
	utils.InitViper()
	// 初始化Logger
	utils.InitLogger()
	// 初始化数据库 DB
	dao.DB = utils.InitMySQLDB()
	// 初始化 Redis
	utils.InitRedis()
	// 初始化Casbin
	utils.InitCasbin(dao.DB)
}

// BackendServer 后台服务
// g.Go() 就是做并发的客户端服务的
func BackendServer() *http.Server {
	bakePort := config.Cfg.Server.BackPort
	log.Printf("后台服务启动于 %s 端口", bakePort)
	return &http.Server{
		Addr:         bakePort,
		Handler:      BackRouter(), // 这个handler的设置必须记住
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
