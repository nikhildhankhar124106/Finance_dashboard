# Finance Dashboard API

A professional-grade financial management backend developed with Go (Golang). This system provides a robust architecture for real-time transaction tracking, multi-dimensional financial analytics, and multi-layered access control.

## Deployment Status

*   **Production URL**: [https://finance-dashboard-t11s.onrender.com](https://finance-dashboard-t11s.onrender.com)
*   **API Documentation**: [https://finance-dashboard-t11s.onrender.com/docs/index.html](https://finance-dashboard-t11s.onrender.com/docs/index.html)
*   **System Health**: [https://finance-dashboard-t11s.onrender.com/health](https://finance-dashboard-t11s.onrender.com/health)

## Core Technical Features

### 1. Advanced Persistence Layer
*   **GORM Integration**: Seamless Postgres ORM mapping with automated schema migrations.
*   **Soft Deletes**: Native support for record recovery via state-based deletion.
*   **Global Search**: High-performance PostgreSQL pattern matching (ILIKE) across categories and notes.
*   **Sequence Synchronization**: Automatic primary key sequence alignment to prevent primary key conflicts in cloud environments.

### 2. Analytics & Reporting
*   **Multi-Level Aggregation**: Dynamic endpoints for summary metrics, category breakdowns, and monthly financial trends.
*   **Performance Cache**: Redis-backed aggregation layer with transparent database fallback for maximum reliability.

### 3. Security & Access Control
*   **JWT Authentication**: Secure Bearer token implementation using industry-standard HMAC-SHA256 signing.
*   **RBAC (Role-Based Access Control)**: Granular permission enforcement at the route level.

| Role | Permissions | Description |
| :--- | :--- | :--- |
| **Admin** | Full Access | Can manage users, transactions, system logs, and view all analytics. |
| **Analyst** | Read + Analytics | Access to all transaction data and full dashboard summaries. |
| **Viewer** | Read Scoped | Restricted to viewing their own transactions and basic summary data. |

## Technical Architecture

### Middleware & Interceptors
The application utilizes a sequence of reusable middleware layers for consistent request processing:
*   **Audit Layer**: Automatically logs all write operations (POST/PUT/PATCH/DELETE) and system mutations to the `ActivityLog` table.
*   **Rate Limiting**: Protects against automated abuse using token bucket algorithms.
*   **Auth Orchestrator**: Extracts JWT claims and populates the request context with verified `user_id` and `role`.

### Validation Strategy
All incoming payloads are strictly validated using the `validator/v10` package via Gin's `binding` tags. This ensures that:
*   Required fields are present.
*   Data types (amounts, dates) conform to financial precision standards.
*   Enums (Transaction Type, Role) are strictly enforced.

## Environment Configuration

The application is fully configurable via environment variables. For production (Render), ensure the following are set:

| Variable | Requirement | Description |
| :--- | :--- | :--- |
| `APP_ENV` | Optional | Set to `production` to enable Release Mode and SSL. |
| `DB_URL` | **Required** | The Neon PostgreSQL connection string (supports `sslmode=require`). |
| `JWT_SECRET` | **Required** | The secure key used for signing JWT tokens. |
| `REDIS_URL` | Optional | Connection string for distributed caching. |
| `PORT` | Dynamic | Render dynamically binds this to the server listener. |

## Local Installation

1. Clone the repository and navigate to the backend directory.
2. Initialize environment: `cp .env.example .env`.
3. Download dependencies: `go mod tidy`.
4. Launch local development server: `go run ./cmd/api/main.go`.
5. Access Swagger UI at `http://localhost:8080/docs/index.html`.
