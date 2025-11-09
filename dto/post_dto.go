package dto

type CreatePostRequest struct {
	Content string `json:"content" validate:"required,max=200"`
}

type PostResponse struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"userid"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdat"`
}
