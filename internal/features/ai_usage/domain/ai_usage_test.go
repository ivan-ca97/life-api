package domain

import (
	"testing"
	"time"
)

func ptr(v float64) *float64 { return &v }

func TestEffectiveLimitUsd(t *testing.T) {
	cases := []struct {
		name      string
		tierLimit *float64
		selfLimit *float64
		want      *float64
	}{
		{"unlimited tier, no self limit", nil, nil, nil},
		{"unlimited tier, self limit binds", nil, ptr(10), ptr(10)},
		{"tier limit, no self limit", ptr(5), nil, ptr(5)},
		{"self limit lower than tier", ptr(20), ptr(8), ptr(8)},
		{"tier lower than self limit", ptr(3), ptr(8), ptr(3)},
		{"equal limits", ptr(5), ptr(5), ptr(5)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Allocation{
				Tier:         Tier{MonthlyLimitUsd: c.tierLimit},
				SelfLimitUsd: c.selfLimit,
			}.EffectiveLimitUsd()
			switch {
			case c.want == nil && got != nil:
				t.Fatalf("expected unlimited (nil), got %v", *got)
			case c.want != nil && got == nil:
				t.Fatalf("expected %v, got unlimited (nil)", *c.want)
			case c.want != nil && got != nil && *got != *c.want:
				t.Fatalf("expected %v, got %v", *c.want, *got)
			}
		})
	}
}

func TestOverLimit(t *testing.T) {
	if OverLimit(5, nil) {
		t.Error("nil limit must never be over")
	}
	if !OverLimit(5, ptr(5)) {
		t.Error("reaching the limit exactly should be over")
	}
	if OverLimit(4.99, ptr(5)) {
		t.Error("below the limit should not be over")
	}
}

func TestPeriodStart(t *testing.T) {
	in := time.Date(2026, 6, 22, 13, 45, 0, 0, time.UTC)
	got := PeriodStart(in)
	want := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
