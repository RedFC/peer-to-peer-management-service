package models

import (
	"github.com/google/uuid"
)

type Subscription struct {
	ID                 uuid.UUID `json:"id"`
	SubscriptionStatus string    `json:"subscription_status"`
	CreatedAt          string    `json:"created_at"`
	UpdatedAt          string    `json:"updated_at"`
}

type SubscriptionWithDetails struct {
	Subscription
	User  UserResponseParams
	Group Group
}

type CreateSubscriptionParams struct {
	UserID             uuid.UUID `json:"user_id"`
	GroupID            uuid.UUID `json:"group_id"`
	SubscriptionStatus string    `json:"subscription_status"`
}

type UpdateSubscriptionParams struct {
	ID                 uuid.UUID `json:"id"`
	UserID             uuid.UUID `json:"user_id,omitempty"`
	GroupID            uuid.UUID `json:"group_id,omitempty"`
	SubscriptionStatus string    `json:"subscription_status,omitempty"`
}
