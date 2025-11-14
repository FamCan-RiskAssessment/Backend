## Backend Maintenance Manual

This document is for future developers who need to maintain or extend the **Backend** service.  
For a structural overview, see `ARCHITECTURE.md`.

---

## 1. Development Environment

### 1.1 Prerequisites

- **Go** (version compatible with `go.mod`)
- **Docker** and **Docker Compose**
- **Git**

Optional but recommended:

- An IDE with Go support (GoLand, VS Code + Go extension).
- `make` (for the simple `Makefile` commands).

### 1.2 Environment Variables

Backend configuration is read via the `bootstrap` package from a `.env` file and/or environment variables.

Minimal required entries (actual names/types are defined in `bootstrap/env.go`):

- **Server**
  - `SERVER_PORT`
- **PostgreSQL**
  - `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- **Redis**
  - `RDB_HOST`, `RDB_PORT`
- **JWT / Auth**
  - Paths or settings for JWT keys (see `internal/infrastructure/jwt` and `bootstrap/constants.go`)
- **External services**
  - SMS gateway configuration
  - S3 storage configuration
  - Calc service URLs (`CalcURL` fields in `bootstrap.Env`)
- **Super admin / pagination / OTP**
  - As defined in `bootstrap.Env` (`SuperAdmin`, `Pagination`, `OTP`, etc.)

Create a `.env` file in `Backend/` by copying from an example (if available) or from deployment configuration.

---

## 2. Running the Backend

### 2.1 Running with Docker (recommended for local)

From `Backend/`:

make upThis will:

- Stop existing containers (`docker compose down`).
- Build and start:
  - `backend` (Go service)
  - `db` (Postgres)
  - `rdb` (Redis)
- Prune unused Docker images.

You can check running containers:

make psThe backend will listen on `:${SERVER_PORT}` as defined in `.env`.

### 2.2 Running Locally (without Docker)

1. Ensure PostgreSQL and Redis are running and accessible.
2. Set required environment variables (or `.env` file).
3. Run:

cd Backend
go run main.goOn startup, the service will:

- Initialize configuration (`bootstrap.Run()`).
- Wire dependencies via `wire.InitializeApplication`.
- Run `AutoMigrate` on entities.
- Run seeders (roles, dummy data).
- Start the Gin HTTP server.

---

## 3. Project Layout Recap (Backend)

For full details, see `ARCHITECTURE.md`, but the key maintenance points:

- `main.go` – process entry point and server startup.
- `bootstrap/` – configuration and constants.
- `internal/`:
  - `presentation/` – controllers, middlewares, routes.
  - `application/` – usecases and service implementations.
  - `domain/` – entities, enums, repository interfaces, exceptions.
  - `infrastructure/` – DB, cache, JWT, SMS, S3, localization, seeds.
- `wire/` – dependency injection configuration (`wire.go`, `wire_gen.go`).
- `docker-compose.yaml` – local stack (backend + db + redis).

---

## 4. Common Maintenance Tasks

### 4.1 Adding a New HTTP Endpoint

1. **Define DTOs (if needed)**

   - Add request/response structs under `internal/application/dto/<module>/`.

2. **Extend Usecase Interface**

   - In `internal/application/usecase/<module>_service.go`, add a new method to the relevant interface (e.g., `FormService`, `UserService`, `CalcService`).

3. **Implement in Service Layer**

   - In `internal/application/service/<module>_service_impl.go`, implement the new method.
   - Use domain repositories (`internal/domain/repository/...`) rather than direct DB calls.

4. **Controller method**

   - In the appropriate controller (`internal/presentation/controller/...`):
     - Add handler function (e.g., `CreateSomething`, `GetSomething`).
     - Validate input (bindings, DTOs).
     - Call the usecase service.
     - Map result to HTTP response format.

5. **Route registration**

   - Update `internal/presentation/routes/*.go` (`admin.go`, `customer.go`, `general.go`) to add the route:
     - Choose the correct group (admin/general/customer).
     - Attach required middlewares (e.g., auth).

6. **Wire changes (if needed)**

   - If you added a _new service_, controller, or repository type:

     - Update `wire/wire.go` provider sets (e.g., `ServiceProviderSet`, `RepositoryProviderSet`, controller sets).
     - Regenerate wire code:

       cd Backend
       wire ./wire 7. **Test the endpoint**

   - Use curl/Postman to verify:
     - Correct request/response formats.
     - Permissions and authentication.
     - Error handling.

---

### 4.2 Modifying the Data Model

1. **Update entity struct**

   - Modify the relevant struct in `internal/domain/entity/*.go`.
   - Add fields with proper tags (e.g., GORM tags, JSON tags).

2. **Adjust enums/logic**

   - If the change affects enums (`internal/domain/enum`) or domain rules, update them accordingly.

3. **Update DTOs and Services**

   - Ensure DTOs and services include the new fields where needed.

4. **AutoMigrate**

   - On next startup, `AutoMigrate` in `main.go` will attempt to align the DB schema.
   - For destructive changes (column drops, type changes), you may need manual migrations or DB scripts.

5. **Seed data**
   - If seed data must change, update `internal/infrastructure/seed/*.go`.

---

### 4.3 Working with Repositories & Database

- Domain repository interfaces live in `internal/domain/repository/postgres` and `internal/domain/repository/redis`.
- Infrastructure implementations are in `internal/infrastructure/repository/postgres` and `internal/infrastructure/repository/redis`.

To add a new repository:

1. **Define the interface** in `internal/domain/repository/postgres` or `redis`.
2. **Implement it** in `internal/infrastructure/repository/...`.
3. **Wire it** in `wire/wire.go` (`RepositoryProviderSet`).
4. **Inject it into services** via constructor parameters and interfaces.

---

### 4.4 Authentication & Authorization

- JWT logic lives under `internal/domain/jwt` and `internal/infrastructure/jwt`.
- Auth middleware: `internal/presentation/middleware/auth.go`.
- Roles and permissions:
  - Entities: `internal/domain/entity/role.go`, `permission.go`.
  - Enums: `internal/domain/enum/role.go`, `permission.go`.

When adding a new **role** or **permission**:

1. Update enums and entities as needed.
2. Update seeders in `internal/infrastructure/seed/role_permission_seed.go`.
3. Adjust auth middleware or controller permission checks.

---

### 4.5 Localization

- Localization code is in `internal/infrastructure/localization` and `internal/domain/localization`.
- Middleware: `internal/presentation/middleware/localization.go`.

To add new messages or languages:

1. Extend language files (e.g., `fa.go`) with new keys/messages.
2. Use the translation service in controllers/services instead of hard-coded strings.

---

### 4.6 Connecting to External Calc Services

- Calc-related logic is in `internal/application/service/calc_service_impl.go` and `internal/application/usecase/calc_service.go`.
- URLs and configuration are sourced from `bootstrap.Env.CalcURL`.

When updating or adding a calc service:

1. Update `bootstrap.Env` (Calc URLs).
2. Adjust `CalcService` implementation to call the new endpoint.
3. Update DTOs and controllers in `internal/presentation/controller/calc`.

---

## 5. Deployment & Operations

### 5.1 Docker Image

- The backend image is defined in `Dockerfile` and referenced in `docker-compose.yaml`:

  image: ghcr.io/famcan-riskassessment/backend:latest
  To build locally:

cd Backend
docker build -t famcan-backend:local .To run with the stack:

docker compose up -dEnsure:

- The external network `famcan-network` exists in your environment (or adjust compose file).
- `.env` in production matches deployment settings.

### 5.2 Logs

- Application logs are printed to stdout/stderr (check container logs).
- For quick view:

docker logs backend-RA -f

### 5.3 Database & Cache Volumes

- PostgreSQL volume: `db-data` → `/var/lib/postgresql`
- Redis volume: `rdb-data` → `/var/lib/redis/data`

These volumes preserve data across container restarts.  
Back up or clean them carefully when doing destructive operations.

---

## 6. Troubleshooting

### 6.1 Backend Container Fails to Start

- Check logs:

  docker logs backend-RA
  Common issues:

- Missing or incorrect `.env` values (DB host, port, credentials).
- Network `famcan-network` not created (create or remove the external network requirement).
- Database not reachable (check `db-RA` container logs).

### 6.2 DB Connection Errors

- Verify Postgres is running:

  docker ps | grep db-RA
  docker logs db-RA

  - Verify credentials in `.env` match what the DB is expecting.

- Check port mappings: `${DB_PORT}:${DB_PORT}`.

### 6.3 Redis Connection Errors

- Verify Redis is running:

  docker ps | grep rdb-RA
  docker logs rdb-RA

  - Check `RDB_HOST`, `RDB_PORT` in `.env`.

### 6.4 Wire / DI Issues

If you modify providers or constructors and see DI errors:

1. Update `wire/wire.go` accordingly (bindings, provider sets).
2. Regenerate wiring:

   cd Backend
   wire ./wire 3. Rebuild/restart backend.

---

## 7. Coding Conventions & Best Practices

- Keep **business logic** in `internal/application/service` and **avoid** putting it in controllers.
- Use **domain repositories** and **interfaces** instead of direct DB access in services.
- Prefer **enum types** for domain concepts instead of magic strings.
- Use **localization** and error types from `internal/domain/exception` instead of ad hoc errors.
- Keep **controllers thin**: input validation, calling services, formatting responses.
- Keep **cross-cutting concerns** (auth, logging, error translation, localization) in middlewares and shared utilities, not duplicated in each handler.

---
