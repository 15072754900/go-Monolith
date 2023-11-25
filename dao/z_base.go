package dao

import (
	"errors"
	"gorm.io/gorm"
)

// DB 数据库指针
var DB *gorm.DB

// 通用CRUD

// Create 创建数据（单条创建和批量创建）
func Create[T any](data *T) {
	err := DB.Create(&data).Error
	if err != nil {
		panic(err)
	}
}

// GetOne [单条]数据查询
func GetOne[T any](data T, query string, args ...any) T {
	err := DB.Where(query, args...).First(&data).Error         // find the first data type that query type of args
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) { // 判断是不是对应的错误
		panic(err)
	}
	return data
}

// Update [单行]更新：传入对应结构体[传递主键用]和 带有对应更新字段值的[结构体]，结构体不断更新零值
func Update[T any](data *T, slt ...string) {
	if len(slt) > 0 {
		DB.Model(&data).Select(slt).Updates(&data)
		return
	}
	err := DB.Model(&data).Updates(&data).Error
	if err != nil {
		panic(err)
	}
}

// List 数据列表
func List[T any](data T, slt, order, query string, args ...any) T {
	db := DB.Model(&data).Select(slt).Order(order)
	if query != "" {
		db = db.Where(query, args...)
	}
	if err := db.Find(&data).Error; err != nil {
		panic(err)
	}
	return data
}

// Count 统计数量
func Count[T any](data T, query string, args ...any) int64 { // 在类型前面加...是将类型转换为可变（数量）参数slice
	var total int64
	db := DB.Model(data)
	if query != "" {
		db = db.Where(query, args...) // 在参数后面加...是将参数slice中的数据解放出来
	}
	if err := db.Count(&total).Error; err != nil {
		panic(err)
	}
	return total
}

// Delete 批量删除数据，通过条件控制可以删除单条数据
func Delete[T any](data T, query string, args ...any) {
	err := DB.Where(query, args...).Delete(&data).Error
	if err != nil {
		panic(err)
	}
}

// UpdatesMap 批量更新：map的字段就是要更新的字段（map可以更新零值），通过条件可以实现单行更新
func UpdatesMap[T any](data *T, maps map[string]any, query string, args ...any) {
	err := DB.Model(&data).Where(query, args...).Updates(maps).Error
	if err != nil {
		panic(err)
	}
}
