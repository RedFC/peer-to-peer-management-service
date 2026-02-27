package tests

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"github.com/sqlc-dev/pqtype"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"p2p-management-service/controllers"
	"p2p-management-service/db/generated"
	"p2p-management-service/services"
	testmocks "p2p-management-service/tests/mocks"
	testutils "p2p-management-service/tests/testutils"
	"p2p-management-service/utils"
)

func newSubscriptionControllerTestApp() (*iris.Application, *testmocks.SubscriptionRepository, *utils.Encryptor) {
	app := iris.New()
	repo := &testmocks.SubscriptionRepository{}
	encryptor := utils.NewEncryptor(testutils.TestEncryptionKey)

	service := services.NewSubscriptionService(repo, encryptor, nil)
	controller := controllers.NewSubscriptionController(service)

	app.Get("/subscriber", controller.GetSubscriptions)

	return app, repo, encryptor
}

func TestSubscriptionController_GetSubscriptions_Success(t *testing.T) {
	app, repo, encryptor := newSubscriptionControllerTestApp()
	e := httptest.New(t, app)

	now := time.Now().UTC()
	userID := uuid.New()
	encFirst, err := encryptor.Encrypt("Dana")
	require.NoError(t, err)
	encLast, err := encryptor.Encrypt("Scully")
	require.NoError(t, err)
	encEmail, err := encryptor.Encrypt("dana@example.com")
	require.NoError(t, err)
	encGroup, err := encryptor.Encrypt("Investigation")
	require.NoError(t, err)

	repo.On("GetTotalSubscriptionsCount", mock.Anything).Return(int64(1), nil)
	repo.On("ListSubscriptions", mock.Anything, generated.ListSubscriptionsParams{
		Limit:  10,
		Offset: 0,
	}).Return([]generated.ListSubscriptionsRow{
		{
			SubscriptionID:        uuid.New(),
			SubscriptionStatus:    generated.NullSubscriptionStatus{SubscriptionStatus: generated.SubscriptionStatusActive, Valid: true},
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
				RawMessage: []byte(`{"tier":"gold"}`),
				Valid:      true,
			},
			GroupCreatedAt: sql.NullTime{Time: now, Valid: true},
			GroupUpdatedAt: sql.NullTime{Time: now, Valid: true},
		},
	}, nil)

	resp := e.GET("/subscriber").Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("message").String().IsEqual("Subscriptions retrieved successfully")
	data := resp.Value("data").Array()
	data.Length().IsEqual(1)
	first := data.Element(0).Object()
	first.Value("subscription_status").String().IsEqual(string(generated.SubscriptionStatusActive))
	first.Value("User").Object().Value("email").String().IsEqual("dana@example.com")
	first.Value("Group").Object().Value("name").String().IsEqual("Investigation")

	repo.AssertExpectations(t)
}

func TestSubscriptionController_GetSubscriptions_Failure(t *testing.T) {
	app, repo, _ := newSubscriptionControllerTestApp()
	e := httptest.New(t, app)

	repo.On("GetTotalSubscriptionsCount", mock.Anything).Return(int64(0), testutils.ErrBoom)

	resp := e.GET("/subscriber").Expect().Status(http.StatusInternalServerError).JSON().Object()
	resp.Value("message").String().IsEqual("Failed to retrieve subscriptions count")

	repo.AssertExpectations(t)
}
