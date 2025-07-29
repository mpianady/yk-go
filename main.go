package main

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/robfig/cron/v3"
	_ "go-blog/docs"
	"go-blog/services"
	"go-blog/services/config"
	"go-blog/utils/validators"
	"log"
)

// @title Go Blog API
// @version 1.0
// @description A blog API built with Gin framework providing features like post-categories, blog posts, comments, and user authentication including login, registration, and refresh token functionality.
// @host https://api.go.tiana-zo.site
// @BasePath /
func main() {
	initializeConfiguration()
	setupCustomValidators()
	setupCronJobs()
	startServer()
}

func initializeConfiguration() {
	config.Init()
	config.InitJWTConfig()
}

func setupCustomValidators() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		panic("Failed to get validator engine")
	}

	if err := v.RegisterValidation("strong_password", validators.ValidatePassword); err != nil {
		panic("Failed to register strong_password validator: " + err.Error())
	}
}

func setupCronJobs() {
	c := cron.New()

	_, err := c.AddFunc("@every 24h", func() {
		log.Println("[CRON] Starting news fetching")
		go fetchNewsAsync()
	})
	if err != nil {
		panic("Error while adding cron task: " + err.Error())
	}

	c.Start()

	// Immediate execution once at startup
	go func() {
		log.Println("[CRON] Immediate news fetch on startup...")
		fetchNewsAsync()
	}()
}

func fetchNewsAsync() {
	if err := services.FetchAndSaveNews(); err != nil {
		log.Printf("[CRON] Error while fetching news: %v", err)
	}
}

func startServer() {
	router := services.InitRoutes()

	if err := router.SetTrustedProxies(nil); err != nil {
		panic("Failed to set trusted proxies: " + err.Error())
	}

	if err := router.Run(":8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
