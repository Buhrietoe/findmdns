# AGENTS.md
## Overview
This repository implements an mDNS discovery service that:
- Discovers services on the local network using `github.com/hashicorp/mdns`.
- Provides a JSON API and a small web UI for viewing discovered devices.
- Written in Go 1.23 with minimal external dependencies.

## Project Structure
```
findmdns/
├─ main.go                 # CLI entrypoint; parses args, runs JSON output or HTTP server
├─ discovery/              # mDNS discovery logic
│  └─ manager.go           # Manager struct, Scan, GetDevices, GetDevice, concurrency safety
├─ api/                    # HTTP API handlers
│  ├─ handlers.go          # HTTP endpoints: /api/devices, /api/device/{id}, /api/scan, /api/last-scan
│  └─ routes.go            # Route configuration using gorilla/mux
├─ models/                 # Domain models
│  └─ device.go            # Device struct representing discovered mDNS services
├─ static/                 # Static assets served by HTTP server
│  ├─ script.js            # Frontend dashboard logic
│  ├─ style.css            # Styling (Midnight Galaxy theme)
│  └─ device-script.js     # Additional script (currently unused)
├─ ui/                     # UI entrypoint
│  └─ index.html           # Single-page app that bootstraps the dashboard
├─ go.mod / go.sum         # Module definitions and dependencies
└─ *.md / *.txt            # Documentation files
```

## Building & Running
- **Build**: `go build -o findmdns` creates the static CLI binary.
- **Run (JSON mode)**: `./findmdns` prints discovered devices as JSON to stdout.
- **Run (HTTP server)**: `./findmdns serve` starts a web server on `$PORT` (default 8080) serving the UI at `/` and API at `/api/*`.
- **Development build**: `go run main.go` or `go run .` for quick iteration.
- **Install**: `go install github.com/Buhrietoe/findmdns/cmd/findmdns@latest` (if packaged as a command).

## API Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/devices` | Return JSON array of all discovered devices |
| GET | `/api/device/{id}` | Return details for a specific device |
| POST | `/api/scan` | Trigger a new discovery scan |
| GET | `/api/last-scan` | Return timestamp of the last scan |

All responses are JSON; error responses use appropriate HTTP status codes.

## UI Interaction
- The UI is served from `/static/` and `index.html` at root.
- `script.js` handles:
  - Loading devices on page load
  - Scanning on button click (`/api/scan`)
  - Updating last scan timestamp
  - Rendering device details in a modal
- Styling is theme‑driven via CSS variables in `style.css`.

## Coding Patterns & Conventions
- **Package organization**: Each logical concern (discovery, API, models) lives in its own package.
- **Concurrency safety**: `Manager` uses `sync.RWMutex` to protect shared device maps.
- **Error handling**: Functions return `error`; HTTP handlers translate errors to HTTP status codes.
- **JSON encoding**: Consistent struct tags (`json:"field"`), indentation with two spaces.
- **Import grouping**: Third‑party imports grouped together, standard library imports grouped separately, sorted alphabetically within groups.
- **Naming**: Exported structs and functions use UpperCase; internal lowerCase. Method names reflect actions (`NewManager`, `Scan`, `GetDevices`).
- **Error messages**: Include original error text when wrapping, but avoid leaking internal details to clients.

## Testing
- Current test suite (`main_test.go`) only verifies compilation.
- Add unit tests under `*_test.go` files:
  - Test `Manager.updateDevices` parsing logic
  - Test `Device` JSON marshaling
  - Integration test for HTTP handlers using `net/http/httptest`
- Run tests with `go test -v ./...`.
- Coverage can be checked with `go test -cover`.

## Linting & Formatting
- `go vet` – static analysis (run before commits)
- `go fmt ./...` – code formatting
- `go mod tidy` – keep dependencies up to date
- Optional: run `staticcheck` or `golangci-lint` for additional checks (not required by CI currently).

## Gotchas & Linter Hints
- The current codebase triggers a linter warning: `interface{} can be replaced by any` in `api/handlers.go:69`. Consider replacing `map[string]interface{}` with a concrete struct when the shape is known.
- Channels are used for discovery entry streaming; ensure no nil channel panics when starting goroutines.
- `github.com/hashicorp/mdns` requires Go 1.20+; verify `go.mod` uses at least Go 1.23.

## Environment Variables
- `PORT` – HTTP server port (default `"8080"` if unset).
- `PWD` – Used to resolve static file paths; defaults to current working directory.

## Common Development Tasks
```bash
# Build & run unit tests
go test -v ./...

# Format code
go fmt ./...

# Update dependencies
go mod tidy

# Build binary
go build -o findmdns

# Run in JSON mode
./findmdns

# Run HTTP server (development)
PORT=8080 ./findmdns serve

# Run with live reload (if using reflex or similar)
```

## Example Workflow
1. **Add a new device field**  
   - Update `models.Device` struct in `models/device.go`.  
   - Adjust JSON marshaling in `api/handlers.go`.  
   - Add unit test covering the new field.  

2. **Expose a new API endpoint**  
   - Add route in `api/routes.go`.  
   - Implement handler in `api/handlers.go`.  
   - Document endpoint in this AGENTS.md.  

3. **Update UI**  
   - Modify `static/script.js` to call the new endpoint.  
   - Add corresponding CSS if needed.  

4. **Run tests**  
   ```bash
   go test -v ./...
   ```

## Additional Resources
- `github.com/hashicorp/mdns` package documentation
- Gorilla mux routing guide
- Go standard library net/http docs

## Development Environment
- Requires Go 1.23 or later.
- `go mod tidy` to ensure dependencies are up‑to‑date.

## Dependency Management
- Run `go mod tidy` after adding new imports.
- Use `go get` to add third‑party packages.
