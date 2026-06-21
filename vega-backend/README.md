# Development Setup

## Database migrations

- Install tern: `go install github.com/jackc/tern/v2@latest`
- Create migration: `tern new migration_name`
- Run migrations: `tern migrate`

## Query generation

- Install sqlc: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
- Generate code: `sqlc generate`
