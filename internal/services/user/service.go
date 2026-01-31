package user

import (
	"context"
	"database/sql"

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

func (s *Service) Register(ctx context.Context, email, password string) error {
	_, err := s.queries.GetUserByEmail(ctx, email)
	if err == nil {
		return types.ErrEmailAlreadyExists
	}
	if err != sql.ErrNoRows {
		utils.Error("Failed to check existing user: " + err.Error())
		return types.ErrInternalServer
	}

	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		utils.Error("Failed to hash password: " + err.Error())
		return types.ErrInternalServer
	}

	id := uuid.New()
	err = s.queries.CreateUser(ctx, sqlc.CreateUserParams{
		UUIDTOBIN: id.String(),
		Email:     email,
		Password:  string(hashed),
	})
	if err != nil {
		utils.Error("Failed to create user: " + err.Error())
		return types.ErrInternalServer
	}

	return nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		bcrypt.CompareHashAndPassword(
			[]byte("$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"),
			[]byte(password),
		)
		if err == sql.ErrNoRows {
			return "", types.ErrInvalidCredentials
		}
		utils.Error("Failed to fetch user: " + err.Error())
		return "", types.ErrInternalServer
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", types.ErrInvalidCredentials
	}

	userID, err := uuid.Parse(user.ID)
	if err != nil {
		utils.Error("Failed to parse user ID: " + err.Error())
		return "", types.ErrInternalServer
	}

	token, err := auth.GenerateJWT(userID, user.Email)
	if err != nil {
		utils.Error("Failed to generate JWT: " + err.Error())
		return "", types.ErrInternalServer
	}

	return token, nil
}
