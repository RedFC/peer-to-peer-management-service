package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"p2p-management-service/config"
	"p2p-management-service/db/generated"
	"p2p-management-service/models"
	"p2p-management-service/utils"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	DB        SubscriptionRepository
	Encryptor *utils.Encryptor
	EmailSvc  EmailService
}

type SubscriptionRepository interface {
	ListSubscriptions(ctx context.Context, arg generated.ListSubscriptionsParams) ([]generated.ListSubscriptionsRow, error)
	GetTotalSubscriptionsCount(ctx context.Context) (int64, error)
	GetSubscriptionByID(ctx context.Context, id uuid.UUID) (generated.GetSubscriptionByIDRow, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (generated.GetUserByIDRow, error)
	GetSubscriptionByUserIDAndGroupIDAndStatus(ctx context.Context, arg generated.GetSubscriptionByUserIDAndGroupIDAndStatusParams) (generated.Subscription, error)
	UpdateSubscription(ctx context.Context, arg generated.UpdateSubscriptionParams) (generated.Subscription, error)
	CreateSubscription(ctx context.Context, arg generated.CreateSubscriptionParams) (generated.Subscription, error)
	GetRoleByName(ctx context.Context, nameHash string) (generated.Role, error)
	AssignRoleToUser(ctx context.Context, arg generated.AssignRoleToUserParams) error
	AssignGroupToUser(ctx context.Context, arg generated.AssignGroupToUserParams) error
	DeleteSubscription(ctx context.Context, arg generated.DeleteSubscriptionParams) error
	GetSubscriptionByUserIDAndGroupID(ctx context.Context, arg generated.GetSubscriptionByUserIDAndGroupIDParams) ([]generated.Subscription, error)
	RevokeSubscription(ctx context.Context, arg generated.RevokeSubscriptionParams) (generated.Subscription, error)
}

func NewSubscriptionService(db SubscriptionRepository, encryptor *utils.Encryptor, emailSvc EmailService) *SubscriptionService {
	return &SubscriptionService{
		DB:        db,
		Encryptor: encryptor,
		EmailSvc:  emailSvc,
	}
}

// List All Subscriptions With Pagination
func (s *SubscriptionService) GetSubscriptions(pageNo string, pageSize string, totalCount int64) ([]models.SubscriptionWithDetails, error) {
	pageNoInt, err := strconv.Atoi(pageNo)
	if err != nil {
		return nil, err
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		return nil, err
	}

	subscriptions, err := s.DB.ListSubscriptions(context.Background(), generated.ListSubscriptionsParams{
		Limit:  int32(pageSizeInt),
		Offset: int32((pageNoInt - 1) * pageSizeInt),
	})
	if err != nil {
		return []models.SubscriptionWithDetails{}, err
	}

	var subscriptionResponses []models.SubscriptionWithDetails
	for _, subscription := range subscriptions {

		DecryptObject := []string{
			subscription.UserFirstName.String,
			subscription.UserLastName.String,
			subscription.UserEmail,
			subscription.GroupName,
		}

		// Decrypt them
		decrypted, err := s.Encryptor.DecryptMultiple(DecryptObject)
		if err != nil {
			return []models.SubscriptionWithDetails{}, err
		}

		subscriptionResponses = append(subscriptionResponses, models.SubscriptionWithDetails{
			Subscription: models.Subscription{
				ID:                 subscription.SubscriptionID,
				SubscriptionStatus: string(subscription.SubscriptionStatus.SubscriptionStatus),
				CreatedAt:          subscription.SubscriptionCreatedAt.Time.UTC().String(),
				UpdatedAt:          subscription.SubscriptionUpdatedAt.Time.UTC().String(),
			},
			User: models.UserResponseParams{
				ID:        subscription.UserID,
				FirstName: decrypted[0],
				LastName:  decrypted[1],
				Email:     decrypted[2],
				CreatedAt: subscription.UserCreatedAt.Time.UTC().String(),
				UpdatedAt: subscription.UserUpdatedAt.Time.UTC().String(),
			},
			Group: models.Group{
				ID:        subscription.GroupID,
				Name:      decrypted[3],
				Metadata:  subscription.GroupMetadata.RawMessage,
				CreatedAt: subscription.GroupCreatedAt.Time.UTC().String(),
				UpdatedAt: subscription.GroupUpdatedAt.Time.UTC().String(),
			},
		})
	}

	return subscriptionResponses, nil
}

// Get Subscription Count
func (s *SubscriptionService) GetTotalSubscriptionsCount(ctx context.Context) (int64, error) {
	totalCount, err := s.DB.GetTotalSubscriptionsCount(ctx)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

// Get Subscription By ID
func (s *SubscriptionService) GetSubscriptionByID(id uuid.UUID) (*models.SubscriptionWithDetails, error) {
	subscription, err := s.DB.GetSubscriptionByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	DecryptObject := []string{
		subscription.UserFirstName.String,
		subscription.UserLastName.String,
		subscription.UserEmail,
		subscription.GroupName,
	}

	// Decrypt them
	decrypted, err := s.Encryptor.DecryptMultiple(DecryptObject)
	if err != nil {
		return nil, err
	}

	subscriptionResponse := models.SubscriptionWithDetails{
		Subscription: models.Subscription{
			ID:                 subscription.SubscriptionID,
			SubscriptionStatus: string(subscription.SubscriptionStatus.SubscriptionStatus),
			CreatedAt:          subscription.SubscriptionCreatedAt.Time.UTC().String(),
			UpdatedAt:          subscription.SubscriptionUpdatedAt.Time.UTC().String(),
		},
		User: models.UserResponseParams{
			ID:        subscription.UserID,
			FirstName: decrypted[0],
			LastName:  decrypted[1],
			Email:     decrypted[2],
			CreatedAt: subscription.UserCreatedAt.Time.UTC().String(),
			UpdatedAt: subscription.UserUpdatedAt.Time.UTC().String(),
		},
		Group: models.Group{
			ID:        subscription.GroupID,
			Name:      decrypted[3],
			Metadata:  subscription.GroupMetadata.RawMessage,
			CreatedAt: subscription.GroupCreatedAt.Time.UTC().String(),
			UpdatedAt: subscription.GroupUpdatedAt.Time.UTC().String(),
		},
	}

	return &subscriptionResponse, nil
}

// Create New Subscription
func (s *SubscriptionService) CreateSubscription(subscription models.CreateSubscriptionParams) (*models.Subscription, error) {

	// validate referenced user exists to avoid FK violation
	if subscription.UserID == (uuid.UUID{}) {
		return nil, fmt.Errorf("user_id is required")
	}

	if _, err := s.DB.GetUserByID(context.Background(), subscription.UserID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	// Check for existing subscription for this user+group
	existing, err := s.DB.GetSubscriptionByUserIDAndGroupIDAndStatus(context.Background(), generated.GetSubscriptionByUserIDAndGroupIDAndStatusParams{
		UserID:             uuid.NullUUID{UUID: subscription.UserID, Valid: true},
		GroupID:            uuid.NullUUID{UUID: subscription.GroupID, Valid: true},
		SubscriptionStatus: generated.NullSubscriptionStatus{SubscriptionStatus: generated.SubscriptionStatus("active"), Valid: true},
	})
	if err == nil {
		// found an existing subscription
		if existing.SubscriptionStatus.SubscriptionStatus == "active" {
			return nil, fmt.Errorf("subscription already active for this user and group")
		}

		// exists but not active -> activate it (update)
		updated, err := s.DB.UpdateSubscription(context.Background(), generated.UpdateSubscriptionParams{
			ID:                 existing.ID,
			SubscriptionStatus: generated.NullSubscriptionStatus{SubscriptionStatus: generated.SubscriptionStatus("active"), Valid: true},
			UpdatedAt:          sql.NullTime{Time: time.Now().UTC(), Valid: true},
		})
		if err != nil {
			return nil, err
		}

		created := &models.Subscription{
			ID: updated.ID,
			// UserID:             updated.UserID.UUID,
			// GroupID:            int(updated.GroupID.Int32),
			SubscriptionStatus: string(updated.SubscriptionStatus.SubscriptionStatus),
			CreatedAt:          updated.CreatedAt.Time.UTC().String(),
			UpdatedAt:          updated.UpdatedAt.Time.UTC().String(),
		}

		return created, nil
	}

	// if error is not sql.ErrNoRows and not nil, propagate
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// create new subscription with active status
	newSubscription, err := s.DB.CreateSubscription(context.Background(), generated.CreateSubscriptionParams{
		UserID:             uuid.NullUUID{UUID: subscription.UserID, Valid: true},
		GroupID:            uuid.NullUUID{UUID: subscription.GroupID, Valid: true},
		SubscriptionStatus: generated.NullSubscriptionStatus{SubscriptionStatus: generated.SubscriptionStatus("active"), Valid: true},
		CreatedAt:          sql.NullTime{Time: time.Now().UTC(), Valid: true},
		UpdatedAt:          sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})

	if err != nil {
		return nil, err
	}

	created := &models.Subscription{
		ID: newSubscription.ID,
		// UserID:             newSubscription.UserID.UUID,
		// GroupID:            newSubscription.GroupID.UUID,
		SubscriptionStatus: string(newSubscription.SubscriptionStatus.SubscriptionStatus),
		CreatedAt:          newSubscription.CreatedAt.Time.UTC().String(),
		UpdatedAt:          newSubscription.UpdatedAt.Time.UTC().String(),
	}

	// If an email service is configured, send a simple invitation email asynchronously.
	if s.EmailSvc != nil {
		// Try to look up the user's email and name from the users table, decrypt and send to that address.
		s.LookupUserAndSendEmail(newSubscription.UserID.UUID, subscription.GroupID, "Welcome to Peer-Peer Communication Platform")
	}

	role, err := s.DB.GetRoleByName(context.Background(), utils.Hash("subscriber"))
	if err != nil {
		return nil, err
	}

	s.DB.AssignRoleToUser(context.Background(), generated.AssignRoleToUserParams{
		UserID: subscription.UserID,
		RoleID: role.ID,
	})

	s.DB.AssignGroupToUser(context.Background(), generated.AssignGroupToUserParams{
		UserID:  subscription.UserID,
		GroupID: subscription.GroupID,
	})

	return created, nil
}

// Update Existing Subscription
func (s *SubscriptionService) UpdateSubscription(id uuid.UUID, subscription models.UpdateSubscriptionParams) (*models.Subscription, error) {

	updatedSubscription, err := s.DB.UpdateSubscription(context.Background(), generated.UpdateSubscriptionParams{
		ID:                 id,
		SubscriptionStatus: generated.NullSubscriptionStatus{SubscriptionStatus: generated.SubscriptionStatus(subscription.SubscriptionStatus), Valid: true},
		UpdatedAt:          sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})

	if err != nil {
		return nil, err
	}

	// Map the updated subscription to the response model
	SubscriptionResponse := &models.Subscription{
		ID: updatedSubscription.ID,
		// UserID:             updatedSubscription.UserID.UUID,
		// GroupID:            updatedSubscription.GroupID.UUID,
		SubscriptionStatus: string(updatedSubscription.SubscriptionStatus.SubscriptionStatus),
		CreatedAt:          updatedSubscription.CreatedAt.Time.UTC().String(),
		UpdatedAt:          updatedSubscription.UpdatedAt.Time.UTC().String(),
	}

	return SubscriptionResponse, nil
}

// Delete Subscription
func (s *SubscriptionService) DeleteSubscription(userId uuid.UUID, groupId uuid.UUID) error {

	subscription, err := s.DB.GetSubscriptionByUserIDAndGroupIDAndStatus(context.Background(), generated.GetSubscriptionByUserIDAndGroupIDAndStatusParams{
		UserID:             uuid.NullUUID{UUID: userId, Valid: true},
		GroupID:            uuid.NullUUID{UUID: groupId, Valid: true},
		SubscriptionStatus: generated.NullSubscriptionStatus{SubscriptionStatus: generated.SubscriptionStatus("active"), Valid: true},
	})
	if err != nil {
		return err
	}

	id := subscription.ID

	err = s.DB.DeleteSubscription(context.Background(), generated.DeleteSubscriptionParams{
		ID:        id,
		UpdatedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})

	return err
}

// revoke Subscription
func (s *SubscriptionService) RevokeSubscription(userId uuid.UUID, groupId uuid.UUID) error {

	subscription, err := s.DB.GetSubscriptionByUserIDAndGroupIDAndStatus(context.Background(), generated.GetSubscriptionByUserIDAndGroupIDAndStatusParams{
		UserID:             uuid.NullUUID{UUID: userId, Valid: true},
		GroupID:            uuid.NullUUID{UUID: groupId, Valid: true},
		SubscriptionStatus: generated.NullSubscriptionStatus{SubscriptionStatus: generated.SubscriptionStatus("active"), Valid: true},
	})
	if err != nil {
		return err
	}

	id := subscription.ID

	_, err = s.DB.RevokeSubscription(context.Background(), generated.RevokeSubscriptionParams{
		ID:        id,
		UpdatedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})

	if err != nil {
		return err
	}

	return nil
}

// OnboardSubscription creates or activates a subscription for the given user and group.
// It only requires userID and groupID and will set the subscription status to active.
func (s *SubscriptionService) OnboardSubscription(userID uuid.UUID, groupID uuid.UUID) (*models.Subscription, error) {
	if userID == (uuid.UUID{}) {
		return nil, fmt.Errorf("user_id is required")
	}

	if _, err := s.DB.GetUserByID(context.Background(), userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	existing, err := s.DB.GetSubscriptionByUserIDAndGroupIDAndStatus(context.Background(), generated.GetSubscriptionByUserIDAndGroupIDAndStatusParams{
		UserID:             uuid.NullUUID{UUID: userID, Valid: true},
		GroupID:            uuid.NullUUID{UUID: groupID, Valid: true},
		SubscriptionStatus: generated.NullSubscriptionStatus{SubscriptionStatus: generated.SubscriptionStatus("active"), Valid: true},
	})
	if err == nil {
		if existing.SubscriptionStatus.SubscriptionStatus == "active" {
			return nil, fmt.Errorf("subscription already active for this user and group")
		}

		updated, err := s.DB.UpdateSubscription(context.Background(), generated.UpdateSubscriptionParams{
			ID:                 existing.ID,
			SubscriptionStatus: generated.NullSubscriptionStatus{SubscriptionStatus: generated.SubscriptionStatus("active"), Valid: true},
			UpdatedAt:          sql.NullTime{Time: time.Now().UTC(), Valid: true},
		})
		if err != nil {
			return nil, err
		}

		return &models.Subscription{
			ID: updated.ID,
			// UserID:             updated.UserID.UUID,
			// GroupID:            updated.GroupID.UUID,
			SubscriptionStatus: string(updated.SubscriptionStatus.SubscriptionStatus),
			CreatedAt:          updated.CreatedAt.Time.UTC().String(),
			UpdatedAt:          updated.UpdatedAt.Time.UTC().String(),
		}, nil
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	newSubscription, err := s.DB.CreateSubscription(context.Background(), generated.CreateSubscriptionParams{
		UserID:             uuid.NullUUID{UUID: userID, Valid: true},
		GroupID:            uuid.NullUUID{UUID: groupID, Valid: true},
		SubscriptionStatus: generated.NullSubscriptionStatus{SubscriptionStatus: generated.SubscriptionStatus("active"), Valid: true},
		CreatedAt:          sql.NullTime{Time: time.Now().UTC(), Valid: true},
		UpdatedAt:          sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	role, err := s.DB.GetRoleByName(context.Background(), utils.Hash("subscriber"))
	if err != nil {
		return nil, err
	}

	s.DB.AssignRoleToUser(context.Background(), generated.AssignRoleToUserParams{
		UserID: newSubscription.UserID.UUID,
		RoleID: role.ID,
	})

	return &models.Subscription{
		ID: newSubscription.ID,
		// UserID:             newSubscription.UserID.UUID,
		// GroupID:            newSubscription.GroupID.UUID,
		SubscriptionStatus: string(newSubscription.SubscriptionStatus.SubscriptionStatus),
		CreatedAt:          newSubscription.CreatedAt.Time.UTC().String(),
		UpdatedAt:          newSubscription.UpdatedAt.Time.UTC().String(),
	}, nil
}

// send invitation
func (s *SubscriptionService) ResendInvitation(userID uuid.UUID, groupID uuid.UUID) (string, error) {

	if userID == (uuid.UUID{}) {
		return "", fmt.Errorf("user_id is required")
	}

	if _, err := s.DB.GetUserByID(context.Background(), userID); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found")
		}
		return "", err
	}

	existing, err := s.DB.GetSubscriptionByUserIDAndGroupIDAndStatus(context.Background(), generated.GetSubscriptionByUserIDAndGroupIDAndStatusParams{
		UserID:             uuid.NullUUID{UUID: userID, Valid: true},
		GroupID:            uuid.NullUUID{UUID: groupID, Valid: true},
		SubscriptionStatus: generated.NullSubscriptionStatus{SubscriptionStatus: generated.SubscriptionStatus("active"), Valid: true},
	})

	if err == nil {
		if existing.SubscriptionStatus.SubscriptionStatus == "active" {
			if s.EmailSvc != nil {
				// Try to look up the user's email and name from the users table, decrypt and send to that address.

				s.LookupUserAndSendEmail(userID, groupID, "Welcome to Peer-Peer Communication Platform")

			}
		} else {
			fmt.Errorf("acive subscription not found")
		}
	}

	return "Invitation resent successfully", nil
}

// lookup user and send email
func (s *SubscriptionService) LookupUserAndSendEmail(userID uuid.UUID, groupID uuid.UUID, Sub string) {
	userRow, err := s.DB.GetUserByID(context.Background(), userID)
	if err == nil {
		DecryptObject := []string{
			userRow.FirstName.String,
			userRow.LastName.String,
			userRow.Email,
		}
		decrypted, err := s.Encryptor.DecryptMultiple(DecryptObject)
		to := ""
		first := ""
		last := ""
		if err == nil && len(decrypted) >= 3 {
			first = decrypted[0]
			last = decrypted[1]
			to = decrypted[2]
		}
		if to != "" {

			ios_url := "https://apps.apple.com/us/app/peer-peer-communication/id6441234567"
			android_url := "https://play.google.com/store/apps/details?id=com.p2pcommunication"

			qr, err := utils.GenerateQRCode(map[string]interface{}{
				"hubs":    config.AppConfig.HUBS,
				"userID":  userID,
				"groupID": groupID,
			})

			if err != nil {
				log.Fatal(err)
			}

			body := utils.EmailBodyTemplate(first, last, qr, ios_url, android_url)
			fmt.Println("TO: ", to)
			fmt.Println("SUBJECT: ", Sub)
			fmt.Println("BODY: ", body)
			s.EmailSvc.SendEmailAsync(context.Background(), EmailRequest{
				To:      to,
				Subject: Sub,
				Body:    body,
			})
		}
	}
}
