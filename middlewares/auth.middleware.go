package middlewares

import (
	"fmt"
	"log"
	"p2p-management-service/config"
	"p2p-management-service/utils"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kataras/iris/v12"
)

// AuthMiddleware is a middleware that checks if the user is authenticated
func AuthMiddleware(ctx iris.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": "unauthorized"})
		ctx.StopExecution()
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	_, err := utils.ValidateJWT(tokenStr)
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": "unauthorized"})
		ctx.StopExecution()
		return
	}

	ctx.Next()
}

func SuperAdminMiddleware(ctx iris.Context) {

	// Initialize the encryptor with a secret key from env
	secretKey := config.AppConfig.ENCRYPTION_SECRET
	if len(secretKey) != 32 {
		log.Fatal("ENCRYPTION_SECRET must be 32 bytes for AES-256")
	}
	encryptor := utils.NewEncryptor(secretKey)

	// Get token from Authorization header
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": "unauthorized"})
		ctx.StopExecution()
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	// Validate token
	token, err := utils.ValidateJWT(tokenStr)
	if err != nil || !token.Valid {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": "unauthorized"})
		ctx.StopExecution()
		return
	}

	// Extract claims as MapClaims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": "invalid token claims"})
		ctx.StopExecution()
		return
	}

	// Check for superadmin role
	roles, ok := claims["role"].([]interface{})
	if !ok {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": "invalid token claims"})
		ctx.StopExecution()
		return
	}

	var roleStrings []string
	for _, r := range roles {
		roleStrings = append(roleStrings, fmt.Sprintf("%v", r))
	}

	// Try to decrypt roles; if decryption fails, assume roles are plain strings already.
	decryptedRoles, err := encryptor.DecryptMultiple(roleStrings)
	var checkedRoles []string
	if err == nil {
		checkedRoles = decryptedRoles
		fmt.Println("Decrypted Roles in Middleware:", decryptedRoles)
	} else {
		checkedRoles = roleStrings
		fmt.Println("Roles in Middleware (plaintext):", roleStrings)
	}

	if !utils.Contains(checkedRoles, "super_admin") {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{"error": "forbidden: insufficient permissions"})
		ctx.StopExecution()
		return
	}

	// All good — continue
	ctx.Next()
}
