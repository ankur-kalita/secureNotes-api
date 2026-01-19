# SecureNotes API

A production-grade REST API for managing personal notes, built with Go and featuring a comprehensive CI/CD pipeline demonstrating DevSecOps best practices.

## Project Overview

SecureNotes API is a simple yet functional notes management service that showcases:
- Clean Go architecture with separation of concerns
- Production-ready CI/CD pipeline with GitHub Actions
- DevSecOps integration (SAST, SCA, Container Scanning)
- Kubernetes-ready containerization

## Tech Stack

- **Language**: Go 1.23
- **Framework**: Gin (HTTP router)
- **Container**: Docker (multi-stage build)
- **Orchestration**: Kubernetes
- **CI/CD**: GitHub Actions

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check endpoint |
| GET | `/api/v1/notes` | Get all notes |
| GET | `/api/v1/notes/:id` | Get note by ID |
| POST | `/api/v1/notes` | Create a new note |
| PUT | `/api/v1/notes/:id` | Update a note |
| DELETE | `/api/v1/notes/:id` | Delete a note |

## Project Structure

```
securenotes-api/
├── .github/workflows/
│   ├── ci.yml              # CI Pipeline
│   └── cd.yml              # CD Pipeline
├── cmd/api/
│   └── main.go             # Application entry point
├── internal/
│   ├── handlers/           # HTTP handlers
│   ├── models/             # Data models
│   ├── repository/         # Data layer
│   └── middleware/         # HTTP middleware
├── tests/                  # Unit tests
├── k8s/                    # Kubernetes manifests
├── Dockerfile              # Multi-stage build
├── .golangci.yml           # Linting configuration
└── README.md
```

## Quick Start

### Prerequisites

- Go 1.23+
- Docker
- kubectl (for Kubernetes deployment)
- Minikube (for local Kubernetes)

### Run Locally

```bash
# Clone the repository
git clone https://github.com/YOUR_USERNAME/securenotes-api.git
cd securenotes-api

# Install dependencies
go mod download

# Run the application
go run cmd/api/main.go

# Test the API
curl http://localhost:8080/health
```

### Run Tests

```bash
# Run all tests with coverage
go test -v -cover ./...

# Run linter
golangci-lint run
```

### Run with Docker

```bash
# Build image
docker build -t securenotes-api .

# Run container
docker run -p 8080:8080 securenotes-api

# Test
curl http://localhost:8080/health
```

### Deploy to Kubernetes (Minikube)

```bash
# Start Minikube
minikube start

# Update image in deployment (replace YOUR_USERNAME)
sed -i 's|IMAGE_PLACEHOLDER|YOUR_DOCKERHUB_USERNAME/securenotes-api:latest|g' k8s/deployment.yaml

# Deploy
kubectl apply -f k8s/

# Access the service
minikube service securenotes-api
```

## CI/CD Pipeline

### CI Pipeline Stages

| Stage | Tool | Purpose |
|-------|------|---------|
| **Lint** | golangci-lint | Enforce code quality standards |
| **SAST** | CodeQL | Detect security vulnerabilities in code |
| **SCA** | govulncheck | Scan dependencies for vulnerabilities |
| **Test** | go test | Run unit tests with coverage |
| **Build** | go build | Compile application binary |
| **Docker Build** | Docker | Create container image |
| **Image Scan** | Trivy | Scan container for CVEs |
| **Runtime Test** | curl | Validate container runs correctly |
| **Push** | Docker | Publish to DockerHub |

### CD Pipeline Stages

| Stage | Purpose |
|-------|---------|
| **Deploy** | Apply Kubernetes manifests |
| **Verify** | Post-deployment health checks |
| **DAST** | Optional runtime security scan |

## GitHub Secrets Configuration

Configure the following secrets in your GitHub repository:

| Secret Name | Description | How to Get |
|-------------|-------------|------------|
| `DOCKERHUB_USERNAME` | Your DockerHub username | hub.docker.com |
| `DOCKERHUB_TOKEN` | DockerHub access token | DockerHub → Account Settings → Security → New Access Token |

### Steps to Configure:

1. Go to your GitHub repository
2. Navigate to **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Add each secret with its value

## Security Features

### Shift-Left Security Implementation

1. **Code Level (SAST)**: CodeQL analyzes source code before build
2. **Dependency Level (SCA)**: govulncheck scans Go modules
3. **Container Level**: Trivy scans Docker images
4. **Runtime Level (DAST)**: Optional ZAP scanning in CD

### Security Best Practices

- Multi-stage Docker build (minimal attack surface)
- Non-root container user
- Read-only root filesystem
- Resource limits in Kubernetes
- Health checks for automatic recovery

## Testing the API

```bash
# Health check
curl http://localhost:8080/health

# Create a note
curl -X POST http://localhost:8080/api/v1/notes \
  -H "Content-Type: application/json" \
  -d '{"title":"My Note","content":"Hello World"}'

# Get all notes
curl http://localhost:8080/api/v1/notes

# Get specific note (replace NOTE_ID)
curl http://localhost:8080/api/v1/notes/NOTE_ID

# Update a note
curl -X PUT http://localhost:8080/api/v1/notes/NOTE_ID \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated Title","content":"Updated Content"}'

# Delete a note
curl -X DELETE http://localhost:8080/api/v1/notes/NOTE_ID
```

## License

MIT License
