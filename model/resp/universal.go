package resp

// PageResult 分页响应结果
type PageResult[T any] struct {
	PageSize int   `json:"pageSize"`
	PageNum  int   `json:"pageNum"`
	Total    int64 `json:"total"`
	List     T     `json:"pageData"` // ！ 作者提示注意这里的别名
}

type OptionVo struct {
	ID   int    `json:"value"`
	Name string `json:"label"`
}
