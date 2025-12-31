This is a REST API written in Go that fetches Twitter user data by combining two services:
- **Nitter** (self-hosted RSS feed service) - for getting tweet IDs from user timelines
- **FxTwitter API** (public API) - for fetching detailed tweet information

API FxTwitter documentation: see @docs/FxTwitterApi.md

Look @infra folder for deployment configs

## Run with Docker

Docker Compose lives in `infra/`.

1) Ensure `infra/.env` exists (you can copy `infra/.env.example` and adjust values).
2) From repo root:
   - Development (local ports): `docker compose -f infra/docker-compose.yml -f infra/docker-compose.override.yml up -d --build`
   - Production (shared proxy network): `docker compose -f infra/docker-compose.yml -f infra/docker-compose.prod.yml up -d --build`
3) Local endpoints (dev override):
   - Nitter: `http://127.0.0.1:8049`
   - API: `http://127.0.0.1:8080`
