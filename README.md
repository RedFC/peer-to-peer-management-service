# P2P Management Service

Go/Iris service that manages users, roles, groups, and subscriber invitations for the peer‑to‑peer communication platform. Data at rest (names, emails, passwords) is encrypted with AES‑256, authentication is handled with RS256 JWTs, and database access is type‑safe via sqlc.

---

## Stack
- Go 1.24.5, Iris v12
- PostgreSQL + `golang-migrate`
- sqlc for query-safe DB access
- JWT (RS256) with generated RSA keys
- AES‑256 encryption (requires 32‑byte `ENCRYPTION_SECRET`)
- Swagger (swag) → OpenAPI 3.0 converter
- Email backends: AWS SES, local file drop, or postfix/sendmail

---

## Quick Start
1) **Install prerequisites**
   - Go 1.24+, `make`, OpenSSL  
   - `golang-migrate` CLI, `sqlc`, `swag`  
   - PostgreSQL 14+ running locally

2) **Create `.env`** (example)
   ```env
   APP_NAME=p2p management service
   PORT=8080

   DB_HOST=127.0.0.1
   DB_PORT=5433
   DB_USER=postgres
   DB_PASSWORD=mysecretpassword
   DB_NAME=postgres
   DB_SSLMODE=disable

   CORS_ALLOWED_ORIGINS=http://localhost:3000
   CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
   CORS_ALLOWED_HEADERS=Authorization,Content-Type
   CORS_EXPOSED_HEADERS=

   ENCRYPTION_SECRET=0123456789abcdef0123456789abcdef  # must be 32 bytes

   AWS_REGION=us-east-1
   AWS_SES_SENDER=no-reply@example.com
   EMAIL_BACKEND=local                # local | ses | postfix
   EMAIL_SAVE_PATH=./logs/emails      # used when EMAIL_BACKEND=local
   HUBS=                              # optional QR hubs payload
   SENDMAIL_PATH=/usr/sbin/sendmail   # only for postfix backend
   ```

3) **Generate JWT keys**
   ```bash
   make generate-jwt-keys
   ```

4) **Run migrations (destructive reset)**
   - The app calls `scripts.RunMigrations` on startup, which executes `make migrate-dev` (drops all tables, then migrates up). **Do not run against a database you need to keep.**
   - For a clean local DB once, run:
     ```bash
     make migrate DATABASE_URL="postgresql://postgres:mysecretpassword@127.0.0.1:5433/postgres?sslmode=disable"
     ```

5) **Start the API**
   ```bash
   make run   # starts on :8080 and seeds data
   ```
   Health check: `GET http://localhost:8080/api/v1/health-check`

6) **Swagger / OpenAPI**
   ```bash
   make swagger
   # then visit http://localhost:8080/swagger/index.html
   # SWAGGER_SERVER_URLS (comma or newline separated) controls generated server URLs
   ```

7) **Run tests**
   ```bash
   go test ./...
   ```
   Tests use mocks; no database required.

---

## Environment Reference
| Variable | Description | Default |
| --- | --- | --- |
| `APP_NAME` | Service label | `p2p management service` |
| `PORT` | HTTP port | `8080` |
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE` | PostgreSQL settings | `127.0.0.1`, `5433`, `postgres`, `mysecretpassword`, `postgres`, `disable` |
| `ENCRYPTION_SECRET` | 32‑byte key for AES‑256 encryption/decryption | **required** |
| `CORS_ALLOWED_ORIGINS`, `CORS_ALLOWED_METHODS`, `CORS_ALLOWED_HEADERS`, `CORS_EXPOSED_HEADERS` | CORS configuration (comma-separated) | empty |
| `AWS_REGION`, `AWS_SES_SENDER` | SES region/sender when using SES backend | none |
| `EMAIL_BACKEND` | `local`, `ses`, or `postfix` | `local` |
| `EMAIL_SAVE_PATH` | Folder for saved emails when `EMAIL_BACKEND=local` | `./logs/emails` |
| `SENDMAIL_PATH` | Path to sendmail when `EMAIL_BACKEND=postfix` | `/usr/sbin/sendmail` |
| `HUBS` | Additional QR hub metadata embedded in invitations | empty |
| `DATABASE_URL` | Full DSN used by Makefile migration targets | `postgresql://postgres:mysecretpassword@127.0.0.1:5433/postgres?sslmode=disable` |

---

## Database & Seeding
- Migrations live in `db/migrations` and are generated with `make migration name=...`.
- Startup flow (`main.go`) connects to Postgres, **runs `migrate-dev` (down then up)**, and then seeds:
  - Roles: `super_admin`, `admin`, `user`, `subscriber`
  - Groups: Super/Admin/User plus sample tiers (B-2 Spirit, C-130 Hercules Operators, etc.)
  - Users: `superadmin@p2pservice.com`, `admin@p2pservice.com`, `user@p2pservice.com` (password for all: `Test@12345`)
- Encryption uses `ENCRYPTION_SECRET`; if it is not exactly 32 bytes the service will exit.

---

## API Surface (base path `/api/v1`)
- `GET /health-check` — readiness probe.
- **Auth**: `POST /auth/login` → JWT (RS256). Provide `email`, `password`.
- **Users** (super_admin only): list/create/get/update/delete; attach role or group.
- **Roles** (super_admin only): CRUD.
- **Groups** (super_admin only): CRUD.
- **Subscriber management** (super_admin only):
  - `GET /subscriber` with pagination.
  - `POST /subscriber` create.
  - `GET /subscriber/{id}` fetch by ID.
  - `PUT /subscriber/{id}` update.
  - `DELETE /subscriber/{group_id}/{user_id}` revoke & delete link.
  - `POST /subscriber/{group_id}/{user_id}/revoke` revoke without delete.
  - `GET /subscriber/{group_id}/{user_id}/resend-email` resend invitation.
- All protected routes require `Authorization: Bearer <token>`; super admin role is enforced by middleware.

---

## Developer Tooling (Make targets)
| Command | Purpose |
| --- | --- |
| `make generate-jwt-keys` | Generate RSA keypair in `keys/`. |
| `make migration name=create_table` | Create timestamped migration files. |
| `make migrate` / `make migrate-down` | Apply or roll back all migrations using `DATABASE_URL`. |
| `make migrate-dev` | Roll back everything then re-apply (destroys data). |
| `make sqlc` | Regenerate Go DB access layer from `db/queries`. |
| `make build` | Compile binary to `bin/app`. |
| `make run` | Run the app (also triggers migrations + seeding via startup). |
| `make run-dev` | Drop DB, then run (via `migrate-dev`). |
| `make swagger` | Regenerate Swagger 2.0 via `swag` and convert to OpenAPI 3.0. |
| `make all` | `sqlc` → `build` → `run` → `swagger` in order. |

---

## Project Layout
- `controllers/` — HTTP handlers (auth, users, roles, groups, subscribers)
- `services/` — business logic + email backends (SES, local, postfix)
- `middlewares/` — auth and logging
- `db/` — migrations, SQLC queries, generated code
- `scripts/` — migrations runner, seeding, OpenAPI converter
- `models/` — request/response DTOs
- `tests/` — unit tests with testify mocks
- `docs/` — generated Swagger/OpenAPI artifacts

---

## Troubleshooting
- **Service exits with “ENCRYPTION_SECRET must be 32 bytes”**: fix the env value length.
- **401/403 on protected routes**: ensure you login as `superadmin@p2pservice.com` / `Test@12345` and pass `Authorization: Bearer <token>`.
- **Data keeps disappearing**: startup runs `migrate-dev` (drops all tables). Point to a disposable DB for local dev; do not run this binary against production data without removing that call.
- **Emails not sending**: set `EMAIL_BACKEND=local` to confirm flow; for SES ensure `AWS_REGION` and credentials are available; for postfix ensure `SENDMAIL_PATH` is correct.
