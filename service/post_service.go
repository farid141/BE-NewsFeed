package service

import (
	"github.com/farid141/go-rest-api/dto"
	"github.com/farid141/go-rest-api/repository"
	"github.com/farid141/go-rest-api/response"
	"github.com/sirupsen/logrus"
)

type PostService interface {
	CreatePost(userID string, post dto.CreatePostRequest) (*dto.PostResponse, error)
	GetFeed(userID string, page, limit, offset int) (*response.PaginatedResponse[dto.PostResponse], error)
}

type postService struct {
	repo   repository.PostRepository
	logger *logrus.Logger
}

func NewPostService(repo repository.PostRepository, logger *logrus.Logger) PostService {
	return &postService{repo, logger}
}

func (p *postService) CreatePost(userID string, req dto.CreatePostRequest) (*dto.PostResponse, error) {
	post, err := p.repo.CreatePost(userID, req)
	if err != nil {
		p.logger.Error("error create post")
		return nil, err
	}
	return &dto.PostResponse{
		ID:        post.ID,
		UserID:    post.UserID,
		Content:   post.Content,
		CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (p *postService) GetFeed(userID string, page, limit, offset int) (*response.PaginatedResponse[dto.PostResponse], error) {
	rows, pagination, err := p.repo.GetFeed(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	posts := make([]dto.PostResponse, 0)
	for _, row := range rows {
		p := dto.PostResponse{
			ID:        row.ID,
			Content:   row.Content,
			UserID:    row.UserID,
			CreatedAt: row.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		posts = append(posts, p)
	}

	return &response.PaginatedResponse[dto.PostResponse]{
		Data: posts,
		Pagination: response.Pagination{
			Page:    page,
			Limit:   limit,
			Total:   pagination.Total,
			HasMore: page*limit < pagination.Total,
		},
	}, nil
}
