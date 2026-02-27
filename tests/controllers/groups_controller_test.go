package tests

import (
	"database/sql"
	"encoding/json"
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

func newGroupControllerTestApp() (*iris.Application, *testmocks.GroupRepository, *utils.Encryptor) {
	app := iris.New()
	groupRepo := &testmocks.GroupRepository{}
	encryptor := utils.NewEncryptor(testutils.TestEncryptionKey)

	groupService := services.NewGroupService(groupRepo, encryptor)
	controller := controllers.NewGroupController(groupService)

	app.Get("/groups", controller.GetGroups)

	return app, groupRepo, encryptor
}

func TestGroupController_GetGroups_Success(t *testing.T) {
	app, groupRepo, encryptor := newGroupControllerTestApp()
	e := httptest.New(t, app)

	now := time.Now().UTC()
	encName, err := encryptor.Encrypt("Core Team")
	require.NoError(t, err)

	metadataBytes, _ := json.Marshal(map[string]any{"region": "NA"})

	groupRepo.On("GetTotalGroupsCount", mock.Anything).Return(int64(1), nil)
	groupRepo.On("ListGroups", mock.Anything, generated.ListGroupsParams{
		Limit:  10,
		Offset: 0,
	}).Return([]generated.Group{
		{
			ID:       uuid.Nil,
			Name:     encName,
			NameHash: utils.Hash("Core Team"),
			Metadata: pqtype.NullRawMessage{
				RawMessage: metadataBytes,
				Valid:      true,
			},
			IsDeleted: sql.NullBool{Bool: false, Valid: true},
			CreatedAt: sql.NullTime{Time: now, Valid: true},
			UpdatedAt: sql.NullTime{Time: now, Valid: true},
		},
	}, nil)

	resp := e.GET("/groups").Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("message").String().IsEqual("Groups retrieved successfully")
	data := resp.Value("data").Array()
	data.Length().IsEqual(1)
	data.Element(0).Object().Value("name").String().IsEqual("Core Team")

	groupRepo.AssertExpectations(t)
}

func TestGroupController_GetGroups_Failure(t *testing.T) {
	app, groupRepo, _ := newGroupControllerTestApp()
	e := httptest.New(t, app)

	groupRepo.On("GetTotalGroupsCount", mock.Anything).Return(int64(0), testutils.ErrBoom)

	resp := e.GET("/groups").Expect().Status(http.StatusInternalServerError).JSON().Object()
	resp.Value("message").String().IsEqual("Failed to retrieve groups count")

	groupRepo.AssertExpectations(t)
}
