# Deployment Guide (Golang Gin + PostgreSQL)

This guide walks you through deploying your Golang Finance Dashboard backend to **Railway** (Recommended for easiest PostgreSQL integration) or **Render**.

---

## 🚀 Option 1: Railway (Recommended)
Railway uses Nixpacks, which automatically detects your `go.mod` and builds your Go backend without needing a `Dockerfile`.

### 1. Database Setup
1. Log into [Railway.app](https://railway.app/).
2. Click **New Project** -> **Provision PostgreSQL**.
3. Railway will immediately spin up a Postgres database. 

### 2. Backend Setup
1. In the same project dashboard, click **New** -> **GitHub Repo**.
2. Select your `Finance_dashboard` repository.
3. Railway will analyze the repo and detect the Golang backend.

### 3. Environment Variables Setup
Once your backend service is added, navigate to its **Variables** tab. Add the following variables to link your application with the provisioned PostgreSQL instance:

| Variable | Value (Reference) |
|----------|-------------------|
| `PORT` | `8080` (Railway provides the public `$PORT` dynamically, Go respects this) |
| `APP_ENV` | `production` |
| `DB_HOST` | `${{Postgres.PGHOST}}` |
| `DB_PORT` | `${{Postgres.PGPORT}}` |
| `DB_USER` | `${{Postgres.PGUSER}}` |
| `DB_PASSWORD` | `${{Postgres.PGPASSWORD}}` |
| `DB_NAME` | `${{Postgres.PGDATABASE}}` |

*Note: Railway's `${{...}}` syntax automatically pulls credentials from your Postgres service!*

### 4. Deploy
Once variables are saved, Railway will automatically trigger a redeployment. Go to the **Settings** tab of your backend service, explicitly hit **Generate Domain**, and your API will be live!

---

## ☁️ Option 2: Render.com

### 1. Database Setup
1. Log into [Render](https://render.com/).
2. Click **New +** -> **PostgreSQL**.
3. Name it `financedb`, select a region, and click **Create Database**.
4. Retrieve the **Internal Database URL** and explicit credentials from the database dashboard once it provisions.

### 2. Backend Setup
1. Click **New +** -> **Web Service**.
2. Connect your GitHub repository.
3. Configure the Build Settings:
   - **Environment:** `Go`
   - **Build Command:** `cd backend && go build -o bin/api ./cmd/api`
   - **Start Command:** `cd backend && ./bin/api`

### 3. Environment Variables Setup
Under **Advanced** -> **Environment Variables**, add the extracted individual credentials from your Postgres database dashboard:
- `PORT`: `8080`
- `APP_ENV`: `production`
- `DB_HOST`: *(From Render DB Dashboard, usually an internal hostname like `dpg-...`)*
- `DB_USER`: *(From Render DB Dashboard)*
- `DB_PASSWORD`: *(From Render DB Dashboard)*
- `DB_NAME`: `financedb`
- `DB_PORT`: `5432`

Click **Save and Deploy**. Render will spin up the server.

---

## 🛑 Common Deployment Issues

#### 1. "Failed to auto migrate database schema / Connection Refused"
- **Cause:** The backend is starting before PostgreSQL is fully initialized, or variables mismatch.
- **Fix:** Double-check your `DB_HOST` variables. Ensure no SSL errors block it (You may need to append `sslmode=disable` or `sslmode=require` if Render forces SSL).

#### 2. Swagger Docs returning 404
- **Cause:** Swaggo runs locally generating `docs/`, which are checked into git. If `docs/` wasn't committed, production servers won't have the files.
- **Fix:** Ensure you run `swag init -g cmd/api/main.go -d .` locally and commit the `backend/docs` folder before pushing.

#### 3. "Port already in use" or "Server forced to shutdown"
- **Cause:** The hosting platform explicitly binds ports on Docker containers through an environment variable (often `$PORT`), but your `config.go` might forcefully default to `8080`.
- **Fix:** Our current architecture leverages `getEnv("PORT", "8080")`. Make sure you aren't hardcoding a string inside `main.go`. Platforms will gracefully bind correctly.

#### 4. Environment Variables Returning Empty
- **Cause:** The `godotenv` package logs a missing `.env` warning and fails to capture OS-level variables if written too strictly.
- **Fix:** Our `config.go` already implements fallback `os.LookupEnv()` tracking—so when the `.env` file logically goes missing on Render (it's ignored securely via `.gitignore`), the platform injected variables seamlessly take over without fault.
