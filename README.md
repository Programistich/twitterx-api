# Twitter API

REST API for fetching tweets via Nitter and FxTwitter.

## Quick Start

```bash
# 1. Setup
cp scripts/.env.example .env
# Fill .env with your credentials

# 2. Create dump file
touch nitter/sessions.jsonl

# 3. Run
cd scripts
docker compose -f scripts/docker-compose.yml --env-file scripts/.env up -d

# Local run
docker network rm caddy && docker network create --driver bridge caddynet
docker compose -f scripts/docker-compose.yml --env-file scripts/.env up -d --scale twitter-api=0
set -a; source scripts/.env; set +a
go run main.go

# 4. API available at http://localhost:8080
```

## API Documentation

**Interactive documentation powered by Swagger UI:**
- 🌐 **Swagger UI**: http://localhost:8080/api/docs/ - Interactive API explorer with "Try it out" functionality
- 📄 **OpenAPI Spec**: http://localhost:8080/api/openapi.yaml - OpenAPI 3.0.3 specification file

The Swagger UI allows you to:
- Browse all available endpoints and their parameters
- View request/response schemas and examples
- Test API endpoints directly from your browser
- Download the OpenAPI specification

## Endpoints

- `GET /api/users/{username}` - user profile
- `GET /api/users/{username}/tweets` - user's tweet list
- `GET /api/users/{username}/tweets/{id}` - specific tweet details

## Stack

- Go 1.21
- Nitter (RSS)
- FxTwitter API
- Docker
- Caddy (reverse proxy in production)

## Documentation

- [AGENTS.md](AGENTS.md) - architecture and details
- [docs/openapi/AGENTS.md](docs/openapi/AGENTS.md) - OpenAPI specification maintenance guide
- [scripts/AGENTS.md](scripts/AGENTS.md) - deployment
- [nitter/AGENTS.md](nitter/AGENTS.md) - nitter architecture and details
