package service

import (
	"context"
	"errors"
	"fmt"

	"EcoNode/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService interface {
	Login(ctx context.Context, username string, password string) (IssuedToken, error)
}

type authService struct {
	users  repository.UserRepository
	tokens TokenService
}

func NewAuthService(users repository.UserRepository, tokens TokenService) AuthService {
	return &authService{
		users:  users,
		tokens: tokens,
	}
}

func (s *authService) Login(ctx context.Context, username string, password string) (IssuedToken, error) {
	user, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return IssuedToken{}, ErrInvalidCredentials
		}

		return IssuedToken{}, fmt.Errorf("find user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return IssuedToken{}, ErrInvalidCredentials
		}

		return IssuedToken{}, fmt.Errorf("compare password hash: %w", err)
	}

	issued, err := s.tokens.Issue(user.ID, user.Username)
	if err != nil {
		return IssuedToken{}, fmt.Errorf("issue access token: %w", err)
	}

	return issued, nil
}
