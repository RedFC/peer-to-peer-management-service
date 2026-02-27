package tests

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"p2p-management-service/db/generated"
	"p2p-management-service/models"
	"p2p-management-service/services"
	testmocks "p2p-management-service/tests/mocks"
	testutils "p2p-management-service/tests/testutils"
	"p2p-management-service/utils"
)

func TestSubscriptionService_GetSubscriptions(t *testing.T) {
	repo := &testmocks.SubscriptionRepository{}
	encryptor := utils.NewEncryptor(testutils.TestEncryptionKey)
	service := services.NewSubscriptionService(repo, encryptor, nil)

	now := time.Now().UTC()
	userID := uuid.New()

	encFirst, err := encryptor.Encrypt("Ava")
	require.NoError(t, err)
	encLast, err := encryptor.Encrypt("Lang")
	require.NoError(t, err)
	encEmail, err := encryptor.Encrypt("ava@example.com")
	require.NoError(t, err)
	encGroup, err := encryptor.Encrypt("Finance")
	require.NoError(t, err)

	repo.On("ListSubscriptions", mock.Anything, generated.ListSubscriptionsParams{
		Limit:  20,
		Offset: 0,
	}).Return([]generated.ListSubscriptionsRow{
		{
			SubscriptionID:        uuid.New(),
			SubscriptionStatus:    generated.NullSubscriptionStatus{SubscriptionStatus: generated.SubscriptionStatusInactive, Valid: true},
			SubscriptionCreatedAt: sql.NullTime{Time: now, Valid: true},
			SubscriptionUpdatedAt: sql.NullTime{Time: now, Valid: true},
			UserID:                userID,
			UserFirstName:         sql.NullString{String: encFirst, Valid: true},
			UserLastName:          sql.NullString{String: encLast, Valid: true},
			UserEmail:             encEmail,
			UserCreatedAt:         sql.NullTime{Time: now, Valid: true},
			UserUpdatedAt:         sql.NullTime{Time: now, Valid: true},
			GroupID:               uuid.New(),
			GroupName:             encGroup,
			GroupMetadata: pqtype.NullRawMessage{
				RawMessage: []byte(`{"domain":"finance"}`),
				Valid:      true,
			},
			GroupCreatedAt: sql.NullTime{Time: now, Valid: true},
			GroupUpdatedAt: sql.NullTime{Time: now, Valid: true},
		},
	}, nil)

	result, err := service.GetSubscriptions("1", "20", 1)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "ava@example.com", result[0].User.Email)
	assert.Equal(t, "Finance", result[0].Group.Name)
	assert.Equal(t, string(generated.SubscriptionStatusInactive), result[0].SubscriptionStatus)

	repo.AssertExpectations(t)
}

func TestSubscriptionService_CreateSubscription_New(t *testing.T) {
	repo := &testmocks.SubscriptionRepository{}
	encryptor := utils.NewEncryptor(testutils.TestEncryptionKey)
	service := services.NewSubscriptionService(repo, encryptor, nil)

	userID := uuid.New()
	now := time.Now().UTC()

	// confirm user exists
	repo.On("GetUserByID", mock.Anything, userID).Return(generated.GetUserByIDRow{}, nil)

	repo.On("GetSubscriptionByUserIDAndGroupIDAndStatus", mock.Anything, mock.Anything).Return(generated.Subscription{}, sql.ErrNoRows)

	repo.On("CreateSubscription", mock.Anything, mock.MatchedBy(func(param generated.CreateSubscriptionParams) bool {
		return param.UserID.Valid && param.GroupID.Valid && param.SubscriptionStatus.Valid
	})).Return(generated.Subscription{
		ID:                 uuid.New(),
		UserID:             uuid.NullUUID{UUID: userID, Valid: true},
		GroupID:            uuid.NullUUID{UUID: uuid.New(), Valid: true},
		SubscriptionStatus: generated.NullSubscriptionStatus{SubscriptionStatus: generated.SubscriptionStatusActive, Valid: true},
		CreatedAt:          sql.NullTime{Time: now, Valid: true},
		UpdatedAt:          sql.NullTime{Time: now, Valid: true},
	}, nil)

	repo.On("GetRoleByName", mock.Anything, utils.Hash("subscriber")).Return(generated.Role{ID: uuid.New()}, nil)
	repo.On("AssignRoleToUser", mock.Anything, generated.AssignRoleToUserParams{
		UserID: userID,
		RoleID: uuid.New(),
	}).Return(nil)

	subscription, err := service.CreateSubscription(models.CreateSubscriptionParams{
		UserID:             userID,
		GroupID:            uuid.New(),
		SubscriptionStatus: "active",
	})
	require.NoError(t, err)
	assert.Equal(t, 44, subscription.ID)
	assert.Equal(t, "active", subscription.SubscriptionStatus)

	repo.AssertExpectations(t)
}
