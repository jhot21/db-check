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

| Flag         | Short | Env        | Description       |
|--------------|-------|------------|-------------------|
| `--host`     |       | `HOST`     | Database host     |
| `--port`     | `-p`  | `PORT`     | Database port     |
| `--name`     | `-n`  | `NAME`     | Database name     |
| `--user`     | `-u`  | `USER`     | Database username |
| `--password` |       | `PASSWORD` | Database password |

All flags can be set via environment variables (uppercase, no dashes).

## Docker Compose Example

Use `db-check` as a sidecar with a healthcheck. Your app declares `depends_on` with `condition: service_healthy` so it waits until the database is reachable.

```yaml
services:
  db-check:
    image: codeberg.org/jhot/db-check:latest
    healthcheck:
      test: ["CMD", "/dbcheck", "mysql"]  # or "postgres"
      interval: 15s
      timeout: 35s
      retries: 5
      start_period: 10s
    environment:
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

See [docker-compose.example.yml](docker-compose.example.yml) for a full working example.

## Image

Pre-built multi-arch images (`linux/amd64`, `linux/arm64`) are available at:

```
codeberg.org/jhot/db-check:latest
```

## License

[MIT](LICENSE)
