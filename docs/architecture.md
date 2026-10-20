# 🏛️ System Architecture

## 1. Architectural Philosophy

The **School Management System (SMS)** is designed as an **offline-first modular monolith** adopting **Clean Architecture** (Hexagonal / Ports & Adapters) principles.

The core motivation behind this architecture is:
1. **Total Offline Autonomy:** A single desktop installation operates with 100% functionality without internet connectivity.
2. **Future-Proof Portability:** The core business domain and application use cases remain completely agnostic of whether they are run inside a desktop webview (via Wails IPC) or hosted in a multi-tenant cloud environment (via Chi HTTP and PostgreSQL).
3. **Deterministic Financial Integrity:** Financial calculations are never subject to floating-point imprecision or distributed state anomalies.

---

## 2. High-Level System Architecture

```text
               ┌────────────────────────────────────────────────────────┐
               │              Frontend Presentation (Vue 3)             │
               │   TypeScript + Composition API + Pinia + Tailwind v4   │
               └───────────────────────────┬────────────────────────────┘
                                           │
                                           ▼
               ┌────────────────────────────────────────────────────────┐
               │           Transport Abstraction Interfaces             │
               │        (StudentService, FeeService, PaymentService)    │
               └───────────┬────────────────────────────────┬───────────┘
                           │                                │
                           ▼                                ▼
              ┌─────────────────────────┐      ┌─────────────────────────┐
              │   Wails Desktop Adapter │      │   HTTP REST API Adapter │
              │      (IPC Bindings)     │      │       (Go Chi Router)   │
              └────────────┬────────────┘      └────────────┬────────────┘
                           │                                │
                           └────────────────┬───────────────┘
                                            │
                                            ▼
               ┌────────────────────────────────────────────────────────┐
               │             Application Layer (CQRS-lite)              │
               │  ┌───────────────────────┐  ┌───────────────────────┐  │
               │  │   Command Handlers    │  │     Query Handlers    │  │
               │  │  (State Modifications)│  │    (Optimized Reads)  │  │
               │  └───────────┬───────────┘  └───────────────┬───────┘  │
               └──────────────┼──────────────────────────────┼──────────┘
                              ▼                              │
               ┌──────────────────────────────┐              │
               │         Domain Layer         │              │
               │  (Entities, Value Objects,   │              │
               │   Aggregates, Policies, Inv) │              │
               └──────────────┬───────────────┘              │
                              ▼                              │
               ┌──────────────────────────────┐              │
               │   Ports (Interfaces / UoW)   │◄─────────────┘
               └──────────────┬───────────────┘
                              │
                              ▼
               ┌────────────────────────────────────────────────────────┐
               │                  Infrastructure Layer                  │
               │  ┌───────────────────────┐  ┌───────────────────────┐  │
               │  │    SQLite Adapter     │  │  PostgreSQL Adapter   │  │
               │  │ (WAL Mode, modernc)   │  │    (Cloud / LAN)      │  │
               │  └───────────────────────┘  └───────────────────────┘  │
               │  ┌───────────────────────┐  ┌───────────────────────┐  │
               │  │ Local Filesystem Blob │  │     S3 Blob Storage   │  │
               │  │ (SHA-256 Addressed)   │  │     (Cloud / LAN)     │  │
               │  └───────────────────────┘  └───────────────────────┘  │
               └────────────────────────────────────────────────────────┘
```

---

## 3. Layer Separation & Responsibilities

### 3.1 Domain Layer (`internal/domain/`)
The domain layer represents the innermost core of the software. It contains:
- **Entities & Aggregates:** `Student`, `Guardian`, `AdmissionApplication`, `FeeStructure`, `Invoice`, `Payment`, `CreditBalance`.
- **Value Objects:** `Money` (integer minor units), `BasisPoints` (percentages as integers), `civil.Date` (timezone-safe dates), `ID` (UUIDv7).
- **Domain Policies:** Automated sibling discount rules, merit concessions, late fee calculations.
- **Invariants & Domain Errors:** Rule enforcement ensuring that payments cannot exceed outstanding balances without entering credit advance, fee structures cannot be modified once published, etc.

> **Rule:** The domain layer has **zero dependencies** on outer layers, database drivers, Wails, or third-party web frameworks. It imports only standard Go packages and its own sub-packages.

---

### 3.2 Application Layer (`internal/application/`)
The application layer coordinates use-case execution and implements **CQRS-lite**:
- **Commands (`internal/application/commands/`):** Handlers that mutate state (e.g., `RegisterStudent`, `PublishFeeStructure`, `GenerateBatchInvoices`, `ProcessPayment`). Every command executes within a transaction managed by the `UnitOfWork`.
- **Queries (`internal/application/queries/`):** Handlers that retrieve data (e.g., `GetStudentLedger`, `ListInvoices`, `GetPaymentReceipt`). Query handlers bypass domain aggregates when necessary to run high-performance relational joins and return tailored read models.
- **Application Services (`internal/application/services/`):** Cross-cutting interfaces such as `auditlog.Recorder`, `clock.Clock`, `storage.BlobStore`, and `idgen.IDGenerator`.

---

### 3.3 Persistence & Unit of Work (`internal/infrastructure/database/`)
To preserve consistency across multi-table aggregate updates (e.g., registering a payment, calculating invoice allocations, updating student credit balance, generating an audit entry, and appending a sync outbox entry), the system uses a **Unit of Work (`uow.UnitOfWork`)** pattern:

```go
err := uow.ExecuteInTx(ctx, func(txCtx context.Context) error {
    // 1. Record payment aggregate
    // 2. Perform FIFO allocations on open invoices
    // 3. Update student credit ledger
    // 4. Record audit log entry
    // 5. Append to sync_outbox
    return nil
})
```

---

## 4. SQLite Concurrency & Database Engine

Desktop SQLite in concurrent environments often suffers from `SQLITE_BUSY` lock contention when writers and readers collide. The system avoids this completely through deliberate architectural choices:

### 4.1 Pure Go Driver
We use `modernc.org/sqlite` instead of `mattn/go-sqlite3`.
- **Zero CGO:** Compiles natively across all platforms without requiring gcc, MinGW, or external C compilers.
- **Self-Contained:** Simplifies automated CI builds and Windows installer generation.

### 4.2 Dual-Pool Architecture
SQLite allows multiple concurrent readers **or** one active writer when in WAL mode:
1. **Dedicated Writer Pool:**
   - `MaxOpenConns = 1`
   - `MaxIdleConns = 1`
   - Connection URL: `_txlock=immediate`
   - All write transactions take immediate exclusive locks, guaranteeing that concurrent write operations queue deterministically in Go memory rather than timing out in SQLite.
2. **Reader Pool:**
   - `MaxOpenConns = Runtime.NumCPU * 2`
   - Concurrently serves queries, dashboard analytics, and reports without waiting on the writer.

### 4.3 SQLite Pragmas
Every SQLite connection initializes with:
```sql
PRAGMA journal_mode = WAL;
PRAGMA busy_timeout = 5000;
PRAGMA foreign_keys = ON;
PRAGMA synchronous = NORMAL;
PRAGMA cache_size = -64000; -- 64MB memory cache
```

---

## 5. Transport Abstraction Layer

The frontend application never couples directly to Wails runtime bindings. Instead, it interacts exclusively with TypeScript service contracts.

### Frontend Service Contract:
```typescript
export interface PaymentService {
  processPayment(input: ProcessPaymentInput): Promise<PaymentReceiptDto>;
  getReceipt(receiptNumber: string): Promise<PaymentReceiptDto>;
  listPayments(filters: PaymentFilters): Promise<PaginatedResult<PaymentSummaryDto>>;
  reversePayment(input: ReversePaymentInput): Promise<void>;
  createAdjustment(input: CreateAdjustmentInput): Promise<FinancialAdjustmentDto>;
}
```

### Transport Implementations:
- **`WailsPaymentService` (Desktop):** Calls Wails runtime IPC bindings (`window.go.main.PaymentService.ProcessPayment`).
- **`HttpPaymentService` (Cloud / LAN):** Issues standard HTTPS JSON requests (`fetch('/api/v1/payments')`).

A centralized **Service Factory** dynamically instantiates the correct service implementation based on the application's runtime configuration (`APP_MODE=desktop` or `APP_MODE=web`).

---

## 6. Offline-First & Future Cloud Synchronization

Although deployed as a standalone desktop application in Phase 1, the data layer includes all building blocks required for multi-device sync (Phase 7):

1. **UUIDv7 Primary Keys:** All entity primary keys are time-ordered UUIDv7 identifiers. They can be created locally on any disconnected client without risk of key collision.
2. **Transactional Outbox (`sync_outbox`):**
   - Every mutation creates a row in `sync_outbox` capturing the entity type, record ID, operation (`INSERT`, `UPDATE`, `DELETE`), and canonical JSON payload.
   - When offline, outbox events accumulate safely.
   - When connected to a central server, a background worker streams unacknowledged events in chronological order.
3. **Optimistic Locking:** Aggregate roots maintain integer `version` numbers to detect concurrent modifications and initiate deterministic conflict resolution workflows.

---

## 7. Security & Session Model

- **Desktop Session:** Authentication states are held in protected Go application memory, not in the Chromium webview's insecure `localStorage` or `sessionStorage`.
- **Cryptographic Passwords:** Passwords are salted and hashed using **Argon2id** (with standard OWASP-recommended memory and iteration parameters).
- **Audit Logging:** Every user action is recorded with an immutable timestamp, user identifier, IP/machine host, action verb, and structured payload.

