package tests

import (
	"net/http"
	"testing"
	"time"
)

func TestVideoEndpoint_E2E(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest("GET", "http://localhost:8080/videos/46325", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer anytoken")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
