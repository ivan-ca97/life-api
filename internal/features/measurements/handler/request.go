package handler

type upsertMeasurementRequest struct {
	Value float64 `json:"value"`
	Notes string  `json:"notes"`
}
