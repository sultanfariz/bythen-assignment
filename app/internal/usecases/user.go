package usecases

import (
	"app/internal/commons"
	"app/internal/entities"
	repositories "app/internal/repositories/user"
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Register(ctx context.Context, req *entities.UserRegisterRequest) (entities.User, error)
	Login(ctx context.Context, req *entities.UserLoginRequest) (string, error)
}

type userUsecase struct {
	repo           repositories.UserRepository
	jwtConfig      commons.ConfigJWT
	contextTimeout time.Duration
}

func NewUserUsecase(repo repositories.UserRepository, jwtConfig commons.ConfigJWT, timeout time.Duration) UserUsecase {
	return &userUsecase{
		repo:           repo,
		jwtConfig:      jwtConfig,
		contextTimeout: timeout,
	}
}

// mocked functions
var (
	hashPassword = func(password string) string {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		return string(hashedPassword)
	}
	generateJWT = func(u *userUsecase, email string) (string, error) {
		return u.jwtConfig.GenerateJWT(email)
	}
)

func (u *userUsecase) Register(ctx context.Context, req *entities.UserRegisterRequest) (entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		return entities.User{}, err
	}

	// Check if user already exists
	_, err := u.repo.FindByEmail(ctx, req.Email)
	if err == nil {
		return entities.User{}, commons.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword := hashPassword(req.Password)

	user := entities.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	err = u.repo.CreateUser(ctx, user)
	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (u *userUsecase) Login(ctx context.Context, req *entities.UserLoginRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		return "", err
	}

	user, err := u.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", commons.ErrInvalidCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return "", commons.ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := generateJWT(u, req.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}
