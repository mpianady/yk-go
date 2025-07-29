package post

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	postModel "go-blog/models/post"
	"go-blog/services/config"
	"go-blog/utils"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strconv"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
)

func DeleteCategoryRecursively(category *postModel.Category) error {
	// Charger récursivement tous les enfants
	for i := range category.Children {
		if err := config.Db.Preload("Children").First(&category.Children[i], category.Children[i].ID).Error; err != nil {
			return err
		}
		if err := DeleteCategoryRecursively(&category.Children[i]); err != nil {
			return err
		}
	}
	// Supprimer la catégorie
	return config.Db.Delete(category).Error
}

// BindAndValidateJSON handles JSON binding and validation for request objects
func BindAndValidateJSON(ctx *gin.Context, request interface{}) bool {
	if err := ctx.ShouldBindJSON(request); err != nil {
		if errors.Is(err, io.EOF) {
			ctx.JSON(http.StatusBadRequest, utils.NewValidationErrorResponse([]map[string]string{
				{"field": "body", "message": "Request body cannot be empty"},
			}))
			return false
		}
		errs := utils.FormatValidationError(err, request)
		ctx.JSON(http.StatusBadRequest, utils.NewValidationErrorResponse(errs))
		return false
	}
	return true
}

// ParsePaginationParams extracts and validates pagination parameters from request
func ParsePaginationParams(ctx *gin.Context) (page, limit, offset int) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", strconv.Itoa(DefaultPage)))
	if err != nil || page < 1 {
		page = DefaultPage
	}

	limit, err = strconv.Atoi(ctx.DefaultQuery("limit", strconv.Itoa(DefaultLimit)))
	if err != nil || limit < 1 {
		limit = DefaultLimit
	}

	offset = (page - 1) * limit
	return page, limit, offset
}

// HandleDatabaseError sends appropriate error response for database operations
func HandleDatabaseError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(message))
}

func GetOrCreateCategory(name string) (postModel.Category, error) {
	var category postModel.Category
	result := config.Db.Where("name = ?", name).First(&category)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		category = postModel.Category{Name: name}
		if err := config.Db.Create(&category).Error; err != nil {
			return category, fmt.Errorf("failed to create category %s: %w", name, err)
		}
		return category, nil
	}
	if result.Error != nil {
		return category, fmt.Errorf("error looking for category %s: %w", name, result.Error)
	}

	return category, nil
}


