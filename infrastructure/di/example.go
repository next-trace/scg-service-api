//go:build examples

// Package di provides dependency injection utilities.
package di

import (
	"context"
	"fmt"

	appconfig "github.com/next-trace/scg-service-api/application/config"
	appgrpc "github.com/next-trace/scg-service-api/application/grpc"
	applogger "github.com/next-trace/scg-service-api/application/logger"
	"github.com/next-trace/scg-service-api/domain/entity"
	"github.com/next-trace/scg-service-api/domain/repository"
	"github.com/next-trace/scg-service-api/domain/service"
	"github.com/next-trace/scg-service-api/infrastructure/config"
	"github.com/next-trace/scg-service-api/infrastructure/grpc"
	"github.com/next-trace/scg-service-api/infrastructure/logger"
)

// This file provides an example of how to use the dependency injection container.
// It shows how to register constructors, resolve dependencies, and invoke functions.

// Example demonstrates how to use the dependency injection container.
func Example() {
	// Create a new container
	container := NewContainer()

	// Register constructors for infrastructure components
	registerInfrastructure(container)

	// Register constructors for domain components
	registerDomain(container)

	// Register constructors for application components
	registerApplication(container)

	// Invoke a function that uses the registered components
	err := container.Invoke(func(itemService *service.ItemService, logger applogger.Logger) {
		ctx := context.Background()
		logger.Info(ctx, "Starting application...")

		// Use the item service
		item, err := itemService.CreateItem(ctx, "Example Item", "This is an example item", []string{"example", "demo"})
		if err != nil {
			logger.Error(ctx, err, "Failed to create item")
			return
		}

		logger.InfoKV(ctx, "Item created", map[string]interface{}{
			"item_id":   item.ID,
			"item_name": item.Name,
		})
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// handleProvideError is a helper function to handle errors from container.Provide
func handleProvideError(err error, componentName string) {
	if err != nil {
		fmt.Printf("Error registering %s: %v\n", componentName, err)
	}
}

// registerInfrastructure registers constructors for infrastructure components.
func registerInfrastructure(container *Container) {
	// Register config loader
	err := container.Provide(func() appconfig.Loader {
		return config.NewViperLoader()
	})
	handleProvideError(err, "config loader")

	// Register logger
	err = container.Provide(func() applogger.Logger {
		return logger.NewSlogAdapter(nil, "info")
	})
	handleProvideError(err, "logger")

	// Register gRPC client
	err = container.Provide(func(logger applogger.Logger) appgrpc.Client {
		return grpc.NewClientAdapter(appgrpc.DefaultClientConfig(), logger)
	})
	handleProvideError(err, "gRPC client")

	// Register gRPC server
	err = container.Provide(func(logger applogger.Logger) appgrpc.Server {
		return grpc.NewServerAdapter(appgrpc.DefaultServerConfig(), logger)
	})
	handleProvideError(err, "gRPC server")

	// Register mock item repository
	err = container.Provide(func() repository.ItemRepository {
		return newMockItemRepository()
	})
	handleProvideError(err, "mock item repository")
}

// registerDomain registers constructors for domain components.
func registerDomain(container *Container) {
	// Register item service
	err := container.Provide(func(repo repository.ItemRepository) *service.ItemService {
		return service.NewItemService(repo)
	})
	handleProvideError(err, "item service")
}

// registerApplication registers constructors for application components.
func registerApplication(container *Container) {
	// Register application components here
}

// mockItemRepository is a mock implementation of the ItemRepository interface.
type mockItemRepository struct {
	items map[string]*entity.Item
}

// newMockItemRepository creates a new mock item repository.
func newMockItemRepository() repository.ItemRepository {
	return &mockItemRepository{
		items: make(map[string]*entity.Item),
	}
}

// GetByID retrieves an item by its ID.
func (r *mockItemRepository) GetByID(ctx context.Context, id string) (*entity.Item, error) {
	item, ok := r.items[id]
	if !ok {
		return nil, fmt.Errorf("item not found")
	}
	return item, nil
}

// FindAll retrieves all items with optional filtering.
func (r *mockItemRepository) FindAll(ctx context.Context, filter repository.ItemFilter) ([]*entity.Item, error) {
	var items []*entity.Item
	for _, item := range r.items {
		items = append(items, item)
	}
	return items, nil
}

// Count returns the number of items matching the filter.
func (r *mockItemRepository) Count(ctx context.Context, filter repository.ItemFilter) (int64, error) {
	return int64(len(r.items)), nil
}

// Save persists an item to the repository.
func (r *mockItemRepository) Save(ctx context.Context, item *entity.Item) error {
	r.items[item.ID] = item
	return nil
}

// Delete removes an item from the repository.
func (r *mockItemRepository) Delete(ctx context.Context, id string) error {
	delete(r.items, id)
	return nil
}
