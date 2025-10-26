# Twitter API

REST API for fetching tweets via Nitter and FxTwitter.

## Quick Start

```bash
# 1. Setup
cp scripts/.env.example .env
# Fill .env with your credentials

# 2. Run
cd scripts
docker-compose up -d

# 3. API available at http://localhost:8080
```

## API Documentation

**Interactive documentation powered by Swagger UI:**
- 🌐 **Swagger UI**: http://localhost:8080/ - Interactive API explorer with "Try it out" functionality
- 📄 **OpenAPI Spec**: http://localhost:8080/openapi.yaml - OpenAPI 3.0.3 specification file

The Swagger UI allows you to:
- Browse all available endpoints and their parameters
- View request/response schemas and examples
- Test API endpoints directly from your browser
- Download the OpenAPI specification

## Endpoints

- `GET /users/{username}/tweets` - user's tweet list
- `GET /users/{username}/tweets/{id}` - specific tweet details

## Stack

- Go 1.21
- Nitter (RSS)
- FxTwitter API
- Docker

## Documentation

- [AGENTS.md](AGENTS.md) - architecture and details
- [docs/openapi/AGENTS.md](docs/openapi/AGENTS.md) - OpenAPI specification maintenance guide
- [scripts/AGENTS.md](scripts/AGENTS.md) - deployment
- [nitter/AGENTS.md](nitter/AGENTS.md) - nitter architecture and details
