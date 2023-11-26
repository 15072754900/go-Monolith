package model

import "reflect"

type Tag struct {
	Universal
	Name     string    `gorm:"type:varchar(20);not null;comment:分类名称" json:"name"`
	Articles []Article `gorm:"many2many:article_tag;" json:"articles"`
}

func (c *Tag) IsEmpty() bool {
	return reflect.DeepEqual(c, &Tag{})
}
