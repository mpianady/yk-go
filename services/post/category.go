package post

import (
	"github.com/gin-gonic/gin"
	categoryDTO "go-blog/dto/post"
	postModel "go-blog/models/post"
	"go-blog/services/config"
	"go-blog/utils"
	categoryUtil "go-blog/utils/post"
	"net/http"
)

const (
	CategoryPath     = "/categories"
	CategoryIDPath   = "/categories/:id"
	CategoryTreePath = "/categories/tree"
	CategoryNotFound = "Category not found"
)

// CreateCategory @Summary Create a new category
// @Description Create a new blog post category
// @Tags Categories
// @Accept json
// @Produce json
// @Param request body post.CategoryRequest true "Category creation request"
// @Success 201 {object} post.CategoryResponse
// @Failure 400 {object} utils.ValidationErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/categories [post]
func CreateCategory(ctx *gin.Context) {
	var request categoryDTO.CategoryRequest
	if !categoryUtil.BindAndValidateJSON(ctx, &request) {
		return
	}
	categoryData := postModel.Category{
		Name:        request.Name,
		Description: request.Description,
		ParentID:    request.ParentID,
	}

	if err := config.Db.Create(&categoryData).Error; err != nil {
		categoryUtil.HandleDatabaseError(ctx, "Error saving category to the database")
		return
	}

	response := categoryDTO.ToCategoryResponse(categoryData)
	ctx.JSON(http.StatusCreated, response)
}

// GetAllCategories @Summary Get paginated list of categories
// @Description Retrieve a paginated list of blog post categories
// @Tags Categories
// @Produce json
// @Param page query int false "Page number (default is 1)"
// @Param limit query int false "Items per page (default is 10)"
// @Success 200 {object} utils.PaginatedResponse[post.CategoryResponse]
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/categories [get]
func GetAllCategories(ctx *gin.Context) {
	page, limit, offset := categoryUtil.ParsePaginationParams(ctx)

	var total int64
	if err := config.Db.Model(&postModel.Category{}).Count(&total).Error; err != nil {
		categoryUtil.HandleDatabaseError(ctx, "Error retrieving total count")
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages > 0 && page > totalPages {
		ctx.JSON(http.StatusNotFound, utils.NewErrorResponse("Page not found"))
		return
	}

	var categories []postModel.Category
	if err := config.Db.
		Limit(limit).
		Offset(offset).
		Order("id ASC").
		Find(&categories).Error; err != nil {
		categoryUtil.HandleDatabaseError(ctx, "Error retrieving categories from the database")
		return
	}

	response := make([]categoryDTO.CategoryResponse, 0, len(categories))
	for _, cat := range categories {
		response = append(response, categoryDTO.ToCategoryResponse(cat))
	}

	ctx.JSON(http.StatusOK, utils.NewPaginatedResponse(response, page, limit, total))
}

// GetCategoryTree @Summary Get category tree
// @Description Retrieve hierarchical tree structure of all categories
// @Tags Categories
// @Produce json
// @Success 200 {array} post.CategoryResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/categories/tree [get]
func GetCategoryTree(ctx *gin.Context) {
	var categories []postModel.Category
	if err := config.Db.Order("name ASC").Find(&categories).Error; err != nil {
		categoryUtil.HandleDatabaseError(ctx, "Error retrieving categories")
		return
	}
	tree := categoryUtil.BuildCategoryTree(categories)
	ctx.JSON(http.StatusOK, tree)
}

// GetCategoryByID @Summary Get category by ID
// @Description Retrieve a specific category by its ID
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} post.CategoryResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/categories/{id} [get]
func GetCategoryByID(ctx *gin.Context) {
	var categoryModel postModel.Category
	id := ctx.Param("id")
	if err := config.Db.Preload("Children").First(&categoryModel, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, utils.NewErrorResponse(CategoryNotFound))
		return
	}
	response := categoryDTO.ToCategoryResponse(categoryModel)
	ctx.JSON(http.StatusOK, response)
}

// UpdateCategory @Summary Update a category
// @Description Update an existing category by its ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param request body post.CategoryRequest true "Category update request"
// @Success 200 {object} post.CategoryResponse
// @Failure 400 {object} utils.ValidationErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/categories/{id} [put]
func UpdateCategory(ctx *gin.Context) {
	var request categoryDTO.CategoryRequest
	if !categoryUtil.BindAndValidateJSON(ctx, &request) {
		return
	}
	var categoryModel postModel.Category
	id := ctx.Param("id")
	if err := config.Db.First(&categoryModel, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, utils.NewErrorResponse(CategoryNotFound))
		return
	}
	categoryModel.Name = request.Name
	categoryModel.Description = request.Description
	if request.ParentID != nil {
		categoryModel.ParentID = request.ParentID
	} else {
		categoryModel.ParentID = nil
	}
	if err := config.Db.Save(&categoryModel).Error; err != nil {
		categoryUtil.HandleDatabaseError(ctx, "Error updating category")
		return
	}
	if err := config.Db.Preload("Children").First(&categoryModel, id).Error; err != nil {
		categoryUtil.HandleDatabaseError(ctx, "Error retrieving updated category")
		return
	}
	response := categoryDTO.ToCategoryResponse(categoryModel)
	ctx.JSON(http.StatusOK, response)
}

// DeleteCategory @Summary Delete a category
// @Description Delete a category and handle its children (recursively or by reassignment)
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param recursive query bool false "Delete children recursively"
// @Param reassign_to query int false "ID of category to reassign children to"
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/categories/{id} [delete]
func DeleteCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	var category postModel.Category
	if err := config.Db.Preload("Children").First(&category, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, utils.NewErrorResponse(CategoryNotFound))
		return
	}
	recursive := ctx.Query("recursive") == "true"
	reassignTo := ctx.Query("reassign_to")
	if len(category.Children) > 0 {
		switch {
		case recursive:
			if err := categoryUtil.DeleteCategoryRecursively(&category); err != nil {
				categoryUtil.HandleDatabaseError(ctx, "Failed to delete children categories")
				return
			}
		case reassignTo != "":
			var newParent postModel.Category
			if err := config.Db.First(&newParent, reassignTo).Error; err != nil {
				ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Reassign target category not found"))
				return
			}
			if err := config.Db.Model(&postModel.Category{}).
				Where("parent_id = ?", category.ID).
				Update("parent_id", newParent.ID).Error; err != nil {
				categoryUtil.HandleDatabaseError(ctx, "Failed to reassign child categories")
				return
			}
		default:
			ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Cannot delete category with children. Use recursive=true or reassign_to={id}"))
			return
		}
	}
	if err := config.Db.Delete(&category).Error; err != nil {
		categoryUtil.HandleDatabaseError(ctx, "Failed to delete category")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
