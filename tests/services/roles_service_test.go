package tests

import (
	"database/sql"
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

func TestRoleService_GetRoles(t *testing.T) {
	repo := &testmocks.RoleRepository{}
	encryptor := utils.NewEncryptor(testutils.TestEncryptionKey)
	service := services.NewRoleService(repo, encryptor)

	now := time.Now().UTC()
	encName, err := encryptor.Encrypt(models.RoleUser)
	require.NoError(t, err)
	encDesc, err := encryptor.Encrypt("Standard user role")
	require.NoError(t, err)

	repo.On("ListRoles", mock.Anything).Return([]generated.Role{
		{
			ID:          uuid.New(),
			Name:        encName,
			NameHash:    utils.Hash(models.RoleUser),
			Description: sql.NullString{String: encDesc, Valid: true},
			IsDeleted:   sql.NullBool{Bool: false, Valid: true},
			CreatedAt:   sql.NullTime{Time: now, Valid: true},
			UpdatedAt:   sql.NullTime{Time: now, Valid: true},
		},
	}, nil)

	roles, err := service.GetRoles()
	require.NoError(t, err)
	require.Len(t, roles, 1)
	assert.Equal(t, models.RoleUser, roles[0].Name)
	assert.Equal(t, "Standard user role", roles[0].Description)

	repo.AssertExpectations(t)
}

func TestRoleService_CreateRole(t *testing.T) {
	repo := &testmocks.RoleRepository{}
	encryptor := utils.NewEncryptor(testutils.TestEncryptionKey)
	service := services.NewRoleService(repo, encryptor)

	now := time.Now().UTC()
	encName, err := encryptor.Encrypt("Reviewer")
	require.NoError(t, err)
	encDesc, err := encryptor.Encrypt("Can review submissions")
	require.NoError(t, err)

	repo.On("CreateRole", mock.Anything, mock.MatchedBy(func(param generated.CreateRoleParams) bool {
		return param.NameHash == utils.Hash("Reviewer") && param.Description.Valid
	})).Return(generated.Role{
		ID:          uuid.New(),
		Name:        encName,
		NameHash:    utils.Hash("Reviewer"),
		Description: sql.NullString{String: encDesc, Valid: true},
		IsDeleted:   sql.NullBool{Bool: false, Valid: true},
		CreatedAt:   sql.NullTime{Time: now, Valid: true},
		UpdatedAt:   sql.NullTime{Time: now, Valid: true},
	}, nil)

	role, err := service.CreateRole(models.CreateRoleParams{
		Name:        "Reviewer",
		Description: "Can review submissions",
	})
	require.NoError(t, err)
	assert.Equal(t, "Reviewer", role.Name)
	assert.Equal(t, "Can review submissions", role.Description)

	repo.AssertExpectations(t)
}

func TestRoleService_DeleteRole(t *testing.T) {
	repo := &testmocks.RoleRepository{}
	encryptor := utils.NewEncryptor(testutils.TestEncryptionKey)
	service := services.NewRoleService(repo, encryptor)

	repo.On("DeleteRole", mock.Anything, mock.MatchedBy(func(param generated.DeleteRoleParams) bool {
		return param.ID == uuid.New() && param.UpdatedAt.Valid
	})).Return(nil)

	err := service.DeleteRole(uuid.New())
	require.NoError(t, err)
	repo.AssertExpectations(t)
}
