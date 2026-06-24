package handler

import (
	"github.com/ivan-ca97/life/internal/applications/health_connect_import/ports"
)

type typeResultResponse struct {
	Created int `json:"created"`
	Skipped int `json:"skipped"`
	Blocked int `json:"blocked"`
}

type importResponse struct {
	Weight    typeResultResponse `json:"weight"`
	Exercise  typeResultResponse `json:"exercise"`
	Steps     typeResultResponse `json:"steps"`
	Sleep     typeResultResponse `json:"sleep"`
	HeartRate typeResultResponse `json:"heart_rate"`
}

func importResponseFromResult(result *ports.ImportResult) *importResponse {
	convert := func(t ports.TypeResult) typeResultResponse {
		return typeResultResponse{
			Created: t.Created,
			Skipped: t.Skipped,
			Blocked: t.Blocked,
		}
	}
	return &importResponse{
		Weight:    convert(result.Weight),
		Exercise:  convert(result.Exercise),
		Steps:     convert(result.Steps),
		Sleep:     convert(result.Sleep),
		HeartRate: convert(result.HeartRate),
	}
}
