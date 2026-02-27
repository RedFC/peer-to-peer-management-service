package models

import "github.com/google/uuid"

type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsDeleted   bool      `json:"is_deleted"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

type CreateRoleParams struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateRoleParams struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}
