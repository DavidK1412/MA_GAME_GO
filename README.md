# MA_GAME_GO

This service exposes minimal HTTP endpoints for managing sessions, matches, and moves while persisting data in Postgres following the schema in `internal/adapters/db/seed.sql`.

## Running locally

1. Export a Postgres connection string as `DATABASE_URL`, for example:
   ```bash
   export DATABASE_URL="postgres://user:pass@localhost:5432/ma_game?sslmode=disable"
   ```
2. Optionally configure pool settings: `DB_MAX_CONNS`, `DB_MIN_CONNS`, `DB_MAX_CONN_LIFETIME`, `DB_MAX_CONN_IDLE`, and `APP_NAME`.
3. Run the server:
   ```bash
   go run ./cmd/server
   ```
   The server listens on `:8080` by default; override with `PORT`.

## Sample requests

```bash
# Create a session (player/device inferred)
curl -X POST http://localhost:8080/sessions \
  -H 'Content-Type: application/json' \
  -d '{"game_id":"11111111-1111-1111-1111-111111111111"}'

# Create a match for a session with a difficulty level
curl -X POST http://localhost:8080/matches \
  -H 'Content-Type: application/json' \
  -d '{"session_id":"<SESSION_ID>","difficulty_id":1}'

# Record a move for an existing match
curl -X POST http://localhost:8080/matches/<MATCH_ID>/moves \
  -H 'Content-Type: application/json' \
  -d '{"movement":[0,1,2]}'
```
