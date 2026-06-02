package model

type CreateMainPageEntryRequest struct {
    ProjectID int    `json:"project_id" binding:"required"`
    Position  int    `json:"position" binding:"required,min=1"`
    IconURL string `json:"icon_url" binding:"required,url"`
}

type MainPageProject struct {
    ID       int     `json:"id"`
    Position int     `json:"position"`
    IconURL  string  `json:"icon_url"`
    Project  Project `json:"project"` // Вграждаме целия съществуващ модел Project
}

type UpdateMainPageEntryRequest struct {
    Position int    `json:"position" binding:"required,min=1"`
    IconURL  string `json:"icon_url" binding:"required,url"`
}