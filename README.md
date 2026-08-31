# Chirpy Server (*Built using boot.dev guide*)

A production-ready, RESTful HTTP API backend for a microblogging service built with Go, PostgreSQL, and SQLC.
**Built using boot.dev guild**

---

## Tech Stack

* **Language:** Go (1.22+)
* **Database:** PostgreSQL
* **Query Builder / Code Generation:** [sqlc](https://sqlc.dev/)
* **Database Migration:** [Goose](https://github.com/pressly/goose)
* **Authentication:** JWT (JSON Web Tokens), Refresh Tokens, bcrypt
* **Routing:** Go Standard Library `net/http` (Go 1.22+ ServeMux routing)

---

## Key Features

* **User Management:** Secure user registration, authentication, and profile updates with hashed passwords.
* **Dual-Token Authentication:**
  * Short-lived JWT Access Tokens (1 hour) for low-latency authorization.
  * Revocable Refresh Tokens (60 days) persisted in PostgreSQL for continuous sessions.
* **Authorization & Access Control:** Strict resource ownership validation (IDOR protection) on sensitive operations.
* **Chirp Management:** Full CRUD support with optional author filtering (`?author_id=...`) and sorting (`?sort=asc|desc`).
* **Profanity Filtering:** Automatic replacement of profane words with `****`.
* **Idempotent Webhooks:** Polka integration for Chirpy Red premium membership upgrades authenticated via API Key headers.
* **Server Metrics:** Internal request counter with admin metrics and dev-reset endpoints.

---

## Environment Variables

Create a `.env` file in the root directory:

```env
PORT=8080
DB_URL=postgres://username:password@localhost:5432/chirpy?sslmode=disable
JWT_SECRET=your_jwt_secret_key_here
POLKA_KEY=your_polka_webhook_api_key_here
```

## Getting Started

### 1. Prerequisites

* Go Go 1.22+ installed
* PostgreSQL database instance running locally or via Docker
* Goose CLI & SQLC CLI installed

### 2. Run database Migrations

```bash
goose -dir sql/schema postgres "postgres://username:password@localhost:5432/chirpy?sslmode=disable" up
```

### 3. Generate Database Code 

```bash
sqlc generate
```

### 4. Run the Server

```bash
go run .
```
*The server will start listening at http://localhost:8080*

## API Reference

### Authentication

* `POST /api/users` - Register a new user
* `POST /api/login` - Login and receive Access & Refresh tokens
* `POST /api/refresh` - Exchange a valid Refresh Token for a new Access Token (Required `Bearer Refresh Token`)
* `POST /api/revoke` - Revoke a Refresh Token (Logout) (Required `Bearer Refresh Token`)

### Users

* `PUT /api/users` - Update user credentials (Required `Bearer Access Token`)

### Chirps

* `POST /api/chirps` - Create a chirp (Requires `Bearer Access Token`)
* `GET /api/chirps` - Get chirps
  * Query parameters: `?author_id=<uuid>&sort=asc|desc`
* `DELETE /api/chirps/{chirpID}` - Delete a chirp (Owner only, requires `Bearer Access Token`)

### Webhooks

* `POST /api/polka/webhooks` - Upgrade user to Chirpy Red (Requires `ApiKey <POLKA_KEY>` header)

### Admin & Metrics
* `GET /admin/metrics` - View server hit count (HTML)
* `POST /admin/reset` - reset metric hit counts (Development only)
