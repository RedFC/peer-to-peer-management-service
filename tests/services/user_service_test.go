package tests

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestUserService_GetUsers(t *testing.T) {
	repo := &testmocks.UserRepository{}
	encryptor := utils.NewEncryptor(testutils.TestEncryptionKey)
	service := services.NewUserService(repo, encryptor)

	roleHashes := []string{utils.Hash(models.RoleUser)}
	now := time.Now().UTC()

	encEmail, err := encryptor.Encrypt("test@example.com")
	require.NoError(t, err)
	encFirst, err := encryptor.Encrypt("Testy")
	require.NoError(t, err)
	encLast, err := encryptor.Encrypt("McTest")
	require.NoError(t, err)
	encRole, err := encryptor.Encrypt(models.RoleUser)
	require.NoError(t, err)
	encGroup, err := encryptor.Encrypt("Core")
	require.NoError(t, err)

	roleJSON, _ := json.Marshal([]string{encRole})
	groupJSON, _ := json.Marshal([]string{encGroup})

	repo.On("ListUsers", mock.Anything, generated.ListUsersParams{
		Limit:   10,
		Offset:  0,
		Column3: roleHashes,
	}).Return([]generated.ListUsersRow{
		{
			ID:        uuid.New(),
			Email:     encEmail,
			EmailHash: utils.Hash("test@example.com"),
			FirstName: sql.NullString{String: encFirst, Valid: true},
			LastName:  sql.NullString{String: encLast, Valid: true},
			IsDeleted: sql.NullBool{Bool: false, Valid: true},
			IsActive:  sql.NullBool{Bool: true, Valid: true},
			LastLogin: sql.NullTime{Time: now, Valid: true},
			CreatedAt: sql.NullTime{Time: now, Valid: true},
			UpdatedAt: sql.NullTime{Time: now, Valid: true},
			Roles:     roleJSON,
			Groups:    groupJSON,
		},
	}, nil)

	users, err := service.GetUsers("1", "10", 1, roleHashes, "", "")
	require.NoError(t, err)
	require.Len(t, users, 1)
	assert.Equal(t, "test@example.com", users[0].Email)
	assert.Equal(t, "Testy", users[0].FirstName)
	assert.Equal(t, []string{models.RoleUser}, users[0].Role)

	repo.AssertExpectations(t)
}

func TestUserService_CreateUser(t *testing.T) {
	repo := &testmocks.UserRepository{}
	encryptor := utils.NewEncryptor(testutils.TestEncryptionKey)
	service := services.NewUserService(repo, encryptor)

	now := time.Now().UTC()
	encEmail, err := encryptor.Encrypt("create@example.com")
	require.NoError(t, err)
	encFirst, err := encryptor.Encrypt("Create")
	require.NoError(t, err)
	encLast, err := encryptor.Encrypt("User")
	require.NoError(t, err)

	repo.On("CreateUser", mock.Anything, mock.MatchedBy(func(param generated.CreateUserParams) bool {
		return param.EmailHash == utils.Hash("create@example.com") &&
			param.FirstName.String != "" &&
			param.LastName.String != ""
	})).Return(generated.User{
		ID:        uuid.New(),
		Email:     encEmail,
		FirstName: sql.NullString{String: encFirst, Valid: true},
		LastName:  sql.NullString{String: encLast, Valid: true},
		IsDeleted: sql.NullBool{Bool: false, Valid: true},
		IsActive:  sql.NullBool{Bool: false, Valid: true},
		LastLogin: sql.NullTime{Time: now, Valid: true},
		CreatedAt: sql.NullTime{Time: now, Valid: true},
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
	}, nil)

	result, err := service.CreateUser(models.CreateUserParams{
		Email:     "create@example.com",
		FirstName: "Create",
		LastName:  "User",
		Password:  "password",
	})
	require.NoError(t, err)
	assert.Equal(t, "create@example.com", result.Email)

	repo.AssertExpectations(t)
}

func TestUserService_DeleteUser(t *testing.T) {
	repo := &testmocks.UserRepository{}
	encryptor := utils.NewEncryptor(testutils.TestEncryptionKey)
	service := services.NewUserService(repo, encryptor)

	id := uuid.New()
	repo.On("DeleteUser", mock.Anything, mock.MatchedBy(func(param generated.DeleteUserParams) bool {
		return param.ID == id && param.UpdatedAt.Valid
	})).Return(nil)

	err := service.DeleteUser(id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUserService_GetTotalUsersCount(t *testing.T) {
	repo := &testmocks.UserRepository{}
	encryptor := utils.NewEncryptor(testutils.TestEncryptionKey)
	service := services.NewUserService(repo, encryptor)

	roleHashes := []string{"hash"}
	repo.On("GetTotalUsersCount", mock.Anything, roleHashes).Return(int64(5), nil)

	count, err := service.GetTotalUsersCount(roleHashes, "", "")
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)

	repo.AssertExpectations(t)
}
