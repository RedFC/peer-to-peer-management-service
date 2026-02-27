package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"p2p-management-service/db/generated"
)

type RelationRepository struct {
	mock.Mock
}

func (m *RelationRepository) GetRoleByName(ctx context.Context, nameHash string) (generated.Role, error) {
	args := m.Called(ctx, nameHash)
	if row, ok := args.Get(0).(generated.Role); ok {
		return row, args.Error(1)
	}
	return generated.Role{}, args.Error(1)
}

func (m *RelationRepository) AssignRoleToUser(ctx context.Context, arg generated.AssignRoleToUserParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *RelationRepository) GetRoleByID(ctx context.Context, id uuid.UUID) (generated.Role, error) {
	args := m.Called(ctx, id)
	if row, ok := args.Get(0).(generated.Role); ok {
		return row, args.Error(1)
	}
	return generated.Role{}, args.Error(1)
}

func (m *RelationRepository) GetGroupByID(ctx context.Context, id uuid.UUID) (generated.Group, error) {
	args := m.Called(ctx, id)
	if row, ok := args.Get(0).(generated.Group); ok {
		return row, args.Error(1)
	}
	return generated.Group{}, args.Error(1)
}

func (m *RelationRepository) AssignGroupToUser(ctx context.Context, arg generated.AssignGroupToUserParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}
