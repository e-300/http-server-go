# Chirpy

A RESTful API server built from scratch in Go — a Twitter-like microblogging backend supporting user registration, authentication (JWT + refresh tokens), CRUD operations on "chirps," and webhook-driven account upgrades.

Built by hand as a learning project to get fluent writing HTTP servers and REST APIs in Go. 
AI was used only to help me generate this README, never to generate code, this is handcrafted artisanal backend logic so it might not be perfect. 


## Features

- **User management** — registration, login, email/password updates with Argon2id password hashing
- **JWT authentication** — short-lived access tokens (1hr) + long-lived refresh tokens (60 days) with revocation support
- **Chirps** — create, read, delete with author-only deletion enforcement, content filtering, 140-char limit, and sort/filter by author
- **Webhook integration** — Polka payment provider webhook for account upgrades with API key validation
- **Admin endpoints** — request metrics and database reset (dev-only)
- **Static file serving** — serves frontend assets via the built-in fileserver with middleware-wrapped metrics tracking

## Tech Stack

- **Go** — standard library `net/http` with `ServeMux` routing (no frameworks)
- **PostgreSQL** — persistent storage
- **sqlc** — type-safe Go code generation from SQL queries
- **goose** — database migration management
- **Argon2id** — password hashing via `alexedwards/argon2id`
- **JWT** — `golang-jwt/jwt/v5` for access token signing and validation
- **godotenv** — environment variable management

## Project Structure

```
├── main.go                          # Server entrypoint, config, route registration
├── json.go                          # JSON response/error helpers
├── metrics.go                       # Request counter middleware + admin metrics
├── readiness.go                     # Health check endpoint
├── reset.go                         # Admin reset (dev-only)
├── handler_user_create.go           # POST /api/users
├── handler_user_login.go            # POST /api/login
├── handler_user_update_login.go     # PUT  /api/users
├── handler_chirp_create.go          # POST /api/chirps
├── handler_chirps_get.go            # GET  /api/chirps, GET /api/chirps/{chirpID}
├── handler_chirps_delete.go         # DELETE /api/chirps/{chirpID}
├── handler_refresh.go               # POST /api/refresh
├── handler_revoke.go                # POST /api/revoke
├── handler_webhook.go               # POST /api/polka/webhooks
├── internal/
│   ├── auth/
│   │   ├── auth.go                  # Password hashing (Argon2id)
│   │   ├── jwt.go                   # JWT creation
│   │   ├── validate_jwt.go          # JWT validation
│   │   ├── get_bearer.go            # Header token extraction
│   │   ├── refresh_token.go         # Refresh token generation (256-bit hex)
│   │   └── auth_test.go             # Table-driven tests
│   └── database/
│       ├── sql/
│       │   ├── schema/              # Goose migrations (001-005)
│       │   └── queries/             # sqlc query definitions
│       ├── db.go                    # sqlc generated DB interface
│       ├── models.go                # sqlc generated models
│       ├── post.sql.go              # sqlc generated chirp queries
│       ├── refresh.sql.go           # sqlc generated refresh token queries
│       └── user.sql.go              # sqlc generated user queries
└── sqlc.yaml                        # sqlc configuration
```

## API Endpoints

### Auth
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/users` | None | Create a new user |
| `POST` | `/api/login` | None | Login, returns access + refresh tokens |
| `PUT` | `/api/users` | Bearer (access) | Update email and password |
| `POST` | `/api/refresh` | Bearer (refresh) | Exchange refresh token for new access token |
| `POST` | `/api/revoke` | Bearer (refresh) | Revoke a refresh token |

### Chirps
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/chirps` | Bearer (access) | Create a chirp (140 char max) |
| `GET` | `/api/chirps` | None | List all chirps (supports `?sort=desc` and `?author_id=`) |
| `GET` | `/api/chirps/{chirpID}` | None | Get a single chirp |
| `DELETE` | `/api/chirps/{chirpID}` | Bearer (access) | Delete own chirp |

### Webhooks
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/polka/webhooks` | ApiKey | Handle Polka payment events |

### Admin
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/admin/metrics` | None | View fileserver hit count |
| `POST` | `/admin/reset` | None (dev-only) | Reset metrics and delete all users |

### Health
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/healthz` | Readiness check |

## Setup

### Prerequisites
- Go 1.25+
- PostgreSQL
- [goose](https://github.com/pressly/goose) for migrations
- [sqlc](https://sqlc.dev/) for query code generation

### Install and run

```bash
git clone https://github.com/e-300/http-server-go.git
cd http-server-go

# Copy .env.example and fill in your values
cp .env.example .env

# Run migrations
goose -dir internal/database/sql/schema postgres "$DB_URL" up

# Generate sqlc code (if modifying queries)
sqlc generate

# Run the server
go run .
```

The server starts on `http://localhost:8080`.

## Database Migrations

Migrations are in `internal/database/sql/schema/` and managed by goose:

1. `001_users.sql` — users table
2. `002_posts.sql` — posts table with foreign key to users
3. `003_passwords.sql` — adds hashed_password column
4. `004_refresh_tokens.sql` — refresh tokens table with expiry and revocation
5. `005_polka_webhook.sql` — adds is_chirpy_red flag for premium users

## Architecture Notes

**Middleware pattern** — `middlewareMetricsInc` wraps the fileserver handler, demonstrating the `http.Handler` interface / `http.HandlerFunc` type conversion pattern for composable request processing.

**Auth flow** — Login returns both a JWT access token (1hr, stateless, used for API auth) and a refresh token (60 days, stored in DB, used only to mint new access tokens). Refresh tokens can be revoked, access tokens cannot — they just expire.

**Thread safety** — `atomic.Int32` for the request hit counter ensures safe concurrent reads/writes without mutex overhead.

**sqlc + goose** — SQL is the source of truth. Goose handles schema evolution, sqlc generates type-safe Go code from query definitions, eliminating hand-written SQL mapping.
