package utils

import (
	"errors"
	"os"
)

// 判断文件目录是否存在
func PathExists(path string) (bool, error) {
	fileinfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, err
		}
	}

	// 表示存在文件夹
	if fileinfo.IsDir() {
		return true, nil
	}
	return false, errors.New("存在同名文件")
}
