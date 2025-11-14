// internal/service/user_service.go
package service

import (
	"github.com/farid141/go-rest-api/dto"
	"github.com/farid141/go-rest-api/helper"
	"github.com/farid141/go-rest-api/repository"
	"github.com/farid141/go-rest-api/response"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserService interface {
	ListUsers(userID int, page, limit, offset int) (response.PaginatedResponse[dto.UserResponse], error)
	FollowUser(follower_id, followed_id string, follow bool) error
}

type userService struct {
	repo   repository.UserRepository
	logger *logrus.Logger
}

func NewUserService(repo repository.UserRepository, logger *logrus.Logger) UserService {
	return &userService{repo, logger}
}

func (s *userService) ListUsers(userID int, page, limit, offset int) (response.PaginatedResponse[dto.UserResponse], error) {
	users, total, err := s.repo.GetUsers(userID, limit, offset)
	if err != nil {
		return response.PaginatedResponse[dto.UserResponse]{}, helper.NewServiceError(fiber.StatusInternalServerError, "Internal Server Error", nil)
	}

	userDTOs := make([]dto.UserResponse, len(users))
	for i, u := range users {
		userDTOs[i] = dto.UserResponse{
			ID:        u.ID,
			Username:  u.Username,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
			Following: u.Following,
		}
	}

	return response.PaginatedResponse[dto.UserResponse]{
		Data: userDTOs,
		Pagination: response.Pagination{
			Page:    page,
			Limit:   limit,
			Total:   total,
			HasMore: page*limit < total,
		},
	}, nil
}

func (s *userService) FollowUser(follower_id, followed_id string, follow bool) error {
	err := s.repo.FollowUser(follower_id, followed_id, follow)
	if err != nil {
		return helper.NewServiceError(fiber.StatusInternalServerError, "Internal Server Error", nil)
	}

	return nil
}
