# Inspacemap Backend

A venue mapping and navigation API built with Go, Fiber, and PostgreSQL.

## 🚀 Quick Start

### Development

```bash
# Start development environment
docker-compose up -d

# View logs
docker-compose logs -f
```

### Testing

```bash
# Run all tests with isolated database
docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from tester
```

## 🏗️ Architecture

- **Framework**: Go + Fiber
- **Database**: PostgreSQL
- **Storage**: MinIO (S3-compatible)
- **Authentication**: JWT
- **Testing**: Comprehensive HTTP integration tests

## 📁 Project Structure

```
├── cmd/                    # Application entrypoints
├── config/                 # Configuration management
├── internal/               # Private application code
│   ├── delivery/http/      # HTTP handlers and routes
│   ├── entity/            # Database entities
│   ├── models/            # API models and DTOs
│   ├── repository/        # Data access layer
│   └── service/           # Business logic layer
├── pkg/                   # Public packages
├── test/                  # Test suites
│   ├── integration/       # Integration tests
│   └── unit/             # Unit tests
└── docker-compose.yml     # Development environment
```

## 🧪 Testing Strategy

### Unit Tests

```bash
go test ./test/unit/... -v
```

### Integration Tests

```bash
go test ./test/integration/... -v
```

### HTTP Integration Tests

```bash
docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from tester
```

## 🔄 CI/CD

This project uses GitHub Actions for continuous integration:

### Workflows

- **CI Pipeline**: Runs on every push and PR
  - Unit tests with race detection
  - Integration tests with real database
  - Docker image build
  - Code linting with golangci-lint
  - Coverage reports

### Coverage Reports

Coverage reports are automatically generated and uploaded as artifacts on each CI run.

## 🛠️ Development Tools

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Git

### Code Quality

```bash
# Run linter
golangci-lint run

# Format code
go fmt ./...

# Run tests with coverage
go test -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📚 API Documentation

API documentation will be available at `/swagger/index.html` when the server is running.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from tester`
5. Create a pull request

## 📄 License

This project is licensed under the MIT License.
