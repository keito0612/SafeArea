package response

// DisasterItem は災害履歴アイテムを表す
type DisasterItem struct {
	ID     int64  `json:"id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Year   int    `json:"year"`
	Damage string `json:"damage"`
}

// DisasterHistoryResponse は災害履歴レスポンスを表す
type DisasterHistoryResponse struct {
	Disasters []DisasterItem `json:"disasters"`
}
