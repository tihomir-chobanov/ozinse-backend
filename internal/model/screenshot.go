package model

// Screenshot represents a screenshot image linked to a project.
// swagger:model Screenshot
type Screenshot struct {
	ID         int    `json:"id" example:"1"`
	ProjectID  int    `json:"project_id" example:"42"`
	URLToImage string `json:"url_to_image" example:"https://example.com/screenshot.png"`
}
