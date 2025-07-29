package dto

type News struct {
	Title string `json:"title"`
	Description string `json:"description"`
	URL string `json:"url"`
	Content string `json:"content"`
	PublishedAt string `json:"published_at"`
}

type NewsResponse struct {
	Status string `json:"status"`
	TotalResults int `json:"totalResults"`
	News []News `json:"articles"`
}
