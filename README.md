# Finance Dashboard API

A professional backend service for financial management built with Go (Golang). This application provides a stable foundation for transaction tracking, financial analytics, and role-based access control.

## Deployment Status

*   **Production API**: [https://finance-dashboard-t11s.onrender.com](https://finance-dashboard-t11s.onrender.com)
*   **Interactive Documentation (Swagger)**: [https://finance-dashboard-t11s.onrender.com/docs/index.html](https://finance-dashboard-t11s.onrender.com/docs/index.html)
*   **System Health**: [https://finance-dashboard-t11s.onrender.com/health](https://finance-dashboard-t11s.onrender.com/health)

## 📖 Project Overview
This is a production-grade backend service built with **Clean Architecture**. It features deep internal auditing, strict role-based access control (RBAC), and professional financial data processing.

### Key Enhancements:
- **Audit Trail**: Every data mutation is recorded for security and compliance.
- **User Status Management**: Real-time account activation/deactivation.
- **Advanced Data Processing**: Search, multi-field sorting, and CSV export.
- **Resilient Infrastructure**: Integrated rate limiting and centralized error handling.

## 🛠️ Tech Stack
- **Language**: Go 1.25+ (Gin Framework)
- **Database**: PostgreSQL (GORM)
- **Caching**: Redis 7.0 (Distributed caching for analytics)
- **Security**: JWT Authentication + custom RBAC Middleware
- **Documentation**: Swagger OpenAPI
- **Logging**: Structured JSON Logging (slog)
- **Testing**: Native Go `testing` + `testify`

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
| **Admin** | `admin@finance.com` |
| **Analyst** | `analyst@finance.com` |
| **Viewer** | `viewer@finance.com` |

## Core Features
1. **User & Role Management**: Secure JWT auth with pre-seeded accounts.
2. **Financial Records (CRUD)**: Full transaction management with advanced filtering (category, type, date).
3. **Analytics & Trends**: Financial summary metrics and spending distribution.
4. **Audit & Stability**: Soft deletes for integrity and automatic sequence synchronization.

## 📚 API Features & Documentation
- **Swagger GUI**: `http://localhost:8080/docs/index.html`
- **Advanced Transactions**: Supports `search`, `sort`, `order`, and `pagination`.
- **Export CSV**: Download complete transaction history.
- **Rate Limit**: Strictly enforced at 100 requests per minute per IP.

## Local Installation

#### 1. Pre-Requisites
- **Go** installed.
- **PostgreSQL** & **Redis** running.

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
go run ./cmd/api/main.go
```
