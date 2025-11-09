package dto

import "time"

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Following bool      `json:"following"`
	CreatedAt time.Time `json:"createdat"`
}
