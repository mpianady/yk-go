package post

import (
	"errors"
	"github.com/gin-gonic/gin"
	commentDTO "go-blog/dto/post"
	commentModel "go-blog/models/post"
	"go-blog/services/config"
	"go-blog/utils"
	commentUtil "go-blog/utils/post"
	"gorm.io/gorm"
	"net/http"
)

const (
	CommentPath         = "/comments"
	CommentAllPath      = "/comments/all"
	CommentByPostIDPath = "/comments/post/:id"
	CommentIdPath       = "/comments/:id"
	CommentNotFound     = "Comment not found"
)

// GetAllComments @Summary Get all comments
// @Description Get list of all comments
// @Tags Comments
// @Produce json
// @Success 200 {array} commentDTO.CommentResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/comments/all [get]
func GetAllComments(ctx *gin.Context) {
	var comments []commentModel.Comment
	if err := config.Db.Preload("User").Find(&comments).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Error retrieving comments"))
		return
	}

	var responses []*commentDTO.CommentResponse
	for _, comment := range comments {
		responses = append(responses, commentDTO.ToCommentResponse(comment))
	}

	ctx.JSON(http.StatusOK, responses)
}

// GetCommentByPostID @Summary Get comments by post ID
// @Description Get paginated comments for a specific post
// @Tags Comments
// @Produce json
// @Param id path string true "Post ID"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} utils.PaginatedResponse[post.CommentResponse]
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/comments/post/{id} [get]
func GetCommentByPostID(ctx *gin.Context) {
	postID := ctx.Param("id")
	if postID == "" {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Missing post ID"))
		return
	}

	// Fetch all comments for the post
	allComments, err := commentUtil.FetchCommentsForPost(postID)
	if err != nil {
		commentUtil.HandleDatabaseError(ctx, "Error retrieving comments")
		return
	}

	// Build a comment tree (all levels)
	commentTree := commentUtil.BuildCommentTree(allComments)

	// Pagination on root level comments
	page, limit, _ := commentUtil.ParsePaginationParams(ctx)
	pagedComments, total, err := commentUtil.PaginateComments(commentTree, page, limit)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.NewErrorResponse("Page not found"))
		return
	}

	// Paginated response
	ctx.JSON(http.StatusOK, utils.NewPaginatedResponse(pagedComments, page, limit, int64(total)))
}

// AddComment @Summary Add new comment
// @Description Create a new comment
// @Tags Comments
// @Accept json
// @Produce json
// @Param comment body commentDTO.CommentRequest true "Comment data"
// @Success 201 {object} commentDTO.CommentResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/comments [post]
func AddComment(ctx *gin.Context) {
	var request commentDTO.CommentRequest
	if !commentUtil.BindAndValidateJSON(ctx, &request) {
		return
	}
	status := commentModel.CommentStatusPending
	if request.Status != nil {
		status = *request.Status
	}
	commentData := commentModel.Comment{
		PostID:  request.PostID,
		UserID:  request.UserID,
		Content: request.Content,
		Status:  status,
	}
	if request.ParentID != nil {
		commentData.ParentID = request.ParentID
	}
	if err := config.Db.Create(&commentData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Error saving comment to the database"))
		return
	}
	if err := config.Db.
		Preload("User").
		Preload("Parent").
		First(&commentData, commentData.ID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Error loading comment"))
		return
	}
	ctx.JSON(http.StatusCreated, commentDTO.ToCommentResponse(commentData))
}

// UpdateComment @Summary Update comment
// @Description Update an existing comment
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID"
// @Param comment body commentDTO.CommentUpdateRequest true "Updated comment data"
// @Success 200 {object} commentDTO.CommentResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/comments/{id} [put]
func UpdateComment(ctx *gin.Context) {
	var request commentDTO.CommentUpdateRequest
	if !commentUtil.BindAndValidateJSON(ctx, &request) {
		return
	}
	id := ctx.Param("id")
	var existingComment commentModel.Comment
	if err := config.Db.First(&existingComment, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, utils.NewErrorResponse(CommentNotFound))
		return
	}
	existingComment.Content = request.Content
	if request.Status != "" {
		existingComment.Status = request.Status
	}
	if err := config.Db.Save(&existingComment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to update comment"))
		return
	}
	if err := config.Db.Preload("User").First(&existingComment, existingComment.ID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Error loading comment author"))
		return
	}
	ctx.JSON(http.StatusOK, commentDTO.ToCommentResponse(existingComment))
}

// DeleteComment @Summary Delete comment
// @Description Delete a comment
// @Tags Comments
// @Produce json
// @Param id path string true "Comment ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/comments/{id} [delete]
func DeleteComment(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Missing comment ID"))
		return
	}
	var comment commentModel.Comment
	if err := config.Db.First(&comment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, utils.NewErrorResponse(CommentNotFound))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Database error while finding comment"))
		}
		return
	}
	if err := config.Db.Delete(&comment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Error deleting comment"))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
