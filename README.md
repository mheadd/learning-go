[![Go CI](https://github.com/mheadd/learning-go/actions/workflows/test.yml/badge.svg)](https://github.com/mheadd/learning-go/actions/workflows/test.yml)

[![Go Security Check](https://github.com/mheadd/learning-go/actions/workflows/gosec.yml/badge.svg)](https://github.com/mheadd/learning-go/actions/workflows/gosec.yml)

# Learning Go: Experiments in AI Assisted Coding

This is a simple REST API built with Go and the [Gin web framework](https://gin-gonic.com/). It demonstrates basic routing, JSON handling, HTTP methods, and PostgreSQL database integration. The app is containerized with Docker and uses a configuration file for easy setup.

## Features
- RESTful API for managing users
- PostgreSQL database for persistent storage
- Static landing page with a simple web interface
- Secure input handling and SQL injection protection
- Automated tests for API endpoints
- Easy configuration via `config.json`
- Containerized with Docker and Docker Compose

## Prerequisites
- Docker and Docker Compose
- Go (for running tests locally)
- PostgreSQL client (optional, for DB inspection)

## Getting Started

### 1. Clone the Repository
```bash
git clone <repo-url>
cd learning-go
```

### 2. Configuration
If desired, edit `config.json` to set database and app port settings (otherwise the app just uses the default values):
```json
{
  "db_host": "db",
  "db_user": "postgres",
  "db_password": "postgres",
  "db_name": "usersdb",
  "db_port": "5432",
  "app_port": "8080"
}
```

### 3. Build and Start the App with Docker Compose
```bash
docker compose up --build
```
- The Go app will be available at [http://localhost:8080](http://localhost:8080)
- The PostgreSQL database will be available at `localhost:5432` (from your host)

### 4. Web Interface
Visit [http://localhost:8080](http://localhost:8080) in your browser to see the landing page. You can:
- View available API endpoints
- Add a new user via the form
- List all users

## API Endpoints

### Health Check
```bash
curl http://localhost:8080/health
```
Response:
```json
{"status": "healthy"}
```

### Users API
- **Get all users:**
  ```bash
  curl http://localhost:8080/api/users
  ```
- **Create a new user:**
  ```bash
  curl -X POST http://localhost:8080/api/users \
    -H "Content-Type: application/json" \
    -d '{"id": "3", "name": "Bob Wilson"}'
  ```

## Database Initialization
- The database schema is defined in `init.sql` and is run automatically when the app starts.
- Example table:
  ```sql
  CREATE TABLE IF NOT EXISTS users (
      id VARCHAR(50) PRIMARY KEY,
      name VARCHAR(100) NOT NULL
  );
  ```

## Running Tests

You can run automated tests for the API endpoints using Go:

1. Make sure the database is running (e.g., via Docker Compose)
2. Run the tests:
   ```bash
   go test -v
   ```
- Tests use the same database as the app (by default `usersdb` on localhost:5432). The `users` table is cleaned before each test.
- Tests cover health check, user creation, and user retrieval.

## Security
- All SQL queries use prepared statements to prevent SQL injection.
- User input is validated for required fields and length.
- Internal errors are logged server-side; generic error messages are returned to clients.

## Customization
- Edit `config.json` to change DB or port settings.
- Edit `init.sql` to change the database schema.
- Edit `static/` for the web interface.

---

For any issues or questions, please open an issue or contact the maintainer.
