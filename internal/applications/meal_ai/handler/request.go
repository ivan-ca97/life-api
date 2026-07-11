package handler

type correctionRequest struct {
	Item       string `json:"item"`
	Correction string `json:"correction"`
}

type estimateMealRequest struct {
	PhotoUrls         []string            `json:"photo_urls"`
	Instructions      string              `json:"instructions"`
	AssumeOnlyVisible bool                `json:"assume_only_visible"`
	Corrections       []correctionRequest `json:"corrections"`
}
