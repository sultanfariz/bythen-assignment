package usecases

import (
	"app/internal/commons"
	"app/internal/entities"
	"app/internal/repositories/user/mocks"
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserUsecase_Register(t *testing.T) {
	type args struct {
		req *entities.UserRegisterRequest
	}

	mockRepo := new(mocks.UserRepository)
	mockJWTConfig := commons.ConfigJWT{
		SecretJWT:       "secret",
		ExpiresDuration: 1,
	}
	timeout := time.Second * 2

	hashPassword = func(password string) string {
		return string(password)
	}

	tests := []struct {
		name    string
		args    args
		want    entities.User
		wantErr bool
		mock    func()
	}{
		{
			name: "success",
			args: args{
				req: &entities.UserRegisterRequest{
					Name:     "John Doe",
					Email:    "john@example.com",
					Password: "password",
				},
			},
			want: entities.User{
				Name:         "John Doe",
				Email:        "john@example.com",
				PasswordHash: hashPassword("password"),
			},
			wantErr: false,
			mock: func() {
				mockRepo.On("FindByEmail", mock.Anything, "john@example.com").Return(entities.User{}, commons.ErrNotFound)
				mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("entities.User")).Return(nil)
			},
		},
		{
			name: "user already exists",
			args: args{
				req: &entities.UserRegisterRequest{
					Name:     "John Doe",
					Email:    "john@example.com",
					Password: "password",
				},
			},
			want:    entities.User{},
			wantErr: true,
			mock: func() {
				mockRepo.On("FindByEmail", mock.Anything, "john@example.com").Return(entities.User{}, nil)
			},
		},
		{
			name: "validation error",
			args: args{
				req: &entities.UserRegisterRequest{
					Name:     "",
					Email:    "invalid-email",
					Password: "short",
				},
			},
			want:    entities.User{},
			wantErr: true,
			mock:    func() {},
		},
		{
			name: "repository error",
			args: args{
				req: &entities.UserRegisterRequest{
					Name:     "John Doe",
					Email:    "john@example.com",
					Password: "password",
				},
			},
			want:    entities.User{},
			wantErr: true,
			mock: func() {
				mockRepo.On("FindByEmail", mock.Anything, "john@example.com").Return(entities.User{}, commons.ErrNotFound)
				mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("entities.User")).Return(errors.New("repository error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil

			tt.mock()
			u := NewUserUsecase(mockRepo, mockJWTConfig, timeout)
			got, err := u.Register(context.TODO(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserUsecase.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserUsecase.Register() = %v, want %v", got, tt.want)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserUsecase_Login(t *testing.T) {
	type args struct {
		req *entities.UserLoginRequest
	}

	mockRepo := new(mocks.UserRepository)
	mockJWTConfig := commons.ConfigJWT{
		SecretJWT:       "secret",
		ExpiresDuration: 1,
	}
	timeout := time.Second * 2

	generateJWT = func(u *UserUsecase, email string) (string, error) {
		return "some-jwt-token", nil
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		mock    func()
	}{
		{
			name: "success",
			args: args{
				req: &entities.UserLoginRequest{
					Email:    "john@example.com",
					Password: "password",
				},
			},
			want:    "some-jwt-token",
			wantErr: false,
			mock: func() {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
				user := entities.User{
					Email:        "john@example.com",
					PasswordHash: string(hashedPassword),
				}
				mockRepo.On("FindByEmail", mock.Anything, "john@example.com").Return(user, nil)
			},
		},
		{
			name: "invalid credentials",
			args: args{
				req: &entities.UserLoginRequest{
					Email:    "john@example.com",
					Password: "wrongpassword",
				},
			},
			want:    "",
			wantErr: true,
			mock: func() {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
				user := entities.User{
					Email:        "john@example.com",
					PasswordHash: string(hashedPassword),
				}
				mockRepo.On("FindByEmail", mock.Anything, "john@example.com").Return(user, nil)
			},
		},
		{
			name: "user not found",
			args: args{
				req: &entities.UserLoginRequest{
					Email:    "john@example.com",
					Password: "password",
				},
			},
			want:    "",
			wantErr: true,
			mock: func() {
				mockRepo.On("FindByEmail", mock.Anything, "john@example.com").Return(entities.User{}, commons.ErrNotFound)
			},
		},
		{
			name: "validation error",
			args: args{
				req: &entities.UserLoginRequest{
					Email:    "invalid-email",
					Password: "short",
				},
			},
			want:    "",
			wantErr: true,
			mock:    func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil

			tt.mock()
			u := NewUserUsecase(mockRepo, mockJWTConfig, timeout)
			got, err := u.Login(context.TODO(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserUsecase.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserUsecase.Login() = %v, want %v", got, tt.want)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
