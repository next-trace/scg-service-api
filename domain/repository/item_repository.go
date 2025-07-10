// Package repository defines interfaces for data access.
package repository

import (
	"context"

	"github.com/next-trace/scg-service-api/domain/entity"
)

// ItemRepository defines the interface for item data access.
type ItemRepository interface {
	// GetByID retrieves an item by its ID.
	GetByID(ctx context.Context, id string) (*entity.Item, error)

	// FindAll retrieves all items with optional filtering.
	FindAll(ctx context.Context, filter ItemFilter) ([]*entity.Item, error)

	// Count returns the number of items matching the filter.
	Count(ctx context.Context, filter ItemFilter) (int64, error)

	// Save persists an item to the repository.
	Save(ctx context.Context, item *entity.Item) error

	// Delete removes an item from the repository.
	Delete(ctx context.Context, id string) error
}

// ItemFilter defines criteria for filtering items.
type ItemFilter struct {
	// Status filters items by their status.
	Status entity.ItemStatus

	// Tags filters items that have all the specified tags.
	Tags []string

	// SearchTerm searches in item name and description.
	SearchTerm string

	// Pagination parameters
	Offset int
	Limit  int
}

// NewItemFilter creates a new filter with default values.
func NewItemFilter() ItemFilter {
	return ItemFilter{
		Limit: 50, // Default limit
	}
}

// WithStatus adds status filtering to the filter.
func (f ItemFilter) WithStatus(status entity.ItemStatus) ItemFilter {
	f.Status = status
	return f
}

// WithTags adds tag filtering to the filter.
func (f ItemFilter) WithTags(tags []string) ItemFilter {
	f.Tags = tags
	return f
}

// WithSearch adds text search to the filter.
func (f ItemFilter) WithSearch(term string) ItemFilter {
	f.SearchTerm = term
	return f
}

// WithPagination adds pagination to the filter.
func (f ItemFilter) WithPagination(offset, limit int) ItemFilter {
	f.Offset = offset
	if limit > 0 {
		f.Limit = limit
	}
	return f
}
