package service

import (
	"strconv"
	"time"

	"github.com/farid141/go-rest-api/dto"
	"github.com/farid141/go-rest-api/helper"
	"github.com/farid141/go-rest-api/model"
	"github.com/farid141/go-rest-api/repository"
	"github.com/farid141/go-rest-api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(username, password string) (*TokenResponse, error)
	Register(dto.CreateUserRequest) (*dto.CreateUserResponse, error)
}

type authService struct {
	repo   repository.UserRepository
	logger *logrus.Logger
}

func NewAuthService(repo repository.UserRepository, logger *logrus.Logger) AuthService {
	return &authService{repo, logger}
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
}

func (a *authService) Login(username, password string) (*TokenResponse, error) {
	var err error

	a.logger.Info(username + " try to login")

	var user *model.User
	user, err = a.repo.GetUserByUsername(username)
	if err != nil {
		return nil, helper.NewServiceError(fiber.StatusUnauthorized, "Invalid Credentials", nil)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, helper.NewServiceError(fiber.StatusUnauthorized, "Invalid Credentials", nil)
	}

	token, err := utils.GenerateJWT(username, strconv.Itoa(int(user.ID)), 15*time.Minute)
	if err != nil {
		return nil, helper.NewServiceError(fiber.StatusInternalServerError, "Internal Server Error", nil)
	}

	refresh_token, err := utils.GenerateJWT(username, strconv.Itoa(int(user.ID)), 7*24*time.Hour)
	if err != nil {
		return nil, helper.NewServiceError(fiber.StatusInternalServerError, "Internal Server Error", nil)
	}

	return &TokenResponse{AccessToken: token, RefreshToken: refresh_token}, nil
}

func (a *authService) Register(req dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	var err error

	// username unique validation
	_, err = a.repo.GetUserByUsername(req.Username)
	if err == nil {
		return nil, helper.NewServiceError(fiber.StatusConflict, "Username already exists", nil)
	}

	// insert and get id
	var id int
	id, err = a.repo.CreateUser(req)
	if err != nil {
		return nil, err
	}

	// query by id
	var user *model.User
	user, err = a.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.CreateUserResponse{ID: user.ID, Username: user.Username, CreatedAt: user.CreatedAt.String()}, nil
}
