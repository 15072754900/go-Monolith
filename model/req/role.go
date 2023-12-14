package req

type SaveOrUpdateRole struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Label       string `json:"label"`
	IsDisable   int    `json:"is_disable"`
	ResourceIds []int  `json:"resource_ids"` //*
	MenuIds     []int  `json:"menu_ids"`     //*
}
