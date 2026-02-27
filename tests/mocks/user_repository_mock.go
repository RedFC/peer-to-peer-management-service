package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"p2p-management-service/db/generated"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) GetTotalUsersCount(ctx context.Context, arg generated.GetTotalUsersCountParams) (int64, error) {
	args := m.Called(ctx, arg)
	if val := args.Get(0); val != nil {
		return val.(int64), args.Error(1)
	}
	return 0, args.Error(1)
}

func (m *UserRepository) ListUsers(ctx context.Context, arg generated.ListUsersParams) ([]generated.ListUsersRow, error) {
	args := m.Called(ctx, arg)
	if rows, ok := args.Get(0).([]generated.ListUsersRow); ok {
		return rows, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (generated.GetUserByIDRow, error) {
	args := m.Called(ctx, id)
	if row, ok := args.Get(0).(generated.GetUserByIDRow); ok {
		return row, args.Error(1)
	}
	return generated.GetUserByIDRow{}, args.Error(1)
}

func (m *UserRepository) CreateUser(ctx context.Context, arg generated.CreateUserParams) (generated.User, error) {
	args := m.Called(ctx, arg)
	if row, ok := args.Get(0).(generated.User); ok {
		return row, args.Error(1)
	}
	return generated.User{}, args.Error(1)
}

func (m *UserRepository) UpdateUser(ctx context.Context, arg generated.UpdateUserParams) (generated.User, error) {
	args := m.Called(ctx, arg)
	if row, ok := args.Get(0).(generated.User); ok {
		return row, args.Error(1)
	}
	return generated.User{}, args.Error(1)
}

func (m *UserRepository) UpdateUserProfile(ctx context.Context, arg generated.UpdateUserProfileParams) (generated.User, error) {
	args := m.Called(ctx, arg)
	if row, ok := args.Get(0).(generated.User); ok {
		return row, args.Error(1)
	}
	return generated.User{}, args.Error(1)
}

func (m *UserRepository) DeleteUser(ctx context.Context, arg generated.DeleteUserParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}
