# Developer & Contributor Guide

## 1. Environment Setup

### 1.1 Toolchain Requirements
Ensure your development machine has the following tools installed:
- **Go:** `1.22` or higher (`1.26` recommended).
- **Node.js:** `18.x` or higher (`v24.x` recommended) and `npm`.
- **Wails CLI v2:** Install via:
  ```powershell
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

### 1.2 Repository Preparation
Clone the repository and install the frontend dependencies:
```powershell
# Install frontend npm dependencies
npm --prefix frontend install
```

---

## 2. Database Management & Seeding

The application stores its SQLite database under `%LOCALAPPDATA%\SchoolManagement\data\school.db` (on Windows).

A CLI utility is provided in `cmd/tools/db/main.go` to manage this database:

```powershell
# Reset database and seed demo dataset
go run ./cmd/tools/db -fresh -seed

# Wipe and rerun migrations only (without demo seed)
go run ./cmd/tools/db -fresh

# Seed additional demo records without wiping
go run ./cmd/tools/db -seed
```

### Default Seed Credentials:
- **Username:** `admin`
- **Password:** `Password123!`

---

## 3. Running the Development Server

### Option A: Wails Live Dev (Recommended)
```powershell
wails dev
```
- Starts the Vite dev server with Hot Module Replacement (HMR).
- Automatically watches Go backend files and rebuilds the desktop container on changes.
- Automatically synchronizes TypeScript bindings in `frontend/wailsjs/`.

### Option B: Direct Go Run
If you prefer running without the Wails CLI dev watcher:
```powershell
# 1. Build frontend bundle
npm --prefix frontend run build

# 2. Run Go desktop entrypoint
go run ./cmd/desktop
```

---

## 4. Testing & Verification

### 4.1 Backend Integration Tests
All integration tests run against isolated in-memory SQLite instances and execute in under 3 seconds:

```powershell
# Run all integration test suites
go test -v ./tests/integration/...

# Run a specific test suite
go test -v ./tests/integration/ -run TestPaymentLifecycleAndCreditEngine
```

### 4.2 Frontend Type Checking & Production Build
```powershell
# Run Vue TypeScript compilation check and Vite build
npm --prefix frontend run build
```

---

## 5. Coding Standards & Architectural Guardrails

1. **Domain Isolation:**
   - Code inside `internal/domain/` must **never** import `internal/infrastructure/`, `internal/transport/`, `wails`, or third-party web frameworks.
2. **Monetary Precision:**
   - **Never** use `float32` or `float64` for money or percentages. Always use `int64` minor units and integer basis points (`bps`).
3. **Transport Decoupling:**
   - Vue components must **never** call `window.go.*` directly. Always use the service abstraction layer in `frontend/src/services/`.
4. **Audit Logging:**
   - Any state-modifying use case executed in `internal/application/commands/` must record an audit entry via `recorder.Record(...)` in the same transaction.

