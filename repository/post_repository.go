package repository

import (
	"database/sql"

	"github.com/farid141/go-rest-api/dto"
	"github.com/farid141/go-rest-api/model"
	"github.com/farid141/go-rest-api/response"
)

type PostRepository interface {
	CreatePost(userID string, post dto.CreatePostRequest) (*model.Post, error)
	GetFeed(userID string, limit, offset int) ([]model.Post, *response.Pagination, error)
}

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{db: db}
}

func (p *postRepository) CreatePost(userID string, post dto.CreatePostRequest) (*model.Post, error) {
	// insert post
	res, err := p.db.Exec(
		`INSERT INTO posts (user_id, content, created_at) VALUES (?,?,NOW())`,
		userID,
		post.Content,
	)
	if err != nil {
		return nil, err
	}

	// get id and data
	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	var createdPost model.Post
	err = p.db.QueryRow("SELECT id, user_id, content, created_at FROM posts WHERE id = ?", lastID).
		Scan(&createdPost.ID, &createdPost.UserID, &createdPost.Content, &createdPost.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &model.Post{
		ID:        createdPost.ID,
		UserID:    createdPost.UserID,
		Content:   createdPost.Content,
		CreatedAt: createdPost.CreatedAt,
	}, nil
}

func (p *postRepository) GetFeed(userID string, limit, offset int) ([]model.Post, *response.Pagination, error) {
	var err error

	// getting total
	var total int
	err = p.db.QueryRow(`
            SELECT COUNT(*)
            FROM posts p
            LEFT JOIN follows f ON f.followed_id = p.user_id
            WHERE f.follower_id = ? OR p.user_id = ?
        `, userID, userID).Scan(&total)
	if err != nil {
		return nil, nil, err
	}

	var rows *sql.Rows
	rows, err = p.db.Query(`
            SELECT p.id, p.user_id, p.content, p.created_at
            FROM posts p
            LEFT JOIN follows f ON f.followed_id = p.user_id
            WHERE f.follower_id = ? OR p.user_id = ?
            ORDER BY p.created_at DESC
            LIMIT ? OFFSET ?
        `, userID, userID, limit, offset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	// store result to var
	posts := make([]model.Post, 0)
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Content, &p.CreatedAt); err != nil {
			return nil, nil, err
		}

		posts = append(posts, p)
	}

	return posts, &response.Pagination{
		Limit: limit,
		Total: total,
	}, nil
}
