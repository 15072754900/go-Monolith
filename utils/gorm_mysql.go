package utils

import (
	"fmt"
	"gin-blog-hufeng/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

func InitMySQLDB() *gorm.DB { // 这个DB是一个在gorm里面定义的结构体，区别于下面的那个是已经配置过了的DB数据实例的指针
	mysqlCfg := config.Cfg.Mysql // 获取对MySQL的配置
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlCfg.Username,
		mysqlCfg.Password,
		mysqlCfg.Host,
		mysqlCfg.Port,
		mysqlCfg.Dbname)

	DB, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		// gorm 日志模式
		Logger: logger.Default.LogMode(getLogMode(config.Cfg.Mysql.LogMode)), // 这里的logger是gorm的日志模式，不是自定义中间件的logger内容
		// 禁用外联约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 禁止默认事务（提高运行速度）
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			// 使用单数表名，启动该选项，此时，'User' 的表名应该是`user`
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatal("MySQL 连接失败， 请检查参数")
	}

	log.Println("MySQL 连接成功")
	// 注意，这里存在问题，就是log.println和log.fatal是不一样的，fatal会在输出内容之后exit(1)退出，导致后面的代码不能执行
	// log.println和fmt 无很大区别，就是加上了日志格式

	// 迁移数据表，此处无需要不执行
	// autoMigrate(DB)

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(10)                  // 设置连接池中最大闲置连接
	sqlDB.SetMaxOpenConns(100)                 // 设置数据库的最大连接数量
	sqlDB.SetConnMaxLifetime(10 * time.Second) // 设置连接的最大可复用时间

	return DB
}

// 根据字符串获取对应 LogLevel
func getLogMode(str string) logger.LogLevel {
	switch str {
	case "silent", "Silent":
		return logger.Silent
	case "error", "Error":
		return logger.Error
	case "warn", "Warn":
		return logger.Warn
	case "info", "Info":
		return logger.Info
	default:
		return logger.Info
	}
}

// 迁移数据表，不建议执行，但是可以参考内容（在后续学习中要看）
