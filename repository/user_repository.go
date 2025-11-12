package repository

import (
	"database/sql"
	"time"
)

type UserRepository interface {
	GetUsers(userID int, limit, offset int) ([]UserWithProfile, int, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

type UserWithProfile struct {
	ID        int64
	Username  string
	Password  string
	CreatedAt time.Time
	Following bool
}

func (r *userRepository) GetUsers(userID int, limit, offset int) ([]UserWithProfile, int, error) {
	var total int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(`
		SELECT 
			u.id,
			u.username,
			CASE WHEN f.follower_id IS NOT NULL THEN TRUE ELSE FALSE END AS is_following,
			u.created_at
		FROM users u
		LEFT JOIN follows f ON f.followed_id = u.id AND f.follower_id = ?
		WHERE u.id != ?
		ORDER BY u.created_at DESC
		LIMIT ? OFFSET ?`,
		userID, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := []UserWithProfile{}
	for rows.Next() {
		var u UserWithProfile
		if err := rows.Scan(&u.ID, &u.Username, &u.Following, &u.CreatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	return users, total, nil
}
