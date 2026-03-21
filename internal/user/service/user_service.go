package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/manav1011/ikatva-be/internal/config"
	sqldb "github.com/manav1011/ikatva-be/internal/db/sqlc"
	"github.com/manav1011/ikatva-be/internal/user/model"
	"github.com/manav1011/ikatva-be/internal/user/repository"
	"github.com/manav1011/ikatva-be/pkg/token"
	"github.com/manav1011/ikatva-be/pkg/utils"
)

// ErrInvalidCredentials is returned when login cannot be satisfied (wrong email/password or inactive user).
var ErrInvalidCredentials = errors.New("invalid credentials")

type UserService struct {
	repo *repository.UserRepository
	cfg  *config.Config
}

func NewUserService(repo *repository.UserRepository, cfg *config.Config) *UserService {
	return &UserService{repo: repo, cfg: cfg}
}

// Login validates credentials, issues JWTs, and persists the refresh token.
func (s *UserService) Login(ctx context.Context, email, password string) (*model.LoginData, error) {
	row, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("get user: %w", err)
	}
	if !row.PasswordHash.Valid || row.PasswordHash.String == "" {
		return nil, ErrInvalidCredentials
	}
	if !utils.CheckPasswordHash(password, row.PasswordHash.String) {
		return nil, ErrInvalidCredentials
	}

	uid := row.ID.String()
	access, err := token.GenerateAccessToken(uid)
	if err != nil {
		return nil, fmt.Errorf("access token: %w", err)
	}
	refresh, err := token.GenerateRefreshToken(uid)
	if err != nil {
		return nil, fmt.Errorf("refresh token: %w", err)
	}

	expiresAt := time.Now().Add(s.cfg.RefreshTokenDuration)
	_, err = s.repo.InsertRefreshToken(ctx, sqldb.InsertRefreshTokenParams{
		UserID:    row.ID,
		Token:     refresh,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return nil, fmt.Errorf("persist refresh token: %w", err)
	}

	expiresIn := int64(s.cfg.AccessTokenDuration.Seconds())
	return &model.LoginData{
		AccessToken:  access,
		RefreshToken: refresh,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User: model.LoginUser{
			ID:    row.ID,
			Email: row.Email,
		},
	}, nil
}
