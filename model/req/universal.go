package req

// KeywordQuery 关键字查询
type KeywordQuery struct {
	Keyword string `json:"keyword"`
}

// PageQuery 获取数据（需要分页）
type PageQuery struct {
	PageSize int    `form:"page_size"`
	PageNum  int    `form:"page_num"`
	KeyWord  string `form:"keyword"`
}

// UpdateReview 修改审核 （批量）
type UpdateReview struct {
	Ids      []int `json:"ids"`
	IsReview *int8 `json:"is_review" validate:"required,min=0,max=1"`
}
