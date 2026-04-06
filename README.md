# Finance Dashboard — Backend API

A production-grade backend service for financial management, built with Go and the Gin framework. Provides REST APIs for transaction tracking, financial analytics, and role-based access control (RBAC).

---

## Live Deployment

| Resource | URL |
| :--- | :--- |
| Production API | https://finance-dashboard-t11s.onrender.com |
| Swagger Documentation | https://finance-dashboard-t11s.onrender.com/docs/index.html |
| Health Check | https://finance-dashboard-t11s.onrender.com/health |

---

## Technology Stack

| Layer | Technology |
| :--- | :--- |
| Language | Go 1.21+ |
| HTTP Framework | Gin-Gonic |
| Database | PostgreSQL 16 via GORM |
| Caching | Redis 7.0 |
| Authentication | JWT (HMAC-SHA256) |
| Documentation | Swagger / OpenAPI 2.0 |
| Logging | Structured JSON (`slog`) |

---

## Quick Start (Local Development)

**Prerequisites:** Go 1.21+, PostgreSQL, Redis

**1. Clone the repository**
```bash
git clone https://github.com/nikhildhankhar124106/Finance_dashboard.git
cd "Finance dashboard/backend"
```

**2. Configure environment variables**

Create a `.env` file in the `backend/` directory:
```env
PORT=8080
APP_ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=financedb

JWT_SECRET=your-secret-key
```

**3. Install dependencies**
```bash
go mod tidy
```

**4. Run the server**
```bash
go run ./cmd/api/main.go
```

**5. Open Swagger UI**

Navigate to: `http://localhost:8080/docs/index.html`

---

## Authentication

This API uses **mock JWT authentication** for demonstration purposes.

**Step 1 — Login to get a token:**

```
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@finance.com"
}
```

**Step 2 — Authorize in Swagger:**

Click the **Authorize** button in Swagger UI and enter:
```
Bearer <your_token_here>
```

**Test Accounts:**

| Role | Email | Access Level |
| :--- | :--- | :--- |
| Admin | admin@finance.com | Full access |
| Analyst | analyst@finance.com | Read-only analytics |
| Viewer | viewer@finance.com | Own transactions only |

---

## API Reference

### Authentication
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :---: |
| POST | `/api/v1/auth/login` | Generate JWT token | No |

### Users
| Method | Endpoint | Description | Min. Role |
| :--- | :--- | :--- | :--- |
| GET | `/api/v1/users` | List all users | Viewer |
| GET | `/api/v1/users/:id` | Get a user by ID | Viewer |
| POST | `/api/v1/users` | Create a new user | Admin |
| PATCH | `/api/v1/users/:id/status` | Activate / deactivate user | Admin |

### Transactions
| Method | Endpoint | Description | Min. Role |
| :--- | :--- | :--- | :--- |
| GET | `/api/v1/transactions` | List transactions (filterable, paginated) | Viewer |
| POST | `/api/v1/transactions` | Create a transaction | Admin |
| PUT | `/api/v1/transactions/:id` | Update a transaction | Admin |
| DELETE | `/api/v1/transactions/:id` | Delete a transaction | Admin |
| GET | `/api/v1/transactions/export` | Export transactions as CSV | Admin |

### Analytics
| Method | Endpoint | Description | Min. Role |
| :--- | :--- | :--- | :--- |
| GET | `/api/v1/summary` | Total income, expenses, balance | Viewer |
| GET | `/api/v1/category-breakdown` | Spending by category | Viewer |
| GET | `/api/v1/monthly-trends` | Income vs expense by month | Viewer |

### System
| Method | Endpoint | Description | Min. Role |
| :--- | :--- | :--- | :--- |
| GET | `/health` | Health check | None |
| DELETE | `/api/v1/system/logs` | Clear system audit logs | Admin |

---

## Query Parameters — Transactions

`GET /api/v1/transactions` supports the following query parameters:

| Parameter | Type | Description | Example |
| :--- | :--- | :--- | :--- |
| `page` | integer | Page number (default: 1) | `?page=2` |
| `page_size` | integer | Results per page (default: 10) | `?page_size=25` |
| `category` | string | Filter by category | `?category=Groceries` |
| `type` | string | Filter by type (`Income` or `Expense`) | `?type=Expense` |
| `date` | string | Filter by exact date (YYYY-MM-DD) | `?date=2026-01-15` |
| `search` | string | Search across category and notes | `?search=rent` |
| `sort` | string | Sort field: `amount`, `date`, `category` | `?sort=amount` |
| `order` | string | Sort direction: `asc` or `desc` | `?order=desc` |

---

## Role-Based Access Control

Access is enforced via JWT middleware. The `role` claim in the token is validated against each route's allowed roles.

| Endpoint Group | Admin | Analyst | Viewer |
| :--- | :---: | :---: | :---: |
| Auth (login) | Yes | Yes | Yes |
| View Users | Yes | Yes | Yes |
| Manage Users | Yes | No | No |
| View Transactions | Yes | Yes | Yes |
| Create / Edit / Delete Transactions | Yes | No | No |
| Export CSV | Yes | No | No |
| View Analytics | Yes | Yes | Yes |
| Delete System Logs | Yes | No | No |

---

## Error Handling

All error responses follow a consistent JSON format:

```json
{
  "error": "description of the error"
}
```

| HTTP Code | Meaning |
| :--- | :--- |
| 200 / 201 | Success / Created |
| 400 | Bad Request — invalid input or missing fields |
| 401 | Unauthorized — missing or invalid JWT token |
| 403 | Forbidden — insufficient role permissions |
| 404 | Not Found — resource does not exist |
| 429 | Too Many Requests — rate limit exceeded (100 req/min per IP) |
| 500 | Internal Server Error |

---

## Project Structure

```
backend/
├── cmd/
│   └── api/              # Application entry point (main.go)
├── config/               # Environment variable loading
├── domain/
│   └── models/           # GORM database models
├── handler/
│   ├── api/              # HTTP request handlers
│   └── middleware/       # Auth, RBAC, CORS, rate limiter, error handler
├── infrastructure/
│   ├── db/               # PostgreSQL connection
│   └── cache/            # Redis connection
├── repository/           # Database queries (data access layer)
├── service/              # Business logic layer
├── pkg/
│   ├── auth/             # JWT utilities
│   ├── apperrors/        # Standardized error types
│   └── logger/           # Structured slog setup
├── docs/                 # Swagger auto-generated files
└── tests/                # Integration tests
```

---

## Production Environment Variables (Render)

| Variable | Description |
| :--- | :--- |
| `PORT` | Server port (auto-set by Render) |
| `APP_ENV` | Set to `production` to enable Gin release mode |
| `DB_URL` | Full Neon PostgreSQL connection string (`sslmode=require`) |
| `JWT_SECRET` | Secret key for signing JWT tokens |

---

## Running Tests

```bash
go test ./tests/...
```
