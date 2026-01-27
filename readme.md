# Go Upload - File Upload Service

A clean architecture Go application for file uploads with image processing capabilities, built with Gin and GORM.

## 🏗️ Architecture

This project follows **Clean Architecture** principles with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                        HTTP Layer                            │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │  Auth    │  │ Upload   │  │  File    │  │  Image   │   │
│  │ Handler  │  │ Handler  │  │ Handler  │  │ Handler  │   │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘   │
└───────┼─────────────┼─────────────┼─────────────┼──────────┘
        │             │             │             │
        │ depends on  │             │             │
        ↓             ↓             ↓             ↓
┌─────────────────────────────────────────────────────────────┐
│                    Service Interfaces                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │   Auth   │  │  Upload  │  │   File   │  │  Image   │   │
│  │ Service  │  │ Service  │  │ Service  │  │ Service  │   │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘   │
└───────┼─────────────┼─────────────┼─────────────┼──────────┘
        │ implemented │             │             │
        ↓     by      ↓             ↓             ↓
┌─────────────────────────────────────────────────────────────┐
│                   Business Logic Layer                       │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │   auth   │  │  upload  │  │   file   │  │  image   │   │
│  │ Service  │  │ Service  │  │ Service  │  │ Service  │   │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘   │
└───────┼─────────────┼─────────────┼─────────────┼──────────┘
        │ depends on  │             │             │
        ↓             ↓             ↓             ↓
┌─────────────────────────────────────────────────────────────┐
│                  Repository Interfaces                       │
│         ┌──────────────────┐  ┌──────────────────┐         │
│         │      User        │  │     Upload       │         │
│         │   Repository     │  │   Repository     │         │
│         └────────┬─────────┘  └────────┬─────────┘         │
└──────────────────┼─────────────────────┼───────────────────┘
                   │ implemented by      │
                   ↓                     ↓
┌─────────────────────────────────────────────────────────────┐
│                   Data Access Layer                          │
│         ┌──────────────────┐  ┌──────────────────┐         │
│         │   PostgreSQL     │  │   PostgreSQL     │         │
│         │     User Repo    │  │   Upload Repo    │         │
│         └────────┬─────────┘  └────────┬─────────┘         │
└──────────────────┼─────────────────────┼───────────────────┘
                   │                     │
                   ↓                     ↓
              ┌─────────────────────────────┐
              │    PostgreSQL Database      │
              │    (GORM + UUID PKs)        │
              └─────────────────────────────┘
```

### Layer Responsibilities

- **HTTP Layer (Handlers)**: Handle HTTP requests/responses, parse inputs, format outputs
- **Service Interfaces**: Define contracts between layers (Dependency Inversion)
- **Business Logic (Services)**: Implement business rules, orchestrate operations
- **Repository Interfaces**: Define data access contracts
- **Data Access (Repositories)**: Implement database operations with GORM

## 📁 Project Structure

```
go-upload/
├── cmd/
│   └── api/
│       └── main.go                    # Application entry point with DI
├── internal/
│   ├── domain/
│   │   ├── entity/                    # Domain entities (clean, no tags)
│   │   │   ├── user.go
│   │   │   └── upload.go
│   │   └── errors/
│   │       └── errors.go              # Custom application errors
│   ├── dto/                           # Data Transfer Objects
│   │   ├── auth_dto.go
│   │   ├── upload_dto.go
│   │   ├── file_dto.go
│   │   └── image_dto.go
│   ├── handler/                       # HTTP handlers (thin layer)
│   │   ├── auth_handler.go
│   │   ├── upload_handler.go
│   │   ├── file_handler.go
│   │   └── image_handler.go
│   ├── service/                       # Business logic
│   │   ├── interfaces.go              # Service interfaces
│   │   ├── auth_service.go
│   │   ├── upload_service.go
│   │   ├── file_service.go
│   │   ├── image_service.go
│   │   └── storage_service.go
│   ├── repository/                    # Data access layer
│   │   ├── user_repository.go         # Repository interface
│   │   ├── upload_repository.go       # Repository interface
│   │   └── postgres/                  # GORM implementations
│   │       ├── models.go              # GORM models with tags
│   │       ├── user_repository.go
│   │       └── upload_repository.go
│   ├── middleware/
│   │   └── auth_middleware.go         # JWT authentication
│   └── router/
│       └── router.go                  # Route configuration
├── pkg/                               # Shared utilities
│   ├── jwt/
│   │   └── jwt.go                     # JWT service
│   └── hash/
│       └── password.go                # Password hashing
├── config/
│   └── config.go                      # Configuration management
├── tests/                             # Test files
│   ├── repository/
│   │   └── user_repository_test.go
│   ├── service/
│   │   └── auth_service_test.go
│   └── integration/
│       └── auth_integration_test.go
├── file-storage/                      # Uploaded files directory
├── .env                               # Environment variables
├── go.mod
├── go.sum
├── dockerfile
└── compose.yaml
```

## 🚀 Features

- ✅ **Clean Architecture** with Handler → Service → Repository pattern
- ✅ **Dependency Injection** - No global variables
- ✅ **Interface-based Design** - Easy to mock and test
- ✅ **JWT Authentication** - Secure token-based auth
- ✅ **File Upload** - With folder organization and conflict resolution
- ✅ **Image Processing** - Resize, quality adjustment, format conversion (WEBP, JPEG, PNG, AVIF)
- ✅ **UUID Primary Keys** - Using PostgreSQL's `gen_random_uuid()`
- ✅ **Soft Deletes** - Audit trail with GORM
- ✅ **Restrictive DTOs** - No password exposure, only essential fields
- ✅ **Comprehensive Tests** - Repository, service, and integration tests

## 📡 API Endpoints

### Authentication
```
POST   /api/auth/signup     - Register new user
POST   /api/auth/signin     - Login and get JWT token
POST   /api/auth/signout    - Logout (requires auth)
GET    /api/auth/user       - Get current user info (requires auth)
```

### File Management
```
POST   /api/upload/         - Upload file (requires auth)
GET    /api/upload/         - List user's uploads (requires auth)
GET    /api/upload/:id      - Get upload details (requires auth)
DELETE /api/upload/:id      - Delete upload (requires auth)
```

### File Serving
```
GET    /file/:id            - Download file
GET    /image/:id           - Get optimized image
       Query params:
       - w=<width>          - Resize to width (pixels)
       - q=<quality>        - Quality 1-100 (default: 80)
       - format=<format>    - webp|jpeg|png|avif
```

## 🛠️ Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 14+
- libvips (for image processing)

### Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd go-upload
```

2. **Install dependencies**
```bash
go mod download
```

3. **Install libvips** (for image processing)
```bash
# macOS
brew install vips

# Ubuntu/Debian
sudo apt-get install libvips-dev

# Alpine (Docker)
apk add vips-dev
```

4. **Set up environment variables**
```bash
cp .env.example .env
# Edit .env with your configuration
```

Example `.env`:
```env
PORT=8888
DATABASE_URL="postgres://user:password@localhost:5432/go-upload"
JWT_SECRET="your_secure_jwt_secret_key_change_in_production"
STORAGE_PATH="file-storage"
```

5. **Run database migrations**
```bash
# Migrations run automatically on startup
# Tables: users, uploads
```

### Running the Application

**Development:**
```bash
go run cmd/api/main.go
```

**Production:**
```bash
# Build
go build -o bin/api ./cmd/api

# Run
./bin/api
```

### Docker Deployment

1. **Start database services**
```bash
docker compose -f compose.yaml up -d
```

2. **Start the Go application**
```bash
docker compose -f compose.yaml -f compose.go.yaml up -d
```

## 🧪 Testing

### Run All Tests
```bash
go test ./tests/... -v
```

### Run Specific Test Suites
```bash
# Repository tests
go test ./tests/repository/... -v

# Service tests (with mocks)
go test ./tests/service/... -v

# Integration tests
go test ./tests/integration/... -v
```

### Test Coverage
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 🔒 Security Features

1. **No Password Exposure**: Passwords are hashed with bcrypt, never returned in responses
2. **JWT Authentication**: Secure token-based authentication
3. **Restrictive DTOs**: Only necessary fields exposed in API responses
4. **Environment-based Secrets**: JWT secret from environment variables
5. **Authorization Checks**: User isolation - users can only access their own uploads
6. **File Path Sanitization**: Prevents directory traversal attacks

## 🏛️ Design Principles

### SOLID Principles
- **Single Responsibility**: Each layer has one reason to change
- **Open/Closed**: Open for extension, closed for modification
- **Liskov Substitution**: Interfaces allow swapping implementations
- **Interface Segregation**: Small, focused interfaces
- **Dependency Inversion**: Depend on abstractions, not concretions

### Clean Architecture Benefits
- **Testability**: Easy to test with mocked dependencies
- **Maintainability**: Clear separation of concerns
- **Flexibility**: Easy to swap implementations (e.g., add S3 storage)
- **Scalability**: Independent scaling of layers
- **Team Collaboration**: Clear boundaries for parallel development

## 📦 Dependencies

```go
require (
    github.com/gin-gonic/gin v1.11.0           // Web framework
    github.com/golang-jwt/jwt/v5 v5.3.0        // JWT
    github.com/google/uuid v1.6.0              // UUID generation
    github.com/h2non/bimg v1.1.9               // Image processing
    github.com/joho/godotenv v1.5.1            // Environment variables
    github.com/rs/xid v1.6.0                   // Unique IDs
    golang.org/x/crypto v0.40.0                // Bcrypt
    gorm.io/driver/postgres v1.6.0             // PostgreSQL driver
    gorm.io/gorm v1.31.1                       // ORM
)
```

## 🔄 Migration from Old Architecture

This project was refactored from a traditional MVC pattern to Clean Architecture:

### What Changed
- ❌ **Removed**: Global `utils.DB`, monolithic controllers
- ✅ **Added**: Service interfaces, dependency injection, DTOs
- ✅ **Improved**: Testability, maintainability, SOLID compliance

### Breaking Changes
- None! All endpoints remain backward compatible

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License.

## 👥 Authors

- Original Project: [Your Name]
- Clean Architecture Refactoring: Claude Code

## 🙏 Acknowledgments

- Clean Architecture by Robert C. Martin
- Gin Web Framework
- GORM ORM
- libvips for image processing
