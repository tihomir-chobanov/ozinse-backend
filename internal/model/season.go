package model

// Season represents a season of a series project.
// swagger:model Season
type Season struct {
	ID           int       `json:"id" example:"1"`
	ProjectID    int       `json:"project_id" example:"42"`
	SeasonNumber int       `json:"season_number" example:"1"`
	Episodes     []Episode `json:"episodes,omitempty"`
}
