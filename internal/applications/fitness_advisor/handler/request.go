package handler

type estimateRequest struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}
