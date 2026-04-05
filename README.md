# Finance Dashboard Service

A production-ready financial management system built with Go (Golang). This service provides a robust architectural foundation for real-time transaction tracking, multi-dimensional analytics, and automated role-based access control.

## Deployment Status

*   **Production Environment**: [https://finance-dashboard-t11s.onrender.com](https://finance-dashboard-t11s.onrender.com)
*   **API Documentation (Swagger)**: [https://finance-dashboard-t11s.onrender.com/docs/index.html](https://finance-dashboard-t11s.onrender.com/docs/index.html)
*   **System Health**: [https://finance-dashboard-t11s.onrender.com/health](https://finance-dashboard-t11s.onrender.com/health)

## Architecture Overview

The system is developed using Domain-Driven Design (DDD) principles, ensuring clear separation of concerns between business logic, persistence, and external communication layers.

### Technical Stack
*   **Runtime**: Go 1.21+ (Gin Web Framework)
*   **Persistence**: PostgreSQL 16 (GORM ORM)
*   **Cache**: Redis 7.0 (Distributed caching for analytics)
*   **Security**: JWT (HMAC-SHA256), token-bucket rate limiting
*   **Observability**: Structured JSON logging (slog)

### Directory Structure
```text
├── cmd/api           # Application entry point
├── domain/models     # GORM schemas and domain entities
├── handler/api       # HTTP handlers and request orchestration
├── handler/middleware # Auth, RBAC, Rate-limiting, Audit Logging
├── infrastructure    # Database and Redis connection adapters
├── repository        # Data persistence implementation
├── service           # Core business logic and orchestration
└── pkg               # Shared utilities (Auth, Logging)
```

## Security & Access Control

Access control is enforced globally via a dedicated JWT middleware that extracts claims and validates roles against a persistent user store.

### Permission Matrix

| Endpoint | Method | Admin | Analyst | Viewer |
| :--- | :--- | :---: | :---: | :---: |
| `/api/v1/users` | POST | Read/Write | x | x |
| `/api/v1/users` | GET | Read | Read | Read |
| `/api/v1/transactions` | POST | Read/Write | x | x |
| `/api/v1/transactions` | GET | Read | Read | Read (Scoped) |
| `/api/v1/summary` | GET | Read | Read | Read |
| `/api/v1/system/logs` | DELETE| Read/Write | x | x |

## Core Technical Features

### Advanced Persistence
*   **State-Based Deletion**: Records utilize GORM `DeletedAt` for soft-deletion recovery and audit integrity.
*   **Global Searching**: Performance-optimized `ILIKE` pattern matching across transaction categories and notes.
*   **Sequence Synchronization**: Automated primary key sequence alignment on startup to prevent ID conflicts in managed cloud environments.

### Transaction Management
*   **Advanced Filtering**: Supports multi-parameter filtering (category, type, date), dynamic sorting, and offset-based pagination.
*   **Analytics Aggregation**: Real-time spending analysis across categories and monthly trends with Redis-backed optimized lookups.

## Environment Configuration

The application implements a robust configuration layer that prioritizes system environment variables for production security while providing local development defaults.

### Production (Neon/Render)
The following variables are mandatory for production environments:

| Variable | Requirement | Description |
| :--- | :--- | :--- |
| `DB_URL` | Required | Standard connection string for PostgreSQL (supports `sslmode=require`). |
| `JWT_SECRET` | Required | Secure key for HMAC-SHA256 token signing. |
| `APP_ENV` | Optional | Set to `production` to enable Release Mode. |
| `PORT` | Dynamic | Port binding for the Render service listener. |

### Local Development
Defaults are provided for local `localhost` development:
*   `DB_PORT`: 5432
*   `DB_NAME`: financedb
*   `REDIS_PORT`: 6379

## Local Installation

1. Clone the repository and navigate to the backend directory.
2. Initialize environment configuration: `cp .env.example .env`.
3. Download dependencies: `go mod tidy`.
4. Execute the binary: `go run ./cmd/api/main.go`.
5. Access the interactive documentation at `http://localhost:8080/docs/index.html`.
