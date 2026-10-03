# AGENTS.md

mDNS discovery service in Go 1.23 (`github.com/Buhrietoe/findmdns`). Discovers local network services via `github.com/hashicorp/mdns`, exposes a JSON API (gorilla/mux) and web UI.

## Commands

```bash
go build -o findmdns          # Build binary
./findmdns                    # JSON mode: one scan, prints devices to stdout
./findmdns serve              # HTTP server on $PORT (default 8080)
go run .                      # Quick iteration
go test -v ./...              # No tests exist yet
go vet ./...                  # Passes clean
go fmt ./...                  # Repo is gofmt-clean
go mod tidy                   # Update dependencies
```

Every scan blocks ~5s (timeout set in `discovery/manager.go`), even with zero results. `serve` runs a blocking initial scan before listening, so wait ~6s before curling endpoints.

## Structure

- `main.go` — CLI entrypoint; dispatches to JSON mode or HTTP server on `serve`
- `discovery/manager.go` — mDNS scan logic, device cache with `sync.RWMutex`; scans are serialized via a dedicated `scanMu` mutex
- `api/handlers.go`, `api/routes.go` — HTTP handlers and gorilla/mux routes
- `models/device.go` — Device struct (ID is a 16-char SHA256 hex hash; no Protocol/TTL/Service fields)
- `ui/index.html` + `static/script.js` — dashboard SPA; script.js renders both the device list and `/device/{id}` detail view inline into `#devicesTableContainer`

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/devices` | All discovered devices |
| GET | `/api/device/{id}` | Single device by hash ID |
| POST | `/api/scan` | Synchronous re-scan (blocks ~5s); returns `{"status":"completed","devices":N}` or 500 with error |
| GET | `/api/last-scan` | `{"last_scan":"<RFC3339 timestamp>"}` |

## Gotchas

- Scans merge results into the existing cache; devices discovered in previous scans persist. There is no TTL/expiry — stale entries remain until manually cleared or the process restarts.
- `Scan()` queries `_services._dns-sd._udp` (DNS-SD advertisement enumeration). Device IDs are hash-based (SHA256 prefix) and URL-safe; display names are separate.
- Static/ and ui/ paths are resolved from the process working directory (`os.Getwd()`) at runtime — run `serve` from the repo root.

## Environment

- `PORT` — HTTP server port (default `"8080"`)
- Scans need working multicast networking; empty results usually mean container/firewall blocking, not a code bug. README notes distros may need mDNS daemons configured (e.g. systemd-resolved) first.
