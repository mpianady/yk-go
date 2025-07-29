package post

type CommentRequest struct {
	PostID   uint    `json:"post_id" binding:"required"`
	UserID   uint    `json:"user_id" binding:"required"`
	Content  string  `json:"content" binding:"required"`
	ParentID *uint   `json:"parent_id,omitempty"`
	Status   *string `json:"status,omitempty"`
}

type CommentUpdateRequest struct {
	Content string `json:"content" binding:"required"`
	Status  string `json:"status" binding:"required"`
}