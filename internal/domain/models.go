package domain

type Identity struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

type AvailabilityInfo struct {
	VideoID            int `json:"video_id"`
	AvailabilityWindow struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"availability_window"`
}

type VideoResponse struct {
	VideoID           int    `json:"video_id"`
	Title             string `json:"title"`
	PlaybackBaseURL   string `json:"playback_baseurl"`
	PlaybackFilename  string `json:"playback_filename"`
	PlaybackExtension string `json:"playback_extension"`
}
