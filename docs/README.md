# Project Documentation

- [Google Docs: Requirements and Instructions](https://docs.google.com/document/d/1pjZqxMr_BAnvM8rAb7DBCY8F4FOSwhQFCD122E8qVrs/edit?tab=t.0)

# Application Architecture

This project follows a classic Three-Layer Architecture: Handler -> Service -> Repository.

## 1. Handler Layer (Presentation)
**Location:** `/internal/handler`
**Core Purpose:** The Handler translates between HTTP (web browsers, Postman) and Go structs. It ensures the inner layers don't handle web-specific details.

**Key Responsibilities:**
*   **Routing:** Listens for HTTP methods and paths (e.g., `DELETE /api/categories/:id`).
*   **Request Binding:** Converts raw JSON into Go Models (using Gin's `ShouldBindJSON`).
*   **Parameter Extraction:** Pulls variables from URLs (e.g., converting string IDs to integers).
*   **Input Validation:** Performs shallow validation (checks formats and required fields).
*   **Response Generation:** Converts Go structs back to JSON.
*   **Status Codes:** Assigns the correct HTTP status (200 OK, 201 Created, 400 Bad Request, etc.).

## 2. Service Layer (Business Logic)
**Location:** `/internal/service`
**Core Purpose:** The Service layer acts as the "Brain" of the application. It contains all the core business logic and rules. It sits squarely between the Presentation layer (Handler) and the Data Access layer (Repository).

**Key Responsibilities:**
*   **Business Rules & Validation:** Enforces the specific rules of the application. For example, before creating a new category, the Service checks if a category with that name already exists (preventing duplicates).
*   **Orchestration:** Coordinates complex operations. If a single request requires multiple database actions (e.g., deleting a category might require checking for linked projects first), the Service manages this flow.
*   **Decoupling:** It separates the HTTP transport logic from the database storage logic. The Service doesn't know about JSON, HTTP headers, or SQL queries—it only cares about Go structs and business requirements.

## 3. Repository Layer (Data Access)
**Location:** `/internal/repository`
**Core Purpose:** The Repository acts as the "Bridge" or "Supplier" between the application's business logic and the underlying database. It isolates all data access logic, ensuring the rest of the application remains database-agnostic. 

**Key Responsibilities:**
*   **Exclusive Database Access:** It is the *only* layer permitted to communicate with the database (e.g., PostgreSQL). All raw SQL queries and database transactions live strictly here.
*   **Abstraction:** It allows the Service layer to request or save data using simple method calls (e.g., "Get Category with ID 5") without needing to know *how* that data is queried or stored.
*   **Data Translation:** It handles the technical execution of fetching database rows and cleanly mapping them back into Go Models for the rest of the application to use.

## 4. Model Layer (Data Definitions)
**Location:** `/internal/model`

**Core Purpose:** The Model is the simplest but most essential component of the application. It defines the exact shape and structure of the data as it travels through the Handlers, Services, and Repositories, serving as the application's "single source of truth."

**Key Responsibilities:**
*   **Data Structure Definition:** It defines the properties of database entities as Go structs (e.g., a `Category` struct containing `ID` and `Name` fields).
*   **Data Tagging & Serialization:** It utilizes struct tags (like `json:"name"`) to instruct Go on how to format, rename, or hide fields (e.g., `omitempty`) when converting the data to and from JSON for the client.
*   **Cross-Layer Consistency:** It ensures that all layers (Presentation, Business Logic, and Data Access) speak the exact same language. If a database schema changes (e.g., adding a "description" field), you update the Model once, and the entire system instantly knows how to handle the new data.

## Flow Example (Creating a Resource)
1. The **Handler** receives raw JSON, binds it to a Model, and passes it to the Service.
2. The **Service** applies business rules and passes the Model to the Repository.
3. The **Repository** executes the SQL `INSERT` to save the Model.
4. The **Handler** returns the updated Model (with its new ID) as JSON.

---

# Project Structure

```
ozinse-backend/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/                     # Configuration management
│   ├── handler/                    # HTTP handlers (presentation layer)
│   │   ├── category_handler.go
│   │   ├── genre_handler.go
│   │   └── project_handler.go
│   ├── model/                      # Data models/entities
│   │   ├── category.go
│   │   ├── genre.go
│   │   ├── project.go
│   │   ├── episode.go
│   │   ├── season.go
│   │   ├── age_category.go
│   │   ├── screenshot.go
│   │   └── projectScreenshot.go
│   ├── repository/                 # Data access layer
│   │   ├── category_repository.go
│   │   ├── genre_repository.go
│   │   ├── project_repository.go
│   │   └── postgres.go
│   └── service/                    # Business logic layer
│       ├── category_service.go
│       ├── genre_service.go
│       └── project_service.go
├── migrations/                     # Database migrations
│   └── 000001_init_schema.up.sql
├── assets/                         # Static files
│   ├── images/
│   └── pdfs/
├── docs/                           # Documentation
└── go.mod                          # Go module file
```

---

# Setup & Installation

## Prerequisites
- **Go** 1.20 or higher
- **PostgreSQL** 12 or higher
- **Gin** web framework (automatically installed via `go mod`)

## Steps

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd ozinse-backend
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   go mod tidy
   ```

3. **Configure environment:**
   Set up environment variables for database connection:
   ```bash
   export DATABASE_URL="postgres://user:password@localhost:5432/ozinse_db"
   export PORT=8080
   ```

4. **Run database migrations:**
   ```bash
   # Use migration tool of your choice (e.g., migrate, goose, or sql-migrate)
   migrate -path migrations -database $DATABASE_URL up
   ```

5. **Start the server:**
   ```bash
   go run cmd/api/main.go
   ```

6. **Verify the server is running:**
   ```bash
   curl http://localhost:8080/api/categories
   ```

---

# API Endpoints

## Categories
- `GET /api/categories` - Get all categories
- `GET /api/categories/:id` - Get category by ID
- `POST /api/categories` - Create a new category
- `PUT /api/categories/:id` - Update a category
- `DELETE /api/categories/:id` - Delete a category

## Genres
- `GET /api/genres` - Get all genres
- `GET /api/genres/:id` - Get genre by ID
- `POST /api/genres` - Create a new genre
- `PUT /api/genres/:id` - Update a genre
- `DELETE /api/genres/:id` - Delete a genre

## Projects
- `GET /api/projects` - Get all projects
- `GET /api/projects/:id` - Get project by ID
- `POST /api/projects` - Create a new project with genres, categories, and age categories
- `PUT /api/projects/:id` - Update a project
- `DELETE /api/projects/:id` - Delete a project

---

# Models Documentation

## Category
```go
type Category struct {
    ID   int    // Unique identifier
    Name string // Category name (must be unique)
}
```
**Purpose:** Represents content categories (e.g., Animation, Comedy, Drama).

## Genre
```go
type Genre struct {
    ID   int    // Unique identifier
    Name string // Genre name (must be unique)
}
```
**Purpose:** Represents content genres (e.g., Action, Adventure, Horror).

## Project
```go
type Project struct {
    ID          int       // Unique identifier
    Title       string    // Project title (must be unique)
    Description string    // Project description
    ReleaseDate time.Time // Release date
    CreatedAt   time.Time // Creation timestamp
    UpdatedAt   time.Time // Last update timestamp
}
```
**Purpose:** Represents multimedia projects (TV shows, movies, documentaries).
**Associations:** Can be linked to multiple Genres, AgeCategories, and Categories.

## Season
```go
type Season struct {
    ID        int       // Unique identifier
    ProjectID int       // Foreign key to Project
    Number    int       // Season number
    Title     string    // Season title
    CreatedAt time.Time // Creation timestamp
}
```
**Purpose:** Represents seasons within a project.

## Episode
```go
type Episode struct {
    ID        int       // Unique identifier
    SeasonID  int       // Foreign key to Season
    Number    int       // Episode number
    Title     string    // Episode title
    Duration  int       // Duration in minutes
    CreatedAt time.Time // Creation timestamp
}
```
**Purpose:** Represents episodes within a season.

## AgeCategory
```go
type AgeCategory struct {
    ID   int    // Unique identifier
    Name string // Age category name (e.g., "10-12", "12-14")
}
```
**Purpose:** Represents age restrictions for content.

## Screenshot & ProjectScreenshot
```go
type Screenshot struct {
    ID   int    // Unique identifier
    URL  string // Screenshot image URL
}

type ProjectScreenshot struct {
    ID           int // Unique identifier
    ProjectID    int // Foreign key to Project
    ScreenshotID int // Foreign key to Screenshot
}
```
**Purpose:** Manages media assets associated with projects.

---

# Database Schema

The application uses PostgreSQL as the primary data store. Key tables include:

## Tables Structure
- **categories** - Stores category data (ID, Name)
- **genres** - Stores genre information (ID, Name)
- **projects** - Stores project details (ID, Title, Description, ReleaseDate, timestamps)
- **seasons** - Stores season information (ID, ProjectID, Number, Title, timestamps)
- **episodes** - Stores episode information (ID, SeasonID, Number, Title, Duration, timestamps)
- **age_categories** - Stores age restriction data (ID, Name)
- **screenshots** - Stores screenshot URLs (ID, URL)
- **project_genres** - Junction table linking Projects to Genres
- **project_categories** - Junction table linking Projects to Categories
- **project_age_categories** - Junction table linking Projects to AgeCategories
- **project_screenshots** - Junction table linking Projects to Screenshots

**Migration Location:** `/migrations/000001_init_schema.up.sql`

For detailed schema, see the migration file.

---

# Error Handling

## Error Strategy
The application follows a consistent error handling pattern across all layers:

### Handler Layer
- **Input Validation:** Returns `400 Bad Request` for invalid JSON or missing required fields.
- **Not Found:** Returns `404 Not Found` when a resource doesn't exist.
- **Server Errors:** Returns `500 Internal Server Error` for unexpected failures.

### Service Layer
- **Business Logic Errors:** Returns descriptive error messages (e.g., "category with name 'X' already exists").
- **Duplicate Prevention:** Checks for existing records before creation to prevent constraint violations.
- **No Exception Throwing:** Uses Go's error return pattern instead of panics.

### Repository Layer
- **Database Errors:** Catches and wraps SQL errors from PostgreSQL.
- **Connection Issues:** Handles database connection failures gracefully.
- **Query Execution:** Returns raw database errors to Service for processing.

## Common Error Codes
| HTTP Code | Meaning | Example |
|-----------|---------|---------|
| 200 OK | Request successful | Category retrieved |
| 201 Created | Resource created | New project created |
| 400 Bad Request | Invalid input | Missing required field |
| 404 Not Found | Resource not found | Category ID doesn't exist |
| 409 Conflict | Resource conflict | Duplicate category name |
| 500 Internal Server Error | Server error | Database connection failed |

## Error Response Format
```json
{
  "error": "Descriptive error message",
  "status": 400
}
```