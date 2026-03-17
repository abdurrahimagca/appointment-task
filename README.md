# appointment

A small appointment booking app. Providers set up time slots, clients pick one and book it. That's about it.

## What's in here

```
services/
  api/          Go API (huma + pgx, generates db code with sqlc)
  apps/web/     React frontend (Vite + Tailwind, served by nginx)
  db/           Postgres migrations, seed data, and an init script
```

The API has three endpoints:

- `GET /providers/{username}` -- get a provider's profile
- `GET /providers/{username}/availability?date=2026-03-20` -- list free slots for a day
- `POST /appointments` -- book a slot

Bot protection with Cloudflare Turnstile. Confirmation emails via Resend.

The web frontend talks to the API through nginx (`/api` path), so it works on any domain without config changes.

## Running locally

You need Docker and Docker Compose.

```sh
cp .env.example .env
docker compose up --build
```

That's it. Postgres starts, migrations and seed data run automatically, then the API and frontend come up.

- Frontend: http://localhost:3000
- API: http://localhost:8080
- API docs (Scalar): http://localhost:8080/docs

The seed data creates a few test providers with slots so you have something to click on right away.

## Building without the repo

The `api` and `web` images are published to GHCR on every push to `main`:

```
ghcr.io/abdurrahimagca/appointment-api:latest
ghcr.io/abdurrahimagca/appointment-web:latest
```

If you want to run this without cloning the repo, you can pull those images and point them at any Postgres 16 instance. You still need the migration and seed files to initialize the database, but the app containers themselves are self-contained.

## Environment variables

Copy `.env.example` and fill in what you need. The defaults work fine for local dev.

| Variable | What it does |
|---|---|
| `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD` | Database credentials |
| `TURNSTILE_SECRET` / `VITE_TURNSTILE_SITE_KEY` | Cloudflare Turnstile keys (optional locally) |
| `RESEND_API_KEY`, `RESEND_FROM_EMAIL` | Email sending via Resend (optional locally) |
| `ALLOWED_ORIGINS` | CORS origins the API accepts |
| `LOG_LEVEL` | `debug`, `info`, `warn`, `error` |

## Dev tooling

If you're working on the Go code, there's a Makefile for common tasks:

```sh
make tools-install    # install migrate + sqlc CLIs
make migrate-new      # create a new migration
make sqlc-generate    # regenerate db code from queries
```

Queries live in `services/db/queries/`, generated Go code ends up in `services/api/internal/db/`.
