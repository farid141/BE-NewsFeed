package model

import "time"

type Post struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"userid" db:"userid"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"createdat" db:"createdat"`
}
