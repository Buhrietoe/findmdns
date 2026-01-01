# findmdns

Web-based mDNS device discovery tool.

## Features

- Web UI for discovering mDNS services
- REST API for programmatic access
- Real-time device listing and details
- Responsive dark-themed UI (midnight galaxy theme)
- Automatic background discovery
- Command-line JSON output (backward compatible)

## Installation

```bash
go install -v github.com/Buhrietoe/findmdns@latest
```

## Usage

### Command-line Mode (Default)
Outputs JSON array of discovered devices:
```bash
findmdns
```

### Web Server Mode
Starts the web UI and API server:
```bash
findmdns serve
```

The application will start on port 8080 by default.

## API Endpoints

- `GET /api/devices` - Return list of discovered devices
- `GET /api/device/:id` - Return details for specific device
- `POST /api/scan` - Trigger new discovery scan
- `GET /api/last-scan` - Get time of last scan

## Web UI

Access the web interface at:
```
http://localhost:8080
```

The UI provides:
- Dashboard showing all discovered devices in a table
- Device detail view with comprehensive information
- Real-time scanning capability

## Why?

This tool simplifies mDNS discovery by providing both command-line and web interfaces, making it easier to investigate services on your local network.

## Caveat

A correctly setup system for mDNS is required. For example, Archlinux using systemd-resolved with /etc/resolv.conf symlinked to /run/systemd/resolve/stub-resolv.conf.