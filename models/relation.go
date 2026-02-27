package models

type AttachRoleParams struct {
	User UserResponseParams `json:"user"`
	Role string             `json:"role"`
}

// roles enums
const (
	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleSubscriber = "subscriber"
	RoleUser       = "user"
)
