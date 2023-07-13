package utils

// viper就是设置配置（config）
import (
	"flag"
	"gin-blog-hufeng/config"
	"github.com/spf13/viper"
	"log"
	"strings"
)

// 优先级：命令行 > 默认值 (flag)
func InitViper() {
	// 根据命令行读取配置文件路径
	var configPath string
	flag.StringVar(&configPath, "c", "", "choose config file.")
	flag.Parse()
	if configPath != "" { // 说明从命令行接收到了参数
		log.Printf("命令行读取参数，配置文件路径为：%s", configPath)
	} else {
		log.Println("命令行参数为空， 默认加载：config/config.toml")
		configPath = "config/config.toml"
		// 在这里把设置的参数文件名赋给了变量
	}

	// 现在读取固定路径配置文件
	v := viper.New()
	v.SetConfigFile(configPath)
	v.AutomaticEnv()                                   // 允许使用环境变量
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // SERVER_APPMODE => SERVER.APPMODE

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		log.Panic("配置文件读取失败：", err)
	}

	// 加载配置文件内容到结构体对象
	if err := v.Unmarshal(&config.Cfg); err != nil { // 就是把配置文件configPath的内容加载到&config.Cfg上
		log.Panic("配置文件内容加载失败：", err)
	}

	// TODO: 配置文件热重载，在一些情景需要

	log.Println("配置文件内容加载成功")
}
