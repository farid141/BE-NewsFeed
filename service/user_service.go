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

	return response.PaginatedResponse[dto.UserResponse]{
		Data: users,
		Pagination: response.Pagination{
			Page:    page,
			Limit:   limit,
			Total:   total,
			HasMore: page*limit < total,
		},
	}, nil
}
