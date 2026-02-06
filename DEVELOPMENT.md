# Development

## Prerequisites
- Go 1.25
- yt-dlp
- ffmpeg
- sqlite3
- Telegram bot token and admin user ID

## Setup
1. Copy `.env.example` to `.env` and fill in the values.
2. Ensure `COOKIES_PATH` points to a real cookies file (for example `./data/cookies.txt`).
3. Ensure `DB_FILE_NAME` points to a writable path (for example `./data/bot.db`).
4. If you run the bot in Docker, set `IS_DOCKERIZED=true` and provide environment variables via the compose/env files.

## Run locally
- `make run` (or `go run cmd/bot/main.go`).

## Tests
- `make test` or `go test ./... -v`.

## Tidy dependencies
- `make tidy`.

## Docker helpers
- `make build`, `make rebuild`, `make down`, `make logs`.
- `make deploy` uses `deploy/deploy.sh` and `deploy/deploy.env`.

## Notes
- Downloads go to `./output` (created on startup).
- Forward mode requires the bot to be an admin in the target group.
- Set `IS_DEBUG=true` for verbose logs.
