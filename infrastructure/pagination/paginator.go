package pagination

import (
	"fmt"
	"math"
)

// Page represents a standardized paginated API response.
type Page[T any] struct {
	Items      []T               `json:"items"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalItems int64             `json:"total_items"`
	TotalPages int               `json:"total_pages"`
	Links      map[string]string `json:"links,omitempty"`
}

// NewPage creates a paginated response object.
func NewPage[T any](items []T, totalItems int64, page, pageSize int, baseURL string) Page[T] {
	totalPages := 0
	if pageSize > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(pageSize)))
	}

	p := Page[T]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
		Links:      make(map[string]string),
	}

	p.Links["self"] = fmt.Sprintf("%s?page=%d&page_size=%d", baseURL, page, pageSize)
	if page > 1 {
		p.Links["first"] = fmt.Sprintf("%s?page=1&page_size=%d", baseURL, pageSize)
		p.Links["prev"] = fmt.Sprintf("%s?page=%d&page_size=%d", baseURL, page-1, pageSize)
	}
	if page < totalPages {
		p.Links["next"] = fmt.Sprintf("%s?page=%d&page_size=%d", baseURL, page+1, pageSize)
		p.Links["last"] = fmt.Sprintf("%s?page=%d&page_size=%d", baseURL, totalPages, pageSize)
	}

	return p
}
