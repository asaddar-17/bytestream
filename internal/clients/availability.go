package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"bytestream/internal/domain"
)

type AvailabilityClient struct {
	baseURL string
	httpc   *http.Client
}

func NewAvailabilityClient(baseURL string, timeout time.Duration) *AvailabilityClient {
	return &AvailabilityClient{
		baseURL: baseURL,
		httpc:   &http.Client{Timeout: timeout},
	}
}

func (c *AvailabilityClient) GetAvailability(bearerToken string, videoID int) (domain.AvailabilityInfo, error) {
	url := fmt.Sprintf("%s/availability/availabilityinfo/%d", c.baseURL, videoID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return domain.AvailabilityInfo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, err := c.httpc.Do(req)
	if err != nil {
		return domain.AvailabilityInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return domain.AvailabilityInfo{}, UpstreamError{Status: resp.StatusCode, Body: string(b)}
	}

	var info domain.AvailabilityInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return domain.AvailabilityInfo{}, err
	}
	return info, nil
}
