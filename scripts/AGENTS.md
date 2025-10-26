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

## CI/CD and Deployment

### GitHub Actions Workflow

The project uses GitHub Actions for continuous deployment to production servers:

**Workflow file**: `.github/workflows/deploy.yml`

**Trigger**: Automatic deployment on push to `main` branch

**Pipeline steps**:
1. **Build**: Docker image is built using `scripts/Dockerfile`
2. **Push**: Image is pushed to GitHub Container Registry (GHCR) at `ghcr.io/programistich/twitter-api:latest`
3. **Deploy**: SSH connection to production server
4. **Update**: Pull latest image and restart services
5. **Cleanup**: Remove old unused Docker images

### Production Deployment

**docker-compose.prod.yml** - Production override configuration:
- Replaces local build with pre-built GHCR image
- Always pulls the latest image version
- Used in combination with main docker-compose.yml

**Production usage**:
```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### Required GitHub Secrets

Configure these secrets in your GitHub repository settings (Settings → Secrets and variables → Actions):

- `SSH_HOST` - Production server IP address or hostname
- `SSH_USER` - SSH username for deployment
- `SSH_PRIVATE_KEY` - SSH private key for authentication (full content, including BEGIN/END markers)

The deployment expects the project to be located at `/root/twitter-api` on the production server.

### Image Registry

Production images are stored in GitHub Container Registry:
- Registry: `ghcr.io`
- Image: `ghcr.io/programistich/twitter-api:latest`
- Authentication: Automatic via `GITHUB_TOKEN`

To manually pull the production image:
```bash
docker pull ghcr.io/programistich/twitter-api:latest
```

## API Documentation

The Twitter API includes built-in Swagger UI documentation:
- **Swagger UI**: http://localhost:8080/ - Interactive API documentation
- **OpenAPI Spec**: http://localhost:8080/openapi.yaml - Machine-readable specification

When adding new endpoints or modifying existing ones:
1. Update the OpenAPI specification in `docs/openapi/openapi.yaml`
2. Rebuild the Docker container to include updated documentation
3. See `docs/openapi/AGENTS.md` for detailed OpenAPI maintenance instructions
