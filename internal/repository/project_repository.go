package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"ozinse-backend/internal/model"
)

// ProjectRepository provides CRUD access to project data and related relations.
type ProjectRepository struct {
	db *sql.DB
}

// NewProjectRepository creates a new ProjectRepository instance.
func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// GetAll retrieves all projects from the database along with their associated
// genres, categories, and screenshots in a single, highly optimized query.
//
// WHY THIS APPROACH WAS CHOSEN (Architecture note):
// We are using PostgreSQL's native JSON aggregation functions (json_agg) instead of
// making separate database calls for each project's relations. Fetching the parent
// records and then looping through them to fetch child records creates the dreaded
// "N+1 query problem", which severely degrades performance as the dataset grows.
// By shifting the aggregation logic to the database layer, we reduce network overhead
// and fetch everything in exactly 1 query, ensuring the API remains fast and scalable.
func (r *ProjectRepository) GetAll() ([]model.Project, error) {

	// Define the SQL query.
	// COALESCE is used to return an empty JSON array '[]' instead of NULL
	// if a project has no associated records in the many-to-many tables.
	// json_build_object creates a JSON object (key-value pairs) for each row.
	// json_agg collects all these objects into a single JSON array.
	query := `
		SELECT 
			p.id, p.title, p.description, p.release_year, p.cover_image_url, 
			p.is_featured, p.type, p.duration, p.keywords, p.director, p.producer,
			
			-- Aggregate genres into a JSON array
			COALESCE((
				SELECT json_agg(json_build_object('id', g.id, 'name', g.name, 'icon_url', g.icon_url))
				FROM genre g
				JOIN project_genre pg ON g.id = pg.genre_id
				WHERE pg.project_id = p.id
			), '[]') AS genres,

						-- Aggregate age_categories into a JSON array
			COALESCE((
				SELECT json_agg(json_build_object('id', ac.id, 'range', ac.range, 'icon_url', ac.icon_url))
				FROM age_category ac
				JOIN project_age_category pac ON ac.id = pac.age_category_id
				WHERE pac.project_id = p.id
			), '[]') AS age_categories,

			-- Aggregate categories into a JSON array
			COALESCE((
				SELECT json_agg(json_build_object('id', c.id, 'name', c.name))
				FROM category c
				JOIN project_category pc ON c.id = pc.category_id
				WHERE pc.project_id = p.id
			), '[]') AS categories,

			-- Aggregate screenshots into a JSON array
			COALESCE((
				SELECT json_agg(json_build_object('id', s.id, 'project_id', s.project_id, 'url_to_image', s.url_to_image))
				FROM project_screenshot s
				WHERE s.project_id = p.id
			), '[]') AS screenshots

		FROM project p
	`

	// Execute the query
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure the database connection is released back to the pool

	var projects []model.Project

	// Iterate over the fetched rows
	for rows.Next() {
		var p model.Project

		// Temporary byte slices to hold the raw JSON data returned by PostgreSQL
		var genresJSON, ageCategoriesJSON, categoriesJSON, screenshotsJSON []byte

		// Scan the row columns into the respective struct fields and byte slices
		err := rows.Scan(
			&p.ID, &p.Title, &p.Description, &p.ReleaseYear, &p.CoverImageUrl,
			&p.IsFeatured, &p.Type, &p.Duration, &p.Keywords, &p.Director, &p.Producer,
			&genresJSON, &ageCategoriesJSON, &categoriesJSON, &screenshotsJSON,
		)
		if err != nil {
			return nil, err
		}

		// Decode (unmarshal) the raw JSON arrays directly into Go slices.
		// Since we used COALESCE(..., '[]') in SQL, these byte slices will never be
		// entirely empty/nil, preventing unmarshal errors.
		json.Unmarshal(genresJSON, &p.Genres)
		json.Unmarshal(ageCategoriesJSON, &p.AgeCategories)
		json.Unmarshal(categoriesJSON, &p.Categories)
		json.Unmarshal(screenshotsJSON, &p.Screenshots)

		// Add the fully populated project to our final slice
		projects = append(projects, p)
	}

	return projects, nil
}

// GetByID retrieves a project by its ID, including its screenshots.
func (r *ProjectRepository) GetByID(id int) (*model.Project, error) {
	var p model.Project

	// 1. Fetch the main project data
	query := `SELECT id, title, description, release_year, cover_image_url, is_featured, type, duration, keywords, director, producer 
              FROM project WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.Title, &p.Description, &p.ReleaseYear, &p.CoverImageUrl,
		&p.IsFeatured, &p.Type, &p.Duration, &p.Keywords, &p.Director, &p.Producer,
	)
	if err != nil {
		return nil, err
	}

	// 2. Fetch screenshots for this specific project
	rows, err := r.db.Query(`SELECT id, project_id, url_to_image FROM project_screenshot WHERE project_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 3. Iterate and add them to the project struct
	for rows.Next() {
		var img model.ProjectScreenshot
		if err := rows.Scan(&img.ID, &img.ProjectID, &img.URLToImage); err != nil {
			return nil, err
		}
		p.Screenshots = append(p.Screenshots, img)
	}

	return &p, nil
}

// Create inserts a new project and all related many-to-many associations.
func (r *ProjectRepository) Create(p *model.Project, genreIDs []int, ageCategoryIDs []int, categoryIDs []int) error {
	// 1. Begin a database transaction to ensure data integrity
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 2. Insert main project details into the 'project' table
	query := `INSERT INTO project (title, description, release_year, cover_image_url, is_featured, type, duration, keywords, director, producer) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	err = tx.QueryRow(query,
		p.Title, p.Description, p.ReleaseYear, p.CoverImageUrl,
		p.IsFeatured, p.Type, p.Duration, p.Keywords, p.Director, p.Producer,
	).Scan(&p.ID)

	if err != nil {
		return err
	}

	// 3. Insert many-to-many relationships for project genres
	for _, genreID := range genreIDs {
		_, err = tx.Exec(`INSERT INTO project_genre (project_id, genre_id) VALUES ($1, $2)`, p.ID, genreID)
		if err != nil {
			return err
		}
	}

	// 4. Insert many-to-many relationships for project age categories
	for _, ageCategoryID := range ageCategoryIDs {
		_, err = tx.Exec(`INSERT INTO project_age_category (project_id, age_category_id) VALUES ($1, $2)`, p.ID, ageCategoryID)
		if err != nil {
			return err
		}
	}

	// 5. Insert many-to-many relationships for project categories
	for _, catID := range categoryIDs {
		_, err = tx.Exec(`INSERT INTO project_category (project_id, category_id) VALUES ($1, $2)`, p.ID, catID)
		if err != nil {
			return err
		}
	}

	// 5. Insert project screenshots
	for _, img := range p.Screenshots {
		_, err = tx.Exec(`INSERT INTO project_screenshot (project_id, url_to_image) VALUES ($1, $2)`, p.ID, img.URLToImage)
		if err != nil {
			return err
		}
	}

	// 6. Insert seasons and episodes if the project type is "series"
	if p.Type == model.ProjectSeries {
		for i := range p.Seasons {
			s := &p.Seasons[i]
			var seasonID int

			err = tx.QueryRow(`INSERT INTO season (project_id, season_number) VALUES ($1, $2) RETURNING id`, p.ID, s.SeasonNumber).Scan(&seasonID)
			if err != nil {
				return err
			}

			for _, e := range s.Episodes {
				_, err = tx.Exec(`INSERT INTO episode (season_id, episode_number, youtube_video_id, duration) VALUES ($1, $2, $3, $4)`,
					seasonID, e.EpisodeNumber, e.YoutubeVideoID, e.Duration)
				if err != nil {
					return err
				}
			}
		}
	}

	// 7. Commit the transaction to apply all changes to the database
	return tx.Commit()
}

// Update modifies an existing project and refreshes its screenshots.
func (r *ProjectRepository) Update(p *model.Project) error {
	// 1. Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	// Rollback in case of an error
	defer tx.Rollback()

	// 2. Update the table
	query := `UPDATE project SET 
                title = $1, description = $2, release_year = $3, 
                cover_image_url = $4, is_featured = $5, type = $6, 
                duration = $7, keywords = $8, director = $9, producer = $10 
              WHERE id = $11`

	result, err := tx.Exec(query,
		p.Title, p.Description, p.ReleaseYear,
		p.CoverImageUrl, p.IsFeatured, p.Type,
		p.Duration, p.Keywords, p.Director, p.Producer,
		p.ID,
	)
	if err != nil {
		return err
	}

	// Check if the update affected any rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("project with id %d not found", p.ID)
	}

	// 3. Update the screenshots, which is a separate table
	_, err = tx.Exec(`DELETE FROM project_screenshot WHERE project_id = $1`, p.ID)
	if err != nil {
		return err
	}

	for _, img := range p.Screenshots {
		_, err = tx.Exec(`INSERT INTO project_screenshot (project_id, url_to_image) VALUES ($1, $2)`, p.ID, img.URLToImage)
		if err != nil {
			return err
		}
	}

	// 4. Commit in case of success
	return tx.Commit()
}

// Delete removes a project from the database by ID.
func (r *ProjectRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM project WHERE id = $1`, id)
	if err != nil {
		return err
	}
	// if we didn't implement RowsAffected(), we could delete a project that doesn't exist and that is missleading/wrong
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project with id %d not found", id)
	}
	return nil
}

// ExistsByName checks whether a project title already exists in the database.
func (r *ProjectRepository) ExistsByName(name string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM project WHERE title = $1)`
	err := r.db.QueryRow(query, name).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, err
}


