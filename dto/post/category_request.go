package post

type CategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	ParentID    *uint  `json:"parent_id"`
}
