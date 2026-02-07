package api

import (
	"testing"
	"time"
)

func TestHasRole(t *testing.T) {
	tests := []struct {
		name  string
		roles []string
		want  string
		ok    bool
	}{
		{"empty", nil, "premium", false},
		{"not present", []string{"standard"}, "premium", false},
		{"present", []string{"standard", "premium"}, "premium", true},
		{"case sensitive", []string{"Premium"}, "premium", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := hasRole(tc.roles, tc.want)
			if got != tc.ok {
				t.Fatalf("expected %v, got %v", tc.ok, got)
			}
		})
	}
}

func TestIsInWindow(t *testing.T) {
	now := time.Date(2026, 2, 7, 12, 0, 0, 0, time.UTC) // 2026-02-07 12:00 UTC

	tests := []struct {
		name string
		from string
		to   string
		ok   bool
	}{
		{"inside range", "2026-02-01", "2026-02-10", true},
		{"before range", "2026-02-08", "2026-02-10", false},
		{"after range", "2026-01-01", "2026-02-06", false},

		// boundary checks
		{"on start date", "2026-02-07", "2026-02-10", true},
		{"on end date (midday)", "2026-02-01", "2026-02-07", true},

		// invalid inputs
		{"bad from date", "bad", "2026-02-10", false},
		{"bad to date", "2026-02-01", "bad", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := isInWindow(now, tc.from, tc.to)
			if got != tc.ok {
				t.Fatalf("expected %v, got %v (from=%s to=%s)", tc.ok, got, tc.from, tc.to)
			}
		})
	}
}
