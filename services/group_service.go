package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"p2p-management-service/db/generated"
	"p2p-management-service/models"
	"p2p-management-service/utils"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type GroupService struct {
	DB        GroupRepository
	Encryptor *utils.Encryptor
}

type GroupRepository interface {
	ListGroups(ctx context.Context, arg generated.ListGroupsParams) ([]generated.Group, error)
	GetTotalGroupsCount(ctx context.Context) (int64, error)
	GetGroupByID(ctx context.Context, id uuid.UUID) (generated.Group, error)
	CreateGroup(ctx context.Context, arg generated.CreateGroupParams) (generated.Group, error)
	UpdateGroup(ctx context.Context, arg generated.UpdateGroupParams) (generated.Group, error)
	DeleteGroup(ctx context.Context, arg generated.DeleteGroupParams) error
}

func NewGroupService(db GroupRepository, encryptor *utils.Encryptor) *GroupService {
	return &GroupService{
		DB:        db,
		Encryptor: encryptor,
	}
}

// GetGroups retrieves all groups from the database.
func (s *GroupService) GetGroups(pageNo string, pageSize string) ([]models.Group, error) {
	pageNoInt, err := strconv.Atoi(pageNo)
	if err != nil {
		return nil, err
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		return nil, err
	}
	groups, err := s.DB.ListGroups(context.Background(), generated.ListGroupsParams{
		Limit:  int32(pageSizeInt),
		Offset: int32((pageNoInt - 1) * pageSizeInt),
	})
	if err != nil {
		return []models.Group{}, err
	}
	var result []models.Group
	for _, group := range groups {

		decryptName, err := s.Encryptor.Decrypt(group.Name)
		if err != nil {
			return []models.Group{}, err
		}

		result = append(result, models.Group{
			ID:        group.ID,
			Name:      decryptName,
			Metadata:  group.Metadata.RawMessage,
			IsDeleted: group.IsDeleted.Bool,
			CreatedAt: group.CreatedAt.Time.UTC().String(),
			UpdatedAt: group.UpdatedAt.Time.UTC().String(),
		})
	}
	return result, nil
}

// GetTotalGroupsCount retrieves the total count of groups.
func (s *GroupService) GetTotalGroupsCount() (int64, error) {
	countResult, err := s.DB.GetTotalGroupsCount(context.Background())
	if err != nil {
		return 0, err
	}

	return countResult, nil
}

// GetGroupByID retrieves a group by its ID.
func (s *GroupService) GetGroupByID(id uuid.UUID) (*models.Group, error) {
	group, err := s.DB.GetGroupByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	decryptName, err := s.Encryptor.Decrypt(group.Name)
	if err != nil {
		return nil, err
	}

	return &models.Group{
		ID:        group.ID,
		Name:      decryptName,
		Metadata:  group.Metadata.RawMessage,
		IsDeleted: group.IsDeleted.Bool,
		CreatedAt: group.CreatedAt.Time.UTC().String(),
		UpdatedAt: group.UpdatedAt.Time.UTC().String(),
	}, nil
}

// CreateGroup creates a new group in the database.
func (s *GroupService) CreateGroup(group models.CreateGroupParams) (*models.Group, error) {

	metaDataJSON, err := json.Marshal(group.Metadata)
	if err != nil {
		log.Println("Encryption failed:", err)
		return nil, err
	}

	encName, err := s.Encryptor.Encrypt(group.Name)
	if err != nil {
		log.Println("Encryption failed:", err)
		return nil, err
	}

	newGroup, err := s.DB.CreateGroup(context.Background(), generated.CreateGroupParams{
		NameHash: utils.Hash(group.Name),
		Name:     encName,
		Metadata: pqtype.NullRawMessage{
			RawMessage: metaDataJSON,
			Valid:      true,
		},
		IsDeleted: sql.NullBool{Bool: false, Valid: true},
		CreatedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &models.Group{
		ID:        newGroup.ID,
		Name:      group.Name,
		Metadata:  newGroup.Metadata.RawMessage,
		IsDeleted: newGroup.IsDeleted.Bool,
		CreatedAt: newGroup.CreatedAt.Time.UTC().String(),
		UpdatedAt: newGroup.UpdatedAt.Time.UTC().String(),
	}, nil
}

// UpdateGroup updates an existing group in the database.
func (s *GroupService) UpdateGroup(id uuid.UUID, group models.UpdateGroupParams) (*models.Group, error) {
	encName, err := s.Encryptor.Encrypt(group.Name)
	if err != nil {
		return nil, err
	}

	metaDataJSON, err := json.Marshal(group.Metadata)
	if err != nil {
		log.Println("Encryption failed:", err)
		return nil, err
	}

	updatedGroup, err := s.DB.UpdateGroup(context.Background(), generated.UpdateGroupParams{
		ID:       id,
		NameHash: utils.Hash(group.Name),
		Name:     encName,
		Metadata: pqtype.NullRawMessage{
			RawMessage: metaDataJSON,
			Valid:      true,
		},
		UpdatedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &models.Group{
		ID:        updatedGroup.ID,
		Name:      group.Name,
		Metadata:  updatedGroup.Metadata.RawMessage,
		IsDeleted: updatedGroup.IsDeleted.Bool,
		CreatedAt: updatedGroup.CreatedAt.Time.UTC().String(),
		UpdatedAt: updatedGroup.UpdatedAt.Time.UTC().String(),
	}, nil
}

// DeleteGroup deletes a group by its ID.
func (s *GroupService) DeleteGroup(id uuid.UUID) error {
	err := s.DB.DeleteGroup(context.Background(), generated.DeleteGroupParams{
		ID:        id,
		UpdatedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
	if err != nil {
		return err
	}
	return nil
}
