package clients

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"bytestream/internal/domain"
)

type IdentityClient struct {
	baseURL string
	httpc   *http.Client
}

func NewIdentityClient(baseURL string, timeout time.Duration) *IdentityClient {
	return &IdentityClient{
		baseURL: baseURL,
		httpc:   &http.Client{Timeout: timeout},
	}
}

func (c *IdentityClient) GetUserInfo(bearerToken string) (domain.Identity, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/identity/userinfo", nil)
	if err != nil {
		return domain.Identity{}, err
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, err := c.httpc.Do(req)
	if err != nil {
		return domain.Identity{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return domain.Identity{}, UpstreamError{Status: resp.StatusCode, Body: string(b)}
	}

	var identity domain.Identity
	if err := json.NewDecoder(resp.Body).Decode(&identity); err != nil {
		return domain.Identity{}, err
	}
	return identity, nil
}
