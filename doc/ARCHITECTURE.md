## Backend Architecture

### Overview

The backend is a Go service that exposes a REST API for user management, questionnaire/forms, risk calculations, and logging.  
It follows a **layered / clean architecture**:

- **Presentation layer**: HTTP handlers, routing, middlewares (Gin).
- **Application layer**: Use cases and service implementations (business workflows).
- **Domain layer**: Core entities, interfaces, enums, and domain rules.
- **Infrastructure layer**: Database, cache, external services (SMS, S3, JWT, localization), seeds.

The backend is designed so that business rules are decoupled from concrete infrastructure (DB, HTTP, etc.) and wired together via dependency injection.

---

## Technology Stack

- **Language**: Go
- **Web framework**: `gin-gonic/gin`
- **DI / wiring**: `google/wire`
- **Database**: PostgreSQL (GORM-style models and AutoMigrate)
- **Cache**: Redis
- **Authentication / Authorization**:
  - JWT (RSA keys in `internal/infrastructure/jwt`)
  - Role and permission entities
- **File storage**: S3-compatible storage
- **Messaging**: SMS provider (Asanak)
- **Localization**: Internal translation service (e.g., Persian messages)

---

## High-Level Architecture

### Main Entry Point

- `main.go` is the main entry point.
- It performs:
  1. **Gin engine creation**: `ginEngine := gin.New()`
  2. **Configuration bootstrap**: `config := bootstrap.Run()`
  3. **Application wiring**: `app, err := wire.InitializeApplication(config)`
  4. **Database migrations**: `AutoMigrate` on core entities.
  5. **Seeding**: Roles and dummy data.
  6. **Route registration**: `routes.Run(ginEngine, app)`
  7. **Server start**: `ginEngine.Run(:PORT)`

### Configuration Bootstrap

- `bootstrap/`:
  - `bootstrap.go`: defines `Config` and `Run()` which initializes:
    - `Constants` (paths, templates, JWT key locations, etc.)
    - `Env` (DB, Redis, OTP, S3, SMS, calc URLs, pagination, etc.)
  - `env.go`, `constants.go`: load environment variables and constants.

The bootstrap configuration is later used by the DI layer (`wire`) to wire databases, external services, and app components.

---

## Dependency Injection (google/wire)

- `wire/wire.go` defines how the application is assembled.
- It uses **provider sets** to wire interfaces to concrete implementations:

  - **DatabaseProviderSet**:
    - Creates Postgres and Redis clients.
    - Binds them to `database.Database` and `database.Cache` interfaces.
  - **RepositoryProviderSet**:
    - Binds infrastructure repositories (Postgres, Redis) to domain repository interfaces.
  - **ServiceProviderSet**:
    - Application services: user, form, calc, JWT, OTP, action log.
    - Each service implements a corresponding `usecase` interface.
  - **ControllerProviderSets**:
    - General, Admin, Customer controllers for different API surfaces.
  - **AdapterProviderSet**:
    - JWT key manager, translation service, S3 storage, SMS service.
  - **MiddlewareProviderSet**:
    - CORS, recovery, localization, auth middlewares.
  - **SeedProviderSet**:
    - Role and dummy data seeders.

- `wire.Application` aggregates:
  - `Database` (DB + cache)
  - `Middlewares`
  - `Controllers`
  - `Seeds`

`InitializeApplication(config)` builds a fully wired `Application` instance that `main.go` uses.

---

## Project Structure (Backend)

Only the backend directories are described here.

- **`bootstrap/`**

  - `bootstrap.go`: top-level `Config` struct and `Run()` entry.
  - `env.go`, `constants.go`: environment and constants configuration.

- **`internal/`**

  - **`application/`**
    - `dto/`: Data transfer objects (request/response structs) organized by domain:
      - `actionLog/`, `calc/`, `form/`, `general/`, `user/`.
    - `service/`: Concrete service implementations (business logic).
      - `user_service_impl.go`, `form_service_impl.go`, `calc_service_impl.go`, `jwt_service_impl.go`, `otp_service_impl.go`, `action_log_impl.go`.
    - `usecase/`: Service interfaces defining use cases.
      - `UserService`, `FormService`, `CalcService`, `JwtService`, `OtpService`, `ActionLogService`.
  - **`domain/`**
    - `entity/`: Core business entities (GORM models).
      - `user.go`, `role.go`, `permission.go`, `form.go`, `action_log.go`, etc.
    - `enum/`: Enumerations for domain concepts.
      - Roles, permissions, cancer types, gender, form status, life status, etc.
    - `repository/`:
      - `postgres/`: Repository interfaces for Postgres.
      - `redis/`: Repository interfaces for Redis cache.
    - `jwt/`: JWT-related domain logic.
    - `communication/`: Abstract SMS interfaces.
    - `localization/`: Interfaces and types for localization.
    - `storage/`: S3 storage interfaces.
    - `exception/`: Domain-level error types (auth, binding, conflict, validation, not-found, etc.).
  - **`infrastructure/`**
    - `database/`:
      - Postgres and Redis clients, configuration, model mapping.
    - `repository/`:
      - `postgres/`: concrete implementations of domain Postgres repositories.
      - `redis/`: concrete implementations of domain Redis repositories.
    - `jwt/`: Key management for JWT (RSA keys).
    - `communication/sms/`: Asanak SMS implementation.
    - `localization/`: Implementation of translation service and language resources (e.g., `fa.go` for Persian).
    - `storage/`: S3 storage implementation.
    - `seed/`: Role/permission seeding and dummy data seeding.
  - **`presentation/`**
    - `controller/`:
      - `user/`: admin/general user endpoints.
      - `form/`: admin/customer/general form endpoints.
      - `calc/`: admin endpoints for calculations (calls Calc services).
      - `action_log/`: admin endpoints.
      - `base_controller.go`, `response.go`, `pagination.go`, `validator.go`: shared HTTP utilities (responses, pagination, validation).
    - `middleware/`:
      - `auth.go`: JWT auth and permission checks.
      - `cors.go`: CORS policy.
      - `localization.go`: localization setup per request.
      - `recovery.go`: panic handling and error responses.
    - `routes/`:
      - `route.go`, `admin.go`, `customer.go`, `general.go`: route registration, grouping, and middleware binding.

- **`wire/`**

  - `wire.go`: DI configuration (providers, bindings).
  - `wire_gen.go`: generated code (by `wire`) that is used by `main.go`.

- **Database and runtime directories**
  - `db-data/`: Postgres data directory (for local Docker / embedded DB).
  - `rdb-data/`: Redis data directory (if used).
  - `tmp/`: Compiled binaries or temporary files.

---

## Request/Response Flow

### 1. Startup Flow

1. `main.go` creates a Gin engine.
2. `bootstrap.Run()` loads configuration (env + constants).
3. `wire.InitializeApplication(config)` builds the application:
   - Connects to Postgres and Redis.
   - Instantiates repositories, services, controllers, middlewares, seeds.
4. Auto-migrations are run on entity structs to keep the DB schema in sync.
5. Seeders run to ensure base data (roles, permissions, dummy records).
6. `routes.Run(ginEngine, app)` registers routes, binding the controllers and middlewares.
7. Gin starts listening on the configured port.

### 2. HTTP Request Lifecycle

1. **Incoming request** hits Gin.
2. **Global middlewares**:
   - CORS
   - Recovery (panic → proper error response)
   - Localization (sets language context)
   - Auth (JWT parsing, permission check; applied per-route/group where needed)
3. **Routing**:
   - Request is dispatched to the appropriate controller method based on URL and HTTP method.
4. **Controller**:
   - Validates input (query params / JSON body).
   - Maps HTTP DTOs to domain/app structures.
   - Calls the appropriate `usecase` service.
5. **Service (application layer)**:
   - Implements business logic.
   - Interacts with repositories (domain interfaces) for DB access.
   - May call external services (Calc API, SMS, S3) through domain interfaces.
6. **Repository (infrastructure)**:
   - Executes DB or cache operations using Postgres/Redis clients.
7. **Response**:
   - Service returns result or error.
   - Controller formats response (using shared `response.go` utilities).
   - Gin serializes and sends JSON (with proper HTTP status).

---

## Data Storage & External Services

### PostgreSQL

- Core entities (user, role, permission, form, logs, various calculation results) are stored in Postgres.
- On startup, `AutoMigrate` is run for all entity structs to ensure tables/columns exist.

### Redis

- Used as a cache (e.g., user cache) to speed up lookups and store short-lived data (such as OTPs, tokens, or frequently used user data).

### S3 Storage

- Files (e.g., attachments or uploads, if used) are stored via S3-compatible storage.
- The domain defines `S3Storage` interface; `internal/infrastructure/storage` implements it with MinIO (`minio-go`).

### SMS Gateway

- `internal/infrastructure/communication/sms` provides integration with Asanak SMS.
- The domain exposes an abstract `SmsService` interface, and the infrastructure layer provides the concrete implementation.

### External Calc Services

- The backend integrates with external calculation APIs (BCRA, Gail, PLCO, Premm5, etc.) through the `CalcService`.
- The URLs and configuration are provided via `bootstrap.Env.CalcURL`.

---

## Cross-Cutting Concerns

### Error Handling

- Domain-level error types (`exception/`) model different HTTP error categories (auth error, binding error, conflict, forbidden, validation, not found, rate limit).
- Controllers and middlewares transform these into HTTP responses (status code + JSON error body).

### Authentication & Authorization

- JWT tokens are signed/verified using RSA keys in `internal/infrastructure/jwt`.
- Auth middleware validates tokens and injects user identity/roles into the request context.
- Role/permission checks are enforced at the controller or middleware level using domain enums.

### Localization

- The localization middleware sets the current language (e.g., via headers).
- Translation service in `internal/infrastructure/localization` provides localized messages keyed by domain codes.

---
