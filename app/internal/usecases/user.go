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

type UserUsecase struct {
	Repo           repositories.UserRepository
	JWTConfig      commons.ConfigJWT
	ContextTimeout time.Duration
}

func NewUserUsecase(repo repositories.UserRepository, jwtConfig commons.ConfigJWT, timeout time.Duration) *UserUsecase {
	return &UserUsecase{
		Repo:           repo,
		JWTConfig:      jwtConfig,
		ContextTimeout: timeout,
	}
}

// mocked functions
var (
	hashPassword = func(password string) string {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		return string(hashedPassword)
	}
	generateJWT = func(u *UserUsecase, email string) (string, error) {
		return u.JWTConfig.GenerateJWT(email)
	}
)

func (u *UserUsecase) Register(ctx context.Context, req *entities.UserRegisterRequest) (entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ContextTimeout)
	defer cancel()

	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		return entities.User{}, err
	}

	// Check if user already exists
	_, err := u.Repo.FindByEmail(ctx, req.Email)
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

	err = u.Repo.CreateUser(ctx, user)
	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (u *UserUsecase) Login(ctx context.Context, req *entities.UserLoginRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ContextTimeout)
	defer cancel()

	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		return "", err
	}

	user, err := u.Repo.FindByEmail(ctx, req.Email)
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

// func (u *UserUsecase) Update(ctx context.Context, user entities.User) error {
// 	ctx, cancel := context.WithTimeout(ctx, u.ContextTimeout)
// 	defer cancel()

// 	// retrieve email from token
// 	email, err := u.JWTConfig.ExtractClaims(user.Token)

// 	// Check if user already exists
// 	_, err := u.Repo.FindByEmail(user.Email)
// 	if err != nil {
// 		return errors.New("user not found")
// 	}

// 	// Hash password
// 	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
// 	user.PasswordHash = string(hashedPassword)

// 	err = u.Repo.UpdateUser(ctx, user)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
