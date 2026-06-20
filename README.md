<div align="center">

# 🎬 Ozinse Backend API

**Robust backend service for the Ozinse Video Platform**

[![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin_Gonic-Framework-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://gin-gonic.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Swagger](https://img.shields.io/badge/Swagger-Docs-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)](https://swagger.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)

Built with **Go (Golang)**, following a clean architectural pattern that separates concerns into handlers, services, and repositories — ensuring scalability and maintainability.

</div>

---

## 🚀 Technology Stack

| Category | Technology |
|---|---|
| 🧠 **Language** | Go (Golang) |
| 🌐 **Web Framework** | Gin Gonic |
| 🗄️ **Database** | PostgreSQL |
| 📖 **Documentation** | Swagger (`swag`) |
| 📝 **Logger** | Custom structured JSON logger |

---

## 📂 Project Architecture

The project is organized into clean logical layers:

```
cmd/api/                 → Application entry point
internal/handler/        → HTTP layer — handles requests and parses inputs
internal/logger/         → Custom structured JSON logger
internal/middleware/     → HTTP middleware (auth guards, logging, recovery, etc.)
internal/model/          → Data structures and DTOs
internal/repository/     → Data access layer — interacts directly with PostgreSQL
internal/service/        → Business logic layer — the "brain" of the app
```

<details>
<summary>📁 <strong>Full folder structure</strong></summary>

```
internal/
├── handler/
│   ├── age_category_handler.go
│   ├── auth_handler.go
│   ├── category_handler.go
│   ├── favorite_handler.go
│   ├── genre_handler.go
│   ├── project_handler.go
│   ├── role_handler.go
│   └── user_handler.go
├── logger/
│   └── logger.go
├── middleware/
├── model/
├── repository/
│   ├── age_category_repository.go
│   ├── category_repository.go
│   ├── favorite_repository.go
│   ├── genre_repository.go
│   ├── postgres.go
│   ├── project_repository.go
│   ├── role_repository.go
│   └── user_repository.go
└── service/
    ├── age_category_service.go
    ├── auth_service.go
    ├── category_service.go
    ├── favorite_service.go
    ├── genre_service.go
    ├── project_service.go
    ├── role_service.go
    └── user_service.go
```

</details>

---

## 🌐 Handler Layer Breakdown

The `internal/handler/` layer parses incoming HTTP requests and delegates to the matching service:

<table>
<tr><td>🔒</td><td><code>age_category_handler.go</code></td><td>Endpoints for managing age categories and content rating rules.</td></tr>
<tr><td>🔐</td><td><code>auth_handler.go</code></td><td>Endpoints for registration, login, and JWT-based authentication.</td></tr>
<tr><td>🗂️</td><td><code>category_handler.go</code></td><td>Endpoints for creating, updating, and listing video categories.</td></tr>
<tr><td>⭐</td><td><code>favorite_handler.go</code></td><td>Endpoints for adding, removing, and listing a user's favorite projects.</td></tr>
<tr><td>🎭</td><td><code>genre_handler.go</code></td><td>Endpoints for managing genres used to tag content.</td></tr>
<tr><td>🎥</td><td><code>project_handler.go</code></td><td>Endpoints for creating and managing movies/series (projects).</td></tr>
<tr><td>🛡️</td><td><code>role_handler.go</code></td><td>Endpoints for managing user roles and permissions.</td></tr>
<tr><td>👤</td><td><code>user_handler.go</code></td><td>Endpoints for user profile management and account settings.</td></tr>
</table>

---

## 🛠 Service Layer Breakdown

The `internal/service/` layer processes all business logic before data reaches the database:

<table>
<tr><td>🔒</td><td><code>age_category_service.go</code></td><td>Logic for content rating and age-appropriate restriction management.</td></tr>
<tr><td>🔐</td><td><code>auth_service.go</code></td><td>User registration, JWT token generation, login validation, and secure password hashing using <code>bcrypt</code>.</td></tr>
<tr><td>🗂️</td><td><code>category_service.go</code></td><td>Manages video categories — create, update, and retrieve category lists to organize content.</td></tr>
<tr><td>⭐</td><td><code>favorite_service.go</code></td><td>Manages users' favorite projects — adding, removing, and retrieving personalized favorite lists.</td></tr>
<tr><td>🎭</td><td><code>genre_service.go</code></td><td>Handles genre management for dynamic content tagging (e.g., Action, Drama, Comedy).</td></tr>
<tr><td>🎥</td><td><code>project_service.go</code></td><td>Core engine for video content — manages movies and series (projects), including metadata, screenshots, and project indexing.</td></tr>
<tr><td>🛡️</td><td><code>role_service.go</code></td><td>Manages user roles and access levels, controlling permissions across the platform.</td></tr>
<tr><td>👤</td><td><code>user_service.go</code></td><td>Handles user profile management, role assignments, and account settings.</td></tr>
</table>

---

## 🗄️ Repository Layer Breakdown

The `internal/repository/` layer talks directly to PostgreSQL:

<table>
<tr><td>🔒</td><td><code>age_category_repository.go</code></td><td>CRUD operations for age category records.</td></tr>
<tr><td>🗂️</td><td><code>category_repository.go</code></td><td>CRUD operations for video category records.</td></tr>
<tr><td>⭐</td><td><code>favorite_repository.go</code></td><td>Persists and retrieves users' favorite project entries.</td></tr>
<tr><td>🎭</td><td><code>genre_repository.go</code></td><td>CRUD operations for genre records.</td></tr>
<tr><td>🔌</td><td><code>postgres.go</code></td><td>Database connection setup and configuration for PostgreSQL.</td></tr>
<tr><td>🎥</td><td><code>project_repository.go</code></td><td>CRUD operations for projects (movies/series), including metadata and screenshots.</td></tr>
<tr><td>🛡️</td><td><code>role_repository.go</code></td><td>CRUD operations for role and permission records.</td></tr>
<tr><td>👤</td><td><code>user_repository.go</code></td><td>CRUD operations for user accounts and profiles.</td></tr>
</table>

---

## 📚 API Documentation (Swagger)

The API documentation is auto-generated. Once the server is running, explore the endpoints in your browser:

> 👉 **[ozinse-backend.onrender.com/swagger/index.html](https://ozinse-backend.onrender.com/swagger/index.html)**

---

## ⚙️ Development Guide

### ▶️ Quick Start

**1. Clone the repository**
```bash
git clone https://github.com/tihomir-chobanov/ozinse-backend.git
```

**2. Set up your `.env` file**

**3. Run with Docker**
```bash
docker compose up --build
```

### 🔄 Updating Documentation

Whenever you add a new service or endpoint, regenerate the Swagger docs:

```bash
swag init -g cmd/api/main.go -o cmd/api/docs --parseDependency --parseInternal
```

---

## 🔐 Authentication Flow

| Step | Endpoint | Description |
|---|---|---|
| 1️⃣ **Register** | `/api/auth/register` | Create an account |
| 2️⃣ **Login** | `/api/auth/login` | Obtain your JWT token |
| 3️⃣ **Authorize** | Swagger UI | Click **Authorize** and enter `Bearer <your_token>` to access protected endpoints |

---

<div align="center">

✨ **Developed by [Tihomir Chobanov](https://github.com/tihomir-chobanov)** ✨

</div>
