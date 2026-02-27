// services/user_service_interface.go
package services

import (
	"context"
	"p2p-management-service/models"

	"github.com/google/uuid"
)

type IUserService interface {
	GetUsers(pageNo, pageSize string, totalCount int64) ([]models.UserResponseParams, error)
	GetUserByID(id uuid.UUID) (*models.UserResponseParams, error)
	CreateUser(user models.CreateUserParams) (*models.UserResponseParams, error)
	UpdateUser(id uuid.UUID, user models.UpdateUserParams) (*models.UserResponseParams, error)
	DeleteUser(id uuid.UUID) error
	GetTotalUsersCount(ctx context.Context) (int64, error)
}
