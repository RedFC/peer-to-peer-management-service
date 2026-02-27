package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"p2p-management-service/db/generated"
	"p2p-management-service/models"
	"p2p-management-service/utils"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	DB        UserRepository
	Encryptor *utils.Encryptor
}

type UserRepository interface {
	GetTotalUsersCount(ctx context.Context, arg generated.GetTotalUsersCountParams) (int64, error)
	ListUsers(ctx context.Context, arg generated.ListUsersParams) ([]generated.ListUsersRow, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (generated.GetUserByIDRow, error)
	CreateUser(ctx context.Context, arg generated.CreateUserParams) (generated.User, error)
	UpdateUserProfile(ctx context.Context, arg generated.UpdateUserProfileParams) (generated.User, error)
	UpdateUser(ctx context.Context, arg generated.UpdateUserParams) (generated.User, error)
	DeleteUser(ctx context.Context, arg generated.DeleteUserParams) error
}

func NewUserService(db UserRepository, encryptor *utils.Encryptor) *UserService {
	return &UserService{
		DB:        db,
		Encryptor: encryptor,
	}
}

// total users count
func (s *UserService) GetTotalUsersCount(roles_hashes []string, roleIdParam string, groupIdParam string) (int64, error) {

	var roleID uuid.UUID
	if roleIdParam != "" {
		idVal, _ := uuid.Parse(roleIdParam)
		roleID = idVal
	} else {
		roleID = uuid.Nil
	}

	var groupID uuid.UUID
	if groupIdParam != "" {
		idVal, _ := uuid.Parse(groupIdParam)
		groupID = idVal
	} else {
		groupID = uuid.Nil
	}

	countResult, err := s.DB.GetTotalUsersCount(context.Background(), generated.GetTotalUsersCountParams{
		Column1: []string{},
		Column2: roleID,
		Column3: groupID,
	})
	if err != nil {
		return 0, err
	}

	return countResult, nil
}

// List All Users With Pagination
func (s *UserService) GetUsers(pageNo string, pageSize string, totalCount int64, roles_hashes []string, roleIdParam string, groupIdParam string) ([]models.UserResponseParams, error) {
	pageNoInt, err := strconv.Atoi(pageNo)
	if err != nil {
		return nil, err
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		return nil, err
	}

	var roleID uuid.UUID
	if roleIdParam != "" {
		idVal, _ := uuid.Parse(roleIdParam)
		roleID = idVal
	} else {
		roleID = uuid.Nil
	}

	var groupID uuid.UUID
	if groupIdParam != "" {
		idVal, _ := uuid.Parse(groupIdParam)
		groupID = idVal
	} else {
		groupID = uuid.Nil
	}

	fmt.Println("params", pageSizeInt, pageNoInt, roleID, groupID)
	users, err := s.DB.ListUsers(context.Background(), generated.ListUsersParams{
		Limit:   int32(pageSizeInt),
		Offset:  int32((pageNoInt - 1) * pageSizeInt),
		Column3: roles_hashes,
		Column4: roleID,
		Column5: groupID,
	})
	if err != nil {
		return []models.UserResponseParams{}, err
	}

	var userResponses []models.UserResponseParams

	// decrypt users
	for _, user := range users {

		var roles []models.Role
		if user.Roles != nil {
			// Unmarshal roles from []byte to []string
			err := json.Unmarshal(user.Roles.([]byte), &roles)
			if err != nil {
				// handle error
				return nil, err
			}

			// Decrypt each role if needed
			for i, role := range roles {
				decryptedName, err := s.Encryptor.Decrypt(role.Name)
				decryptedDescription, err := s.Encryptor.Decrypt(role.Description)
				if err == nil {
					roles[i].Name = decryptedName
					roles[i].Description = decryptedDescription
					roles[i].ID = role.ID
				}
			}
		}

		var groups []models.Group
		if user.Groups != nil {
			// Unmarshal groups from []byte to []string
			err := json.Unmarshal(user.Groups.([]byte), &groups)
			if err != nil {
				// handle error
				return nil, err
			}

			// Decrypt each group if needed
			for i, group := range groups {
				decryptedGroupName, err := s.Encryptor.Decrypt(group.Name)
				if err == nil {
					groups[i].Name = decryptedGroupName
					groups[i].Metadata = group.Metadata
					groups[i].ID = group.ID
				}
			}
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
			return []models.UserResponseParams{}, err
		}

		userResponses = append(userResponses, models.UserResponseParams{
			ID:              user.ID,
			Email:           decrypted[2],
			FirstName:       decrypted[0],
			LastName:        decrypted[1],
			IsDeleted:       user.IsDeleted.Bool,
			IsActive:        user.IsActive.Bool,
			IsPasswordReset: user.IsPasswordReset.Bool,
			Role:            roles,
			Group:           groups,
			LastLogin:       user.LastLogin.Time.UTC().String(),
			CreatedAt:       user.CreatedAt.Time.UTC().String(),
			UpdatedAt:       user.UpdatedAt.Time.UTC().String(),
		})

	}

	return userResponses, nil
}

// Get User By ID
func (s *UserService) GetUserByID(id uuid.UUID) (*models.UserWithProfileResponseParams, error) {
	user, err := s.DB.GetUserByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	var roles []models.Role
	if user.Roles != nil {
		// Unmarshal roles from []byte to []string
		err := json.Unmarshal(user.Roles.([]byte), &roles)
		if err != nil {
			// handle error
			return nil, err
		}

		// Decrypt each role if needed
		for i, role := range roles {
			decryptedName, err := s.Encryptor.Decrypt(role.Name)
			decryptedDescription, err := s.Encryptor.Decrypt(role.Description)
			if err == nil {
				roles[i].Name = decryptedName
				roles[i].Description = decryptedDescription
				roles[i].ID = role.ID
			}
		}
	}

	var groups []models.Group
	if user.Groups != nil {
		// Unmarshal groups from []byte to []string
		err := json.Unmarshal(user.Groups.([]byte), &groups)
		if err != nil {
			// handle error
			return nil, err
		}

		// Decrypt each group if needed
		for i, group := range groups {
			decryptedGroupName, err := s.Encryptor.Decrypt(group.Name)
			if err == nil {
				groups[i].Name = decryptedGroupName
				groups[i].Metadata = group.Metadata
				groups[i].ID = group.ID
			}
		}
	}

	// profile
	var profile map[string]interface{}
	err = json.Unmarshal(user.Profile.([]byte), &profile)
	if err != nil {
		// handle error
		return nil, err
	}

	DecryptObject := []string{
		user.FirstName.String,
		user.LastName.String,
		user.Email,
	}

	// Decrypt them
	decrypted, err := s.Encryptor.DecryptMultiple(DecryptObject)
	if err != nil {
		return &models.UserWithProfileResponseParams{}, err
	}

	userResponse := models.UserWithProfileResponseParams{
		ID:              user.ID,
		Email:           decrypted[2],
		FirstName:       decrypted[0],
		LastName:        decrypted[1],
		IsDeleted:       user.IsDeleted.Bool,
		IsActive:        user.IsActive.Bool,
		IsPasswordReset: user.IsPasswordReset.Bool,
		Role:            roles,
		Group:           groups,
		Profile:         profile,
		LastLogin:       user.LastLogin.Time.UTC().String(),
		CreatedAt:       user.CreatedAt.Time.UTC().String(),
		UpdatedAt:       user.UpdatedAt.Time.UTC().String(),
	}

	return &userResponse, nil
}

// Create New User
func (s *UserService) CreateUser(user models.CreateUserParams) (*models.UserResponseParams, error) {
	// encryption
	encEmail, err := s.Encryptor.Encrypt(user.Email)
	if err != nil {
		return nil, err
	}

	encFirstName, err := s.Encryptor.Encrypt(user.FirstName)
	if err != nil {
		return nil, err
	}

	encLastName, err := s.Encryptor.Encrypt(user.LastName)
	if err != nil {
		return nil, err
	}
	// create user
	newUser, err := s.DB.CreateUser(context.Background(), generated.CreateUserParams{
		Email:           encEmail,
		EmailHash:       utils.Hash(user.Email),
		FirstName:       sql.NullString{String: encFirstName, Valid: true},
		LastName:        sql.NullString{String: encLastName, Valid: true},
		Password:        utils.Hash(user.Password), // Password should be encrypted before storing
		IsDeleted:       sql.NullBool{Valid: true, Bool: false},
		IsActive:        sql.NullBool{Valid: true, Bool: false},
		IsPasswordReset: sql.NullBool{Valid: true, Bool: true},
		CreatedAt:       sql.NullTime{Time: time.Now().UTC(), Valid: true},
		UpdatedAt:       sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &models.UserResponseParams{
		ID:              newUser.ID,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		IsDeleted:       newUser.IsDeleted.Bool,
		IsActive:        newUser.IsActive.Bool,
		IsPasswordReset: newUser.IsPasswordReset.Bool,
		LastLogin:       newUser.LastLogin.Time.UTC().String(),
		CreatedAt:       newUser.CreatedAt.Time.UTC().String(),
		UpdatedAt:       newUser.UpdatedAt.Time.UTC().String(),
	}, nil
}

// Update Existing User
func (s *UserService) UpdateUser(id uuid.UUID, user models.UpdateUserParams) (*models.UserResponseParams, error) {

	//encryption
	encFirstName, err := s.Encryptor.Encrypt(user.FirstName)
	if err != nil {
		return nil, err
	}

	encLastName, err := s.Encryptor.Encrypt(user.LastName)
	if err != nil {
		return nil, err
	}

	updatedUser, err := s.DB.UpdateUserProfile(context.Background(), generated.UpdateUserProfileParams{
		ID:        id,
		FirstName: sql.NullString{String: encFirstName, Valid: true},
		LastName:  sql.NullString{String: encLastName, Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})

	if err != nil {
		return nil, err
	}

	updatedUser.Email, err = s.Encryptor.Decrypt(updatedUser.Email)
	if err != nil {
		return nil, err
	}

	// Map the updated user to the response model
	userResponse := &models.UserResponseParams{
		ID:              updatedUser.ID,
		Email:           updatedUser.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		IsDeleted:       updatedUser.IsDeleted.Bool,
		IsActive:        updatedUser.IsActive.Bool,
		IsPasswordReset: updatedUser.IsPasswordReset.Bool,
		LastLogin:       updatedUser.LastLogin.Time.UTC().String(),
		CreatedAt:       updatedUser.CreatedAt.Time.UTC().String(),
		UpdatedAt:       updatedUser.UpdatedAt.Time.UTC().String(),
	}

	return userResponse, nil
}

// Delete User
func (s *UserService) DeleteUser(id uuid.UUID) error {
	err := s.DB.DeleteUser(context.Background(), generated.DeleteUserParams{
		ID:        id,
		UpdatedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
	if err != nil {
		return err
	}
	return nil
}
