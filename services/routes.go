package services

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "go-blog/docs"
	"go-blog/services/auth"
	"go-blog/services/post"
)

func InitRoutes() *gin.Engine {
	router := gin.Default()

	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Main group with prefix /v1
	v1 := router.Group("/v1")

	setupPublicRoutes(v1)
	setupProtectedRoutes(v1)

	return router
}

func setupPublicRoutes(v1 *gin.RouterGroup) {
	v1.POST("/register", auth.Register)
	v1.POST("/login", auth.Login)
	v1.POST("/refresh-token", auth.RefreshToken)

	// Post routes
	v1.GET(post.Path, post.GetAllPosts)
	v1.GET(post.IdPath, post.GetPostByID)
	v1.GET(post.CategoryPath, post.GetAllCategories)
	v1.GET(post.CategoryTreePath, post.GetCategoryTree)
	v1.GET(post.CategoryIDPath, post.GetCategoryByID)
	v1.GET(post.CommentByPostIDPath, post.GetCommentByPostID)
}

func setupProtectedRoutes(v1 *gin.RouterGroup) {
	protected := v1.Group("/", auth.AuthenticationMiddleWare)

	// Auth routes
	protected.POST("/logout", auth.Logout)
	protected.GET("/me", auth.Me)

	// Routes only accessible to ADMIN
	adminOnly := protected.Group("/", auth.AuthorizeRoles("ADMIN"))
	{
		adminOnly.POST(post.CategoryPath, post.CreateCategory)
		adminOnly.PUT(post.CategoryIDPath, post.UpdateCategory)
		adminOnly.DELETE(post.CategoryIDPath, post.DeleteCategory)
	}

	// Routes accessible to ADMIN and AUTHOR
	authorOrAdmin := protected.Group("/", auth.AuthorizeRoles("ADMIN", "AUTHOR"))
	{
		authorOrAdmin.POST(post.Path, post.CreatePost)
		authorOrAdmin.PUT(post.IdPath, post.UpdatePost)
		authorOrAdmin.DELETE(post.IdPath, post.DeletePost)

		protected.GET(post.CommentAllPath, post.GetAllComments)
	}

	// Public routes for all authenticated users
	protected.POST(post.CommentPath, post.AddComment)
	protected.PUT(post.CommentIdPath, post.UpdateComment)
	protected.DELETE(post.CommentIdPath, post.DeleteComment)
}
