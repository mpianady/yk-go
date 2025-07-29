package post

type PostRequest struct {
	Title  string `json:"title" binding:"required"`
	Excerpt string `json:"excerpt"`
	Content string `json:"content" binding:"required"`
	CategoryIDs []uint `json:"category_ids"`
}
