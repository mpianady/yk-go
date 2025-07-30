package post

import (
	"errors"
	commentDTO "go-blog/dto/post"
	commentModel "go-blog/models/post"
	"go-blog/services/config"
)

// FetchCommentsForPost Helper function to fetch comments for a post
func FetchCommentsForPost(postID string) ([]commentModel.Comment, error) {
	var comments []commentModel.Comment
	err := config.Db.Preload("User").Where("post_id = ?", postID).Order("created_at ASC").Find(&comments).Error
	return comments, err
}

// BuildCommentTree Helper function to build a comment tree structure
func BuildCommentTree(allComments []commentModel.Comment) []*commentDTO.CommentResponse {
	commentMap := make(map[uint]*commentDTO.CommentResponse)

	// Étape 1 : créer tous les commentaires (vides de children)
	for _, c := range allComments {
		cr := commentDTO.ToCommentResponse(c)
		cr.Children = []*commentDTO.CommentResponse{}
		commentMap[c.ID] = cr
	}

	var roots []*commentDTO.CommentResponse

	// Étape 2 : remplir la hiérarchie
	for _, c := range allComments {
		if c.ParentID != nil {
			parent, exists := commentMap[*c.ParentID]
			if exists {
				child := commentMap[c.ID]
				parent.Children = append(parent.Children, child)
			}
		} else {
			roots = append(roots, commentMap[c.ID])
		}
	}

	return roots
}

// PaginateComments Helper function to paginate comments
func PaginateComments(roots []*commentDTO.CommentResponse, page, limit int) ([]*commentDTO.CommentResponse, int, error) {
	total := len(roots)
	totalPages := (total + limit - 1) / limit

	if totalPages > 0 && page > totalPages {
		return nil, 0, errors.New("page not found")
	}

	start := (page - 1) * limit
	end := start + limit
	if end > total {
		end = total
	}

	return roots[start:end], total, nil
}
