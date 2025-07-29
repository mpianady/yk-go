package auth

import (
	"github.com/gin-gonic/gin"
	"go-blog/dto/auth"
	"go-blog/models/user"
	"go-blog/services/config"
	"go-blog/utils"
	authUtils "go-blog/utils/auth"
	"net/http"
	"strings"
)

var tokenService = NewTokenService()

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body auth.RegisterRequest true "User registration payload"
// @Success 200 {object} auth.RegisterResponse
// @Failure 400 {object} utils.ValidationErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/register [post]
func Register(ctx *gin.Context) {
	var input auth.RegisterRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		errs := utils.FormatValidationError(err, input)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewValidationErrorResponse(errs))
		return
	}

	var existingUser user.User
	config.Db.Where("email = ?", input.Email).First(&existingUser)

	if existingUser.ID != 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewValidationErrorResponse([]map[string]string{
			{"field": "email", "message": "This email is already registered"},
		}))
		return
	}

	hashedPassword, err := authUtils.HashPassword(input.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.NewErrorResponse("Unable to process password"))
		return
	}

	newUser := user.User{
		Email:     input.Email,
		Password:  hashedPassword,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Role:      "READER", // valeur par défaut
		Status:    "ACTIVE", // valeur par défaut
	}

	if err := config.Db.Create(&newUser).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.NewErrorResponse("Error saving user to the database"))
		return
	}

	accessToken, err := tokenService.GenerateAccessToken(newUser.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to generate access token"))
		return
	}

	refreshToken, err := tokenService.GenerateRefreshToken(newUser.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to generate refresh token"))
		return
	}

	ctx.JSON(http.StatusOK, auth.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return access and refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body auth.LoginRequest true "User login credentials"
// @Success 200 {object} auth.RegisterResponse
// @Failure 400 {object} utils.ValidationErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/login [post]
func Login(ctx *gin.Context) {
	var input auth.LoginRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		errs := utils.FormatValidationError(err, input)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewValidationErrorResponse(errs))
		return
	}

	userInfos, err := authUtils.ValidateCredentials(input.Email, input.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.NewErrorResponse("Invalid credentials"))
		return
	}

	accessToken, err := tokenService.GenerateAccessToken(userInfos.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to generate access token"))
		return
	}

	refreshToken, err := tokenService.GenerateRefreshToken(userInfos.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to generate refresh token"))
	}

	ctx.JSON(http.StatusOK, auth.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.RefreshRequest true "Refresh token request"
// @Success 200 {object} auth.RefreshResponse
// @Failure 400 {object} utils.ValidationErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/refresh-token [post]
func RefreshToken(ctx *gin.Context) {
	var input auth.RefreshRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		errs := utils.FormatValidationError(err, input)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewValidationErrorResponse(errs))
		return
	}

	refreshToken, err := tokenService.ValidateRefreshToken(input.RefreshToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.NewErrorResponse("Invalid refresh token"))
		return
	}

	// Generate a new access token
	userInfos := user.User{}
	if err := config.Db.First(&userInfos, refreshToken.UserID).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.NewErrorResponse("Error retrieving user from database"))
		return
	}
	accessToken, err := tokenService.GenerateAccessToken(userInfos.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to generate access token"))
		return
	}

	ctx.JSON(http.StatusOK, auth.RefreshResponse{
		AccessToken: accessToken,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate refresh token and logout user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.LogoutRequest true "Logout request with refresh token"
// @Success 200 {object} auth.LogoutResponse
// @Failure 400 {object} utils.ValidationErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /v1/logout [post]
func Logout(ctx *gin.Context) {
	var input auth.LogoutRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		errs := utils.FormatValidationError(err, input)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewValidationErrorResponse(errs))
		return
	}

	if err := tokenService.DeleteRefreshToken(input.RefreshToken); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to delete refresh token"))
		return
	}

	ctx.JSON(http.StatusOK, auth.LogoutResponse{Message: "Logout successful"})
}

// Me godoc
// @Summary Get current user information
// @Description Retrieve detailed information about the currently authenticated user
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} auth.UserResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /v1/me [get]
func Me(ctx *gin.Context) {
	userAny, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.NewErrorResponse("User not found in context"))
		return
	}

	userModel, ok := userAny.(user.User)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.NewErrorResponse("Invalid user data in context"))
		return
	}

	ctx.JSON(http.StatusOK, auth.NewUserResponse(userModel))
}

// AuthenticationMiddleWare is a middleware function that validates Authorization tokens and sets user claims in the Gin context.
func AuthenticationMiddleWare(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	tokenString, ok := authUtils.ExtractBearerToken(authHeader)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.NewErrorResponse("Invalid authorization header"))
		return
	}

	claims, err := tokenService.ParseAndValidateAccessToken(tokenString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.NewErrorResponse("Invalid access token"))
		return
	}

	var userModel user.User
	if err := config.Db.Where("email = ?", claims.Email).First(&userModel).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.NewErrorResponse("User not found"))
		return
	}

	ctx.Set("user", userModel)
	ctx.Next()
}

func AuthorizeRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userAny, exists := ctx.Get("user")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.NewErrorResponse("Unauthorized: no user in context"))
			return
		}

		userModel, ok := userAny.(user.User)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.NewErrorResponse("Invalid user in context"))
			return
		}

		for _, role := range allowedRoles {
			if strings.ToUpper(userModel.Role) == strings.ToUpper(role) {
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, utils.NewErrorResponse("Forbidden: insufficient role"))
	}
}
