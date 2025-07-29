# Go Blog API

A blog API built with Go (Gin framework) that provides features for managing blog posts, categories, user
authentication, and automated news fetching.

## Features

- **Blog Posts Management**: Create, read, update, and delete blog posts
- **Category Management**: Hierarchical category system with parent-child relationships
- **Automated News Fetching**: Periodic fetching of news articles from `https://newsapi.org/v2/everything` every 24
  hours via cron job and goroutines
- **Multi-category Support**: Fetch and categorize news from multiple categories
- **Duplicate Prevention**: Automatic detection and prevention of duplicate posts

## Pre-requisites

- Docker installed and running
- Docker Compose installed
- Git

## Technical Stack

- Go 1.24
- Gin Web Framework
- GORM (Object Relational Mapper)
- JWT Authentication
- News API Integration

## Installation

1. Clone the repository
2. Copy `.env.sample` to `.env`:
   ```bash
   cp .env.sample .env
   ```
3. Fill all the blank values in `.env` file
4. Start the containers:
   ```bash
   ./bin/start.sh
   ```
5. To stop the containers:
   ```bash
   ./bin/stop.sh
   ```

## API Documentation

The API runs on `localhost:8080` and includes the following main endpoints:

- Posts management
- Category management
- User authentication
- News fetching

For detailed API documentation, visit [Swagger Documentation](https://api.go.tiana-zo.site/swagger/index.html) after
starting the server.
