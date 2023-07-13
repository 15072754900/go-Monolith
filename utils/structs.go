package utils

import (
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

// 结构体 转 map[string]any, 需要配合 `ampstructure`

// CopyProperties 拷贝属性， 一般用于 vo -> po
func CopyProperties[T any](from any) (to T) {
	if err := copier.Copy(&to, from); err != nil {
		Logger.Error("CopyProperties: ", zap.Error(err))
		panic(err)
	}
	return to
}
