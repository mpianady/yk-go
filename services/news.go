package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"go-blog/dto"
	postModel "go-blog/models/post"
	"go-blog/services/config"
	postUtils "go-blog/utils/post"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type NewsService struct {
	apiKey     string
	apiUrl     string
	categories []string
}

func NewNewsService() *NewsService {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found or unable to load it. Proceeding with existing environment variables.")
	}
	
	// List of required environment variables
	requiredEnv := []string{
		"NEWS_API_KEY", "NEWS_API_URL", "NEWS_CATEGORIES",
	}

	// Verification
	for _, key := range requiredEnv {
		if os.Getenv(key) == "" {
			log.Fatalf("Missing required env var: %s", key)
		}
	}

	return &NewsService{
		apiKey:     os.Getenv("NEWS_API_KEY"),
		apiUrl:     os.Getenv("NEWS_API_URL"),
		categories: strings.Split(os.Getenv("NEWS_CATEGORIES"), ","),
	}
}

type fetchResult struct {
	category string
	posts    []dto.News
}

// FetchAndSaveNews fetches news for each category and saves them to the database.
func FetchAndSaveNews() error {
	service := NewNewsService()
	return service.fetchAndSaveAllCategories()
}

func (ns *NewsService) fetchAndSaveAllCategories() error {
	resultsCh := make(chan fetchResult, len(ns.categories))
	var wg sync.WaitGroup

	for _, cat := range ns.categories {
		wg.Add(1)
		go ns.fetchCategoryAsync(cat, &wg, resultsCh)
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	return ns.processResults(resultsCh)
}

func (ns *NewsService) fetchCategoryAsync(category string, wg *sync.WaitGroup, resultsCh chan<- fetchResult) {
	defer wg.Done()
	posts, err := ns.fetchCategory(category)
	if err != nil {
		log.Printf("Error fetching category %s : %v", category, err)
		resultsCh <- fetchResult{category: category, posts: nil}
		return
	}
	resultsCh <- fetchResult{category: category, posts: posts}
}

func (ns *NewsService) processResults(resultsCh <-chan fetchResult) error {
	for res := range resultsCh {
		if res.posts == nil {
			log.Printf("No posts for category %s", res.category)
			continue
		}

		if err := ns.savePostsForCategory(res.category, res.posts); err != nil {
			log.Printf("Error saving category %s: %v", res.category, err)
		}
	}
	return nil
}

func (ns *NewsService) savePostsForCategory(categoryName string, posts []dto.News) error {
	category, err := postUtils.GetOrCreateCategory(categoryName)
	if err != nil {
		return fmt.Errorf("error creating/retrieving category: %w", err)
	}

	for _, post := range posts {
		if err := ns.savePost(post, category); err != nil {
			log.Printf("Error saving post %s: %v", post.Title, err)
		}
	}
	return nil
}

func (ns *NewsService) savePost(newsPost dto.News, category postModel.Category) error {
	// Check if article already exists
	var existingPost postModel.Post
	err := config.Db.Where("title = ?", newsPost.Title).First(&existingPost).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error searching existing post: %w", err)
	}

	if err == nil {
		log.Printf("Post already exists: %s", newsPost.Title)
		return nil
	}

	// Create post
	postData := postModel.Post{
		Title:   newsPost.Title,
		Excerpt: newsPost.Description,
		Content: newsPost.Content + "<br><br><a href=\"" + newsPost.URL + "\">Read All...</a>",
	}

	// Save post
	if err := config.Db.Create(&postData).Error; err != nil {
		return fmt.Errorf("error DB insertion: %w", err)
	}

	// Associate category
	if err := config.Db.Model(&postData).Association("Categories").Append(&category); err != nil {
		return fmt.Errorf("error category association: %w", err)
	}

	log.Printf("[%s] Post inserted: %s", category.Name, newsPost.Title)
	return nil
}

// fetchCategory retrieves posts for a given category and returns the slice of posts or an error
func (ns *NewsService) fetchCategory(category string) ([]dto.News, error) {
	url := fmt.Sprintf("%s?q=%s&language=en&sortBy=publishedAt&apiKey=%s", ns.apiUrl, category, ns.apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing body (%s): %v", category, err)
		}
	}()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}

	var result dto.NewsResponse
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("error parsing json: %w", err)
	}

	return result.News, nil
}
