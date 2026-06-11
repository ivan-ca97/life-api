package handler

type uploadURLRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
}
