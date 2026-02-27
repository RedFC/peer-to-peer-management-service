package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(email string, roles []string) (string, error) {
	privateKeyData, err := os.ReadFile("keys/jwt_private.pem")
	if err != nil {
		return "", fmt.Errorf("failed to read private key file: %w", err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	claims := jwt.MapClaims{
		"email": email,
		"role":  roles,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signedToken, nil
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	publicKeyData, err := os.ReadFile("keys/jwt_public.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}
	return token, nil
}
