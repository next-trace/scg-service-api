package pagination_test

import (
	"testing"

	"github.com/next-trace/scg-service-api/infrastructure/pagination"
	"github.com/stretchr/testify/assert"
)

type TestItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestNewPage(t *testing.T) {
	// Test data
	items := []TestItem{
		{ID: 1, Name: "Item 1"},
		{ID: 2, Name: "Item 2"},
		{ID: 3, Name: "Item 3"},
	}
	baseURL := "/api/items"

	// Test first page with more pages available
	t.Run("First page with more pages", func(t *testing.T) {
		page := pagination.NewPage(items, 10, 1, 3, baseURL)

		// Check basic pagination info
		assert.Equal(t, items, page.Items)
		assert.Equal(t, 1, page.Page)
		assert.Equal(t, 3, page.PageSize)
		assert.Equal(t, int64(10), page.TotalItems)
		assert.Equal(t, 4, page.TotalPages) // Ceiling of 10/3 = 4

		// Check links
		assert.Equal(t, "/api/items?page=1&page_size=3", page.Links["self"])
		assert.NotContains(t, page.Links, "prev")
		assert.NotContains(t, page.Links, "first")
		assert.Equal(t, "/api/items?page=2&page_size=3", page.Links["next"])
		assert.Equal(t, "/api/items?page=4&page_size=3", page.Links["last"])
	})

	// Test middle page
	t.Run("Middle page", func(t *testing.T) {
		page := pagination.NewPage(items, 10, 2, 3, baseURL)

		// Check basic pagination info
		assert.Equal(t, items, page.Items)
		assert.Equal(t, 2, page.Page)
		assert.Equal(t, 3, page.PageSize)
		assert.Equal(t, int64(10), page.TotalItems)
		assert.Equal(t, 4, page.TotalPages)

		// Check links
		assert.Equal(t, "/api/items?page=2&page_size=3", page.Links["self"])
		assert.Equal(t, "/api/items?page=1&page_size=3", page.Links["prev"])
		assert.Equal(t, "/api/items?page=1&page_size=3", page.Links["first"])
		assert.Equal(t, "/api/items?page=3&page_size=3", page.Links["next"])
		assert.Equal(t, "/api/items?page=4&page_size=3", page.Links["last"])
	})

	// Test last page
	t.Run("Last page", func(t *testing.T) {
		page := pagination.NewPage(items, 10, 4, 3, baseURL)

		// Check basic pagination info
		assert.Equal(t, items, page.Items)
		assert.Equal(t, 4, page.Page)
		assert.Equal(t, 3, page.PageSize)
		assert.Equal(t, int64(10), page.TotalItems)
		assert.Equal(t, 4, page.TotalPages)

		// Check links
		assert.Equal(t, "/api/items?page=4&page_size=3", page.Links["self"])
		assert.Equal(t, "/api/items?page=3&page_size=3", page.Links["prev"])
		assert.Equal(t, "/api/items?page=1&page_size=3", page.Links["first"])
		assert.NotContains(t, page.Links, "next")
		assert.NotContains(t, page.Links, "last")
	})

	// Test empty items
	t.Run("Empty items", func(t *testing.T) {
		emptyItems := []TestItem{}
		page := pagination.NewPage(emptyItems, 0, 1, 10, baseURL)

		// Check basic pagination info
		assert.Empty(t, page.Items)
		assert.Equal(t, 1, page.Page)
		assert.Equal(t, 10, page.PageSize)
		assert.Equal(t, int64(0), page.TotalItems)
		assert.Equal(t, 0, page.TotalPages)

		// Check links
		assert.Equal(t, "/api/items?page=1&page_size=10", page.Links["self"])
		assert.NotContains(t, page.Links, "prev")
		assert.NotContains(t, page.Links, "first")
		assert.NotContains(t, page.Links, "next")
		assert.NotContains(t, page.Links, "last")
	})

	// Test zero page size
	t.Run("Zero page size", func(t *testing.T) {
		page := pagination.NewPage(items, 10, 1, 0, baseURL)

		// Check basic pagination info
		assert.Equal(t, items, page.Items)
		assert.Equal(t, 1, page.Page)
		assert.Equal(t, 0, page.PageSize)
		assert.Equal(t, int64(10), page.TotalItems)
		assert.Equal(t, 0, page.TotalPages) // Should be 0 when page size is 0

		// Check links
		assert.Equal(t, "/api/items?page=1&page_size=0", page.Links["self"])
		assert.NotContains(t, page.Links, "prev")
		assert.NotContains(t, page.Links, "first")
		assert.NotContains(t, page.Links, "next")
		assert.NotContains(t, page.Links, "last")
	})
}
