package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"go-blog/models/auth"
	"go-blog/models/post"
	"go-blog/models/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var Db *gorm.DB

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found or unable to load it. Proceeding with existing environment variables.")
	}

	// List of required environment variables
	requiredEnv := []string{
		"DB_USER", "DB_PASSWORD", "DB_HOST", "DB_PORT",
		"DB_NAME", "DB_CHARSET", "DB_PARSE_TIME", "DB_LOC",
	}
	fmt.Println(os.Getenv("DB_HOST"))
	// Validation
	for _, key := range requiredEnv {
		if os.Getenv(key) == "" {
			log.Fatalf("Missing required env var: %s", key)
		}
	}

	connectionString := buildConnectionString()

	var err error
	Db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database connection error: %s", err)
	}

	if err := Db.AutoMigrate(
		&user.User{},
		&auth.RefreshToken{},
		&post.Post{},
		&post.Category{},
		&post.Comment{},
	); err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	fmt.Println("Connected!")
}

func buildConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_CHARSET"),
		os.Getenv("DB_PARSE_TIME"),
		os.Getenv("DB_LOC"),
	)
}
