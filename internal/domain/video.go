package domain

type VideoMeta struct {
	Title        string
	BaseURL      string
	StandardName string
	PremiumName  string
	Extension    string
}

var videoStore = map[int]VideoMeta{
	46325: {
		Title:        "Example Video 001",
		BaseURL:      "https://s3.eu-west-1.amazonaws.com/bytestreamfake",
		StandardName: "example001",
		PremiumName:  "example001-premium",
		Extension:    ".mp4",
	},
	77777: {
		Title:        "Example Video 002",
		BaseURL:      "https://s3.eu-west-1.amazonaws.com/bytestreamfake",
		StandardName: "example002",
		PremiumName:  "example002-premium",
		Extension:    ".mp4",
	},
}

func LookupVideo(videoID int) (VideoMeta, bool) {
	v, ok := videoStore[videoID]
	return v, ok
}

func BuildVideoResponse(videoID int, meta VideoMeta, isPremium bool) VideoResponse {
	filename := meta.StandardName
	if isPremium {
		filename = meta.PremiumName
	}
	return VideoResponse{
		VideoID:           videoID,
		Title:             meta.Title,
		PlaybackBaseURL:   meta.BaseURL,
		PlaybackFilename:  filename,
		PlaybackExtension: meta.Extension,
	}
}
