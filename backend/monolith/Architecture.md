internal/
├── geolocation/                    # Bounded context: geolocation
│   ├── domain/                      # Enterprise business rules (pure Go, no external dependencies)
│   │   ├── location.go               # Location entity / value object
│   │   ├── track.go                  # Track aggregate (e.g., list of points)
│   │   ├── repository.go             # Repository interfaces (e.g., LocationRepository)
│   │   └── service.go                # Domain services (if any, e.g., DistanceCalculator)
│   ├── application/                  # Use cases / application services
│   │   ├── command/                   # Command DTOs and handlers
│   │   │   ├── update_location.go      # UpdateDriverLocation command + handler
│   │   │   └── get_history.go          # GetLocationHistory command + handler
│   │   ├── query/                      # Query DTOs and handlers
│   │   │   ├── current_location.go     # GetCurrentLocation query + handler
│   │   │   └── track_history.go        # GetTrackHistory query + handler
│   │   ├── dto.go                      # Shared data transfer objects (optional)
│   │   └── service.go                  # Orchestration of commands/queries (optional facade)
│   ├── infrastructure/                # Frameworks, drivers, external concerns
│   │   ├── repositories/                # Implementations of domain repository interfaces
│   │   │   ├── postgres_location_repo.go  # LocationRepository using PostgreSQL
│   │   │   └── redis_track_cache.go       # Optional cache for tracks
│   │   ├── clients/                      # Clients to external services
│   │   │   └── graphhopper_client.go      # Calls GraphHopper map matching API
│   │   └── pubsub/                        # Event publishers (if using domain events)
│   │       └── event_publisher.go
│   ├── interfaces/                    # Adapters to the outside world (presentation layer)
│   │   ├── grpc/                        # gRPC service implementation
│   │   │   ├── server.go                 # Implements the generated gRPC interface
│   │   │   └── dto_converter.go          # Converts between protobuf and domain/application DTOs
│   │   ├── http/                         # HTTP handlers (if REST exposed)
│   │   │   ├── handler.go                 # HTTP handlers for location endpoints
│   │   │   └── middleware.go              # Geolocation‑specific middleware
│   │   └── amqp/                         # Message consumer (e.g., RabbitMQ) for incoming location events
│   │       └── consumer.go
│   └── test/                           # Tests for this context (optional, can be alongside code)
│       ├── domain_test/
│       ├── application_test/
│       └── infrastructure_test/
