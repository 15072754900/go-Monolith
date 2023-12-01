package req

type GetMessages struct {
	PageQuery
	Nickname string `form:"nickname"`
	IsReview *int8  `form:"is_review"`
}
