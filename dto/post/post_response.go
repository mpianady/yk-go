package post

import (
	"go-blog/models/post"
	"time"
)

type PostResponse struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Excerpt    string    `json:"excerpt"`
	Content    string    `json:"content"`
	Categories []uint    `json:"category_ids"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ToPostResponse(post post.Post) PostResponse {
	categoryIDs := make([]uint, len(post.Categories))

	for i, cat := range post.Categories {
		categoryIDs[i] = cat.ID
	}

	return PostResponse{
		ID:         post.ID,
		Title:      post.Title,
		Excerpt:    post.Excerpt,
		Content:    post.Content,
		Categories: categoryIDs,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
	}
}
