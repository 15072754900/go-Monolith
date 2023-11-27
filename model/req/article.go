package req

// 获取文章ID，获取ID对应的文章内容这些都是不同函数执行的结果

// GetArts 条件获取文章，在文章标题中存在的内容
type GetArts struct {
	PageQuery
	Title      string `form:"title"`
	CategoryId int    `form:"category_id"`
	TagId      int    `form:"tag_id"`
	Type       int8   `form:"type" validate:"required,min=1,max=3"`
	Status     int8   `form:"status" validate:"required,min=1,max=3"`
	IsDelete   *int8  `form:"is_delete" validate:"required,min=1,max=3"`
}
