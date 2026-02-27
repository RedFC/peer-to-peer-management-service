package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Group struct {
	ID        uuid.UUID       `json:"id"`
	Name      string          `json:"name"`
	Metadata  json.RawMessage `json:"metadata"`
	IsDeleted bool            `json:"is_deleted"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

type CreateGroupParams struct {
	Name     string                 `json:"name,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type UpdateGroupParams struct {
	Name     string                 `json:"name,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
