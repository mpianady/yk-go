package post

import (
	"github.com/gin-gonic/gin"
	postDTO "go-blog/dto/post"
	postModel "go-blog/models/post"
	"go-blog/services/config"
	"go-blog/utils"
	postUtil "go-blog/utils/post"
	"net/http"
	"strconv"
	"strings"
)

const (
	Path   = "/posts"
	IdPath = "/posts/:id"
	NotFound = "Post not found"
)

// CreatePost @Summary Create a new post
// @Description Create a new blog post with optional category assignments
// @Tags Posts
// @Accept json
// @Produce json
// @Param request body post.PostRequest true "Post creation request"
// @Success 201 {object} post.PostResponse
// @Failure 400 {object} utils.ValidationErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/posts [post]
func CreatePost(ctx *gin.Context) {
	var request postDTO.PostRequest
	if !postUtil.BindAndValidateJSON(ctx, &request) {
		return
	}

	postData := postModel.Post{
		Title:   request.Title,
		Excerpt: request.Excerpt,
		Content: request.Content,
	}

	if len(request.CategoryIDs) > 0 {
		var categories []postModel.Category
		if err := config.Db.Where("id IN ?", request.CategoryIDs).Find(&categories).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading categories"})
			return
		}
		postData.Categories = categories
	}

	if err := config.Db.Create(&postData).Error; err != nil {
		postUtil.HandleDatabaseError(ctx, "Error saving post to the database")
		return
	}

	response := postDTO.ToPostResponse(postData)
	ctx.JSON(http.StatusCreated, response)
}

// GetAllPosts @Summary Get paginated list of posts
// @Description Retrieve a paginated list of blog posts with optional category filtering
// @Tags Posts
// @Produce json
// @Param page query int false "Page number (default is 1)"
// @Param limit query int false "Items per page (default is 10)" 
// @Param category_ids query string false "Comma-separated list of category IDs to filter by"
// @Success 200 {object} utils.PaginatedResponse[post.PostResponse]
// @Failure 400 {object} utils.ErrorResponse "Invalid category_ids"
// @Failure 404 {object} utils.ErrorResponse "Page not found"
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/posts [get]
func GetAllPosts(ctx *gin.Context) {
	page, limit, offset := postUtil.ParsePaginationParams(ctx)

	// 1. Get category_ids from a query string
	categoryIDsParam := ctx.Query("category_ids")
	var categoryIDs []uint
	if categoryIDsParam != "" {
		parts := strings.Split(categoryIDsParam, ",")
		for _, part := range parts {
			id, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_ids"})
				return
			}
			categoryIDs = append(categoryIDs, uint(id))
		}
	}

	var posts []postModel.Post
	query := config.Db.Model(&postModel.Post{}).Preload("Categories")

	// 2. Apply category filter if requested
	if len(categoryIDs) > 0 {
		query = query.
			Joins("JOIN post_categories pc ON pc.post_id = posts.id").
			Where("pc.category_id IN ?", categoryIDs).
			Group("posts.id")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		postUtil.HandleDatabaseError(ctx, "Error retrieving total count")
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages > 0 && page > totalPages {
		ctx.JSON(http.StatusNotFound, gin.H{"error": NotFound})
		return
	}

	if err := query.
		Limit(limit).
		Offset(offset).
		Order("posts.id ASC").
		Find(&posts).Error; err != nil {
		postUtil.HandleDatabaseError(ctx, "Error retrieving posts from the database")
		return
	}

	response := make([]postDTO.PostResponse, 0, len(posts))
	for _, post := range posts {
		response = append(response, postDTO.ToPostResponse(post))
	}

	ctx.JSON(http.StatusOK, utils.NewPaginatedResponse(response, page, limit, total))
}

// GetPostByID @Summary Get post by ID
// @Description Retrieve a specific post by its ID
// @Tags Posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} post.PostResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /v1/posts/{id} [get]
func GetPostByID(ctx *gin.Context) {
	var model postModel.Post
	id := ctx.Param("id")

	if err := config.Db.Preload("Categories").First(&model, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, utils.NewErrorResponse(NotFound))
		return
	}

	response := postDTO.ToPostResponse(model)
	ctx.JSON(http.StatusOK, response)
}

// UpdatePost @Summary Update a post
// @Description Update an existing post by its ID
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param request body post.PostRequest true "Post update request"
// @Success 200 {object} post.PostResponse
// @Failure 400 {object} utils.ValidationErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/posts/{id} [put]
func UpdatePost(ctx *gin.Context) {
	id := ctx.Param("id")

	var request postDTO.PostRequest
	if !postUtil.BindAndValidateJSON(ctx, &request) {
		return
	}

	var post postModel.Post
	if err := config.Db.Preload("Categories").First(&post, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, utils.NewErrorResponse(NotFound))
		return
	}

	post.Title = request.Title
	post.Excerpt = request.Excerpt
	post.Content = request.Content

	// Update associated categories
	if len(request.CategoryIDs) > 0 {
		var categories []postModel.Category
		if err := config.Db.Where("id IN ?", request.CategoryIDs).Find(&categories).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading categories"})
			return
		}
		// Replace existing associations
		if err := config.Db.Model(&post).Association("Categories").Replace(&categories); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating categories"})
			return
		}
	} else {
		// Clear categories if none are sent
		if err := config.Db.Model(&post).Association("Categories").Clear(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error clearing categories"})
			return
		}
	}

	if err := config.Db.Save(&post).Error; err != nil {
		postUtil.HandleDatabaseError(ctx, "Error updating post")
		return
	}

	response := postDTO.ToPostResponse(post)
	ctx.JSON(http.StatusOK, response)
}

// DeletePost @Summary Delete a post
// @Description Delete a post by its ID
// @Tags Posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} string "Post deleted successfully"
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/posts/{id} [delete]
func DeletePost(ctx *gin.Context) {
	id := ctx.Param("id")

	var post postModel.Post
	if err := config.Db.First(&post, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, utils.NewErrorResponse(NotFound))
		return
	}

	// Clear associations manually if needed
	if err := config.Db.Model(&post).Association("Categories").Clear(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error clearing categories"})
		return
	}

	if err := config.Db.Delete(&post).Error; err != nil {
		postUtil.HandleDatabaseError(ctx, "Error deleting post")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
