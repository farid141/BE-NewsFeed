package model

import "time"

type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `json:"createdat" db:"created_at"`
}
