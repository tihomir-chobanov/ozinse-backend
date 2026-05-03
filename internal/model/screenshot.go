package model

type Screenshot struct {
    ID           int    `json:"id"`
    ProjectID    int    `json:"project_id"`
    URLToImage   string `json:"url_to_image"`
}