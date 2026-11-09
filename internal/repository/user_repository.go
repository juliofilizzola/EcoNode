package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username     string    `gorm:"size:64;not null;uniqueIndex"`
	PasswordHash string    `gorm:"size:255;not null"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime"`
}

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (User, error)
}

type gormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) FindByUsername(ctx context.Context, username string) (User, error) {
	if username == "" {
		return User{}, ErrUserNotFound
	}

	var user User

	result := r.db.WithContext(ctx).Where(&User{Username: username}).Limit(1).Find(&user)
	if result.Error != nil {
		return User{}, fmt.Errorf("find user by username: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return User{}, ErrUserNotFound
	}

	return user, nil
}
