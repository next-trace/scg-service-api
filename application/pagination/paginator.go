// Package pagination provides interfaces and types for pagination functionality.
package pagination

// PageInfo contains metadata about a paginated response.
type PageInfo struct {
	// Page is the current page number.
	Page int `json:"page"`

	// PageSize is the number of items per page.
	PageSize int `json:"page_size"`

	// TotalItems is the total number of items across all pages.
	TotalItems int64 `json:"total_items"`

	// TotalPages is the total number of pages.
	TotalPages int `json:"total_pages"`

	// Links contains navigation links for the paginated response.
	Links map[string]string `json:"links,omitempty"`
}

// PaginationOptions contains options for paginating a collection.
type PaginationOptions struct {
	// Page is the requested page number (1-based).
	Page int

	// PageSize is the requested number of items per page.
	PageSize int

	// BaseURL is the base URL used to generate navigation links.
	BaseURL string
}

// DefaultPageSize is the default number of items per page if not specified.
const DefaultPageSize = 20

// DefaultPaginationOptions returns default pagination options.
func DefaultPaginationOptions() PaginationOptions {
	return PaginationOptions{
		Page:     1,
		PageSize: DefaultPageSize,
		BaseURL:  "",
	}
}
