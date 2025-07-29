package post

import "go-blog/models/post"

type CommentResponse struct {
	ID        uint             `json:"id"`
	PostID    uint             `json:"post_id"`
	UserID    uint             `json:"user_id"`
	Author    string           `json:"author"`
	Content   string           `json:"content"`
	Status    string           `json:"status"`
	CreatedAt string           `json:"created_at"`
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
