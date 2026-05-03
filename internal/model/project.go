package model

type Project struct {
    ID   int    `json:"id"`
    Title string `json:"title"`
    Description string `json:"description"`
    ReleaseYear int `json:"release_year"`
    CoverImageUrl string `json:"cover_image_url"`
    IsFeatured bool `json:"is_featured"`
    Type string `json:"type"`
    Duration int `json:"duration"`
    Keywords string `json:"keywords"`
    Director string `json:"director"`
    Producer string `json:"producer"`
    Seasons     []Season            `json:"seasons,omitempty"`   
    Screenshots []ProjectScreenshot `json:"screenshots,omitempty"`
}

const (
    ProjectSeries = "series"
    ProjectMovie  = "movie"
)