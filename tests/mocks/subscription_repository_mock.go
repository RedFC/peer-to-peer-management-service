package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"p2p-management-service/db/generated"
)

type SubscriptionRepository struct {
	mock.Mock
}

func (m *SubscriptionRepository) ListSubscriptions(ctx context.Context, arg generated.ListSubscriptionsParams) ([]generated.ListSubscriptionsRow, error) {
	args := m.Called(ctx, arg)
	if rows, ok := args.Get(0).([]generated.ListSubscriptionsRow); ok {
		return rows, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SubscriptionRepository) GetTotalSubscriptionsCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	if val := args.Get(0); val != nil {
		return val.(int64), args.Error(1)
	}
	return 0, args.Error(1)
}

func (m *SubscriptionRepository) GetSubscriptionByID(ctx context.Context, id uuid.UUID) (generated.GetSubscriptionByIDRow, error) {
	args := m.Called(ctx, id)
	if row, ok := args.Get(0).(generated.GetSubscriptionByIDRow); ok {
		return row, args.Error(1)
	}
	return generated.GetSubscriptionByIDRow{}, args.Error(1)
}

func (m *SubscriptionRepository) GetUserByID(ctx context.Context, id uuid.UUID) (generated.GetUserByIDRow, error) {
	args := m.Called(ctx, id)
	if row, ok := args.Get(0).(generated.GetUserByIDRow); ok {
		return row, args.Error(1)
	}
	return generated.GetUserByIDRow{}, args.Error(1)
}

func (m *SubscriptionRepository) GetSubscriptionByUserIDAndGroupIDAndStatus(ctx context.Context, arg generated.GetSubscriptionByUserIDAndGroupIDAndStatusParams) (generated.Subscription, error) {
	args := m.Called(ctx, arg)
	if row, ok := args.Get(0).(generated.Subscription); ok {
		return row, args.Error(1)
	}
	return generated.Subscription{}, args.Error(1)
}

func (m *SubscriptionRepository) UpdateSubscription(ctx context.Context, arg generated.UpdateSubscriptionParams) (generated.Subscription, error) {
	args := m.Called(ctx, arg)
	if row, ok := args.Get(0).(generated.Subscription); ok {
		return row, args.Error(1)
	}
	return generated.Subscription{}, args.Error(1)
}

func (m *SubscriptionRepository) CreateSubscription(ctx context.Context, arg generated.CreateSubscriptionParams) (generated.Subscription, error) {
	args := m.Called(ctx, arg)
	if row, ok := args.Get(0).(generated.Subscription); ok {
		return row, args.Error(1)
	}
	return generated.Subscription{}, args.Error(1)
}

func (m *SubscriptionRepository) GetRoleByName(ctx context.Context, nameHash string) (generated.Role, error) {
	args := m.Called(ctx, nameHash)
	if row, ok := args.Get(0).(generated.Role); ok {
		return row, args.Error(1)
	}
	return generated.Role{}, args.Error(1)
}

func (m *SubscriptionRepository) AssignRoleToUser(ctx context.Context, arg generated.AssignRoleToUserParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *SubscriptionRepository) DeleteSubscription(ctx context.Context, arg generated.DeleteSubscriptionParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *SubscriptionRepository) GetSubscriptionByUserIDAndGroupID(ctx context.Context, arg generated.GetSubscriptionByUserIDAndGroupIDParams) ([]generated.Subscription, error) {
	args := m.Called(ctx, arg)
	if rows, ok := args.Get(0).([]generated.Subscription); ok {
		return rows, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SubscriptionRepository) RevokeSubscription(ctx context.Context, arg generated.RevokeSubscriptionParams) (generated.Subscription, error) {
	args := m.Called(ctx, arg)
	if row, ok := args.Get(0).(generated.Subscription); ok {
		return row, args.Error(1)
	}
	return generated.Subscription{}, args.Error(1)
}
