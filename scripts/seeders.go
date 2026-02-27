package scripts

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/sqlc-dev/pqtype"

	"p2p-management-service/config"
	"p2p-management-service/db/generated"
	"p2p-management-service/utils"
)

func seedRoles(q *generated.Queries, encryptor *utils.Encryptor) {
	ctx := context.Background()

	// Seed roles
	roles := []struct {
		Name        string
		Description string
	}{
		{"super_admin", "Super Administrator role with all privileges"},
		{"admin", "Administrator role with limited privileges"},
		{"user", "Regular user role"},
		{"subscriber", "Subscriber role"},
	}

	for _, role := range roles {
		encNameHash := utils.Hash(role.Name)

		encName, err := encryptor.Encrypt(role.Name)
		if err != nil {
			log.Println("Encryption failed:", err)
			continue
		}

		encDesc, err := encryptor.Encrypt(role.Description)
		if err != nil {
			log.Println("Encryption failed:", err)
			continue
		}

		_, err = q.GetRoleByName(ctx, encNameHash)
		if err == nil {
			log.Printf("Role with name %s already exists, skipping...\n", role.Name)
			continue
		}
		if err != sql.ErrNoRows {
			log.Println("Error checking role existence:", err)
			continue
		}

		_, err = q.CreateRole(ctx, generated.CreateRoleParams{
			Name:        encName,
			NameHash:    encNameHash,
			Description: sql.NullString{String: encDesc, Valid: true},
			CreatedAt:   sql.NullTime{Time: time.Now().UTC(), Valid: true},
			UpdatedAt:   sql.NullTime{Time: time.Now().UTC(), Valid: true},
			IsDeleted:   sql.NullBool{Bool: false, Valid: true},
		})
		if err != nil {
			log.Println("Seeding error:", err)
		}
	}

}

func seedGroups(q *generated.Queries, encryptor *utils.Encryptor) {
	ctx := context.Background()

	// Seed roles
	groups := []struct {
		Name     string
		MetaData map[string]interface{}
	}{
		{
			Name: "Super Admin Group",
			MetaData: map[string]interface{}{
				"access_level": "all",
			},
		},
		{
			Name: "Admin Group",
			MetaData: map[string]interface{}{
				"access_level": "all",
			},
		},
		{
			Name: "Support Group",
			MetaData: map[string]interface{}{
				"access_level": "limited",
			},
		},
		{
			Name: "User Group",
			MetaData: map[string]interface{}{
				"access_level": "limited",
			},
		},
		{
			Name: "B-2 Spirit",
			MetaData: map[string]interface{}{
				"tier": "platinum",
			},
		},
		{
			Name: "C-130 Hercules Operators",
			MetaData: map[string]interface{}{
				"tier": "gold",
			},
		},
		{
			Name: "C-17 Globemaster III",
			MetaData: map[string]interface{}{
				"tier": "silver",
			},
		},
		{
			Name: "C-5 Galaxy",
			MetaData: map[string]interface{}{
				"tier": "bronze",
			},
		},
		{
			Name: "Core Team",
			MetaData: map[string]interface{}{
				"region": "NA",
			},
		},
		{
			Name: "Engineering Team",
			MetaData: map[string]interface{}{
				"region": "NA",
			},
		},
		{
			Name: "Marketing Team",
			MetaData: map[string]interface{}{
				"region": "NA",
			},
		},
	}

	for _, group := range groups {
		encNameHash := utils.Hash(group.Name)

		encName, err := encryptor.Encrypt(group.Name)
		if err != nil {
			log.Println("Encryption failed:", err)
			continue
		}

		_, err = q.GetGroupByName(ctx, encNameHash)
		if err == nil {
			log.Printf("Group with name %s already exists, skipping...\n", group.Name)
			continue
		}
		if err != sql.ErrNoRows {
			log.Println("Error checking group existence:", err)
			continue
		}

		metaDataJSON, err := json.Marshal(group.MetaData)
		if err != nil {
			log.Println("Encryption failed:", err)
			continue
		}

		_, err = q.CreateGroup(ctx, generated.CreateGroupParams{
			Name:     encName,
			NameHash: encNameHash,
			Metadata: pqtype.NullRawMessage{
				RawMessage: metaDataJSON,
				Valid:      true,
			},
			CreatedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
			UpdatedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
			IsDeleted: sql.NullBool{Bool: false, Valid: true},
		})
		if err != nil {
			log.Println("Seeding error:", err)
		}
	}
}

func seedUsers(q *generated.Queries, encryptor *utils.Encryptor) {
	ctx := context.Background()

	// Seed users
	users := []struct {
		Email     string
		Password  string
		FirstName string
		LastName  string
		Role      string
		Group     string
	}{
		{"superadmin@p2pservice.com", "Test@12345", "Super", "Admin", "super_admin", "Super Admin Group"},
		{"admin@p2pservice.com", "Test@12345", "Admin", "Admin", "admin", "Admin Group"},
		{"user@p2pservice.com", "Test@12345", "Test", "User", "user", "User Group"},
	}

	for _, user := range users {

		emailHash := utils.Hash(user.Email)
		_, err := q.GetUserByEmailHash(ctx, emailHash)
		if err == nil {
			log.Printf("User with email %s already exists, skipping...\n", user.Email)
			continue
		}
		if err != sql.ErrNoRows {
			log.Println("Error checking user existence:", err)
			continue
		}

		encEmail, err := encryptor.Encrypt(user.Email)
		if err != nil {
			log.Println("Encryption failed:", err)
			continue
		}

		encPassword, err := encryptor.Encrypt(user.Password)
		if err != nil {
			log.Println("Encryption failed:", err)
			continue
		}

		encFirstName, err := encryptor.Encrypt(user.FirstName)
		if err != nil {
			log.Println("Encryption failed:", err)
			continue
		}

		encLastName, err := encryptor.Encrypt(user.LastName)
		if err != nil {
			log.Println("Encryption failed:", err)
			continue
		}

		data, err := q.CreateUser(ctx, generated.CreateUserParams{
			Email:           encEmail,
			EmailHash:       emailHash,
			Password:        encPassword,
			FirstName:       sql.NullString{String: encFirstName, Valid: true},
			LastName:        sql.NullString{String: encLastName, Valid: true},
			IsDeleted:       sql.NullBool{Bool: false, Valid: true},
			IsActive:        sql.NullBool{Valid: true, Bool: true},
			IsPasswordReset: sql.NullBool{Bool: false, Valid: true},
			LastLogin:       sql.NullTime{Valid: false},
			CreatedAt:       sql.NullTime{Time: time.Now().UTC(), Valid: true},
			UpdatedAt:       sql.NullTime{Time: time.Now().UTC(), Valid: true},
		})

		if err != nil {
			log.Println("Seeding error:", err)
		}

		fmt.Println(utils.Hash(user.Role))

		Role, err := q.GetRoleByName(ctx, utils.Hash(user.Role))
		if err != nil {
			log.Println("Seeding error:", err)
			continue
		}

		Group, err := q.GetGroupByName(ctx, utils.Hash(user.Group))
		if err != nil {
			log.Println("Seeding error:", err)
			continue
		}

		err = q.AssignRoleToUser(ctx, generated.AssignRoleToUserParams{
			UserID: data.ID,
			RoleID: Role.ID,
		})

		err = q.AssignGroupToUser(ctx, generated.AssignGroupToUserParams{
			UserID:  data.ID,
			GroupID: Group.ID,
		})

		if err != nil {
			log.Println("Seeding error:", err)
		}
	}

}

func Seed(db *sql.DB) {
	q := generated.New(db)

	// Initialize the encryptor with a secret key from env
	secretKey := config.AppConfig.ENCRYPTION_SECRET
	if len(secretKey) != 32 {
		log.Fatal("ENCRYPTION_SECRET must be 32 bytes for AES-256")
	}
	encryptor := utils.NewEncryptor(secretKey)

	// seeding roles
	seedRoles(q, encryptor)

	// seeding groups
	seedGroups(q, encryptor)

	// seeding Users
	seedUsers(q, encryptor)

	log.Println("✅ Seeding complete.")
}
