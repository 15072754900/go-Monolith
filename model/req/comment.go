package req

type GetComments struct {
	PageQuery
	Nickname string `json:"nickname"`
	IsReview *int8  `form:"is_review"`
	Type     int    `form:"type"`
}
