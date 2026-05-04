package repository

import (
	"database/sql"
	"fmt"
	"ozinse-backend/internal/model"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) GetAll() ([]model.Project, error) {
	// 1. Fetch all basic project data
	rows, err := r.db.Query(`SELECT id, title, description, release_year, cover_image_url, is_featured, type, duration, keywords, director, producer FROM project`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var p model.Project
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.ReleaseYear, &p.CoverImageUrl, &p.IsFeatured, &p.Type, &p.Duration, &p.Keywords, &p.Director, &p.Producer); err != nil {
			return nil, err
		}

		// 2. For EACH project, fetch its screenshots from the separate table
		screenshotRows, err := r.db.Query(`SELECT id, project_id, url_to_image FROM project_screenshot WHERE project_id = $1`, p.ID)
		if err != nil {
			// If we fail to fetch screenshots, we can either return an error or just continue with empty screenshots
			return nil, err
		}

		// Important: Close screenshot rows after we are done with them for this specific project
		for screenshotRows.Next() {
			var img model.ProjectScreenshot
			if err := screenshotRows.Scan(&img.ID, &img.ProjectID, &img.URLToImage); err != nil {
				screenshotRows.Close()
				return nil, err
			}
			// Append the screenshot to the current project's slice
			p.Screenshots = append(p.Screenshots, img)
		}
		screenshotRows.Close()

		projects = append(projects, p)
	}

	return projects, nil
}

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

func (r *ProjectRepository) Create(p *model.Project) error {
	// 1. Begin a new database transaction to ensure atomicity.
	// If any step fails, we can roll back all changes to keep the database consistent.
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Ensure that if the function returns early due to an error, all changes are undone.
	defer tx.Rollback()

	// 2. Insert the main project data into the 'project' table.
	// We use the RETURNING id clause to get the auto-generated ID for the new project.
	query := `INSERT INTO project (title, description, release_year, cover_image_url, is_featured, type, duration, keywords, director, producer) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	err = tx.QueryRow(query,
		p.Title, p.Description, p.ReleaseYear, p.CoverImageUrl,
		p.IsFeatured, p.Type, p.Duration, p.Keywords, p.Director, p.Producer,
	).Scan(&p.ID)

	if err != nil {
		return err
	}

	// 3. Iterate through the screenshots array and insert each record into the 'project_screenshot' table.
	// Each screenshot is linked to the project using the newly generated project ID.
	for _, img := range p.Screenshots {
		_, err = tx.Exec(`INSERT INTO project_screenshot (project_id, url_to_image) VALUES ($1, $2)`, p.ID, img.URLToImage)
		if err != nil {
			return err
		}
	}

	// 4. Handle series-specific data: seasons and episodes.
	if p.Type == model.ProjectSeries {
		for i := range p.Seasons {
			s := &p.Seasons[i]
			var seasonID int

			// Insert the season record and retrieve its generated ID.
			err = tx.QueryRow(`INSERT INTO season (project_id, season_number) VALUES ($1, $2) RETURNING id`, p.ID, s.SeasonNumber).Scan(&seasonID)
			if err != nil {
				return err
			}

			// Insert all episodes belonging to this specific season.
			for _, e := range s.Episodes {
				_, err = tx.Exec(`INSERT INTO episode (season_id, episode_number, youtube_video_id, duration) VALUES ($1, $2, $3, $4)`,
					seasonID, e.EpisodeNumber, e.YoutubeVideoID, e.Duration)
				if err != nil {
					return err
				}
			}
		}
	}

	// 5. If all insertions were successful, commit the transaction to make the changes permanent.
	return tx.Commit()
}

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

func (r *ProjectRepository) ExistsByName(name string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM project WHERE title = $1)`
	err := r.db.QueryRow(query, name).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, err
}

/* The Repository acts as a bridge between your application's business logic and the data source (database). It is a design pattern that isolates data access logic. Instead of your code knowing exactly how to write SQL queries, it simply asks the repository for an object (e.g., "Give me Category with ID 5"). The repository handles the technical execution with the database and returns a clean Model.

Repository (The Supplier): The only layer allowed to communicate with the database (PostgreSQL) using raw SQL queries.

*/
