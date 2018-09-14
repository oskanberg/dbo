package model

// Film deliberately omits daily gross
type Film struct {
	ID    string  `json:"id"`
	BomID *string `json:"bomID"`
	Title *string `json:"title"`
}
