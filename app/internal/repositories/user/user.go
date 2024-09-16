package user

import (
	"app/internal/commons"
	"app/internal/entities"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (entities.User, error)
	CreateUser(ctx context.Context, user entities.User) error
	UpdateUser(ctx context.Context, user entities.User) error
}

type userRepository struct {
	db             *gorm.DB
	ContextTimeout time.Duration
}

func NewUserRepository(db *gorm.DB, timeout time.Duration) UserRepository {
	return &userRepository{db: db, ContextTimeout: timeout}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ContextTimeout)
	defer cancel()

	var user entities.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, commons.ErrNotFound
		}
		// Check if the context was canceled
		if ctx.Err() == context.DeadlineExceeded {
			return entities.User{}, commons.ErrTimeout
		}
		return entities.User{}, err
	}
	return user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user entities.User) error {
	ctx, cancel := context.WithTimeout(ctx, r.ContextTimeout)
	defer cancel()

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		// Check if the context was canceled
		if ctx.Err() == context.DeadlineExceeded {
			return commons.ErrTimeout
		}
		return err
	}

	return nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user entities.User) error {
	ctx, cancel := context.WithTimeout(ctx, r.ContextTimeout)
	defer cancel()

	user.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		// Check if the context was canceled
		if ctx.Err() == context.DeadlineExceeded {
			return commons.ErrTimeout
		}
		return err
	}

	return nil
}
