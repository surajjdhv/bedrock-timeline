# Artisan's Timeline

A real-time player activity monitor for Minecraft Bedrock servers. Reads `journalctl` logs and displays a live timeline of player joins/leaves via a web interface.

## Features

- Real-time streaming of player connection events via WebSocket
- SQLite database for persistent history
- Clean, responsive Svelte frontend
- Single binary deployment (embedded static files)
- Designed for Cloudflare Tunnel

## Requirements

- Go 1.21+ (with CGO for SQLite)
- Node.js 18+ (for frontend build only)
- Systemd journal access (for reading bedrock service logs)

## Quick Start

### Install Dependencies

```bash
make install-deps
```

### Build

```bash
make build
```

This compiles the frontend and builds the Go binary into `build/bedrock-timeline`.

### Run

```bash
PORT=3000 DB_PATH=./data/bedrock.db JOURNAL_UNIT=bedrock ./build/bedrock-timeline
```

Environment variables:
- `PORT` - Server port (default: `8080`)
- `DB_PATH` - SQLite database path (default: `data/bedrock.db`)
- `JOURNAL_UNIT` - Systemd unit name for journalctl (default: `bedrock`)
- `TITLE` - Application title displayed in header (default: `Artisan's Timeline`)

## Deployment

### Systemd Service

```bash
sudo make install
sudo systemctl enable --now bedrock-timeline
```

### Cloudflare Tunnel

1. Install `cloudflared` on your server
2. Create a tunnel:

```bash
cloudflared tunnel create bedrock-timeline
```

3. Configure routing:

```yaml
# ~/.cloudflared/config.yml
tunnel: <your-tunnel-id>
ingress:
  - hostname: timeline.yourdomain.com
    service: http://localhost:8080
  - service: http_status:404
```

4. Run the tunnel:

```bash
cloudflared tunnel run bedrock-timeline
```

For production, run cloudflared as a systemd service:

```bash
cloudflared service install
systemctl enable --now cloudflared
```

## Development

### Run Frontend Dev Server

```bash
make dev
```

Frontend runs on http://localhost:5173 with hot reload.

### Run Backend

```bash
make run
```

Backend runs on http://localhost:8080.

## API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /api/events` | Recent events (optional `?player=` and `?type=` filters) |
| `GET /api/players` | List of all recorded players |
| `GET /api/stats` | Aggregate statistics |
| `WS /ws` | WebSocket for real-time event stream |

## Log Format Support

The parser recognizes these Bedrock log patterns:

```
[timestamp] Player connected: Steve, xuid: 12345...
[timestamp] Player disconnected: Steve, xuid: 12345...
[timestamp] Steve joined the game
[timestamp] Steve left the game
```

## Project Structure

```
bedrock-timeline/
├── cmd/server/main.go       # Entry point
├── internal/
│   ├── journal/             # journalctl stream reader
│   ├── parser/              # Bedrock log parser
│   ├── store/               # SQLite operations
│   └── ws/                  # WebSocket hub
├── web/                     # Svelte frontend
├── build/                   # Compiled binary output
└── Makefile
```

## License

MIT