# Gator

Gator is a small command-line RSS reader/aggregator written in Go. It lets you:
- Register/login a user (stored in PostgreSQL)
- Add RSS feeds and follow/unfollow them
- Periodically aggregate new posts from followed feeds
- Browse your latest posts right in the terminal

This README documents how to set up the project locally, the tech stack, configuration, commands, and how to run tests.

## Stack
- Language: Go (go 1.25.1 as per go.mod)
- Database: PostgreSQL
- SQL data access: sqlc (code generated into internal/database)
- RSS parsing: encoding/xml + net/http (custom rss package)

Package manager/build tool: Go modules (go.mod/go.sum).

Entry point: main.go (invokes the CLI in internal/cli).

## Requirements
- Go 1.25.1 or newer (see go.mod)
- PostgreSQL (version not pinned; 13+ recommended) — TODO: confirm minimal supported version.
- sqlc (for regenerating database access code) — optional at runtime.

## Database schema and migrations
SQL schema migration files live under:
- sql/schema/*.sql
- The files include goose-style annotations (e.g., `-- +goose Up/Down`).

You can apply the schema in one of two ways:
1) With a migration tool (preferred):
   - The files appear compatible with goose, but no goose config is provided in this repo. TODO: confirm the intended migration tool and add instructions.
2) Manually via psql, applying in numeric order:
   - 001_users.sql
   - 002_feeds.sql
   - 003_feed_follows.sql
   - 004_feed_last_fetched.sql
   - 005_posts.sql

Example (psql):
- psql "$DB_URL" -f sql/schema/001_users.sql
- psql "$DB_URL" -f sql/schema/002_feeds.sql
- ... and so on.

The SQL queries consumed by sqlc live under sql/queries/*.sql.

## Configuration (env/vars)
Gator reads its configuration from a JSON file in your home directory:
- Path: ~/.gatorconfig.json
- Fields:
  - db_url: PostgreSQL connection string
  - current_user_name: set automatically by the app when you register/login

Example ~/.gatorconfig.json:
{
  "db_url": "postgres://user:pass@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}

Note:
- The application connects to PostgreSQL using the value of db_url.
- The current_user_name is written by the app when you run register/login.

No other environment variables are required by the app. If you prefer not to use a config file, you could extend config.Read() to read from env vars — TODO: consider adding this feature.

## Setup
1) Install dependencies
- Install Go (matching go.mod)
- Install and start PostgreSQL
- (Optional) Install sqlc if you plan to regenerate database code

2) Create a PostgreSQL database
- createdb gator

3) Create config file
- Create ~/.gatorconfig.json (see example above) and set db_url to your database connection string

4) Apply database schema
- Use goose (TODO) or psql to apply the files in sql/schema in order

5) Build
- go build -o gator ./

## Running
You can run the app via go run or the compiled binary.

- go run ./ main.go <command> [args]
- Or after building: ./gator <command> [args]

Examples:
- Initialize a user
  - ./gator register alice
  - ./gator login alice
- Show users and current user
  - ./gator users
- Manage feeds
  - ./gator addfeed "My Blog" https://example.com/rss.xml
  - ./gator feeds
  - ./gator follow https://example.com/rss.xml
  - ./gator following
  - ./gator unfollow https://example.com/rss.xml
- Browse latest posts (default: 2 posts; pass an optional limit)
  - ./gator browse
  - ./gator browse 10
- Run the aggregator periodically (duration uses Go time format like 30s, 5m, 1h)
  - ./gator agg 30s

Notes:
- reset command deletes all users (and may cascade-delete related data depending on FK constraints). Use with caution: ./gator reset

## CLI commands
The CLI commands are registered in internal/cli/cli.go and implemented in internal/cli/handlers.go.
- login <username>
- register <username>
- reset
- users
- agg <duration>
- addfeed <name> <url>
- feeds
- follow <url>
- following
- unfollow <url>
- browse [limit]

Some commands require you to be logged in (middlewareLoggedIn), e.g., addfeed, follow, following, unfollow, browse.

## Scripts and tooling
- sqlc generate code (requires sqlc installed):
  - sqlc generate
  - Configuration: sqlc.yaml
  - Generated package: internal/database

- Tests:
  - go test ./...
  - Current tests target internal/config (file I/O for ~/.gatorconfig.json). Running tests may write to or depend on a temp location; review tests before running in sensitive environments.

## Project structure
- main.go — application entry point
- internal/cli — CLI framework, command handlers, and RSS scraper logic
- internal/config — config file read/write (~/.gatorconfig.json)
- internal/database — sqlc-generated models and query methods
- internal/rss — RSS fetch and parse utilities
- sql/schema — database schema (with goose-style annotations)
- sql/queries — SQL queries used by sqlc
- go.mod / go.sum — dependencies
- sqlc.yaml — sqlc configuration

## License
No license file is present in this repository. TODO: add a LICENSE and update this section (e.g., MIT, Apache-2.0).

## Contributing
- Open PRs with small, focused changes.
- For database changes, update sql/schema and sql/queries and run sqlc generate.
- Please add/update tests where applicable.

## Troubleshooting
- Database connection errors at startup usually mean db_url is missing/incorrect in ~/.gatorconfig.json.
- Aggregator (agg) does not fetch posts: ensure at least one feed exists and that it has new items; the scraper skips duplicate URLs silently.
- Character entities in titles/descriptions show up escaped: the rss package unescapes strings before saving, but sources vary.

---
Last updated: 2025-10-16
