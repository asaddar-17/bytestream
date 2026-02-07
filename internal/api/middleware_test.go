package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExtractBearer(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		wantToken  string
		wantErr    bool
	}{
		{"missing header", "", "", true},
		{"wrong scheme", "Basic abc", "", true},
		{"empty token", "Bearer   ", "", true},
		{"valid token", "Bearer abc123", "abc123", false},
		{"case-insensitive bearer", "bEaReR xyz", "xyz", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}

			got, err := extractBearer(req)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if got != tc.wantToken {
				t.Fatalf("expected token %q, got %q", tc.wantToken, got)
			}
		})
	}
}
