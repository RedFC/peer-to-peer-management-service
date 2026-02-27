package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"p2p-management-service/db/generated"
)

type RoleRepository struct {
	mock.Mock
}

func (m *RoleRepository) ListRoles(ctx context.Context) ([]generated.Role, error) {
	args := m.Called(ctx)
	if rows, ok := args.Get(0).([]generated.Role); ok {
		return rows, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *RoleRepository) GetRoleByID(ctx context.Context, id uuid.UUID) (generated.Role, error) {
	args := m.Called(ctx, id)
	if row, ok := args.Get(0).(generated.Role); ok {
		return row, args.Error(1)
	}
	return generated.Role{}, args.Error(1)
}

func (m *RoleRepository) CreateRole(ctx context.Context, arg generated.CreateRoleParams) (generated.Role, error) {
	args := m.Called(ctx, arg)
	if row, ok := args.Get(0).(generated.Role); ok {
		return row, args.Error(1)
	}
	return generated.Role{}, args.Error(1)
}

func (m *RoleRepository) UpdateRole(ctx context.Context, arg generated.UpdateRoleParams) (generated.Role, error) {
	args := m.Called(ctx, arg)
	if row, ok := args.Get(0).(generated.Role); ok {
		return row, args.Error(1)
	}
	return generated.Role{}, args.Error(1)
}

func (m *RoleRepository) DeleteRole(ctx context.Context, arg generated.DeleteRoleParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}
