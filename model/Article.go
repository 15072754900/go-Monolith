package model

import "reflect"

const (
	PUBLIC = 1 + iota // 公开
	SECRET            // 私密
	DRAFT             // 草稿
)

// Article belongTo: 一个文章 属于 一个分类
// belongTo: 一个文章 属于 一个用户
// manyTomany: 一个文章 可以拥有 多个标签，多个文章 可以使用 一个标签
type Article struct {
	Universal

	CategoryId int `gorm:"type:bigint;not null;comment:分类 ID" json:"category_id"`

	UserId int `gorm:"type:int;not null;comment:用户 ID" json:"user_id"`

	Title       string `gorm:"type:varchar(100);not null;comment:文章标题" json:"title"`
	Desc        string `gorm:"type:varchar(200);comment:文章描述" json:"desc"`
	Content     string `gorm:"type:longtext;comment:文章内容" json:"content"`
	Img         string `gorm:"type:varchar(100);comment:封面图片地址" json:"img"`
	Type        int8   `gorm:"type:tinyint;comment:类型(1-原创 2-转载 3-翻译)" json:"type"` // tinyint 比较小的类型 小于255的整数
	Status      int8   `gorm:"type:tinyint;comment:状态(1-公开 2-私密)" json:"status"`
	IsTop       *int8  `gorm:"type:tinyint;not null;default:0;comment:是否置顶(0-否 1-是)" json:"is_top"`
	IsDelete    *int8  `gorm:"type:tinyint;not null;default:0;comment:是否放到回收站(0-否 1-是)" json:"is_delete"` // 这个在后面学习数据库的时候还要好好思考
	OriginalUrl string `gorm:"type:varchar(100);comment:源链接" json:"original_url"`
}

type ArticleTag struct {
	ArticleId int
	TagId     int
}

func (a *Article) IsEmpty() bool {
	return reflect.DeepEqual(a, &Article{})
}
