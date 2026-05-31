package model

type FavoriteRequest struct {
	ProjectID int `json:"project_id" binding:"required"`
}
