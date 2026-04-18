# GoStack architecture

This document is a **structural overview** of the framework as implemented in this repository. For product intent and feature depth, see [PRD.md](../PRD.md).

---

## 1. System context

How a typical GoStack deployment sits next to clients and backing services.

```mermaid
flowchart LR
  subgraph clients [Clients]
    B[Browser / HTMX]
    M[Mobile / SPA]
    API[API clients]
  end

  subgraph app [GoStack process]
    S[HTTP server]
    K[App kernel]
  end

  subgraph data [Data and infra]
    SQL[(SQL DB)]
    R[(Redis)]
    O[Object storage]
    MAI[AI HTTP APIs]
  end

  B --> S
  M --> S
  API --> S
  S --> K
  K --> SQL
  K --> R
  K --> O
  K --> MAI
```

---

## 2. HTTP request path

Incoming requests pass through **global middleware**, then the **router** dispatches to a **handler**. Handlers receive a `*Context` and return `error`; the kernel maps failures to HTTP responses.

```mermaid
flowchart TB
  IN[HTTP request] --> MW[Global middleware chain]
  MW --> R[Router match]
  R --> C[New Context]
  C --> H[HandlerFunc]
  H --> OUT[Response]
  H -->|error| ERR[Error response]

  subgraph mw_examples [Typical middleware]
    M1[Recover]
    M2[RequestID / Logger]
    M3[CORS / Timeout / Gzip]
    M4[JWT / CSRF]
  end
  MW -.-> mw_examples
```

**Route groups** wrap handlers with **group-local** middleware before the handler runs (same `Context` model).

---

## 3. Application kernel

Core types in package `github.com/rohitdas13595/go-stack` and how they relate.

```mermaid
flowchart TB
  subgraph kernel [gostack package]
    App[App]
    Rout[router.Router]
    Mwchain[middleware chain]
    RE[RenderEngine]
    VP[validator + policy]
    App --> Rout
    App --> Mwchain
    App --> RE
    App --> VP
  end

  Ctx[Context per request]
  HF[HandlerFunc]

  App --> HF
  HF --> Ctx
  Ctx -.-> App
```

- **`App`** — Registers routes, groups, and resources; owns global middleware; optional renderer and validation/policy hooks.
- **`Context`** — Per-request wrapper: `Param`, `Query`, `Bind`, `JSON`, `Render`, `RenderPartial`, `User`, `Authorize`, etc.
- **`RenderEngine`** — `html/template` over a filesystem `fs.FS` (e.g. `views/`).

---

## 4. Package map

Major **library packages** under the module (simplified dependency view).

```mermaid
flowchart TB
  subgraph core [Core]
    GS[gostack]
    RT[router]
    MW[middleware]
  end

  subgraph app_layer [App composition]
    CFG[config]
    DBP[db]
    MIG[migrate]
    ORM[orm]
    AUTH[auth]
  end

  subgraph realtime [Realtime]
    SSEP[sse]
    WSP[ws]
  end

  subgraph integrations [Integrations]
    JOB[jobs]
    CACHE[cache]
    STO[storage]
    MAIL[mail]
    AI[ai]
    OBS[observability]
    OAU[oauth]
    OAPI[openapi]
    WBN[webauthn]
  end

  GS --> RT
  GS --> MW
  GS --> CFG
  GS --> DBP
  GS --> ORM
  GS --> AUTH
  MIG --> DBP
  ORM --> DBP
  GS --> SSEP
  GS --> WSP
  GS --> JOB
  GS --> CACHE
  GS --> STO
  GS --> MAIL
  GS --> AI
  GS --> OBS
  GS --> OAU
  GS --> OAPI
  GS --> WBN
```

Handlers usually import **`gostack`** and **`middleware`**; data access uses **`db`**, **`migrate`**, **`orm`**, or **`gostack.DB()` / `gostack.Query`** when the DB manager is registered.

---

## 5. Bridges and globals

Some capabilities are configured once at startup and reached through **package-level helpers** (optional sugar; you can also pass dependencies explicitly in your own code).

```mermaid
flowchart LR
  subgraph boot [Startup]
    MGR[db.Manager]
    LD[config.Load]
  end

  subgraph bridges [gostack bridges]
    SetDB[SetDBManager]
    DBfn[DB]
    Q[Query / Find]
    CFG[Config / ConfigInt]
  end

  MGR --> SetDB
  SetDB --> DBfn
  DBfn --> Q
  LD --> CFG
```

- **`SetDBManager`** — Registers named pools; **`DB()`** returns `*sql.DB` for ORM and raw SQL.
- **`config.Global()`** — Filled by **`config.Load`**; **`gostack.Config*`** reads dotted paths.

---

## 6. CLI and generated apps

The **`gostack`** binary is a thin entrypoint; commands live under **`internal/cli`**.

```mermaid
flowchart TB
  MAIN[cmd/gostack/main.go] --> CLI[internal/cli]
  CLI --> NEW[new scaffold]
  CLI --> SERVE[serve]
  CLI --> DB[db migrate / rollback / status]
  CLI --> RT[routes / env / work / docs]

  NEW --> APP[Generated app module]
  APP --> SRVR[cmd/server]
  APP --> RU[routes + handlers]
  APP --> VWS[views]
  APP --> MIGF[db/migrations]
```

Generated applications depend on the **same module** `github.com/rohitdas13595/go-stack` (often with a **`replace`** to a local checkout during development).

---

## 7. Examples workspace

Local **examples** are separate **Go modules** in `examples/*`, each replacing the framework path so builds work without publishing.

```mermaid
flowchart TB
  W[go.work] --> ROOT[module root]
  W --> E1[examples/hello]
  W --> E2[examples/ssr]
  W --> E3[examples/api]
  W --> EM[more examples/...]

  subgraph loc [Local replace]
    REP["replace go-stack => ../.."]
  end

  E1 --> REP
  E2 --> REP
  E3 --> REP
  EM --> REP
  REP --> ROOT
```

---

## 8. Shutdown and server lifecycle

`App.ListenAndServe` uses the standard library `http.Server`. **`ListenAndServeContext`** ties server shutdown to a context and runs **`App.Shutdown`** hooks (registered by integrations that need cleanup).

```mermaid
sequenceDiagram
  participant Main
  participant App
  participant Srv as http.Server
  participant Hooks as shutdown hooks

  Main->>App: ListenAndServeContext(ctx, addr)
  App->>Srv: ListenAndServe
  Note over Main, Srv: until ctx cancelled or error
  Main->>App: ctx done
  App->>Srv: Shutdown
  App->>Hooks: Shutdown(ctx)
```

---

## Reading order

1. **`gostack.go`** — routing, groups, `ServeHTTP`, `ListenAndServe*`
2. **`context.go`** — request API surface
3. **`router/router.go`** — matching semantics
4. **`middleware/middleware.go`** — common middleware
5. **`internal/cli/`** — how scaffolding and DB commands work

For file-level detail, see [DEVELOPMENT.md](DEVELOPMENT.md).
