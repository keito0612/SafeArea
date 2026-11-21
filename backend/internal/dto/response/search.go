package response

// Coordinates は緯度経度を表す
type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// RiskItem は個別のリスク情報を表す
type RiskItem struct {
	Type        string `json:"type"`
	Level       string `json:"level"`
	Description string `json:"description"`
}

// SearchResponse は検索結果レスポンスを表す
type SearchResponse struct {
	Address              string      `json:"address"`
	Coordinates          Coordinates `json:"coordinates"`
	OverallRisk          string      `json:"overall_risk"`
	Risks                []RiskItem  `json:"risks"`
	Tags                 []string    `json:"tags"`
	DisasterHistoryCount int         `json:"disaster_history_count"`
}
