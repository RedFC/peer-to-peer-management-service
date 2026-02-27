package tests

import (
	"database/sql"
	"encoding/json"
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

func TestGroupService_GetGroups(t *testing.T) {
	repo := &testmocks.GroupRepository{}
	encryptor := utils.NewEncryptor(testutils.TestEncryptionKey)
	service := services.NewGroupService(repo, encryptor)

	now := time.Now().UTC()
	encName, err := encryptor.Encrypt("Analytics")
	require.NoError(t, err)

	metaBytes, _ := json.Marshal(map[string]string{"tier": "gold"})

	repo.On("ListGroups", mock.Anything, generated.ListGroupsParams{
		Limit:  5,
		Offset: 5,
	}).Return([]generated.Group{
		{
			ID:       uuid.New(),
			Name:     encName,
			NameHash: utils.Hash("Analytics"),
			Metadata: pqtype.NullRawMessage{
				RawMessage: metaBytes,
				Valid:      true,
			},
			IsDeleted: sql.NullBool{Bool: false, Valid: true},
			CreatedAt: sql.NullTime{Time: now, Valid: true},
			UpdatedAt: sql.NullTime{Time: now, Valid: true},
		},
	}, nil)

	groups, err := service.GetGroups("2", "5")
	require.NoError(t, err)
	require.Len(t, groups, 1)
	assert.Equal(t, "Analytics", groups[0].Name)
	assert.Equal(t, `{"tier":"gold"}`, groups[0].Metadata)

	repo.AssertExpectations(t)
}

func TestGroupService_CreateGroup(t *testing.T) {
	repo := &testmocks.GroupRepository{}
	encryptor := utils.NewEncryptor(testutils.TestEncryptionKey)
	service := services.NewGroupService(repo, encryptor)

	now := time.Now().UTC()
	encName, err := encryptor.Encrypt("Creators")
	require.NoError(t, err)

	repo.On("CreateGroup", mock.Anything, mock.MatchedBy(func(param generated.CreateGroupParams) bool {
		return param.NameHash == utils.Hash("Creators") &&
			param.Metadata.Valid &&
			len(param.Metadata.RawMessage) > 0
	})).Return(generated.Group{
		ID:       uuid.New(),
		Name:     encName,
		NameHash: utils.Hash("Creators"),
		Metadata: pqtype.NullRawMessage{
			RawMessage: []byte(`{"feature":"beta"}`),
			Valid:      true,
		},
		IsDeleted: sql.NullBool{Bool: false, Valid: true},
		CreatedAt: sql.NullTime{Time: now, Valid: true},
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
	}, nil)

	group, err := service.CreateGroup(models.CreateGroupParams{
		Name:     "Creators",
		Metadata: map[string]interface{}{"feature": "beta"},
	})
	require.NoError(t, err)
	assert.Equal(t, "Creators", group.Name)

	repo.AssertExpectations(t)
}
