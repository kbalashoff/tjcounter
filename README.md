# tjcounter

`tjcounter` is now a complete minimal real-time counter app:

- **Backend**: Go HTTP server exposing JSON endpoints and a Server-Sent Events stream.
- **Frontend**: Embedded single-page HTML + JavaScript client.
- **Realtime sync**: Multiple browser tabs stay in sync via SSE.

## Run

```bash
go run ./cmd/tjcounter
```

Open http://localhost:8080.

## API

- `GET /api/state` → current value
- `POST /api/increment` → increments value
- `POST /api/decrement` → decrements value
- `POST /api/reset` → resets to `0`
- `GET /events` → SSE stream with `value` events

## Project structure

- `cmd/tjcounter/main.go` - HTTP server and SSE broker
- `internal/counter/counter.go` - concurrency-safe counter domain logic
- `web/index.html` - client UI

## Test

```bash
go test ./...
```
