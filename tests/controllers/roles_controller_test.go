package tests

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"p2p-management-service/controllers"
	"p2p-management-service/db/generated"
	"p2p-management-service/models"
	"p2p-management-service/services"
	testmocks "p2p-management-service/tests/mocks"
	testutils "p2p-management-service/tests/testutils"
	"p2p-management-service/utils"
)

func newRoleControllerTestApp() (*iris.Application, *testmocks.RoleRepository, *utils.Encryptor) {
	app := iris.New()
	roleRepo := &testmocks.RoleRepository{}
	encryptor := utils.NewEncryptor(testutils.TestEncryptionKey)

	roleService := services.NewRoleService(roleRepo, encryptor)
	controller := controllers.NewRoleController(roleService)

	app.Get("/roles", controller.GetRoles)

	return app, roleRepo, encryptor
}

func TestRoleController_GetRoles_Success(t *testing.T) {
	app, roleRepo, encryptor := newRoleControllerTestApp()
	e := httptest.New(t, app)

	now := time.Now().UTC()
	encName, err := encryptor.Encrypt(models.RoleAdmin)
	require.NoError(t, err)
	encDesc, err := encryptor.Encrypt("Platform administrator")
	require.NoError(t, err)

	roleRepo.On("ListRoles", mock.Anything).Return([]generated.Role{
		{
			ID:          uuid.Nil,
			Name:        encName,
			NameHash:    utils.Hash(models.RoleAdmin),
			Description: sql.NullString{String: encDesc, Valid: true},
			IsDeleted:   sql.NullBool{Bool: false, Valid: true},
			CreatedAt:   sql.NullTime{Time: now, Valid: true},
			UpdatedAt:   sql.NullTime{Time: now, Valid: true},
		},
	}, nil)

	resp := e.GET("/roles").Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("message").String().IsEqual("Roles retrieved successfully")
	data := resp.Value("data").Array()
	data.Length().IsEqual(1)
	data.Element(0).Object().Value("name").String().IsEqual(models.RoleAdmin)

	roleRepo.AssertExpectations(t)
}

func TestRoleController_GetRoles_Failure(t *testing.T) {
	app, roleRepo, _ := newRoleControllerTestApp()
	e := httptest.New(t, app)

	roleRepo.On("ListRoles", mock.Anything).Return(nil, testutils.ErrBoom)

	resp := e.GET("/roles").Expect().Status(http.StatusInternalServerError).JSON().Object()
	resp.Value("message").String().IsEqual("Failed to retrieve roles")

	roleRepo.AssertExpectations(t)
}
