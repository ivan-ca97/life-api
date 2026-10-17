package handler

type uploadUrlRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
}
