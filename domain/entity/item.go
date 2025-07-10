// Package entity contains the domain entities for the application.
package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ItemStatus represents the possible states of an item.
type ItemStatus string

const (
	// ItemStatusActive indicates the item is active and available.
	ItemStatusActive ItemStatus = "active"

	// ItemStatusInactive indicates the item is inactive but can be reactivated.
	ItemStatusInactive ItemStatus = "inactive"

	// ItemStatusDeleted indicates the item has been deleted and cannot be recovered.
	ItemStatusDeleted ItemStatus = "deleted"
)

// Item represents a domain entity in the system.
type Item struct {
	// ID is the unique identifier for the item.
	ID string

	// Name is the display name of the item.
	Name string

	// Description provides details about the item.
	Description string

	// Tags are labels associated with the item for categorization and searching.
	Tags []string

	// Status indicates the current state of the item.
	Status ItemStatus

	// CreatedAt is the timestamp when the item was created.
	CreatedAt time.Time

	// UpdatedAt is the timestamp when the item was last updated.
	UpdatedAt time.Time
}

// NewItem creates a new item with the given name and description.
// It generates a new UUID for the ID and sets the status to active.
func NewItem(name, description string, tags []string) (*Item, error) {
	if name == "" {
		return nil, fmt.Errorf("item name cannot be empty")
	}

	now := time.Now().UTC()
	return &Item{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Tags:        tags,
		Status:      ItemStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Validate checks if the item is valid according to business rules.
func (i *Item) Validate() error {
	if i.ID == "" {
		return fmt.Errorf("item ID cannot be empty")
	}

	if i.Name == "" {
		return fmt.Errorf("item name cannot be empty")
	}

	if i.Status == "" {
		return fmt.Errorf("item status cannot be empty")
	}

	return nil
}

// Update updates the item's fields with the provided values.
// Only non-empty values are used for the update.
func (i *Item) Update(name, description string, tags []string, status ItemStatus) {
	if name != "" {
		i.Name = name
	}

	if description != "" {
		i.Description = description
	}

	if len(tags) > 0 {
		i.Tags = tags
	}

	if status != "" {
		i.Status = status
	}

	i.UpdatedAt = time.Now().UTC()
}

// Activate sets the item's status to active.
func (i *Item) Activate() {
	i.Status = ItemStatusActive
	i.UpdatedAt = time.Now().UTC()
}

// Deactivate sets the item's status to inactive.
func (i *Item) Deactivate() {
	i.Status = ItemStatusInactive
	i.UpdatedAt = time.Now().UTC()
}

// Delete sets the item's status to deleted.
func (i *Item) Delete() {
	i.Status = ItemStatusDeleted
	i.UpdatedAt = time.Now().UTC()
}

// IsActive returns true if the item is active.
func (i *Item) IsActive() bool {
	return i.Status == ItemStatusActive
}

// HasTag returns true if the item has the specified tag.
func (i *Item) HasTag(tag string) bool {
	for _, t := range i.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// AddTag adds a tag to the item if it doesn't already exist.
func (i *Item) AddTag(tag string) {
	if tag == "" || i.HasTag(tag) {
		return
	}
	i.Tags = append(i.Tags, tag)
	i.UpdatedAt = time.Now().UTC()
}

// RemoveTag removes a tag from the item if it exists.
func (i *Item) RemoveTag(tag string) {
	if tag == "" {
		return
	}

	for j, t := range i.Tags {
		if t == tag {
			i.Tags = append(i.Tags[:j], i.Tags[j+1:]...)
			i.UpdatedAt = time.Now().UTC()
			return
		}
	}
}
