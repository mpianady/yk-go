package post

import (
	"go-blog/models/post"
)

// swagger:model CommentResponseDoc
type CommentResponseDoc struct {
	ID        uint                      `json:"id"`
	PostID    uint                      `json:"post_id"`
	UserID    uint                      `json:"user_id"`
	Author    string                    `json:"author"`
	Content   string                    `json:"content"`
	Status    string                    `json:"status"`
	CreatedAt string                    `json:"created_at"`
	Children  []CommentResponseChildDoc `json:"children,omitempty"`
}

// swagger:model CommentResponseChildDoc
type CommentResponseChildDoc struct {
	ID        uint   `json:"id"`
	PostID    uint   `json:"post_id"`
	UserID    uint   `json:"user_id"`
	Author    string `json:"author"`
	Content   string `json:"content"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// swagger:model PaginatedCommentResponse
type PaginatedCommentResponse struct {
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
	Total int64                `json:"total"`
	Data  []CommentResponseDoc `json:"data"`
}

type CommentResponse struct {
	ID        uint               `json:"id"`
	PostID    uint               `json:"post_id"`
	UserID    uint               `json:"user_id"`
	Author    string             `json:"author"`
	Content   string             `json:"content"`
	Status    string             `json:"status"`
	CreatedAt string             `json:"created_at"`
	Children  []*CommentResponse `json:"children,omitempty"`
}

func ToCommentResponse(comment post.Comment) *CommentResponse {
	authorName := comment.User.FirstName + " " + comment.User.LastName

	children := make([]*CommentResponse, 0, len(comment.Children))
	for _, child := range comment.Children {
		children = append(children, ToCommentResponse(child)) // récursion
	}

	return &CommentResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		Author:    authorName,
		Content:   comment.Content,
		Status:    comment.Status,
		CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"), // format plus lisible
		Children:  children,
	}
}
