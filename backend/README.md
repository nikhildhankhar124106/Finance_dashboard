# Finance Dashboard API

A professional-grade backend for financial management, built with Go (Golang). This system provides a robust foundation for tracking transactions, analyzing spending trends, and managing user roles with high-performance caching and structured logging.

## Technical Overview

The application follows a clean, Domain-Driven Design (DDD) architecture, separating concerns between handlers, services, and repositories to ensure maintainability and testability.

### Core Features

*   **Role-Based Access Control (RBAC)**: Integrated security layer supporting Admin, Analyst, and Viewer roles with specific endpoint permissions.
*   **Soft Deletes**: Native support for GORM-based soft deletes across User and Transaction models, ensuring data recoverability.
*   **Global Search & Filtering**: Optimized PostgreSQL queries using ILIKE for keyword searches across categories and notes, with dynamic sorting and pagination.
*   **Distributed Caching**: Optional Redis integration for analytics aggregation, with a "Direct-DB" fallback for high availability.
*   **Audit Logging**: Automatic tracking of all write operations (POST/PUT/PATCH/DELETE) and system activities.
*   **Automated Schema Management**: Built-in sequence synchronization and auto-migration to prevent primary key conflicts in cloud environments.

## Technology Stack

*   **Core**: Go 1.21+ (Gin Framework)
*   **Database**: PostgreSQL 16+ (GORM ORM)
*   **Caching**: Redis (go-redis)
*   **Authentication**: JWT (JSON Web Tokens)
*   **Documentation**: Swagger (swaggo)
*   **Logging**: Structured JSON Logging (slog)

## Project Structure

```text
├── cmd/api           # Application entry point (main.go)
├── config            # Environment and configuration mapping
├── domain/models     # Core business entities and GORM schemas
├── docs              # Auto-generated Swagger documentation
├── handler/api       # HTTP request/response orchestration
├── handler/middleware # Auth, Rate-Limiting, and Audit layers
├── infrastructure    # Database and Cache connection logic
├── repository        # Data persistence and SQL logic
├── service           # Pure business logic and orchestration
└── pkg               # Shared utilities (logging, errors)
```

## Getting Started

### Prerequisites

*   Go 1.21 or higher installed.
*   PostgreSQL 15+ instance.
*   Redis (optional, but recommended).

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/nikhildhankhar124106/Finance_dashboard.git
   cd Finance_dashboard/backend
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Configure environment variables:
   Create a `.env` file in the `backend/` directory:
   ```env
   PORT=8080
   DB_HOST=localhost
   DB_USER=postgres
   DB_PASSWORD=yourpassword
   DB_NAME=financedb
   DB_PORT=5432
   APP_ENV=development
   ```

4. Run the application:
   ```bash
   go run ./cmd/api
   ```

## API Documentation

Once the server is running, you can access the interactive API documentation at:
`http://localhost:8080/docs/index.html`

**Authentication Note**: To test protected endpoints, click the "Authorize" button and enter your JWT in the format: `Bearer <your_token>`.

## Deployment

The project is optimized for deployment on **Render** using **Neon PostgreSQL**.

1. Connect your GitHub repository to a Render Web Service.
2. Set `APP_ENV` to `production` (this automatically enables SSL requirements for Neon).
3. Provide the `DATABASE_URL` or individual `DB_*` variables in the Render Dashboard.
4. The application will automatically handle schema migrations and sequence synchronization on startup.
