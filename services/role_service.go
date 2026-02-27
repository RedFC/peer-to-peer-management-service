package services

import (
	"context"
	"database/sql"
	"p2p-management-service/db/generated"
	"p2p-management-service/models"
	"p2p-management-service/utils"
	"time"

	"github.com/google/uuid"
)

type RoleService struct {
	DB        RoleRepository
	Encryptor *utils.Encryptor
}

type RoleRepository interface {
	ListRoles(ctx context.Context) ([]generated.Role, error)
	GetRoleByID(ctx context.Context, id uuid.UUID) (generated.Role, error)
	CreateRole(ctx context.Context, arg generated.CreateRoleParams) (generated.Role, error)
	UpdateRole(ctx context.Context, arg generated.UpdateRoleParams) (generated.Role, error)
	DeleteRole(ctx context.Context, arg generated.DeleteRoleParams) error
}

func NewRoleService(db RoleRepository, encryptor *utils.Encryptor) *RoleService {
	return &RoleService{
		DB:        db,
		Encryptor: encryptor,
	}
}

// GetRoles retrieves all roles from the database.
func (s *RoleService) GetRoles() ([]models.Role, error) {
	roles, err := s.DB.ListRoles(context.Background())
	if err != nil {
		return []models.Role{}, err
	}
	var result []models.Role
	for _, role := range roles {

		DecryptObject := []string{
			role.Name,
			role.Description.String,
		}

		// Decrypt them
		decrypted, err := s.Encryptor.DecryptMultiple(DecryptObject)
		if err != nil {
			return []models.Role{}, err
		}

		result = append(result, models.Role{
			ID:          role.ID,
			Name:        decrypted[0],
			Description: decrypted[1],
			IsDeleted:   role.IsDeleted.Bool,
			CreatedAt:   role.CreatedAt.Time.UTC().String(),
			UpdatedAt:   role.UpdatedAt.Time.UTC().String(),
		})
	}
	return result, nil
}

// GetRoleByID retrieves a role by its ID.
func (s *RoleService) GetRoleByID(id uuid.UUID) (*models.Role, error) {
	role, err := s.DB.GetRoleByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	DecryptObject := []string{
		role.Name,
		role.Description.String,
	}

	// Decrypt them
	decrypted, err := s.Encryptor.DecryptMultiple(DecryptObject)
	if err != nil {
		return nil, err
	}

	return &models.Role{
		ID:          role.ID,
		Name:        decrypted[0],
		Description: decrypted[1],
		IsDeleted:   role.IsDeleted.Bool,
		CreatedAt:   role.CreatedAt.Time.UTC().String(),
		UpdatedAt:   role.UpdatedAt.Time.UTC().String(),
	}, nil
}

// CreateRole creates a new role in the database.
func (s *RoleService) CreateRole(role models.CreateRoleParams) (*models.Role, error) {

	encName, err := s.Encryptor.Encrypt(role.Name)
	if err != nil {
		return nil, err
	}

	encDesc, err := s.Encryptor.Encrypt(role.Description)
	if err != nil {
		return nil, err
	}

	newRole, err := s.DB.CreateRole(context.Background(), generated.CreateRoleParams{
		NameHash:    utils.Hash(role.Name),
		Name:        encName,
		Description: sql.NullString{Valid: true, String: encDesc},
		IsDeleted:   sql.NullBool{Valid: true, Bool: false},
		CreatedAt:   sql.NullTime{Time: time.Now().UTC(), Valid: true},
		UpdatedAt:   sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &models.Role{
		ID:          newRole.ID,
		Name:        role.Name,
		Description: role.Description,
		IsDeleted:   newRole.IsDeleted.Bool,
		CreatedAt:   newRole.CreatedAt.Time.UTC().String(),
		UpdatedAt:   newRole.UpdatedAt.Time.UTC().String(),
	}, nil
}

// UpdateRole updates an existing role in the database.
func (s *RoleService) UpdateRole(id uuid.UUID, role models.UpdateRoleParams) (*models.Role, error) {

	encName, err := s.Encryptor.Encrypt(role.Name)
	if err != nil {
		return nil, err
	}

	encDesc, err := s.Encryptor.Encrypt(role.Description)
	if err != nil {
		return nil, err
	}

	updatedRole, err := s.DB.UpdateRole(context.Background(), generated.UpdateRoleParams{
		ID:          id,
		NameHash:    utils.Hash(role.Name),
		Name:        encName,
		Description: sql.NullString{Valid: true, String: encDesc},
		UpdatedAt:   sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &models.Role{
		ID:          updatedRole.ID,
		Name:        role.Name,
		Description: role.Description,
		IsDeleted:   updatedRole.IsDeleted.Bool,
		CreatedAt:   updatedRole.CreatedAt.Time.UTC().String(),
		UpdatedAt:   updatedRole.UpdatedAt.Time.UTC().String(),
	}, nil
}

// DeleteRole deletes a role by its ID.
func (s *RoleService) DeleteRole(id uuid.UUID) error {
	err := s.DB.DeleteRole(context.Background(), generated.DeleteRoleParams{
		ID:        id,
		UpdatedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
	if err != nil {
		return err
	}
	return nil
}
