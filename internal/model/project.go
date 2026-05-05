package model

// Project represents a movie or series with metadata and relations.
// swagger:model Project
type Project struct {
	ID            int                 `json:"id" example:"1"`
	Title         string              `json:"title" example:"Example Movie"`
	Description   string              `json:"description" example:"A short description of the project."`
	ReleaseYear   int                 `json:"release_year" example:"2024"`
	CoverImageUrl string              `json:"cover_image_url" example:"https://example.com/poster.png"`
	IsFeatured    bool                `json:"is_featured" example:"false"`
	Type          string              `json:"type" binding:"required,oneof=movie series" example:"movie"`
	Duration      int                 `json:"duration" example:"120"`
	Keywords      string              `json:"keywords" example:"action,thriller"`
	Director      string              `json:"director" example:"Jane Doe"`
	Producer      string              `json:"producer" example:"John Smith"`
	Seasons       []Season            `json:"seasons,omitempty"`
	Genres        []Genre             `json:"genres,omitempty"`
	AgeCategories []Age_Category      `json:"age_categories,omitempty"`
	Categories    []Category          `json:"categories,omitempty"`
	Screenshots   []ProjectScreenshot `json:"screenshots,omitempty"`
}

const (
	ProjectSeries = "series"
	ProjectMovie  = "movie"
)
