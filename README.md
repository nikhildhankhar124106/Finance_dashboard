<div align="center">
  <h1>📊 Finance Dashboard API</h1>
  <p>A production-ready Golang REST architecture backing a high-performance Financial Analytics Dashboard.</p>
</div>

---

## 📖 Project Overview
This is a production-grade backend service built with **Clean Architecture**. It features deep internal auditing, strict role-based access control (RBAC), and professional financial data processing.

### Key Enhancements:
- **Audit Trail**: Every data mutation is recorded for security and compliance.
- **User Status Management**: Real-time account activation/deactivation.
- **Advanced Data Processing**: Search, multi-field sorting, and CSV export.
- **Resilient Infrastructure**: Integrated rate limiting and centralized error handling.

---

## 🛠️ Tech Stack
- **Language**: Go 1.21+
- **Framework**: [Gin-Gonic](https://github.com/gin-gonic/gin)
- **Database**: PostgreSQL (GORM)
- **Security**: JWT Authentication + custom RBAC Middleware
- **Documentation**: Swagger OpenAPI
- **Testing**: Native Go `testing` + `testify`

---

## 🔐 Roles & Access Control
Access is determined by the `role` payload in the JWT.

| Role | Permissions |
| :--- | :--- |
| **Admin** | Full system control. User management, total transaction visibility, data export. |
| **Analyst** | Read-only access to global financial metrics, summaries, and trends. |
| **Viewer** | Personal finance tracking. Restricted to their own `user_id` context. |

### 🧪 Test Credentials
| Role | Email |
| :--- | :--- |
| **Admin** | `admin@test.com` |
| **Analyst** | `analyst@test.com` |
| **Viewer** | `viewer@test.com` |

---

## 🚀 Setup Instructions

#### 1. Pre-Requisites
- **Go** installed.
- **PostgreSQL** running.

#### 2. Local Setup
```bash
git clone https://github.com/nikhildhankhar124106/Finance_dashboard.git
cd "Finance dashboard/backend"
go mod tidy
```

#### 3. Environment Config (`.env`)
```env
PORT=8080
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=financedb
DB_PORT=5432
APP_ENV=development
```

#### 4. Run Application
```bash
# Optional: Seed the database with sample data
go run ./cmd/seed

# Start the API server
go run ./cmd/api
```

---

## 📚 API Features & Documentation

### Documentation
- **Swagger GUI**: `http://localhost:8080/docs/index.html`

### Advanced Transactions (`GET /api/v1/transactions`)
Supports complex queries:
- `search`: Searches category and notes.
- `sort`: `amount`, `date`, or `category`.
- `order`: `asc` or `desc`.
- `page` & `page_size`: Paginated results with `total_pages`.

### Audit & Security
- **Export CSV**: `GET /api/v1/transactions/export` (Download CSV history).
- **User Status**: `PATCH /api/v1/users/:id/status` (Deactivated users are blocked from all endpoints instantly).
- **Rate Limit**: Strictly enforced at **100 requests per minute** per IP.

---

## 🧪 Testing
Run the integration test suite to verify auth and status logic:
```bash
go test ./tests/...
```
