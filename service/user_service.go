// internal/service/user_service.go
package service

import (
	"github.com/farid141/go-rest-api/dto"
	"github.com/farid141/go-rest-api/repository"
	"github.com/farid141/go-rest-api/response"
)

type UserService interface {
	ListUsers(userID int, page, limit, offset int) (response.PaginatedResponse[dto.UserResponse], error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) ListUsers(userID int, page, limit, offset int) (response.PaginatedResponse[dto.UserResponse], error) {
	users, total, err := s.repo.GetUsers(userID, limit, offset)
	if err != nil {
		return response.PaginatedResponse[dto.UserResponse]{}, err
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
