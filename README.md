# TwitterX API

A Go REST API for fetching Twitter data without the official API.

## About

TwitterX API combines two services to fetch Twitter data:

- **Nitter** — self-hosted RSS service for getting tweet IDs from user timelines
- **FxTwitter API** — public API for fetching detailed tweet information

### Why?

The official Twitter API became paid and has strict limitations. This project provides a free alternative for:
- Getting a user's tweet list
- Fetching detailed tweet information (text, media, statistics)
- Getting user profile information

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/users/{username}` | User profile information |
| GET | `/api/users/{username}/tweets` | List of user's tweet IDs |
| GET | `/api/users/{username}/tweets/{id}` | Detailed tweet information |
| GET | `/api/docs/` | Swagger UI documentation |

## Quick Start

### Requirements

- Docker and Docker Compose

### Running

```bash
cd infra
docker compose up -d
```

Services will be available at:
- API: `http://localhost:8080`
- Swagger UI: `http://localhost:8080/api/docs/`
- Nitter: `http://localhost:8049`

### Configuration

Environment variables in `infra/.env`:

| Variable | Description | Default |
|----------|-------------|---------|
| `NITTER_URL` | Nitter instance URL | `http://nitter:8049` |
| `DEBUG` | Enable debug logs | `False` |
| `NITTER_IMAGE` | Nitter Docker image | `zedeus/nitter:latest` |

## Development

```bash
# Install dependencies
go mod download

# Generate Swagger documentation
go install github.com/swaggo/swag/cmd/swag@latest
swag init -o docs/api

# Run (requires running Nitter instance)
NITTER_URL=http://localhost:8049 go run main.go
```

## Production Deployment

```bash
cd infra
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
```
