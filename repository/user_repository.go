package repository

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/farid141/go-rest-api/dto"
	"github.com/farid141/go-rest-api/helper"
	"github.com/farid141/go-rest-api/model"
	"github.com/gofiber/fiber/v2"
)

type UserRepository interface {
	GetUsers(userID int, limit, offset int) ([]UserWithProfile, int, error)
	CreateUser(dto.CreateUserRequest) (int, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByID(id int) (*model.User, error)
	FollowUser(follower_id, followed_id string, follow bool) error
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

func (r *userRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(
		`SELECT id, username, password_hash, created_at 
		FROM users WHERE username=?`,
		username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByID(id int) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(
		`SELECT id, username, password_hash, created_at 
		FROM users WHERE id=?`,
		id,
	).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) CreateUser(req dto.CreateUserRequest) (int, error) {
	// username validation
	exists, err := helper.CoulmnValueExists(r.db, "users", "username", req.Username)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, helper.NewServiceError(fiber.StatusConflict, "Username already exists", nil)
	}

	// insert new user
	res, err := r.db.Exec(
		`INSERT INTO users (username, password_hash, created_at) VALUES (?,?,NOW())`,
		req.Username,
		req.Password,
	)
	if err != nil {
		return 0, err
	}

	// get id of new user
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *userRepository) FollowUser(follower_id, followed_id string, follow bool) error {
	var err error

	var followed_str int
	followed_str, err = strconv.Atoi(followed_id)
	if err != nil {
		return err
	}

	// validasi user yang diikuti ada
	_, err = r.GetUserByID(followed_str)
	if err != nil {
		return err
	}

	// follow/unfollow
	queryString := func() string {
		if follow {
			return `INSERT INTO follows (follower_id, followed_id) VALUES (?, ?)`
		} else {
			return `DELETE FROM follows WHERE follower_id=? AND followed_id=?`
		}
	}()

	_, err = r.db.Exec(
		queryString,
		follower_id,
		followed_id,
	)
	if err != nil {
		return err
	}

	return nil
}
