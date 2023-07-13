package model

import "time"

// 包含逻辑删除的模型（作者感觉不好用，所以他就注释了）
// type DelUniversal struct {...}

// Universal 不包含逻辑删除的模型
type Universal struct {
	ID        int       `gorm:"primary_key;auto_increment" json:"id" mapstruct:"id"`
	CreatedAt time.Time `json:"created_at" mapstructure:"-"` // mapstructure 是让在改数据指定为其他类型数据结构不会被映射
	UpdatedAt time.Time `json:"updated_at" mapstructure:"-"` // json 指的是在序列化为JSON类型是该字段名称为"id"
}

// tag 元数据
