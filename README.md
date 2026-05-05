# stackprobe

Lightweight HTTP health-check orchestrator that aggregates service status from multiple endpoints into a single dashboard.

---

## Installation

```bash
go install github.com/youruser/stackprobe@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/stackprobe.git && cd stackprobe && go build -o stackprobe .
```

---

## Usage

Define your endpoints in a `config.yaml` file:

```yaml
interval: 30s
services:
  - name: API Gateway
    url: https://api.example.com/health
  - name: Auth Service
    url: https://auth.example.com/ping
  - name: Database Proxy
    url: https://db-proxy.internal/status
```

Start the dashboard:

```bash
stackprobe --config config.yaml --port 8080
```

Then open `http://localhost:8080` in your browser to view the aggregated health status of all registered services.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `config.yaml` | Path to configuration file |
| `--port` | `8080` | Port to serve the dashboard on |
| `--timeout` | `5s` | Per-request probe timeout |

---

## License

[MIT](LICENSE)