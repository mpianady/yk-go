package utils

import "math"

type PaginatedResponse[T any] struct {
	Data       []T  `json:"data"`
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int  `json:"totalPages"`
	Empty      bool  `json:"empty"`
	First        bool  `json:"first"`
	Last         bool  `json:"last"`
	HasNext      bool  `json:"hasNext"`
	HasPrevious  bool  `json:"hasPrevious"`
}

func NewPaginatedResponse[T any](items []T, page, limit int, total int64) PaginatedResponse[T] {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 && total > 0 {
		totalPages = 1
	}

	if items == nil {
		items = []T{}
	}

	return PaginatedResponse[T]{
		Data:         items,
		Page:         page,
		Limit:        limit,
		Total:        total,
		TotalPages:   totalPages,
		Empty:        len(items) == 0,
		First:        page <= 1,
		Last:         page >= totalPages,
		HasNext:      page < totalPages,
		HasPrevious:  page > 1,
	}
}


