package clients

import "fmt"

type UpstreamError struct {
	Status int
	Body   string
}

func (e UpstreamError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("upstream status=%d", e.Status)
	}
	return fmt.Sprintf("upstream status=%d body=%s", e.Status, e.Body)
}
