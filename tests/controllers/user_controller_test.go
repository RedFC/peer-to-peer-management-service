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

func newUserControllerTestApp() (*iris.Application, *testmocks.UserRepository, *testmocks.RelationRepository, *utils.Encryptor) {
	app := iris.New()
	userRepo := &testmocks.UserRepository{}
	relationRepo := &testmocks.RelationRepository{}
	encryptor := utils.NewEncryptor(testutils.TestEncryptionKey)

	userService := services.NewUserService(userRepo, encryptor)
	relationService := services.NewRelationService(relationRepo, encryptor)
	controller := controllers.NewUserController(userService, relationService)

	app.Get("/users", controller.GetUsers)
	app.Get("/users/{id}", controller.GetUserByID)
	app.Post("/users", controller.CreateUser)

	return app, userRepo, relationRepo, encryptor
}

func TestUserController_GetUsers_Success(t *testing.T) {
	app, userRepo, _, encryptor := newUserControllerTestApp()
	e := httptest.New(t, app)

	roleHashes := []string{utils.Hash(models.RoleSuperAdmin), utils.Hash(models.RoleAdmin)}
	userRepo.On("GetTotalUsersCount", mock.Anything, roleHashes).Return(int64(1), nil)

	now := time.Now().UTC()
	encEmail, err := encryptor.Encrypt("alice@example.com")
	require.NoError(t, err)
	encFirstName, err := encryptor.Encrypt("Alice")
	require.NoError(t, err)
	encLastName, err := encryptor.Encrypt("Doe")
	require.NoError(t, err)
	encRole, err := encryptor.Encrypt(models.RoleUser)
	require.NoError(t, err)
	encGroup, err := encryptor.Encrypt("Core Team")
	require.NoError(t, err)

	roleJSON, _ := json.Marshal([]string{encRole})
	groupJSON, _ := json.Marshal([]string{encGroup})

	userRepo.On(
		"ListUsers",
		mock.Anything,
		generated.ListUsersParams{
			Limit:   10,
			Offset:  0,
			Column3: roleHashes,
		},
	).Return([]generated.ListUsersRow{
		{
			ID:        uuid.New(),
			Email:     encEmail,
			EmailHash: utils.Hash("alice@example.com"),
			Password:  utils.Hash("password"),
			FirstName: sql.NullString{String: encFirstName, Valid: true},
			LastName:  sql.NullString{String: encLastName, Valid: true},
			IsDeleted: sql.NullBool{Bool: false, Valid: true},
			IsActive:  sql.NullBool{Bool: true, Valid: true},
			LastLogin: sql.NullTime{Time: now, Valid: true},
			CreatedAt: sql.NullTime{Time: now, Valid: true},
			UpdatedAt: sql.NullTime{Time: now, Valid: true},
			Roles:     roleJSON,
			Groups:    groupJSON,
		},
	}, nil)

	resp := e.GET("/users").Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("message").String().IsEqual("Users retrieved successfully")
	data := resp.Value("data").Array()
	data.Length().IsEqual(1)
	firstUser := data.Element(0).Object()
	firstUser.Value("first_name").String().IsEqual("Alice")
	firstUser.Value("last_name").String().IsEqual("Doe")
	firstUser.Value("email").String().IsEqual("alice@example.com")

	roleArray := firstUser.Value("role").Array()
	roleArray.Length().IsEqual(1)
	roleArray.Element(0).String().IsEqual(models.RoleUser)

	userRepo.AssertExpectations(t)
}

func TestUserController_GetUsers_Failure(t *testing.T) {
	app, userRepo, _, _ := newUserControllerTestApp()
	e := httptest.New(t, app)

	roleHashes := []string{utils.Hash(models.RoleSuperAdmin), utils.Hash(models.RoleAdmin)}
	userRepo.On("GetTotalUsersCount", mock.Anything, roleHashes).Return(int64(0), testutils.ErrBoom)

	resp := e.GET("/users").Expect().Status(http.StatusInternalServerError).JSON().Object()
	resp.Value("message").String().IsEqual("Failed to retrieve users count")

	userRepo.AssertExpectations(t)
}

func TestUserController_GetUserByID_Success(t *testing.T) {
	app, userRepo, _, encryptor := newUserControllerTestApp()
	e := httptest.New(t, app)

	userID := uuid.New()
	now := time.Now().UTC()
	encEmail, err := encryptor.Encrypt("bob@example.com")
	require.NoError(t, err)
	encFirstName, err := encryptor.Encrypt("Bob")
	require.NoError(t, err)
	encLastName, err := encryptor.Encrypt("Builder")
	require.NoError(t, err)
	encRole, err := encryptor.Encrypt(models.RoleAdmin)
	require.NoError(t, err)
	rolesJSON, _ := json.Marshal([]string{encRole})
	profileJSON, _ := json.Marshal(map[string]interface{}{"department": "Engineering"})

	userRepo.On("GetUserByID", mock.Anything, userID).Return(generated.GetUserByIDRow{
		ID:        userID,
		Email:     encEmail,
		EmailHash: utils.Hash("bob@example.com"),
		Password:  utils.Hash("secret"),
		FirstName: sql.NullString{String: encFirstName, Valid: true},
		LastName:  sql.NullString{String: encLastName, Valid: true},
		IsDeleted: sql.NullBool{Bool: false, Valid: true},
		IsActive:  sql.NullBool{Bool: true, Valid: true},
		LastLogin: sql.NullTime{Time: now, Valid: true},
		CreatedAt: sql.NullTime{Time: now, Valid: true},
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
		Roles:     rolesJSON,
		Groups:    []byte("[]"),
		Profile:   profileJSON,
	}, nil)

	resp := e.GET("/users/" + userID.String()).Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("message").String().IsEqual("User retrieved successfully")
	data := resp.Value("data").Object()
	data.Value("email").String().IsEqual("bob@example.com")
	data.Value("first_name").String().IsEqual("Bob")
	data.Value("last_name").String().IsEqual("Builder")

	userRepo.AssertExpectations(t)
}

func TestUserController_GetUserByID_InvalidUUID(t *testing.T) {
	app, _, _, _ := newUserControllerTestApp()
	e := httptest.New(t, app)

	resp := e.GET("/users/not-a-uuid").Expect().Status(http.StatusBadRequest).JSON().Object()
	resp.Value("message").String().IsEqual("Invalid user ID")
}

func TestUserController_CreateUser_Success(t *testing.T) {
	app, userRepo, relationRepo, encryptor := newUserControllerTestApp()
	e := httptest.New(t, app)

	now := time.Now().UTC()
	encEmail, err := encryptor.Encrypt("carol@example.com")
	require.NoError(t, err)
	encFirstName, err := encryptor.Encrypt("Carol")
	require.NoError(t, err)
	encLastName, err := encryptor.Encrypt("Danvers")
	require.NoError(t, err)

	newID := uuid.New()
	userRepo.On("CreateUser", mock.Anything, mock.Anything).Return(generated.User{
		ID:        newID,
		Email:     encEmail,
		FirstName: sql.NullString{String: encFirstName, Valid: true},
		LastName:  sql.NullString{String: encLastName, Valid: true},
		IsDeleted: sql.NullBool{Bool: false, Valid: true},
		IsActive:  sql.NullBool{Bool: false, Valid: true},
		LastLogin: sql.NullTime{Time: now, Valid: true},
		CreatedAt: sql.NullTime{Time: now, Valid: true},
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
	}, nil)

	relationRepo.On("GetRoleByName", mock.Anything, utils.Hash(models.RoleUser)).Return(generated.Role{ID: uuid.New()}, nil)
	relationRepo.On("AssignRoleToUser", mock.Anything, mock.Anything).Return(nil)

	payload := models.CreateUserParams{
		Email:     "carol@example.com",
		FirstName: "Carol",
		LastName:  "Danvers",
		Password:  "password",
	}

	resp := e.POST("/users").WithJSON(payload).Expect().Status(http.StatusCreated).JSON().Object()
	resp.Value("message").String().IsEqual("User created successfully")

	userRepo.AssertExpectations(t)
	relationRepo.AssertExpectations(t)
}
