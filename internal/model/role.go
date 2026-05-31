package model


type Role struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Permissions []Permission `json:"permissions"`
}

type Permission struct {
	ID          int    `json:"id"`
	RoleID      int    `json:"role_id"`
	Module      string `json:"module"`
	AccessLevel string `json:"access_level"`
}