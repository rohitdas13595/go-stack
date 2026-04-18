# GoStack — Product Requirements Document (PRD)

**Version:** 1.0.0  
**Status:** Draft  
**Author:** Framework Design Team  
**Last Updated:** 2026-04-18

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Vision & Goals](#2-vision--goals)
3. [Target Audience](#3-target-audience)
4. [System Architecture Overview](#4-system-architecture-overview)
5. [Core Modules](#5-core-modules)
   - 5.1 [CLI & Project Scaffolding](#51-cli--project-scaffolding)
   - 5.2 [HTTP Router & Middleware](#52-http-router--middleware)
   - 5.3 [Rendering Engine (SSR / CSR / ISR)](#53-rendering-engine-ssr--csr--isr)
   - 5.4 [Database Layer](#54-database-layer)
   - 5.5 [Migration System](#55-migration-system)
   - 5.6 [ORM](#56-orm)
   - 5.7 [Authentication & Authorization](#57-authentication--authorization)
   - 5.8 [AI SDK](#58-ai-sdk)
   - 5.9 [Configuration & Environment](#59-configuration--environment)
   - 5.10 [Job Queue & Background Workers](#510-job-queue--background-workers)
   - 5.11 [Caching](#511-caching)
   - 5.12 [File Storage](#512-file-storage)
   - 5.13 [Mailer](#513-mailer)
   - 5.14 [WebSocket & Real-Time](#514-websocket--real-time)
   - 5.15 [Observability (Logging, Metrics, Tracing)](#515-observability-logging-metrics-tracing)
   - 5.16 [Testing Utilities](#516-testing-utilities)
6. [API Design Principles](#6-api-design-principles)
7. [Directory Structure](#7-directory-structure)
8. [Data Models](#8-data-models)
9. [Security Requirements](#9-security-requirements)
10. [Performance Requirements](#10-performance-requirements)
11. [Developer Experience (DX)](#11-developer-experience-dx)
12. [Deployment & Infrastructure](#12-deployment--infrastructure)
13. [Versioning & Compatibility](#13-versioning--compatibility)
14. [Roadmap & Milestones](#14-roadmap--milestones)
15. [Open Questions & Decisions](#15-open-questions--decisions)
16. [Glossary](#16-glossary)

---

## 1. Executive Summary

**GoStack** is a batteries-included, opinionated, full-stack web framework for Go. It provides everything a production team needs — routing, SSR/CSR rendering, a database-agnostic ORM, schema migrations, authentication, an AI SDK, background jobs, real-time communication, and observability — in a single cohesive package with a consistent API.

GoStack is designed to feel as ergonomic as Laravel, Rails, or Next.js, while preserving Go's hallmark traits: type safety, concurrency, and raw performance.

---

## 2. Vision & Goals

### 2.1 Vision Statement

> _"Build production-grade web applications in Go without stitching together dozens of third-party libraries."_

### 2.2 Primary Goals

| #   | Goal                                                         | Success Metric                         |
| --- | ------------------------------------------------------------ | -------------------------------------- |
| G1  | Ship a working CRUD app in under 10 minutes                  | Measured via new-user onboarding study |
| G2  | Zero-boilerplate for 80% of common web patterns              | Feature coverage audit                 |
| G3  | SSR page cold-start ≤ 10 ms, API response p99 ≤ 5 ms (local) | Benchmark suite                        |
| G4  | Type-safe ORM with zero `interface{}` leakage                | Static analysis CI gate                |
| G5  | First-class AI integration (streaming, RAG, agents)          | SDK unit test coverage ≥ 90%           |

### 2.3 Non-Goals

- GoStack is **not** a replacement for bare `net/http` when you want zero opinions.
- GoStack does **not** compile Go to WASM for the browser; CSR is served as JavaScript bundles (HTMX / Alpine.js / optional React adapter).
- GoStack is **not** a microservices orchestration platform (though it can power microservices).

---

## 3. Target Audience

| Persona                | Description                                                                    |
| ---------------------- | ------------------------------------------------------------------------------ |
| **Solo Developer**     | Wants a single tool to ship SaaS products fast, no DevOps expertise required   |
| **Go Backend Team**    | Wants unified conventions across projects, tired of glue-code between packages |
| **Polyglot Developer** | Comes from Rails/Laravel/Django; wants equivalent productivity in Go           |
| **AI Product Builder** | Building LLM-powered features and needs tight server-side AI integration       |

---

## 4. System Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                          GoStack Application                    │
│                                                                 │
│  ┌────────────┐  ┌─────────────┐  ┌─────────────────────────┐  │
│  │  CLI Tool  │  │  Config &   │  │   Plugin / Extension    │  │
│  │ (gostack)  │  │  Env Mgmt   │  │       System            │  │
│  └────────────┘  └─────────────┘  └─────────────────────────┘  │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                    HTTP Kernel                          │    │
│  │   Router ─► Middleware Chain ─► Handler ─► Response    │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌─────────────────────┐   │
│  │ Render Engine│  │  Auth Layer  │  │     AI SDK          │   │
│  │ SSR/CSR/ISR  │  │ JWT/Session/ │  │ Completion/Stream/  │   │
│  │ Templates    │  │ OAuth/RBAC   │  │ Embeddings/Agents   │   │
│  └──────────────┘  └──────────────┘  └─────────────────────┘   │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌─────────────────────┐   │
│  │   ORM &      │  │  Migration   │  │   Query Builder     │   │
│  │  Model Layer │  │   Engine     │  │   (type-safe)       │   │
│  └──────────────┘  └──────────────┘  └─────────────────────┘   │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌─────────────────────┐   │
│  │  Job Queue   │  │   Caching    │  │   File Storage      │   │
│  │  & Workers   │  │ (Redis/Mem)  │  │   (S3 / Local)      │   │
│  └──────────────┘  └──────────────┘  └─────────────────────┘   │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌─────────────────────┐   │
│  │  WebSocket   │  │   Mailer     │  │  Observability      │   │
│  │  & SSE       │  │  (SMTP/SES)  │  │  Logs/Metrics/Trace │   │
│  └──────────────┘  └──────────────┘  └─────────────────────┘   │
│                                                                 │
│         ┌──────────────────────────────────────┐               │
│         │         Database Drivers             │               │
│         │  PostgreSQL │ MySQL │ SQLite │ Others │               │
│         └──────────────────────────────────────┘               │
└─────────────────────────────────────────────────────────────────┘
```

### 4.1 Request Lifecycle

```
Incoming Request
      │
      ▼
 [Global Middleware]   ← CORS, RateLimiter, RequestID, Logger
      │
      ▼
   [Router]            ← Matches method + path pattern
      │
      ▼
 [Route Middleware]    ← Auth, CSRF, ACL
      │
      ▼
   [Handler]           ← Controller method or inline func
      │
      ▼
 [Render Engine]       ← SSR template / JSON / SSE / WS upgrade
      │
      ▼
 [Response Writer]     ← Compression, header injection, streaming
```

---

## 5. Core Modules

---

### 5.1 CLI & Project Scaffolding

#### Overview

The `gostack` binary is the primary developer interface. It provides code generation, migration management, development server with hot reload, and build tooling.

#### Commands

```
gostack new <project-name>          # Scaffold a new application
gostack new <project-name> --api   # API-only (no rendering)
gostack new <project-name> --spa   # SPA mode (React/Vue adapter)

gostack make:controller <Name>      # Generate controller
gostack make:model <Name>           # Generate model + migration
gostack make:middleware <Name>      # Generate middleware
gostack make:migration <name>       # Generate blank migration file
gostack make:job <Name>             # Generate background job
gostack make:mailer <Name>          # Generate mailer

gostack db:migrate                  # Run pending migrations
gostack db:rollback                 # Rollback last migration batch
gostack db:seed                     # Run seeders
gostack db:reset                    # Drop, migrate, seed
gostack db:status                   # Show migration status

gostack serve                       # Dev server with hot reload
gostack serve --port 4000

gostack build                       # Production build
gostack build --docker              # Output Dockerfile + binary

gostack test                        # Run test suite
gostack test --coverage

gostack routes                      # List all registered routes
gostack env                         # Print resolved config
```

#### Scaffolded Project Structure

See [Section 7 — Directory Structure](#7-directory-structure).

#### Hot Reload

- Uses `fsnotify` to watch `.go` and template files.
- Recompiles and restarts the binary in-process using a subprocess model.
- Template changes re-render without recompile.

---

### 5.2 HTTP Router & Middleware

#### Overview

GoStack ships its own high-performance radix-tree router, wrapping `net/http`. It supports named params, wildcard segments, route groups, and inline/named middleware.

#### Routing API

```go
app := gostack.New()

// Basic routes
app.GET("/", handlers.Home)
app.POST("/users", handlers.CreateUser)
app.PUT("/users/:id", handlers.UpdateUser)
app.DELETE("/users/:id", handlers.DeleteUser)
app.PATCH("/users/:id", handlers.PatchUser)

// Resource shorthand (generates CRUD routes)
app.Resource("/posts", handlers.PostHandler{})

// Route groups with shared middleware
api := app.Group("/api/v1", middleware.Auth(), middleware.RateLimit(100))
api.GET("/me", handlers.Me)

// Nested groups
admin := app.Group("/admin", middleware.Auth(), middleware.RequireRole("admin"))
admin.GET("/dashboard", handlers.AdminDashboard)
admin.Resource("/users", handlers.AdminUserHandler{})

// Named routes
app.GET("/profile", handlers.Profile).Name("user.profile")
url := app.RouteURL("user.profile") // "/profile"
```

#### Built-in Middleware

| Middleware                | Description                                    |
| ------------------------- | ---------------------------------------------- |
| `middleware.Logger()`     | Structured request logging                     |
| `middleware.Recover()`    | Panic recovery → 500                           |
| `middleware.CORS(opts)`   | Configurable CORS headers                      |
| `middleware.RateLimit(n)` | In-memory token bucket; Redis-backed option    |
| `middleware.Auth()`       | Validates JWT or session; populates `ctx.User` |
| `middleware.CSRF()`       | CSRF token generation and validation           |
| `middleware.RequestID()`  | Injects `X-Request-ID`                         |
| `middleware.Compress()`   | Gzip / Brotli response compression             |
| `middleware.Timeout(d)`   | Request timeout via context cancellation       |
| `middleware.Cache(d)`     | HTTP cache headers + optional response cache   |

#### Context (`ctx`)

```go
func (h *UserHandler) Show(ctx *gostack.Context) error {
    id   := ctx.Param("id")           // route param
    page := ctx.QueryInt("page", 1)   // query string with default
    var body CreateUserRequest
    if err := ctx.Bind(&body); err != nil { // JSON/form decode + validate
        return ctx.BadRequest(err)
    }
    user := ctx.User()                 // authenticated user (or nil)
    ctx.Set("key", value)             // per-request store
    return ctx.JSON(200, user)
    // or: ctx.Render("users/show", gostack.Data{"user": user})
    // or: ctx.Redirect("/login")
    // or: ctx.Stream(streamFunc)
}
```

---

### 5.3 Rendering Engine (SSR / CSR / ISR)

#### Rendering Modes

| Mode                     | Abbreviation | Description                                      |
| ------------------------ | ------------ | ------------------------------------------------ |
| Server-Side Rendering    | SSR          | Full HTML generated per request on the server    |
| Client-Side Rendering    | CSR          | Server sends shell HTML; JS takes over rendering |
| Incremental Static Regen | ISR          | Page cached for TTL, revalidated in background   |
| Partial Page Rendering   | PPR          | Component-level SSR (for HTMX / partial swaps)   |

#### Template Engine

- Default: Go `html/template` with an ergonomic layout/component system.
- Optional: `templ` adapter (type-safe Go templates compiled to functions).
- Optional: React/JSX adapter via build pipeline (Vite integration).

```
views/
  layouts/
    app.html          # root layout
    auth.html         # auth pages layout
  components/
    nav.html          # reusable partials
    flash.html
  pages/
    home.html
    users/
      index.html
      show.html
      edit.html
```

#### Template Rendering

```go
// Controller
return ctx.Render("users/show", gostack.Data{
    "user":  user,
    "title": "User Profile",
})

// Partial (for HTMX swap)
return ctx.RenderPartial("components/user-card", gostack.Data{"user": user})

// ISR with 60-second TTL
return ctx.RenderISR("pages/home", gostack.Data{...}, 60*time.Second)
```

#### Template Helpers

```go
// Available inside templates
{{ url "user.profile" }}
{{ asset "app.css" }}       // fingerprinted asset path
{{ csrf_token }}
{{ flash "success" }}
{{ partial "nav" . }}
{{ truncate .Bio 120 }}
{{ timeAgo .CreatedAt }}
{{ json .Data }}
```

#### CSR Mode

- Server renders an HTML shell with a `<div id="app">` mount point.
- JavaScript bundle (Vite-built) mounts and takes over routing.
- API routes are automatically available; auth tokens injected via a `window.__INIT__` payload.
- Supports React, Vue, Svelte, or vanilla JS adapter packages.

#### SSE (Server-Sent Events)

```go
app.GET("/stream", func(ctx *gostack.Context) error {
    return ctx.SSE(func(send gostack.SSESender) {
        for event := range eventCh {
            send("update", event)
        }
    })
})
```

---

### 5.4 Database Layer

#### Supported Drivers

| Database        | Driver                | Status                  |
| --------------- | --------------------- | ----------------------- |
| PostgreSQL      | `pgx/v5`              | ✅ Primary              |
| MySQL / MariaDB | `go-sql-driver/mysql` | ✅ Supported            |
| SQLite          | `modernc.org/sqlite`  | ✅ Supported (dev/test) |
| CockroachDB     | `pgx/v5` (compat)     | 🟡 Planned              |
| TiDB            | MySQL driver          | 🟡 Planned              |

#### Connection Configuration

```yaml
# config/database.yaml
default: postgres

connections:
  postgres:
    driver: postgres
    host: ${DB_HOST:localhost}
    port: ${DB_PORT:5432}
    database: ${DB_NAME:myapp}
    username: ${DB_USER:postgres}
    password: ${DB_PASS:}
    ssl_mode: prefer
    pool:
      max_open: 25
      max_idle: 5
      max_lifetime: 5m

  sqlite:
    driver: sqlite
    database: ./storage/dev.db
```

#### Multi-DB Support

```go
db := gostack.DB()                    // default connection
readDB := gostack.DB("read_replica")  // named connection
```

#### Connection Pool

- Wraps `database/sql` pool.
- Exposes pool stats via metrics endpoint.
- Automatic reconnect with exponential backoff.

---

### 5.5 Migration System

#### Overview

File-based, version-controlled SQL migrations with a Go API for programmatic migrations. Inspired by Flyway and Rails Active Record Migrations.

#### Migration Files

```
db/migrations/
  20240101000001_create_users.sql
  20240101000002_add_email_index.sql
  20240102000001_create_posts.sql
```

#### File Format

```sql
-- 20240101000001_create_users.sql

-- +gostack:up
CREATE TABLE users (
  id          BIGSERIAL PRIMARY KEY,
  email       VARCHAR(255) NOT NULL UNIQUE,
  name        VARCHAR(255) NOT NULL,
  password    VARCHAR(255) NOT NULL,
  role        VARCHAR(50)  NOT NULL DEFAULT 'user',
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_users_email ON users(email);

-- +gostack:down
DROP TABLE IF EXISTS users;
```

#### Go-based Migrations

```go
// db/migrations/20240101000001_create_users.go
func init() {
    gostack.Migrations.Register(&gostack.Migration{
        Version: "20240101000001",
        Name:    "create_users",
        Up: func(db *gostack.DB) error {
            return db.Schema().CreateTable("users", func(t *gostack.Table) {
                t.BigIncrements("id")
                t.String("email", 255).Unique()
                t.String("name", 255)
                t.String("password", 255)
                t.String("role", 50).Default("user")
                t.Timestamps()
                t.SoftDeletes()
            })
        },
        Down: func(db *gostack.DB) error {
            return db.Schema().DropTable("users")
        },
    })
}
```

#### Migration Tracking

- Tracks applied migrations in a `_gostack_migrations` table.
- Supports batching: rollback reverts the most recent batch atomically.
- Dry-run mode prints SQL without executing.
- Locking: uses advisory locks (PostgreSQL) to prevent concurrent migrations.

---

### 5.6 ORM

#### Design Philosophy

- Generic-based, fully type-safe. No `interface{}` in the public API.
- Generates no code — all type safety comes from Go generics (1.21+).
- Supports eager loading, scopes, soft deletes, timestamps automatically.
- Raw SQL escape hatch always available.

#### Model Definition

```go
// models/user.go
type User struct {
    gostack.Model                         // ID, CreatedAt, UpdatedAt, DeletedAt

    Email    string        `db:"email"`
    Name     string        `db:"name"`
    Password string        `db:"password" json:"-"`
    Role     string        `db:"role"`

    Posts    []Post        `db:"-" has_many:"posts,foreign_key:user_id"`
    Profile  *UserProfile  `db:"-" has_one:"user_profiles,foreign_key:user_id"`
}

func (u *User) TableName() string { return "users" }

func (u *User) BeforeCreate(ctx context.Context) error {
    hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
    u.Password = string(hashed)
    return err
}
```

#### Query API

```go
// Find by primary key
user, err := gostack.Find[User](ctx, 42)

// First match
user, err := gostack.Query[User]().Where("email = ?", email).First(ctx)

// All with conditions, ordering, pagination
users, err := gostack.Query[User]().
    Where("role = ?", "admin").
    Where("deleted_at IS NULL").
    OrderBy("created_at DESC").
    Limit(20).
    Offset(page * 20).
    All(ctx)

// Eager loading
users, err := gostack.Query[User]().
    With("Posts", "Profile").
    All(ctx)

// Create
user := &User{Email: "alice@example.com", Name: "Alice"}
err = gostack.Create(ctx, user)

// Update
err = gostack.Update(ctx, user, gostack.Changes{"name": "Alice Smith"})

// Delete (soft delete if model embeds gostack.Model)
err = gostack.Delete(ctx, user)

// Hard delete
err = gostack.ForceDelete(ctx, user)

// Upsert
err = gostack.Upsert(ctx, user, []string{"email"})

// Count
count, err := gostack.Query[User]().Where("role = ?", "admin").Count(ctx)

// Exists
exists, err := gostack.Query[User]().Where("email = ?", email).Exists(ctx)

// Aggregates
avg, err := gostack.Query[Order]().Avg(ctx, "total")
sum, err := gostack.Query[Order]().Where("user_id = ?", uid).Sum(ctx, "total")

// Raw SQL
rows, err := gostack.DB().QueryRaw(ctx, "SELECT * FROM users WHERE ...", args...)
```

#### Scopes

```go
func ActiveUsers() gostack.Scope[User] {
    return func(q *gostack.QueryBuilder[User]) *gostack.QueryBuilder[User] {
        return q.Where("deleted_at IS NULL").Where("role != ?", "banned")
    }
}

users, err := gostack.Query[User]().Scope(ActiveUsers()).All(ctx)
```

#### Transactions

```go
err := gostack.Transaction(ctx, func(tx *gostack.TX) error {
    if err := tx.Create(ctx, order); err != nil {
        return err // auto-rollback
    }
    return tx.Update(ctx, inventory, gostack.Changes{"stock": stock - 1})
})
```

#### Relationships

| Type         | Tag Example                                        |
| ------------ | -------------------------------------------------- |
| Belongs To   | `belongs_to:"users,foreign_key:user_id"`           |
| Has One      | `has_one:"user_profiles,foreign_key:user_id"`      |
| Has Many     | `has_many:"posts,foreign_key:user_id"`             |
| Many-to-Many | `many_to_many:"post_tags,join_table:post_tag_map"` |
| Polymorphic  | `polymorphic:"commentable"`                        |

---

### 5.7 Authentication & Authorization

#### Authentication Strategies

| Strategy               | Description                                               |
| ---------------------- | --------------------------------------------------------- |
| **Session-based**      | Server-side sessions stored in Redis or DB                |
| **JWT**                | Stateless tokens; access + refresh token pair             |
| **OAuth2 / OIDC**      | Built-in providers: Google, GitHub, Discord, generic OIDC |
| **API Key**            | Hash-stored keys with scopes                              |
| **Magic Link**         | Passwordless email login                                  |
| **Passkey (WebAuthn)** | FIDO2 / WebAuthn support                                  |

#### Auth Configuration

```go
app.Use(auth.New(auth.Config{
    Strategy: auth.JWT,
    Secret:   gostack.Env("JWT_SECRET"),
    TTL:      15 * time.Minute,
    Refresh: auth.RefreshConfig{
        TTL:      7 * 24 * time.Hour,
        Rotation: true,          // rotate refresh token on use
        CookieOnly: true,        // store refresh token in HttpOnly cookie
    },
    UserLoader: func(ctx context.Context, id string) (auth.User, error) {
        return gostack.Find[User](ctx, id)
    },
}))
```

#### OAuth2

```go
app.Use(oauth.New(oauth.Config{
    Providers: []oauth.Provider{
        oauth.Google(gostack.Env("GOOGLE_CLIENT_ID"), gostack.Env("GOOGLE_SECRET")),
        oauth.GitHub(gostack.Env("GITHUB_CLIENT_ID"), gostack.Env("GITHUB_SECRET")),
    },
    OnSuccess: func(ctx *gostack.Context, profile oauth.Profile) error {
        user, _ := gostack.Query[User]().Where("email = ?", profile.Email).First(ctx)
        if user == nil {
            user = &User{Email: profile.Email, Name: profile.Name}
            gostack.Create(ctx, user)
        }
        return auth.Login(ctx, user)
    },
}))

// Auto-registers:
//   GET  /auth/google
//   GET  /auth/google/callback
//   GET  /auth/github
//   GET  /auth/github/callback
```

#### Authorization (RBAC + Policies)

```go
// Define policy
type PostPolicy struct{}

func (p PostPolicy) Update(user *User, post *Post) bool {
    return post.UserID == user.ID || user.Role == "admin"
}

// Register
app.Policy(&Post{}, PostPolicy{})

// Use in handler
func (h *PostHandler) Update(ctx *gostack.Context) error {
    post, _ := gostack.Find[Post](ctx, ctx.Param("id"))
    if err := ctx.Authorize("update", post); err != nil {
        return ctx.Forbidden("Cannot edit this post")
    }
    // ...
}

// Route-level authorization
admin := app.Group("/admin", middleware.RequireRole("admin"))
```

#### Built-in Auth Endpoints (optional)

```
POST   /auth/register
POST   /auth/login
POST   /auth/logout
POST   /auth/refresh
POST   /auth/forgot-password
POST   /auth/reset-password
GET    /auth/me
```

---

### 5.8 AI SDK

#### Overview

First-class integration with major LLM providers. Supports streaming, function calling / tool use, embeddings, vector search, and multi-step agents.

#### Provider Support

| Provider      | Models                         |
| ------------- | ------------------------------ |
| **Anthropic** | Claude 3.5/4 series            |
| **OpenAI**    | GPT-4o, o1, o3                 |
| **Google**    | Gemini 1.5/2.0                 |
| **Mistral**   | Mistral Large, Codestral       |
| **Ollama**    | Any locally-hosted model       |
| **Custom**    | Any OpenAI-compatible endpoint |

#### Configuration

```yaml
# config/ai.yaml
default: anthropic

providers:
  anthropic:
    api_key: ${ANTHROPIC_API_KEY}
    default_model: claude-sonnet-4-20250514
  openai:
    api_key: ${OPENAI_API_KEY}
    default_model: gpt-4o
  ollama:
    base_url: http://localhost:11434
    default_model: llama3
```

#### Completion API

```go
ai := gostack.AI()

// Simple completion
resp, err := ai.Complete(ctx, &ai.CompletionRequest{
    Model:  "claude-sonnet-4-20250514",
    System: "You are a helpful assistant.",
    Messages: []ai.Message{
        {Role: "user", Content: "Summarize this document: " + doc},
    },
    MaxTokens: 1024,
})
fmt.Println(resp.Content)
```

#### Streaming

```go
func (h *ChatHandler) Stream(ctx *gostack.Context) error {
    return ctx.SSE(func(send gostack.SSESender) {
        stream, err := gostack.AI().Stream(ctx, &ai.CompletionRequest{
            Messages: messages,
        })
        for chunk := range stream {
            send("delta", chunk.Delta)
        }
    })
}
```

#### Tool Use / Function Calling

```go
tools := []ai.Tool{
    {
        Name:        "get_weather",
        Description: "Get current weather for a city",
        Schema: ai.Schema{
            Properties: map[string]ai.Property{
                "city": {Type: "string", Description: "City name"},
            },
            Required: []string{"city"},
        },
        Handler: func(input map[string]any) (any, error) {
            city := input["city"].(string)
            return weather.Fetch(city)
        },
    },
}

resp, err := ai.Complete(ctx, &ai.CompletionRequest{
    Messages: messages,
    Tools:    tools,
    ToolChoice: "auto",
})
```

#### Embeddings & Vector Search

```go
// Generate embedding
vec, err := gostack.AI().Embed(ctx, "The quick brown fox")

// Store in DB (requires pgvector extension for Postgres)
doc := &Document{Content: text, Embedding: vec}
gostack.Create(ctx, doc)

// Similarity search
results, err := gostack.Query[Document]().
    NearestNeighbor("embedding", queryVec, 10). // top-10
    All(ctx)
```

#### Agent Loop

```go
agent := ai.NewAgent(ai.AgentConfig{
    Model:      "claude-sonnet-4-20250514",
    System:     "You are a research assistant.",
    Tools:      []ai.Tool{searchTool, calculatorTool},
    MaxSteps:   10,
    OnStep: func(step ai.Step) {
        log.Info("agent step", "tool", step.ToolName)
    },
})

result, err := agent.Run(ctx, "Research the top 5 Go web frameworks and compare them.")
```

#### Structured Output

```go
type SentimentResult struct {
    Score    float64 `json:"score"`
    Label    string  `json:"label"`
    Reasoning string `json:"reasoning"`
}

result, err := ai.Extract[SentimentResult](ctx, &ai.ExtractionRequest{
    Model:  "claude-sonnet-4-20250514",
    Prompt: "Analyze sentiment: " + reviewText,
})
// result is *SentimentResult, fully typed
```

#### Middleware & Hooks

```go
gostack.AI().Use(
    ai.RateLimit(1000),                    // requests/minute
    ai.Logger(),                            // log all requests
    ai.Cache(5*time.Minute),               // cache identical prompts
    ai.Retry(3, ai.ExponentialBackoff),    // retry on 5xx
    ai.CostTracker(budget.Track),          // track token costs
)
```

---

### 5.9 Configuration & Environment

```yaml
# config/app.yaml
app:
  name: MyApp
  env: ${APP_ENV:development}
  url: ${APP_URL:http://localhost:3000}
  key: ${APP_KEY} # 32-byte secret for encryption/signing
  debug: ${APP_DEBUG:true}

server:
  host: 0.0.0.0
  port: ${PORT:3000}
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 120s
```

```go
// Access
gostack.Config("app.name")             // string
gostack.ConfigInt("server.port")       // int
gostack.ConfigBool("app.debug")        // bool
gostack.Env("DATABASE_URL")            // raw env variable
gostack.Env("STRIPE_KEY", "fallback") // with default
```

- Config files loaded in order: `config/app.yaml` → `config/{env}.yaml` → environment variables.
- Secrets can be sourced from Vault, AWS SSM, or GCP Secret Manager via adapters.

---

### 5.10 Job Queue & Background Workers

#### Overview

Redis-backed (or DB-backed for simple cases) priority job queue with reliable delivery, retries, scheduling, and concurrency control.

```go
// Define a job
type WelcomeEmailJob struct {
    UserID int64
}

func (j *WelcomeEmailJob) Handle(ctx context.Context) error {
    user, err := gostack.Find[User](ctx, j.UserID)
    return mailer.Send("welcome", user)
}

// Dispatch
gostack.Dispatch(ctx, &WelcomeEmailJob{UserID: user.ID})

// Dispatch with options
gostack.Dispatch(ctx, &WelcomeEmailJob{UserID: user.ID},
    jobs.Delay(5*time.Minute),
    jobs.Queue("emails"),
    jobs.MaxRetries(3),
    jobs.RetryAfter(time.Minute),
)

// Scheduled (cron-style)
app.Schedule("0 8 * * *", &DailyDigestJob{})

// Run workers
gostack serve --workers          // in dev
// or separate process:
gostack work --queues=default,emails --concurrency=10
```

---

### 5.11 Caching

```go
cache := gostack.Cache()

// Set / Get
cache.Set(ctx, "user:42", user, 10*time.Minute)
cached, ok := cache.Get[User](ctx, "user:42")

// Remember (compute if missing)
user, err := gostack.Remember(ctx, "user:42", 10*time.Minute, func() (*User, error) {
    return gostack.Find[User](ctx, 42)
})

// Invalidate
cache.Delete(ctx, "user:42")
cache.Flush(ctx)                     // clear all

// Tags
cache.Tags("users").Set(ctx, "user:42", user, 0)
cache.Tags("users").Flush(ctx)      // invalidate all user cache
```

Backends: In-memory (default for dev), Redis, Memcached, DynamoDB.

---

### 5.12 File Storage

```go
storage := gostack.Storage()       // default disk
s3 := gostack.Storage("s3")       // named disk

// Upload
url, err := storage.Put(ctx, "avatars/"+filename, fileReader, storage.Options{
    Visibility: "public",
    ContentType: "image/jpeg",
})

// Read
reader, err := storage.Get(ctx, "avatars/profile.jpg")

// Signed URL (private files)
url, err := storage.SignedURL(ctx, "invoices/123.pdf", 24*time.Hour)

// Delete
err = storage.Delete(ctx, "old-avatar.jpg")
```

Backends: Local disk, AWS S3, GCS, Cloudflare R2, MinIO (S3-compatible).

---

### 5.13 Mailer

```go
// Define mailer
type WelcomeMailer struct {
    User *User
}

func (m *WelcomeMailer) Build() *mail.Message {
    return mail.New().
        To(m.User.Email, m.User.Name).
        Subject("Welcome to MyApp!").
        Template("emails/welcome", mail.Data{"user": m.User}).
        Attach("terms.pdf", termsFile)
}

// Send
gostack.Mail().Send(ctx, &WelcomeMailer{User: user})

// Queue for async delivery
gostack.Mail().Queue(ctx, &WelcomeMailer{User: user})
```

Backends: SMTP, AWS SES, Postmark, Resend, Mailgun. Dev preview via `gostack mail:preview`.

---

### 5.14 WebSocket & Real-Time

```go
// WebSocket handler
app.GET("/ws/chat/:room", func(ctx *gostack.Context) error {
    return ctx.WebSocket(func(conn *gostack.WSConn) {
        room := ctx.Param("room")
        gostack.Broadcast.Join(conn, "room:"+room)

        for msg := range conn.Receive() {
            gostack.Broadcast.To("room:"+room).Emit("message", msg)
        }
    })
})

// Broadcast from anywhere (e.g., after a DB write)
gostack.Broadcast.To("room:general").Emit("new-post", post)
gostack.Broadcast.ToUser(userID).Emit("notification", notif)
```

Channels are backed by Redis pub/sub for horizontal scaling.

---

### 5.15 Observability (Logging, Metrics, Tracing)

#### Logging

```go
log := gostack.Log()
log.Info("user created", "user_id", user.ID, "email", user.Email)
log.Warn("rate limit approaching", "ip", ip, "count", count)
log.Error("payment failed", "error", err, "order_id", order.ID)
```

Structured JSON by default. Adapters for `slog`, `zap`, `zerolog`.

#### Metrics

Prometheus-compatible metrics exposed at `/metrics`:

- HTTP request count, latency, error rate (by route, method, status)
- DB query count and latency
- Job queue depth and processing time
- Cache hit/miss ratio
- AI token usage

#### Distributed Tracing

OpenTelemetry-native. Exporters: Jaeger, Tempo, Datadog, Honeycomb.

```go
ctx, span := gostack.Trace(ctx, "process-payment")
defer span.End()
span.SetAttribute("order.id", orderID)
```

---

### 5.16 Testing Utilities

```go
// HTTP testing
func TestCreateUser(t *testing.T) {
    app := testutil.NewApp(t)
    defer app.Teardown()

    res := app.POST("/users", testutil.JSON{
        "email": "test@example.com",
        "name":  "Test User",
    })

    res.AssertStatus(201)
    res.AssertJSON("email", "test@example.com")
}

// Auth helpers
res := app.AsUser(adminUser).DELETE("/users/42")
res.AssertStatus(204)

// DB helpers
testutil.Seed[User](t, &User{Email: "alice@example.com"})
testutil.AssertRowCount[User](t, 1)
testutil.AssertNotExists[User](t, "email = ?", "deleted@example.com")

// Mail assertion
testutil.AssertMailSent(t, "welcome@example.com", "Welcome to MyApp!")

// Job assertion
testutil.AssertJobDispatched(t, &WelcomeEmailJob{})
```

---

## 6. API Design Principles

1. **Explicit over implicit** — No magic struct tags driving hidden behavior beyond what is documented.
2. **Errors as values** — No panics in the public API. All functions return `error`.
3. **Context propagation** — Every function accepting I/O takes a `context.Context`.
4. **Generics for type safety** — ORM, cache, and AI extraction use Go generics; no runtime casts.
5. **Progressive disclosure** — Simple defaults work out-of-the-box; advanced options via functional option pattern.
6. **Composability** — Each module is independently importable; you don't have to use all of GoStack.
7. **No global state pollution** — App instance carries all state; safe for testing in parallel.

---

## 7. Directory Structure

```
myapp/
├── cmd/
│   └── server/
│       └── main.go               # entry point
├── config/
│   ├── app.yaml
│   ├── database.yaml
│   ├── ai.yaml
│   └── production.yaml
├── db/
│   ├── migrations/               # .sql or .go migration files
│   └── seeds/
│       └── seed.go
├── internal/
│   ├── handlers/                 # HTTP handlers
│   │   ├── user_handler.go
│   │   └── post_handler.go
│   ├── models/                   # ORM models
│   │   ├── user.go
│   │   └── post.go
│   ├── middleware/               # Custom middleware
│   ├── jobs/                     # Background jobs
│   ├── mailers/                  # Email builders
│   ├── policies/                 # Authorization policies
│   └── services/                 # Business logic
├── views/
│   ├── layouts/
│   ├── components/
│   ├── pages/
│   └── emails/
├── public/                       # Static assets (served at /public)
│   ├── css/
│   ├── js/
│   └── images/
├── storage/                      # Local file storage (gitignored)
├── tests/
│   ├── integration/
│   └── fixtures/
├── routes/
│   └── web.go                    # Route definitions
├── go.mod
├── go.sum
├── .env
├── .env.example
└── Makefile
```

---

## 8. Data Models

### Core Framework Tables

```sql
-- Migration tracking
CREATE TABLE _gostack_migrations (
  id          SERIAL PRIMARY KEY,
  version     VARCHAR(255) NOT NULL UNIQUE,
  name        VARCHAR(255) NOT NULL,
  batch       INT NOT NULL,
  applied_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Session storage (optional)
CREATE TABLE _gostack_sessions (
  id          VARCHAR(128) PRIMARY KEY,
  user_id     BIGINT REFERENCES users(id) ON DELETE CASCADE,
  payload     JSONB NOT NULL,
  expires_at  TIMESTAMPTZ NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Job queue (DB backend)
CREATE TABLE _gostack_jobs (
  id           BIGSERIAL PRIMARY KEY,
  queue        VARCHAR(100) NOT NULL DEFAULT 'default',
  payload      JSONB NOT NULL,
  attempts     INT NOT NULL DEFAULT 0,
  max_attempts INT NOT NULL DEFAULT 3,
  reserved_at  TIMESTAMPTZ,
  available_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Password reset tokens
CREATE TABLE _gostack_password_resets (
  email       VARCHAR(255) NOT NULL,
  token       VARCHAR(255) NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## 9. Security Requirements

| Requirement              | Implementation                                                    |
| ------------------------ | ----------------------------------------------------------------- |
| SQL Injection prevention | Parameterized queries only; no string interpolation in ORM        |
| XSS prevention           | Auto-escaping in HTML templates; CSP header middleware            |
| CSRF protection          | Synchronizer token pattern; double-submit cookie for SPAs         |
| Password hashing         | bcrypt with cost ≥ 12 by default; Argon2id adapter                |
| Secrets management       | Never logged; redacted in error reports; Vault integration        |
| HTTPS enforcement        | `middleware.HTTPS()` redirects HTTP; HSTS header                  |
| Rate limiting            | Per-IP and per-user token bucket on auth endpoints                |
| Input validation         | `ctx.Bind()` runs struct validation via `go-playground/validator` |
| Audit logging            | Auth events, admin actions logged with actor and IP               |
| Dependency scanning      | `govulncheck` in CI                                               |

---

## 10. Performance Requirements

| Metric                            | Target         | Notes                               |
| --------------------------------- | -------------- | ----------------------------------- |
| HTTP routing overhead             | < 1 µs/request | Radix tree, zero allocs on hot path |
| SSR render time (simple page)     | < 5 ms         | Compiled templates, buffer pool     |
| ORM single-row fetch              | < 2 ms         | Postgres localhost                  |
| ORM list query (100 rows)         | < 10 ms        | With eager loading                  |
| Cold start (binary)               | < 50 ms        | Compiled Go binary                  |
| Memory baseline (idle)            | < 30 MB        | Single instance, no workers         |
| Throughput (echo server baseline) | ≥ 100k req/s   | Benchmarked on 4-core VM            |

---

## 11. Developer Experience (DX)

### 11.1 Hot Reload

- Go source recompiles in background using incremental build.
- Templates re-render on save without recompile.
- Browser auto-refresh via injected SSE script in dev mode.

### 11.2 Error Pages

- Rich development error pages showing stack trace, request details, DB query log, and template context.
- Production mode returns clean JSON or HTML error pages.

### 11.3 Interactive Debug Console

- `gostack console` launches a REPL with the full app context loaded.
- Query the DB, inspect routes, dispatch jobs interactively.

### 11.4 Documentation Generation

- `gostack docs` generates OpenAPI 3.1 spec from route annotations.
- Serves Swagger UI at `/docs` in dev mode.

### 11.5 Database Studio

- `gostack db:studio` opens a browser-based DB explorer (tables, query runner, migration history).

---

## 12. Deployment & Infrastructure

### 12.1 Build

```bash
gostack build
# Outputs: ./bin/server (single static binary with embedded assets)

gostack build --docker
# Outputs: Dockerfile with multi-stage build + fly.toml / docker-compose.yaml
```

### 12.2 Deployment Targets

| Target              | Support                                |
| ------------------- | -------------------------------------- |
| Bare metal / VPS    | ✅ Single binary                       |
| Docker / Kubernetes | ✅ Official Dockerfile template        |
| Fly.io              | ✅ `fly.toml` generator                |
| Railway             | ✅ `railway.json` generator            |
| AWS Lambda          | 🟡 Adapter planned (cold-start caveat) |
| Google Cloud Run    | ✅ Containerized binary                |

### 12.3 Health Checks

```
GET /health          # liveness probe  → {"status":"ok"}
GET /health/ready    # readiness probe → checks DB + Redis connectivity
```

### 12.4 Graceful Shutdown

- On `SIGTERM` / `SIGINT`: stops accepting new connections, drains in-flight requests (configurable timeout, default 30s), closes DB connections, flushes job workers.

---

## 13. Versioning & Compatibility

| Policy               | Detail                                                                |
| -------------------- | --------------------------------------------------------------------- |
| Go version           | Minimum Go 1.22; tested on latest two stable releases                 |
| Module versioning    | Semantic versioning via Go modules (`github.com/yourorg/gostack/v2`)  |
| Stability            | Public API stable after v1.0; breaking changes only in major versions |
| Deprecation          | Features deprecated for one minor cycle before removal                |
| Migration guides     | Published in `CHANGELOG.md` for every minor and major release         |
| CLI backwards compat | CLI flags stable within a major version                               |

---

## 14. Roadmap & Milestones

### Phase 1 — Foundation (Months 1–3)

- [ ] CLI scaffolding (`new`, `make:*`, `serve`)
- [ ] HTTP router + middleware chain
- [ ] Config & environment management
- [ ] ORM (CRUD, relations, transactions)
- [ ] Migration engine (SQL files)
- [ ] JWT + Session auth
- [ ] SSR templates

### Phase 2 — Full Stack (Months 4–6)

- [ ] CSR mode (HTMX + Alpine adapter)
- [ ] OAuth2 providers
- [ ] Redis-backed job queue
- [ ] Cache layer (Redis)
- [ ] Mailer
- [ ] File storage (local + S3)
- [ ] WebSocket / SSE
- [ ] Prometheus metrics

### Phase 3 — AI & Advanced (Months 7–9)

- [ ] AI SDK (Anthropic + OpenAI)
- [ ] Streaming completions
- [ ] Tool use / function calling
- [ ] Embeddings + vector search
- [ ] Agent loop
- [ ] ISR rendering mode

### Phase 4 — Production Hardening (Months 10–12)

- [ ] OpenTelemetry tracing
- [ ] RBAC + Policy system
- [ ] WebAuthn / Passkeys
- [ ] Rate limiting (Redis-backed)
- [ ] DB Studio
- [ ] OpenAPI docs generator
- [ ] React/Vue SPA adapter
- [ ] v1.0.0 release

---

## 15. Open Questions & Decisions

| #    | Question                      | Options                             | Decision                                   |
| ---- | ----------------------------- | ----------------------------------- | ------------------------------------------ |
| OQ-1 | Template engine default       | `html/template` vs `templ`          | TBD                                        |
| OQ-2 | CSR frontend default          | HTMX, Alpine, React                 | HTMX + Alpine (React optional)             |
| OQ-3 | Job queue backing store       | Redis vs DB vs NATS                 | Redis primary, DB fallback                 |
| OQ-4 | ORM code generation           | Generics only vs `sqlc` integration | Generics-first; `sqlc` adapter planned     |
| OQ-5 | Multi-tenancy model           | Row-level vs schema-level           | Row-level with tenant scope                |
| OQ-6 | Vector search                 | pgvector vs external Qdrant         | pgvector for Postgres; Qdrant adapter      |
| OQ-7 | Realtime horizontal scale     | Redis pub/sub vs NATS               | Redis pub/sub                              |
| OQ-8 | Module monorepo vs multi-repo | Monorepo                            | Monorepo with separate `go.mod` per module |

---

## 16. Glossary

| Term                | Definition                                                                               |
| ------------------- | ---------------------------------------------------------------------------------------- |
| **ISR**             | Incremental Static Regeneration — cached pages revalidated in the background after a TTL |
| **ORM**             | Object-Relational Mapper — maps Go structs to database rows                              |
| **SSR**             | Server-Side Rendering — HTML generated on the server per request                         |
| **CSR**             | Client-Side Rendering — HTML rendered in the browser via JavaScript                      |
| **PPR**             | Partial Page Rendering — server renders only a component fragment (used with HTMX)       |
| **Scope**           | A reusable query constraint function applied to ORM queries                              |
| **Soft Delete**     | Logical deletion via a `deleted_at` timestamp; row remains in DB                         |
| **Migration Batch** | A group of migrations applied together; rollback reverts the entire batch                |
| **Context (ctx)**   | GoStack's request context object wrapping `*http.Request` and `http.ResponseWriter`      |
| **Policy**          | An authorization rule object defining who can perform what actions on a resource         |
| **Agent**           | An AI loop that autonomously selects and invokes tools to complete a goal                |
| **Embedding**       | A dense vector representation of text used for semantic search                           |

---

_GoStack PRD — End of Document_  
_For questions or contributions, open an issue at `github.com/yourorg/gostack`._
