# Repository Guidelines

## Project Structure & Module Organization

This is a Go module for a simple banking backend: `github.com/NatdanaiKhe/simplebank`.

- `db/migration/` contains PostgreSQL migration files. Keep paired `*.up.sql` and `*.down.sql` files with matching sequence numbers.
- `db/schema/` contains schema input used by sqlc.
- `db/query/` contains sqlc query definitions grouped by domain, such as `account.sql`, `entry.sql`, and `transfer.sql`.
- `db/sqlc/` contains generated Go data-access code plus database tests.
- `sqlc.yaml` configures sqlc generation into package `db` under `db/sqlc`.

## Build, Test, and Development Commands

- `make migrate-up` applies database migrations to `postgres://postgres:postgres@localhost:5432/bank?sslmode=disable`.
- `make migrate-down` rolls migrations back against the same local database.
- `make migrate-create name=add_table_name` creates a new migration pair. Verify the output path before use.
- `make sqlc` regenerates Go code from `db/query` and `db/schema`.
- `go test ./...` runs the Go test suite. Tests expect a local PostgreSQL `bank` database matching the migration schema.

## Coding Style & Naming Conventions

Use standard Go formatting with `gofmt`. Keep package names short and lowercase; generated database code uses package `db`. Prefer explicit error handling and small focused functions. SQL query files should stay grouped by aggregate or table, and sqlc query names should use clear action names such as `CreateAccount` or `GetTransfer`.

## Testing Guidelines

Tests use Go's `testing` package and `github.com/stretchr/testify/require`. Place database tests next to the generated data-access package in `db/sqlc`. Name tests as `TestXxx`, for example `TestCreateAccount`. New query behavior should include coverage for successful execution and important error cases where practical.

## Commit & Pull Request Guidelines

The current history uses short, imperative commit messages, for example `Create README.md for Go learning repository`. Keep commits focused and describe the change in plain language.

Pull requests should include a concise summary, any database migration impact, commands run for validation, and linked issues when applicable. Include screenshots only when a future change adds user-facing UI.

## Security & Configuration Tips

Do not commit real credentials. The local database URL in the `Makefile` is for development only. Keep schema changes reversible by maintaining accurate down migrations.

## Git Commit Messages

Always use Conventional Commits.

Format:

<type>(<scope>): <summary>

Rules:

- Types: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert
- Use imperative mood.
- Keep the summary under 72 characters.
- Do not end the summary with a period.
- Infer the scope from the changed files when possible.
- Output only the commit message unless the user requests otherwise.
