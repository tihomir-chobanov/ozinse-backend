package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"ozinse-backend/internal/model"
)

// RoleRepository provides CRUD access to roles and their associated module permissions.
type RoleRepository struct {
	db *sql.DB
}

// NewRoleRepository creates a new RoleRepository instance.
func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// Create inserts a new role and all its module permissions using a database transaction.
func (r *RoleRepository) Create(role *model.Role) error {
	// Start the transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	// Ensure a rollback is executed if any step fails before committing
	defer tx.Rollback()

	// 1. Insert the role name and scan the auto-generated ID
	roleQuery := `INSERT INTO role (name) VALUES ($1) RETURNING id`
	err = tx.QueryRow(roleQuery, role.Name).Scan(&role.ID)
	if err != nil {
		return err
	}

	// 2. Insert permissions for each module in a loop
	permQuery := `INSERT INTO role_permission (role_id, module, access_level) VALUES ($1, $2, $3)`
	for i := range role.Permissions {
		p := &role.Permissions[i]
		p.RoleID = role.ID // Link the permission to the newly generated role ID

		_, err = tx.Exec(permQuery, p.RoleID, p.Module, p.AccessLevel)
		if err != nil {
			return err
		}
	}

	// Commit the transaction to apply all changes permanently
	return tx.Commit()
}

// GetAll retrieves all roles along with their permissions aggregated as a JSON array.
// role has fields: id, name, and permissions (array of module permissions)
// permissions have fields: id, role_id, module, access_level
func (r *RoleRepository) GetAll() ([]model.Role, error) {
	query := `
		SELECT r.id, r.name, 
		       COALESCE(json_agg(json_build_object(
		           'id', rp.id, 
		           'role_id', rp.role_id, 
		           'module', rp.module, 
		           'access_level', rp.access_level
		       )) FILTER (WHERE rp.id IS NOT NULL), '[]') as permissions
		FROM role r
		LEFT JOIN role_permission rp ON r.id = rp.role_id
		GROUP BY r.id, r.name
		ORDER BY r.id ASC;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var role model.Role
		var permissionsJSON []byte

		// Scan row columns into struct fields and raw JSON byte slice
		err := rows.Scan(&role.ID, &role.Name, &permissionsJSON)
		if err != nil {
			return nil, err
		}

		// Unmarshal the raw JSON array directly into the Go slice
		if err := json.Unmarshal(permissionsJSON, &role.Permissions); err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (r *RoleRepository) GetByID(id int) (*model.Role, error) {
	query := `
		SELECT r.id, r.name,
		       COALESCE(json_agg(json_build_object(
		           'id', rp.id, 
		           'role_id', rp.role_id, 
		           'module', rp.module, 
		           'access_level', rp.access_level
		       )) FILTER (WHERE rp.id IS NOT NULL), '[]') as permissions
		FROM role r
		LEFT JOIN role_permission rp ON r.id = rp.role_id
		WHERE r.id = $1
		GROUP BY r.id, r.name
	`
	var role model.Role
	var permissionsJSON []byte

	err := r.db.QueryRow(query, id).Scan(&role.ID, &role.Name, &permissionsJSON)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(permissionsJSON, &role.Permissions); err != nil {
		return nil, err
	}

	return &role, nil
}

// Update modifies an existing role name and replaces all its module permissions within a transaction.
func (r *RoleRepository) Update(role *model.Role) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Update the role name
	roleQuery := `UPDATE role SET name = $1 WHERE id = $2`
	result, err := tx.Exec(roleQuery, role.Name, role.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("role with id %d not found", role.ID)
	}

	// 2. Delete all old permissions associated with this role ID
	_, err = tx.Exec(`DELETE FROM role_permission WHERE role_id = $1`, role.ID)
	if err != nil {
		return err
	}

	// 3. Insert the new updated permissions
	permQuery := `INSERT INTO role_permission (role_id, module, access_level) VALUES ($1, $2, $3)`
	for i := range role.Permissions {
		p := &role.Permissions[i]
		p.RoleID = role.ID

		_, err = tx.Exec(permQuery, p.RoleID, p.Module, p.AccessLevel)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Delete removes a role by ID. Cascade constraints in the database automatically delete its permissions.
func (r *RoleRepository) Delete(id int) error {
	query := `DELETE FROM role WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("role with id %d not found", id)
	}

	return nil
}