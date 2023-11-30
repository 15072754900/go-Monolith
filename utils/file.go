package utils

import (
	"bufio"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"os"
)

var File fileUtils

type fileUtils struct{}

// WriteFile 将文件写入到本地：已有同名文件则覆盖：写入文件这种工作需要时时回忆
func (*fileUtils) WriteFile(name, path, content string) {
	//  指定模式打开文件： 读写|创建
	file, err := os.OpenFile(path+name, os.O_RDWR|os.O_CREATE, 0666) // 0666指的是编辑的文件的权限控制，0表示8进制，666表示都可以编辑/读取的权限
	if err != nil {
		Logger.Error("文件写入， 目标地址错误：", zap.Error(err))
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)      // 文件写入对象
	_, err = writer.WriteString(content) // 将内容写到缓存
	writer.Flush()                       // 将缓存写到文件
	if err != nil {
		Logger.Error("文件写入失败: ", zap.Error(err))
	}
}

func (*fileUtils) ReadFromFileHeader(file *multipart.FileHeader) string {
	open, err := file.Open()
	if err != nil {
		Logger.Error("文件读取，目标地址错误： ", zap.Error(err))
		return ""
	}
	defer open.Close()
	all, err := io.ReadAll(open)
	if err != nil {
		Logger.Error("文件读取失败：", zap.Error(err))
		return ""
	}
	return string(all)
}
