package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/PhantomXD-nepal/goauthtemplate/db/generated/sqlc"
	"github.com/PhantomXD-nepal/goauthtemplate/internal/services/auth"
	"github.com/PhantomXD-nepal/goauthtemplate/internal/types"
	"github.com/PhantomXD-nepal/goauthtemplate/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	queries *sqlc.Queries
}

func NewService(queries *sqlc.Queries) *Service {
	return &Service{
		queries: queries,
	}
}

func (s *Service) Register(ctx context.Context, email, password string) (string, string, error) {
	_, err := s.queries.GetUserByEmail(ctx, email)
	if err == nil {
		return "", "", types.ErrEmailAlreadyExists
	}
	if err != sql.ErrNoRows {
		utils.Error("Failed to check existing user: " + err.Error())
		return "", "", types.ErrInternalServer
	}

	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		utils.Error("Failed to hash password: " + err.Error())
		return "", "", types.ErrInternalServer
	}
	//Generate access token
	token, err := auth.GenerateJWT(uuid.New(), email)
	if err != nil {
		utils.Error("Error when generating jwt when registering" + err.Error())
		return "", "", err
	}

	//Generate refresh token
	refreshTokenValue := utils.GenerateRandomString(32)
	hashedRefreshToken := utils.HashString(refreshTokenValue)
	refreshTokenID := uuid.New()

	id := uuid.New()

	err = s.queries.CreateUser(ctx, sqlc.CreateUserParams{
		UUIDTOBIN: id.String(),
		Email:     email,
		Password:  string(hashed),
	})
	if err != nil {
		utils.Error("Failed to create user: " + err.Error())
		return "", "", types.ErrInternalServer
	}

	err = s.queries.CreateRefreshToken(ctx, sqlc.CreateRefreshTokenParams{
		UUIDTOBIN:   refreshTokenID.String(),
		UUIDTOBIN_2: id.String(),
		TokenHash:   hashedRefreshToken,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		utils.Error("Failed to create refresh token: " + err.Error())
		return "", "", types.ErrInternalServer
	}

	return token, refreshTokenValue, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		bcrypt.CompareHashAndPassword(
			[]byte("$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"),
			[]byte(password),
		)
		if err == sql.ErrNoRows {
			return "", "", types.ErrInvalidCredentials
		}
		utils.Error("Failed to fetch user: " + err.Error())
		return "", "", types.ErrInternalServer
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", types.ErrInvalidCredentials
	}

	userID, err := uuid.Parse(user.ID)
	if err != nil {
		utils.Error("Failed to parse user ID: " + err.Error())
		return "", "", types.ErrInternalServer
	}
	//Generate access token
	token, err := auth.GenerateJWT(userID, user.Email)
	if err != nil {
		utils.Error("Failed to generate JWT: " + err.Error())
		return "", "", types.ErrInternalServer
	}

	//Generate refresh token
	refreshTokenValue := utils.GenerateRandomString(32)
	hashedRefreshToken := utils.HashString(refreshTokenValue)
	refreshTokenID := uuid.New()

	//Delete previous refresh tokens
	err = s.queries.DeleteRefreshToken(ctx, userID.String())
	if err != nil {
		utils.Error("Failed to delete: " + err.Error())
		return "", "", types.ErrInternalServer
	}

	//Create new refresh token

	err = s.queries.CreateRefreshToken(ctx, sqlc.CreateRefreshTokenParams{
		UUIDTOBIN:   refreshTokenID.String(),
		UUIDTOBIN_2: userID.String(),
		TokenHash:   hashedRefreshToken,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	})

	if err != nil {
		utils.Error("Failed to create refresh token: " + err.Error())
		return "", "", types.ErrInternalServer
	}

	return token, refreshTokenValue, nil
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	//Check if refresh token is valid
	hashedRefreshToken := utils.HashString(refreshToken)
	//Fetch refresh token to see if its valid
	res, err := s.queries.GetRefreshTokenFromTokenHash(ctx, hashedRefreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", types.ErrInvalidCredentials
		}
	}
	userID, err := uuid.Parse(res.UserID)
	if err != nil {
		return "", "", err
	}
	//Generate new access token
	token, err := auth.GenerateJWT(userID, res.UserID)
	if err != nil {
		utils.Error("Failed to generate access token" + err.Error())
		return "", "", err
	}

	//Generate refresh token
	refreshTokenValue := utils.GenerateRandomString(32)
	newHashedRefreshToken := utils.HashString(refreshTokenValue)
	refreshTokenID := uuid.New()

	//Delete previous refresh tokens
	err = s.queries.DeleteRefreshToken(ctx, userID.String())
	if err != nil {
		utils.Error("Failed to delete: " + err.Error())
		return "", "", types.ErrInternalServer
	}

	//Create new refresh token

	err = s.queries.CreateRefreshToken(ctx, sqlc.CreateRefreshTokenParams{
		UUIDTOBIN:   refreshTokenID.String(),
		UUIDTOBIN_2: userID.String(),
		TokenHash:   newHashedRefreshToken,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	})

	if err != nil {
		utils.Error("Failed to create refresh token: " + err.Error())
		return "", "", types.ErrInternalServer
	}

	return token, refreshTokenValue, nil

}
