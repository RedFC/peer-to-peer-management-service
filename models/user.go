package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string `json:"email"`
	EmailHash string `json:"emailHash"`
	Password  string `json:"password"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	CreatedAt string `json:"createdAt"`
}

// is_deleted, is_active, is_password_reset, last_login
type UserResponseParams struct {
	ID              uuid.UUID `json:"id"`
	Email           string    `json:"email"`
	FirstName       string    `json:"first_name,omitempty"`
	LastName        string    `json:"last_name,omitempty"`
	IsDeleted       bool      `json:"is_deleted"`
	IsActive        bool      `json:"is_active"`
	IsPasswordReset bool      `json:"is_password_reset"`
	Role            []Role    `json:"role,omitempty"`
	Group           []Group   `json:"group,omitempty"`
	LastLogin       string    `json:"last_login,omitempty"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
}

// is_deleted, is_active, is_password_reset, last_login
type UserWithProfileResponseParams struct {
	ID              uuid.UUID   `json:"id"`
	Email           string      `json:"email"`
	FirstName       string      `json:"first_name,omitempty"`
	LastName        string      `json:"last_name,omitempty"`
	IsDeleted       bool        `json:"is_deleted"`
	IsActive        bool        `json:"is_active"`
	IsPasswordReset bool        `json:"is_password_reset"`
	Role            interface{} `json:"role,omitempty"`
	Group           interface{} `json:"group,omitempty"`
	Profile         interface{} `json:"profile,omitempty"`
	LastLogin       string      `json:"last_login"`
	CreatedAt       string      `json:"created_at"`
	UpdatedAt       string      `json:"updated_at"`
}

type CreateUserParams struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

type UpdateUserParams struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}
