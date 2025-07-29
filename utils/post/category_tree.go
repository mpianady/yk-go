package post

import (
	"go-blog/dto/post"
	postModel "go-blog/models/post"
)

func BuildCategoryTree(categories []postModel.Category) []post.CategoryResponse {
	mapped := convertCategoriesToMap(categories)
	return buildHierarchy(categories, mapped)
}

func convertCategoriesToMap(categories []postModel.Category) map[uint]*post.CategoryResponse {
	mapped := make(map[uint]*post.CategoryResponse)
	for _, cat := range categories {
		response := post.ToCategoryResponse(cat)
		response.Children = []post.CategoryResponse{} // Reset children for tree building
		mapped[cat.ID] = &response
	}
	return mapped
}

func buildHierarchy(categories []postModel.Category, mapped map[uint]*post.CategoryResponse) []post.CategoryResponse {
	var roots []post.CategoryResponse
	for _, cat := range categories {
		if cat.ParentID != nil {
			if parent, ok := mapped[*cat.ParentID]; ok {
				parent.Children = append(parent.Children, *mapped[cat.ID])
			}
		} else {
			roots = append(roots, *mapped[cat.ID])
		}
	}
	return roots
}
