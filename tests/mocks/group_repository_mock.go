package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"p2p-management-service/db/generated"
)

type GroupRepository struct {
	mock.Mock
}

func (m *GroupRepository) ListGroups(ctx context.Context, arg generated.ListGroupsParams) ([]generated.Group, error) {
	args := m.Called(ctx, arg)
	if rows, ok := args.Get(0).([]generated.Group); ok {
		return rows, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *GroupRepository) GetTotalGroupsCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	if val := args.Get(0); val != nil {
		return val.(int64), args.Error(1)
	}
	return 0, args.Error(1)
}

func (m *GroupRepository) GetGroupByID(ctx context.Context, id uuid.UUID) (generated.Group, error) {
	args := m.Called(ctx, id)
	if row, ok := args.Get(0).(generated.Group); ok {
		return row, args.Error(1)
	}
	return generated.Group{}, args.Error(1)
}

func (m *GroupRepository) CreateGroup(ctx context.Context, arg generated.CreateGroupParams) (generated.Group, error) {
	args := m.Called(ctx, arg)
	if row, ok := args.Get(0).(generated.Group); ok {
		return row, args.Error(1)
	}
	return generated.Group{}, args.Error(1)
}

func (m *GroupRepository) UpdateGroup(ctx context.Context, arg generated.UpdateGroupParams) (generated.Group, error) {
	args := m.Called(ctx, arg)
	if row, ok := args.Get(0).(generated.Group); ok {
		return row, args.Error(1)
	}
	return generated.Group{}, args.Error(1)
}

func (m *GroupRepository) DeleteGroup(ctx context.Context, arg generated.DeleteGroupParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}
