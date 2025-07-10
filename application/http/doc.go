// Package http provides HTTP abstractions and helpers (ports) for decoding requests,
// encoding responses, and running HTTP servers with graceful shutdown.
//
// Design
//   - Code-to-interface: Microservices depend on these small interfaces (RequestDecoder, ResponseWriter)
//     instead of concrete implementations.
//   - Transport-agnostic contracts live in the application layer; concrete adapters live under infrastructure.
//
// Highlights
//   - RequestDecoder abstracts deserialization concerns.
//   - ResponseWriter standardizes success and error payloads.
//   - Run helper starts an http.Server and performs graceful shutdown upon context cancel or SIGINT/SIGTERM.
//
// Quickstart
//
//	logger := infralog.NewSlogAdapter(os.Stdout, "info")
//	mux := http.NewServeMux()
//	srv := &nethttp.Server{Addr: ":8080", Handler: mux}
//	ctx := context.Background()
//	if err := apphttp.Run(ctx, srv, logger); err != nil { /* handle */ }
//
// See infrastructure/serializer for a JSON adapter implementing both ports.
package http
