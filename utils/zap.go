package utils

import (
	"gin-blog-hufeng/config"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger() {
	// 生成日志文件目录, 判断文件夹是否存在，存在则过，不存在则创建
	if ok, _ := PathExists(config.Cfg.Zap.Directory); !ok {
		log.Printf("create %v directory\n%", config.Cfg.Zap.Directory)
		_ = os.Mkdir(config.Cfg.Zap.Directory, os.ModePerm)
	}
	core := zapcore.NewCore(getEnCoder(), getWriterSyncer(), getLevelPriority())
	Logger = zap.New(core)

	if config.Cfg.Zap.ShowLine {
		// 获取调用的文件，函数名称，行号
		Logger = Logger.WithOptions(zap.AddCaller())
	}

	log.Println("Zap Logger 初始化成功")
}

// 编码器：如何写入日志 写给zap.New(core)的东西: Encoder
func getEnCoder() zapcore.Encoder {
	// 参考：zap.NewProductionEncoderConfig()
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     customTimeEncoder, // 自定义的时间输出格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 使用面向字节的API来节省存储分配
	}

	if config.Cfg.Zap.Format == "json" {
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig) // 在这里不是一样的吗，在后续我要删除一部分
}

// 日志输出路径：文件、控制台、双向输出
func getWriterSyncer() zapcore.WriteSyncer {
	file, _ := os.Create(config.Cfg.Zap.Directory + "/log.log") // 创建文件

	// 双向输出，一定输出到文件夹，但是不一定输出到控制台（LoginConsole）
	if config.Cfg.Zap.LogInConsole {
		fileWriter := zapcore.AddSync(file)
		consoleWriter := zapcore.AddSync(os.Stdout)                   // 输出到控制台的os.Stdout
		return zapcore.NewMultiWriteSyncer(fileWriter, consoleWriter) // NewMultiWriteSync创建新的输出方式
	}

	// 只输出到文件，上面双向输出（控制台和文件）
	return zapcore.AddSync(file)
}

// 日志输出级别
func getLevelPriority() zapcore.LevelEnabler {
	switch config.Cfg.Zap.Level {
	case "debug", "Debug":
		return zap.DebugLevel // 这里小写zap，是库zap里的
	case "info", "Info":
		return zap.InfoLevel
	case "warn", "Warn":
		return zap.WarnLevel
	case "error", "Error":
		return zap.ErrorLevel
	case "dpanic", "Depanic":
		return zap.ErrorLevel
	case "panic", "Panic":
		return zap.PanicLevel
	case "fatal", "Fatal":
		return zap.FatalLevel
	}
	return zap.InfoLevel
}

// 自定义日志输出时间格式
func customTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(config.Cfg.Zap.Prefix + t.Format("2023/04/14 - 15:38:20"))
}
