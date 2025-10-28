# AGENTS.md

**After all changes, please update any relevant files with the new information.**

## Project Overview

This is a REST API written in Go that fetches Twitter user data by combining two services:
- **Nitter** (self-hosted RSS feed service) - for getting tweet IDs from user timelines
- **FxTwitter API** (public API) - for fetching detailed tweet information

## API Documentation

The API is documented using **OpenAPI 3.0.3** specification:
- **Location**: `docs/openapi/openapi.yaml` - The source of truth for API documentation
- **Swagger UI**: Integrated at `/` endpoint for interactive documentation
- **Spec endpoint**: `/openapi.yaml` serves the specification file

### When Adding New Features
**IMPORTANT**: When adding new endpoints or modifying existing ones:
1. Update `docs/openapi/openapi.yaml` with new endpoints, schemas, or changes
2. Update code in `main.go` and related files
3. Test that Swagger UI reflects the changes correctly
4. See `docs/openapi/AGENTS.md` for detailed OpenAPI maintenance instructions

## Environment Configuration

Required environment variables (create `.env` from `scripts/.env.example`):
- `NITTER_URL` - URL of your running Nitter instance
  - Docker: `http://nitter:8049` (use container name)
  - Local: `http://127.0.0.1:8049` (use localhost)
- `TWITTER_USERNAME` - Twitter account username for Nitter session generation
- `TWITTER_PASSWORD` - Twitter account password
- `TWITTER_OTP_SECRET` - 2FA secret key for Twitter account

For deployment and CI/CD configuration, see `scripts/AGENTS.md`

## Architecture

### Request Flow

**User Profile Flow:**
1. Client requests user profile via `/users/{username}`
2. API queries FxTwitter API with the username
3. User profile data is returned to client

**Tweets Flow:**
1. Client requests tweets from a user via `/users/{username}/tweets`
2. API fetches the user's RSS feed from Nitter instance
3. RSS parser extracts tweet IDs from the feed
4. Client requests specific tweet details via `/users/{username}/tweets/{id}`
5. API queries FxTwitter API with the tweet ID
6. Tweet data is returned to client

### Module Structure
- `go.mod` - Dependencies (includes github.com/swaggo/http-swagger for Swagger UI)
- `main.go` - HTTP server setup with Gorilla Mux router, handler factories, and Swagger UI integration
- `docs/openapi/` - OpenAPI documentation directory
  - `openapi.yaml` - OpenAPI 3.0.3 specification (manually maintained)
  - `AGENTS.md` - OpenAPI maintenance guide
- `docs/fx/API_FxTwitter.md` - External FxTwitter API reference
- `internal/service/nitter.go` - NitterService fetches and processes RSS feeds from Nitter
- `internal/service/fxtwitter.go` - FxTwitterService fetches tweet data from FxTwitter API
- `internal/parser/rss.go` - RSS XML parser that extracts tweet IDs from Nitter feed GUIDs
- `internal/models/fxtwitter.go` - Complete data models for FxTwitter API responses (tweets and users) with custom time parsing

### Key Implementation Details
- RSS GUID format: `http://127.0.0.1:8049/{username}/status/{tweetID}#m`
- Tweet IDs extracted using regex: `/status/(\d+)`
- FxTwitter API endpoints:
  - User profile: `https://api.fxtwitter.com/{username}`
  - Tweet details: `https://api.fxtwitter.com/{username}/status/{tweetID}`
- HTTP clients have 10-15 second timeouts
- Custom TwitterTime type handles Twitter's RFC1123 time format

## FxTwitter API
API FxTwitter documentation: see [docs/fx/API_FxTwitter.md](docs/fx/API_FxTwitter.md)

## Deployment

### Docker Deployment
All Docker and deployment files are located in `scripts/` directory:
- `scripts/Dockerfile` - Docker build configuration
- `scripts/docker-compose.yml` - Complete service orchestration
- `scripts/.env.example` - Environment variables template

See `scripts/AGENTS.md` for detailed deployment instructions.

## Nitter NitterService
Nitter Configuration in nitter/ folder
