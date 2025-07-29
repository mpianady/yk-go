package post

import (
	postModel "go-blog/models/post"
	"time"
)

type CategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parentId,omitempty"`
	Children    []CategoryResponse `json:"children,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func ToCategoryResponse(cat postModel.Category) CategoryResponse {
	children := make([]CategoryResponse, len(cat.Children))
	for i, child := range cat.Children {
		children[i] = ToCategoryResponse(child)
	}

	return CategoryResponse{
		ID:          cat.ID,
		Name:        cat.Name,
		Description: cat.Description,
		ParentID:    cat.ParentID,
		Children:    children,
		CreatedAt:   cat.CreatedAt,
		UpdatedAt:   cat.UpdatedAt,
	}
}
