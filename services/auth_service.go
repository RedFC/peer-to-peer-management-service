package services

import (
	"context"
	"errors"
	"p2p-management-service/db/generated"
	"p2p-management-service/models"
	"p2p-management-service/utils"
)

type AuthService struct {
	DB        *generated.Queries
	Encryptor *utils.Encryptor
}

func NewAuthService(db *generated.Queries, encryptor *utils.Encryptor) *AuthService {
	return &AuthService{
		DB:        db,
		Encryptor: encryptor,
	}
}

func (s *AuthService) Authenticate(email, password string) (models.LoginResponse, error) {

	// hash email
	email = utils.Hash(email)

	// check for user existence
	user, err := s.DB.GetUserByEmailHash(context.Background(), email) // This is a placeholder, implement actual logic
	if err != nil {
		return models.LoginResponse{}, err
	}

	roles, err := s.DB.GetUserRoles(context.Background(), user.ID)
	if err != nil {
		return models.LoginResponse{}, err
	}

	// Decrypt roles
	if len(roles) == 0 {
		return models.LoginResponse{}, errors.New("user has no roles assigned")
	}

	decryptPassword, err := s.Encryptor.Decrypt(user.Password)
	if err != nil {
		return models.LoginResponse{}, err
	}

	if decryptPassword != password {
		return models.LoginResponse{}, errors.New("invalid credentials")
	}

	var roleClaims = make([]string, len(roles))
	var Roles = make([]models.Role, len(roles))
	for i, role := range roles {
		decryptedRoleName, err := s.Encryptor.Decrypt(role.Name)
		if err != nil {
			// skip this role if name decryption fails
			continue
		}
		decryptedRoleDescription, err := s.Encryptor.Decrypt(role.Description.String)
		if err != nil {
			// skip this role if description decryption fails
			continue
		}
		roleClaims[i] = decryptedRoleName
		Roles[i].Name = decryptedRoleName
		Roles[i].ID = role.ID
		Roles[i].Description = decryptedRoleDescription
		Roles[i].IsDeleted = role.IsDeleted.Bool
		Roles[i].CreatedAt = role.CreatedAt.Time.UTC().String()
		Roles[i].UpdatedAt = role.UpdatedAt.Time.UTC().String()
	}

	token, err := utils.GenerateJWT(email, roleClaims)
	if err != nil {
		return models.LoginResponse{}, err
	}

	// Collect encrypted fields
	DecryptObject := []string{
		user.FirstName.String,
		user.LastName.String,
		user.Email,
	}

	// Decrypt them
	decrypted, err := s.Encryptor.DecryptMultiple(DecryptObject)
	if err != nil {
		return models.LoginResponse{}, err
	}

	return models.LoginResponse{Token: token,
		User: models.UserResponseParams{
			ID:              user.ID,
			FirstName:       decrypted[0],
			LastName:        decrypted[1],
			Email:           decrypted[2],
			IsDeleted:       user.IsDeleted.Bool,
			CreatedAt:       user.CreatedAt.Time.UTC().String(),
			UpdatedAt:       user.UpdatedAt.Time.UTC().String(),
			IsActive:        user.IsActive.Bool,
			IsPasswordReset: user.IsPasswordReset.Bool,
			Role:            Roles,
		}}, nil
}
