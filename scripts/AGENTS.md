# AGENTS.md - Scripts Directory

**After all changes, please update any relevant files with the new information.**

## Overview

This directory contains all CI/CD, deployment, and configuration files for the Twitter API project.

## Files

### Environment Configuration
- `.env.example` - Template for environment variables required by the application
  - `TWITTER_USERNAME` - Twitter account username for Nitter session generation
  - `TWITTER_PASSWORD` - Twitter account password
  - `TWITTER_OTP_SECRET` - 2FA secret key for Twitter account
  - `NITTER_URL` - URL of your running Nitter instance (e.g., `http://127.0.0.1:8049`)

### Docker Configuration
- `Dockerfile` - Multi-stage Docker build configuration for the Twitter API
  - Build stage: Uses golang:1.21-alpine to compile the Go application
  - Runtime stage: Uses alpine:latest with minimal dependencies
  - Copies `docs/` directory (includes `docs/openapi/` with OpenAPI specification)
  - Exposes port 8080

- `docker-compose.yml` - Complete orchestration for all services
  - **session-generator** - Python container that generates Twitter sessions using credentials
  - **nitter-redis** - Redis cache for Nitter instance
  - **nitter** - Self-hosted Nitter instance for RSS feed access
  - **twitter-api** - Main API service built from Dockerfile

## Usage

### Running with Docker Compose

From the scripts directory:
```bash
cd scripts
docker-compose up -d
```

Or from the project root:
```bash
docker-compose -f scripts/docker-compose.yml up -d
```

### Environment Setup

1. Copy the example environment file:
```bash
cp scripts/.env.example .env
```

2. Edit `.env` with your credentials:
```bash
TWITTER_USERNAME=your_twitter_username
TWITTER_PASSWORD=your_twitter_password
TWITTER_OTP_SECRET=your_2fa_secret_key
NITTER_URL=http://127.0.0.1:8049
```

### Service Dependencies

The services start in this order:
1. **nitter-redis** - Must be healthy before Nitter starts
2. **session-generator** - Must complete successfully before Nitter starts
3. **nitter** - Must be healthy before twitter-api starts
4. **twitter-api** - Starts last and depends on Nitter being ready

## Path References

All paths in docker-compose.yml are relative to the scripts directory:
- `../nitter/` - References the nitter configuration directory in project root
- `..` - References the project root for build context
- `scripts/Dockerfile` - References this directory's Dockerfile

## Security Notes

- Nitter services run with minimal privileges (no-new-privileges, capability drops)
- Redis and Nitter use read-only filesystems
- Services run as non-root users (998:998, 999:1000)
- Network isolation via shared Docker network

## Port Mappings

- `8080` - Twitter API (exposed publicly)
  - `/` - Swagger UI interactive documentation
  - `/openapi.yaml` - OpenAPI specification file
  - `/users/{username}/tweets` - API endpoint
  - `/users/{username}/tweets/{id}` - API endpoint
- `8049` - Nitter instance (bound to localhost only)
- Redis runs on internal network only

## API Documentation

The Twitter API includes built-in Swagger UI documentation:
- **Swagger UI**: http://localhost:8080/ - Interactive API documentation
- **OpenAPI Spec**: http://localhost:8080/openapi.yaml - Machine-readable specification

When adding new endpoints or modifying existing ones:
1. Update the OpenAPI specification in `docs/openapi/openapi.yaml`
2. Rebuild the Docker container to include updated documentation
3. See `docs/openapi/AGENTS.md` for detailed OpenAPI maintenance instructions
