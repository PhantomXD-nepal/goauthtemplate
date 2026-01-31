package auth

import (
	"fmt"
	"time"

	"github.com/PhantomXD-nepal/goauthtemplate/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJWT(userID uuid.UUID, email string) (string, error) {
	jwtSecret := config.Envs.JWTSecret

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"email":   email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func ValidateJWT(tokenString string, secret []byte) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return &claims, nil
}

func GetUserIDFromToken(tokenString string, secret []byte) (string, error) {
	claims, err := ValidateJWT(tokenString, secret)
	if err != nil {
		return "", err
	}

	userID, ok := (*claims)["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id not found in token")
	}

	return userID, nil
}

func IsTokenExpired(tokenString string, secret []byte) bool {
	claims, err := ValidateJWT(tokenString, secret)
	if err != nil {
		return true
	}

	exp, ok := (*claims)["exp"].(float64)
	if !ok {
		return true
	}

	expirationTime := time.Unix(int64(exp), 0)
	return time.Now().After(expirationTime)
}

func RefreshToken(tokenString string, secret []byte) (string, error) {
	claims, err := ValidateJWT(tokenString, secret)
	if err != nil {
		return "", fmt.Errorf("invalid token for refresh: %w", err)
	}

	userIDStr, ok := (*claims)["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id not found in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", fmt.Errorf("invalid user_id format: %w", err)
	}

	email, ok := (*claims)["email"].(string)
	if !ok {
		return "", fmt.Errorf("email not found in token")
	}

	return GenerateJWT(userID, email)
}
