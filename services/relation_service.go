package services

import (
	"context"
	"p2p-management-service/db/generated"
	"p2p-management-service/models"
	"p2p-management-service/utils"

	"github.com/google/uuid"
)

type RelationService struct {
	DB        RelationRepository
	Encryptor *utils.Encryptor
}

type RelationRepository interface {
	GetRoleByName(ctx context.Context, nameHash string) (generated.Role, error)
	AssignRoleToUser(ctx context.Context, arg generated.AssignRoleToUserParams) error
	GetRoleByID(ctx context.Context, id uuid.UUID) (generated.Role, error)
	GetGroupByID(ctx context.Context, id uuid.UUID) (generated.Group, error)
	AssignGroupToUser(ctx context.Context, arg generated.AssignGroupToUserParams) error
}

func NewRelationService(db RelationRepository, encryptor *utils.Encryptor) *RelationService {
	return &RelationService{
		DB:        db,
		Encryptor: encryptor,
	}
}

func (s *RelationService) RolesToUsers(payload models.AttachRoleParams) (bool, error) {
	// get role by name hash
	role, err := s.DB.GetRoleByName(context.Background(), utils.Hash(payload.Role))
	if err != nil {
		return false, err
	}

	err = s.DB.AssignRoleToUser(context.Background(), generated.AssignRoleToUserParams{
		UserID: payload.User.ID,
		RoleID: role.ID,
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *RelationService) RolesToUsersByID(userId uuid.UUID, roleId uuid.UUID) (bool, error) {
	// get role by name hash
	role, err := s.DB.GetRoleByID(context.Background(), roleId)
	if err != nil {
		return false, err
	}

	err = s.DB.AssignRoleToUser(context.Background(), generated.AssignRoleToUserParams{
		UserID: userId,
		RoleID: role.ID,
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *RelationService) GroupToUsersByID(userId uuid.UUID, groupId uuid.UUID) (bool, error) {
	// get role by name hash
	group, err := s.DB.GetGroupByID(context.Background(), groupId)
	if err != nil {
		return false, err
	}

	err = s.DB.AssignGroupToUser(context.Background(), generated.AssignGroupToUserParams{
		UserID:  userId,
		GroupID: group.ID,
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
