# Finance Dashboard API

A professional backend service for financial management built with Go (Golang). This application provides a stable foundation for transaction tracking, financial analytics, and role-based access control.

## Deployment Status

*   **Production API**: [https://finance-dashboard-t11s.onrender.com](https://finance-dashboard-t11s.onrender.com)
*   **Interactive Documentation (Swagger)**: [https://finance-dashboard-t11s.onrender.com/docs/index.html](https://finance-dashboard-t11s.onrender.com/docs/index.html)
*   **System Health**: [https://finance-dashboard-t11s.onrender.com/health](https://finance-dashboard-t11s.onrender.com/health)

## Technical Stack

*   **Language**: Go 1.21+ (Gin Framework)
*   **Database**: PostgreSQL 16 (GORM ORM)
*   **Caching**: Redis 7.0 (Distributed caching for analytics)
*   **Authentication**: JWT (HMAC-SHA256)
*   **Logging**: Structured JSON Logging (slog)

## Core Features

### 1. User & Role Management
*   Secure authentication using JWT.
*   Pre-seeded accounts for testing (Admin, Analyst, Viewer).
*   Role-based permissions enforced across all endpoints.

### 2. Financial Records (CRUD)
*   Create, Read, Update, and Delete transactions.
*   **Filtering**: Search by category, transaction type, and specific dates.
*   **Pagination**: Efficiently manage large datasets with offset-based paging.
*   **Global Search**: Keyword-based searching across categories and notes using PostgreSQL ILIKE.

### 3. Analytics & Trends
*   **Summary**: Total income, total expenses, and current balance.
*   **Breakdown**: Spending distribution categorized by transaction type.
*   **Monthly Trends**: Periodical data showing financial movement over time.

### 4. System Stability
*   **Soft Deletes**: Records are marked as deleted but preserved for audit integrity.
*   **Sequence Sync**: Automatic database sequence alignment on startup to prevent ID conflicts.

## Security & Access Control

### Role-Based Access Control (RBAC)
Role-based permissions are enforced via custom middleware. The system extracts the user role from the JWT claim and validates it against the allowed roles for each specific route.

| Feature | Endpoint | Admin | Analyst | Viewer |
| :--- | :--- | :---: | :---: | :---: |
| Create User | `POST /api/v1/users` | Yes | No | No |
| List Users | `GET /api/v1/users` | Yes | Yes | Yes |
| Create Transaction | `POST /api/v1/transactions` | Yes | No | No |
| List Transactions | `GET /api/v1/transactions` | Yes | Yes | Yes (Scoped) |
| View Analytics | `GET /api/v1/summary` | Yes | Yes | Yes |
| Delete Logs | `DELETE /api/v1/system/logs` | Yes | No | No |

## Validation & Error Handling

*   **Input Validation**: Request payloads are strictly validated using Gin's `binding:"required"` and the Go Playground validator.
*   **HTTP Status Codes**:
    *   **200/201**: Successful request/creation.
    *   **400**: Bad Request (Invalid input).
    *   **401**: Unauthorized (Missing or invalid token).
    *   **403**: Forbidden (Insufficient role permissions).
    *   **500**: Internal Server Error.
*   **Structured Errors**: All failures return a standardized JSON response: `{"error": "description"}`.

## Project Structure

```text
├── cmd/api           # Application entry point
├── domain/models     # Database schemas and entities
├── handler/api       # HTTP handlers (Request/Response)
├── handler/middleware # Auth, RBAC, and Rate-limiting
├── infrastructure    # Database and Redis adapters
├── repository        # Data persistence (SQL logic)
├── service           # Business logic and processing
└── pkg               # Shared utilities (Auth, Logging)
```

## Environment Configuration

### Production (Render & Neon)
*   **DB_URL**: The Neon PostgreSQL connection string (supports `sslmode=require`).
*   **JWT_SECRET**: The secure signing key for tokens.
*   **APP_ENV**: Set to `production` for release mode.

### Local Development
Defaults are provided for local development:
*   `DB_PORT`: 5432
*   `DB_NAME`: financedb
*   `PORT`: 8080

## Local Installation

1. Clone the repository and navigate to the backend directory.
2. Initialize environment: `cp .env.example .env`.
3. Download dependencies: `go mod tidy`.
4. Run: `go run ./cmd/api/main.go`.
5. Access Swagger at `http://localhost:8080/docs/index.html`.

## Sample Users

| Role    | Email               |
|--------|---------------------|
| Admin   | admin@finance.com   |
| Analyst | analyst@finance.com |
| Viewer  | viewer@finance.com  |

---

## How to Test the API

1. Open Swagger:
https://finance-dashboard-t11s.onrender.com/docs/index.html  

2. Login:

POST /api/v1/auth/login

Example:
{
  "email": "admin@finance.com"
}

3. Copy JWT token

4. Click "Authorize" in Swagger and paste:
Bearer <token>

5. Test endpoints:
- GET /api/v1/transactions
- POST /api/v1/transactions
- GET /api/v1/summary

---

## Deployment Notes

- Uses Neon PostgreSQL with sslmode=require
- Deployed on Render using dynamic PORT binding
