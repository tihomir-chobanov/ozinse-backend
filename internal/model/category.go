package model

// Category represents a content category used to tag projects.
// swagger:model Category
type Category struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"Drama"`
}


// CategoryMovieCount represents the aggregated statistics for a category
type CategoryMovieCount struct {
	CategoryID   int    `json:"category_id" db:"category_id"`
	CategoryName string `json:"category_name" db:"category_name"`
	TotalMovies  int    `json:"total_movies" db:"total_movies"`
}
