// Package grpc provides gRPC server and client implementations.
//
// Note: This package requires the following dependencies:
// - google.golang.org/grpc
// - google.golang.org/grpc/health
// - google.golang.org/grpc/health/grpc_health_v1
// - google.golang.org/grpc/reflection
//
// See docs/dependencies.md for more information.
package grpc

import (
	"context"
	"fmt"
	"net"

	appgrpc "github.com/hbttundar/scg-service-base/application/grpc"
	applogger "github.com/hbttundar/scg-service-base/application/logger"
)

// Type aliases for gRPC types
// These are placeholders for the actual types from the gRPC packages
// They will be replaced with the actual types when the dependencies are added
type (
	// grpcServer represents the gRPC server
	grpcServer struct{}

	// healthServer represents the gRPC health server
	healthServer struct{}

	// serverOption represents a gRPC server option
	serverOption struct{}

	// serviceRegistrar represents a gRPC service registrar
	serviceRegistrar interface{}

	// healthCheckResponse represents a health check response
	healthCheckResponse int
)

// Constants for health check response
const (
	// HealthCheckResponseServing indicates the service is serving
	HealthCheckResponseServing healthCheckResponse = 1

	// HealthCheckResponseNotServing indicates the service is not serving
	HealthCheckResponseNotServing healthCheckResponse = 2
)

// Mock methods for the grpcServer type
func (s *grpcServer) Serve(listener net.Listener) error {
	// In a real implementation, this would start the gRPC server
	return nil
}

func (s *grpcServer) GracefulStop() {
	// In a real implementation, this would gracefully stop the gRPC server
}

// Mock methods for the healthServer type
func (s *healthServer) SetServingStatus(service string, status healthCheckResponse) {
	// In a real implementation, this would set the serving status of a service
}

// Mock namespace for health check responses
var healthpb = struct {
	HealthCheckResponse_SERVING    healthCheckResponse
	HealthCheckResponse_NOT_SERVING healthCheckResponse
}{
	HealthCheckResponse_SERVING:    HealthCheckResponseServing,
	HealthCheckResponse_NOT_SERVING: HealthCheckResponseNotServing,
}

// serverAdapter implements the appgrpc.Server interface using the gRPC library.
type serverAdapter struct {
	server     *grpcServer
	config     appgrpc.ServerConfig
	log        applogger.Logger
	healthSvc  *healthServer
	registered bool
}

// NewServerAdapter creates a new gRPC server adapter.
func NewServerAdapter(config appgrpc.ServerConfig, log applogger.Logger) appgrpc.Server {
 // Create server options
	// In a real implementation, these would be:
	// opts := []grpc.ServerOption{
	//     grpc.MaxConcurrentStreams(config.MaxConcurrentStreams),
	//     grpc.MaxRecvMsgSize(config.MaxRecvMsgSize),
	//     grpc.MaxSendMsgSize(config.MaxSendMsgSize),
	// }
	// We're not using options in this mock implementation

	// Create gRPC server with options
	// In a real implementation, this would be:
	// server := grpc.NewServer(opts...)
	server := &grpcServer{}

	// Create health service if enabled
	// In a real implementation, this would be:
	// var healthSvc *health.Server
	// if config.EnableHealthCheck {
	//     healthSvc = health.NewServer()
	//     healthpb.RegisterHealthServer(server, healthSvc)
	// }
	var healthSvc *healthServer
	if config.EnableHealthCheck {
		healthSvc = &healthServer{}
	}

	// Enable reflection if configured
	// In a real implementation, this would be:
	// if config.EnableReflection {
	//     reflection.Register(server)
	// }
	// No-op for now

	return &serverAdapter{
		server:     server,
		config:     config,
		log:        log,
		healthSvc:  healthSvc,
		registered: false,
	}
}

// Start starts the gRPC server on the given listener.
func (s *serverAdapter) Start(ctx context.Context, listener net.Listener) error {
	if !s.registered {
		s.log.Warn(ctx, "starting gRPC server with no registered services")
	}

	// Set all services to SERVING status if health check is enabled
	if s.healthSvc != nil {
		s.healthSvc.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	}

	// Log server start
	s.log.InfoKV(ctx, "starting gRPC server", map[string]interface{}{
		"address": listener.Addr().String(),
	})

	// Start server (this is blocking)
	return s.server.Serve(listener)
}

// Stop gracefully stops the gRPC server.
func (s *serverAdapter) Stop(ctx context.Context) error {
	// Set all services to NOT_SERVING status if health check is enabled
	if s.healthSvc != nil {
		s.healthSvc.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	}

	// Log server stop
	s.log.Info(ctx, "stopping gRPC server")

	// Gracefully stop the server
	s.server.GracefulStop()
	return nil
}

// RegisterService registers a gRPC service with the server.
func (s *serverAdapter) RegisterService(service interface{}) error {
	// The service should implement a Register method that takes a serviceRegistrar
	// In a real implementation, this would be:
	// if registrar, ok := service.(interface {
	//     Register(grpc.ServiceRegistrar)
	// }); ok {
	//     registrar.Register(s.server)
	//     s.registered = true
	//     return nil
	// }

	// For this mock implementation, we'll just check if the service has a Register method
	if registrar, ok := service.(interface {
		Register(interface{})
	}); ok {
		registrar.Register(s.server)
		s.registered = true
		return nil
	}

	return fmt.Errorf("service does not implement Register method")
}
