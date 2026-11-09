package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidToken = errors.New("invalid token")

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

type IssuedToken struct {
	Value     string
	ExpiresAt time.Time
}

type TokenService interface {
	Issue(userID uuid.UUID, username string) (IssuedToken, error)
	Parse(tokenString string) (Claims, error)
}

type jwtTokenService struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTTokenService(secret string, ttl time.Duration) TokenService {
	return &jwtTokenService{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (s *jwtTokenService) Issue(userID uuid.UUID, username string) (IssuedToken, error) {
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(s.ttl)

	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(s.secret)
	if err != nil {
		return IssuedToken{}, fmt.Errorf("sign token: %w", err)
	}

	return IssuedToken{Value: signed, ExpiresAt: expiresAt}, nil
}

func (s *jwtTokenService) Parse(tokenString string) (Claims, error) {
	var claims Claims

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(*jwt.Token) (any, error) {
			return s.secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return Claims{}, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !token.Valid {
		return Claims{}, ErrInvalidToken
	}

	return claims, nil
}
