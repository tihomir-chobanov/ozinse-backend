package model

// Episode represents a single episode of a season.
// swagger:model Episode
type Episode struct {
	ID             int    `json:"id" example:"1"`
	SeasonID       int    `json:"season_id" example:"10"`
	EpisodeNumber  int    `json:"episode_number" example:"1"`
	YoutubeVideoID string `json:"youtube_video_id" example:"dQw4w9WgXcQ"`
	Duration       int    `json:"duration" example:"45"`
}
