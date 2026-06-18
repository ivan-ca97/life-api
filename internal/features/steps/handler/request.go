package handler

type upsertStepsRequest struct {
	Steps  int    `json:"steps"`
	Source string `json:"source"`
}
