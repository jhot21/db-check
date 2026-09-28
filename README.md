# DB Check

A small command line utility to check database connectivity. Intended to be used as a Docker Compose sidecar container that blocks dependent services from starting until a remote database is reachable.

## Usage

```
dbcheck <command> [flags]
```

### Commands

| Command    | Description                        |
|------------|------------------------------------|
| `mysql`    | Check connection to MySQL/MariaDB  |
| `postgres` | Check connection to PostgreSQL     |

### Flags

| Flag         | Short | Env        | Default                     | Description                          |
|--------------|-------|------------|-----------------------------|--------------------------------------|
| `--host`     |       | `HOST`     | *required*                  | Database host                        |
| `--port`     | `-p`  | `PORT`     | `3306` (mysql), `5432` (postgres) | Database port                  |
| `--name`     | `-n`  | `NAME`     |                             | Database name                        |
| `--user`     | `-u`  | `USER`     |                             | Database username                    |
| `--password` |       | `PASSWORD` |                             | Database password                    |

All flags can be set via environment variables (uppercase, no dashes). Passwords and other values may contain any characters; no URL-escaping is needed.

> **Note:** `USER` is also set by most shells to your login name. When running `dbcheck` outside a container, pass `--user` explicitly (flags take precedence over env vars) or you may connect as the wrong user.

Subcommand-specific flags:

| Command    | Flag        | Env       | Default | Description |
|------------|-------------|-----------|---------|-------------|
| `postgres` | `--sslmode` | `SSLMODE` | `allow` | `disable`, `allow`, `prefer`, `require`, `verify-ca` or `verify-full` |
| `mysql`    | `--tls`     | `TLS`     | `false` | `true`, `false`, `skip-verify` or `preferred` |

The defaults perform no certificate verification. For databases reachable over an untrusted network, use `--sslmode verify-full` (postgres) or `--tls true` (mysql).

Run `dbcheck --version` to print the build version.

### Exit status

| Status | Meaning |
|--------|---------|
| `0`    | The database accepted the connection (`Success` is printed). |
| `1`    | Anything else: missing/invalid settings, no subcommand, connection refused, authentication failure, or the 10 second timeout. The error is printed to stderr. |

The tool makes a single attempt. Retrying is left to Docker's healthcheck settings (`--interval`, `--retries`, `--start-period`).

### Docker-only: `TYPE`

When run as a container, the `TYPE` environment variable controls which database subcommand the built-in healthcheck runs (`mysql` or `postgres`, default: `mysql`). This is not a CLI flag. An empty or invalid `TYPE` makes the healthcheck fail rather than pass.

The built-in healthcheck runs every 15s, allows 5 retries and a 30s start period. Override it per service in Compose with the `healthcheck:` key if you need different timing.

## Docker Compose Example

Use `db-check` as a sidecar with a healthcheck. Your app declares `depends_on` with `condition: service_healthy` so it waits until the database is reachable.

```yaml
services:
  db-check:
    image: ghcr.io/jhot21/db-check:latest
    environment:
      - TYPE=mysql # or postgres
      - HOST=your-db-host
      - PORT=3306
      - USER=myuser
      - PASSWORD=mypassword
      - NAME=mydb

  app:
    image: your-app
    depends_on:
      db-check:
        condition: service_healthy
```

For PostgreSQL, set `TYPE=postgres` and, if needed, `PORT` and `SSLMODE`:

```yaml
  db-check:
    image: ghcr.io/jhot21/db-check:latest
    environment:
      - TYPE=postgres
      - HOST=your-db-host
      - USER=myuser
      - PASSWORD=mypassword
      - NAME=mydb
      - SSLMODE=verify-full
```

See [docker-compose.example.yml](docker-compose.example.yml) for a full working example.

## Image

Pre-built multi-arch images (`linux/amd64`, `linux/arm64`) are available at:

```
ghcr.io/jhot21/db-check:latest
```

## Development

```sh
GOWORK=off go vet ./... && GOWORK=off go test ./...
```

## License

[MIT](LICENSE)
