# Go Clean Architecture Template

A robust and scalable Go project template following the principles of **Clean Architecture**. This template is designed to provide a solid foundation for building maintainable and testable web applications.

## 🏗️ Architecture Overview

The project is structured into layers to ensure a separation of concerns and maintain a one-way dependency flow towards the core business logic.

```mermaid
graph TD
    subgraph "Infrastructure & Adapters"
        H[HTTP Handlers]
        R[Repositories]
        C[Config & DB]
    end

    subgraph "Application Logic"
        UC[Use Cases]
    end

    subgraph "Core Domain"
        D[Domain Entities]
        RI[Repository Interfaces]
    end

    H --> UC
    UC --> RI
    R -- Implements --> RI
    R --> C
```

### 📁 Project Structure

```text
.
├── cmd/api/main.go          # Application entry point & dependency injection
├── internal/
│   ├── adapter/             # Interface Adapters
│   │   ├── http/            # REST API Handlers & Routers
│   │   └── repository/      # Database implementations (GORM)
│   ├── core/                # Infrastructure & Cross-cutting concerns
│   │   ├── config/          # Configuration loading
│   │   └── database/        # Database connection setup
│   ├── domain/              # Core Business Logic (Entities & Interfaces)
│   │   └── user/            
│   └── usecase/             # Application Business Logic (Orchestrators)
│       └── user/
├── pkg/                     # Shared utilities (Logger, Errors, Responses)
└── .env.example             # Environment variables template
```

---

## 🔄 Communication Flow

This template follows a strict communication flow. Below are diagrams illustrating how a typical request (e.g., "Register User") moves through the system.

### 1. Sequence of Interaction
This diagram shows the execution flow from the moment a client sends a request.

```mermaid
sequenceDiagram
    participant Client
    participant Router as HTTP Router (Adapter)
    participant Handler as User Handler (Adapter)
    participant UseCase as User UseCase (Application)
    participant Repo as User Repository (Adapter)
    participant DB as PostgreSQL (Infrastructure)

    Client->>Router: POST /api/v1/users/register
    Router->>Handler: Register(http.ResponseWriter, *http.Request)
    Handler->>UseCase: Register(ctx, RegisterInput)
    UseCase->>Repo: FindByEmail(ctx, email)
    Repo->>DB: SELECT * FROM users WHERE email = ...
    DB-->>Repo: Result (User or Not Found)
    Repo-->>UseCase: Result
    
    UseCase->>Repo: Create(ctx, newUser)
    Repo->>DB: INSERT INTO users ...
    DB-->>Repo: Success
    Repo-->>UseCase: error (nil)

    UseCase-->>Handler: (*User, error)
    Handler-->>Client: 201 Created (JSON Response)
```

### 2. Dependency & Connection Flow
This diagram shows how the components are physically connected via interfaces.

```mermaid
graph LR
    subgraph "External"
        Client((Client))
    end

    subgraph "Adapter Layer (HTTP)"
        Router[Router]
        Handler[User Handler]
    end

    subgraph "Application Layer"
        UC[User UseCase]
    end

    subgraph "Domain Layer (Core)"
        URI[Repository Interface]
        Entity[User Entity]
    end

    subgraph "Adapter Layer (DB)"
        Repo[User Repository Implementation]
    end

    subgraph "Infrastructure"
        DB[(PostgreSQL)]
    end

    Client -- HTTP Request --> Router
    Router -- Call --> Handler
    Handler -- Invoke --> UC
    UC -- Uses --> URI
    UC -- Returns --> Entity
    Repo -- Implements --> URI
    Repo -- Queries --> DB
```

### Layer Responsibilities:

1.  **Request Layer (Router)**: Handles routing and middleware (Auth, Logging).
2.  **Adapter Layer (Handler)**: Unmarshals JSON request body, validates input, calls the UseCase, and marshals the JSON response.
3.  **UseCase Layer**: Contains the "What" of the application. It orchestrates domain entities and repository interfaces to fulfill a business requirement.
4.  **Domain Layer**: The heart of the app. It defines the business rules (Entities) and the contracts (Interfaces) that other layers must follow.
5.  **Adapter Layer (Repository)**: Handles the "How" of data persistence. It translates domain-agnostic calls into GORM/SQL queries.

---

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- PostgreSQL (or your preferred database supported by GORM)

### Setup
1. Clone the repository.
2. Copy `.env.example` to `.env` and fill in your database credentials.
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

## 🛠️ Features Included
- **Clean Architecture** structure.
- **Dependency Injection** in `main.go`.
- **GORM** for database operations.
- **Graceful Shutdown** of the HTTP server.
- **Custom Error Handling** and standardized JSON responses.
- **Structured Logging** using `slog`.
- **JWT** authentication support.
