# School Management System (SMS)

> **Offline-first, cloud-ready desktop school management platform built with Go, Wails v2, SQLite, and Vue 3.**

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Wails](https://img.shields.io/badge/Wails-v2.15+-DF0000?style=flat)](https://wails.io)
[![Vue 3](https://img.shields.io/badge/Vue-3.5+-4FC08D?style=flat&logo=vue.js)](https://vuejs.org)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.9+-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org)
[![Tailwind CSS](https://img.shields.io/badge/TailwindCSS-v4.0-38B2AC?style=flat&logo=tailwind-css)](https://tailwindcss.com)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

---

## Overview

The **School Management System (SMS)** is an offline-first Windows desktop application architected for schools that require complete, uninterrupted operational capability without dependable internet connectivity, while remaining structurally prepared for future multi-device LAN and cloud synchronization (PostgreSQL).

Engineered following **Clean / Hexagonal Architecture** and **CQRS-lite** principles, the domain and application layers remain entirely decoupled from the presentation layer (Wails / Vue 3) and persistence mechanism (SQLite / PostgreSQL).

---

## Key Features

### Identity, Campus & Academics
- **Multi-Campus & Term Management:** Support for campuses, academic years, terms/semesters, classes, and sections.
- **Role-Based Access Control (RBAC):** Granular permissions for Superadmins, School Admins, Accountants, and Registrars.
- **Audit Logging:** Every state-modifying action records actor, IP, timestamp, entity type, and before/after payloads.

### Students, Guardians & Dynamic Admissions
- **Comprehensive Student Profiles:** Demographic details, blood group, medical alerts, admission numbers, and enrollment records.
- **Multi-Guardian Support:** Link multiple guardians per student (parents, emergency contacts) with prioritized billing and pickup authorizations.
- **Configurable Admission Form Engine:** Visual drag-and-drop schema definition for admission forms with field-level validations, versioning, drafts, publishing, and dynamic response storage.
- **Content-Addressed Document Storage:** Student documents and identity proofs are cryptographically hashed (`SHA-256`) and stored deduplicated on the local filesystem.

### Fees Engine & Batch Invoicing
- **Modular Fee Heads:** Configurable recurring (Tuition, Facilities) and one-off (Admission, Lab, Sports) heads with strict currency minor unit precision (`int64`).
- **Immutable Fee Schedules:** Versioned fee structures that protect historical financial records from retroactive corruption.
- **Multi-Tier Discount Policies:** Automated evaluation of sibling discounts (e.g., 20% off for 2nd+ child), academic merit scholarships, and staff concessions.
- **Batch Invoicing:** Deterministic mass generation of invoices for entire classes/terms with idempotent execution.
- **Calculated Overdue Tracking:** Overdue statuses are derived on the fly based on due dates and school timezone to prevent midnight status mutation race conditions.

### Payments, Sequential Receipts & Advance Credit
- **Strict Sequential Receipt Numbering:** Unbroken sequence receipts per calendar year (e.g., `RCP-2026-000001`) enforced via atomic database sequence counters.
- **Automated FIFO Allocations:** Payments are intelligently allocated against the student's oldest overdue invoices first before servicing newer charges.
- **Advance Credit Ledger:** Overpayments automatically flow into a dedicated student credit balance, instantly available for allocation against future invoices.
- **Financial Adjustments:** Formal workflows for fee waivers, bad-debt write-offs, and full/partial payment reversals with auditable reason codes.

### Offline-First Architecture
- **Pure Go SQLite:** Powered by `modernc.org/sqlite` (no CGO or MinGW required on Windows).
- **Concurrent WAL Mode:** Write-Ahead Logging with dedicated single-writer connection (`MaxOpenConns=1`, `_txlock=immediate`) and concurrent reader pool to guarantee zero database locks (`SQLITE_BUSY`).
- **Transactional Outbox:** Every database mutation writes to a `sync_outbox` table in the exact same transaction, establishing the foundation for bidirectional cloud synchronization.

---

## System Architecture

```text
Vue 3 Frontend (Composition API + TypeScript + Tailwind v4)
      │
      ▼
Transport Abstraction (Frontend Service Interfaces)
 ┌──────────────────────────┴──────────────────────────┐
 │                                                     │
 ▼                                                     ▼
Wails Desktop Adapter (IPC)               HTTP REST Adapter (Chi)
 └──────────────────────────┬──────────────────────────┘
                            │
                            ▼
           Go Application Layer (CQRS-lite)
         ┌──────────────────┴──────────────────┐
         ▼                                     ▼
      Commands                              Queries
 (Write Use Cases)                      (Read Use Cases)
         │                                     │
         ▼                                     │
    Domain Layer                               │
 (Aggregates, Policies,                        │
  Money & Basis Points)                        │
         │                                     │
         ▼                                     │
   Unit of Work & Repository Interfaces ◄──────┘
         │
         ▼
Persistence Adapters (Embedded Goose Migrations)
 ┌──────────────────────────┴──────────────────────────┐
 │                                                     │
 ▼                                                     ▼
SQLite Adapter (WAL Mode)               PostgreSQL Adapter (Cloud)
```

For complete architectural details, see [docs/architecture.md](docs/architecture.md) and [docs/database.md](docs/database.md).

---

## Technology Stack

| Layer | Technologies |
|---|---|
| **Desktop Shell** | [Wails v2](https://wails.io) (v2.15+) |
| **Backend Core** | Go 1.26+, Hexagonal Architecture, CQRS-lite |
| **Database** | SQLite 3 (WAL mode via `modernc.org/sqlite`), Goose Migrations |
| **Frontend Framework** | [Vue 3](https://vuejs.org) (Composition API, `<script setup lang="ts">`), Vite 5 |
| **Styling & UI** | [Tailwind CSS v4](https://tailwindcss.com), Lucide Icons, shadcn-vue patterns |
| **State & Data** | [Pinia](https://pinia.vuejs.org), [Vue Router](https://router.vuejs.org) |
| **Validation** | [Zod](https://zod.dev) & [VeeValidate](https://vee-validate.logaretm.com) (Frontend), Go Domain Invariants (Backend) |
| **Data Integrity** | UUIDv7 primary keys, `int64` minor units (Zero float arithmetic), integer basis points |

---

## Quick Start Guide

### Prerequisites
Ensure the following tools are installed on your machine:
- **Go:** `1.22+` (Go `1.26+` tested)
- **Node.js:** `18.x` or higher (`v24.x` recommended) & `npm`
- **Wails CLI v2:** Install via:
  ```powershell
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

---

### 1. Database Setup: Fresh Migration & Seeding
The project includes a built-in CLI tool to wipe, migrate, and seed the SQLite database in the standard OS AppData directory (`%LOCALAPPDATA%\SchoolManagement\data\school.db`):

```powershell
go run ./cmd/tools/db -fresh -seed
```

**Output will confirm:**
- Existing database and WAL/SHM files cleared.
- Goose applied all database migrations (00001 through 00004).
- Comprehensive demo school seeded with:
  - **School:** Oakridge International Academy
  - **Superadmin:** Username: `admin` | Password: `Password123!`
  - **Academic Structure:** 2026-2027 Academic Year, Fall Term, Grades 1-3, Sections A & B.
  - **Students & Guardians:** Alice, Bob, Charlie, and Diana with active enrollments and primary guardians.
  - **Fees & Discounts:** Tuition & Lab fee schedules, 20% sibling discount policy.
  - **Billing & Invoices:** September 2026 batch invoices generated.
  - **Payments & Receipts:** Sample receipts `RCP-2026-000001` and `RCP-2026-000002` (with credit balance advance).

---

### 2. Running in Development Mode
Run the live development environment with hot module replacement (HMR) for the frontend and automatic recompilation of Go code:

```powershell
wails dev
```

*The application window will open automatically. Frontend changes update instantly.*

---

### 3. Compiling Standalone Executables

#### Direct Go Compilation:
```powershell
# 1. Build frontend distribution bundle
npm --prefix frontend run build

# 2. Build standalone Windows binary
go build -o desktop.exe ./cmd/desktop

# 3. Run
.\desktop.exe
```

#### Production Wails Build (with Windows Installer):
```powershell
wails build
```
*The compiled binary is placed in `build/bin/school-management.exe`.*

---

## Testing

The system features complete end-to-end integration test suites using in-memory isolated SQLite instances:

```powershell
# Run all integration tests
go test -v ./tests/integration/...
```

Test coverage includes:
- `TestAuthAndUserLifecycle`: User creation, authentication, password hashing, and session verification.
- `TestAcademicStudentAdmissionLifecycle`: Dynamic admission forms, published revisions, student enrollments, and status transitions.
- `TestFeeEngineAndInvoicingLifecycle`: Fee structures, multi-condition discount policies, and batch invoice generation.
- `TestPaymentLifecycleAndCreditEngine`: Sequential receipt generation, automated FIFO allocations, credit advances, and financial adjustments.

---

## Repository Structure

```text
school-management/
├── cmd/
│   ├── desktop/          # Desktop application entry point (Wails runtime)
│   ├── server/           # Future HTTP server entry point (LAN/Cloud)
│   └── tools/            # Developer tools & utilities (db fresh & seed)
├── database/
│   ├── migrations/       # Versioned SQL migrations (Goose)
│   │   ├── sqlite/       # SQLite migrations embedded in Go binary
│   │   └── postgres/     # PostgreSQL migrations for cloud sync
│   └── queries/          # Raw SQL and sqlc query files
├── docs/                 # In-depth architectural and module documentation
│   ├── architecture.md   # Hexagonal design, CQRS, and sync outbox
│   ├── database.md       # Schema design, UUIDv7, monetary rules, indexing
│   ├── modules/          # Deep-dive module documentation
│   │   ├── academic_admissions.md
│   │   ├── fees_and_billing.md
│   │   └── payments_and_credits.md
│   └── guides/           # Operational guides
│       ├── development.md
│       └── deployment.md
├── frontend/             # Vue 3 + TypeScript + Tailwind v4 application
│   ├── src/
│   │   ├── components/   # Reusable UI components & layouts
│   │   ├── modules/      # Business logic components
│   │   ├── pages/        # View routes (Auth, Academics, Students, Fees, Payments)
│   │   ├── services/     # Transport abstraction layer (Wails / HTTP)
│   │   ├── stores/       # Pinia reactive state stores
│   │   └── types/        # TypeScript DTOs and entity models
│   └── package.json
├── internal/             # Private application source code
│   ├── domain/           # Pure enterprise business logic & aggregates
│   ├── application/      # Use-case handlers (Commands & Queries)
│   ├── infrastructure/   # SQLite database adapters, blob storage, logging
│   ├── platform/         # Configuration, cryptography, and OS paths
│   └── transport/        # Wails IPC and HTTP presentation adapters
├── tests/
│   └── integration/      # End-to-end SQLite integration tests
├── wails.json            # Wails v2 project configuration
└── go.mod
```

---

## =Documentation Index

- [System Architecture](docs/architecture.md) — Hexagonal design, layers, outbox pattern, and thread-safe SQLite concurrency.
- [Database Design](docs/database.md) — Complete relational schema, UUIDv7 indexing, and minor currency units.
- [Academics & Dynamic Admissions](docs/modules/academic_admissions.md) — Campus, terms, students, and form schema builder.
- [Fees Engine & Batch Invoicing](docs/modules/fees_and_billing.md) — Versioned fees, discount evaluation, and billing periods.
- [Payments, Receipts & Credits](docs/modules/payments_and_credits.md) — Atomic receipt sequence, FIFO allocations, and credit ledger.
- [Developer Guide](docs/guides/development.md) — Detailed instructions on contributing, testing, and debugging.
- [Deployment & Backups](docs/guides/deployment.md) — Installer packaging, local app data locations, and data safety.

---

## 👨‍💻 Author & Engineering

Designed, architected, and built with precision by **Syed Faizan Hassaan**.

- **GitHub:** [@zfhassaan](https://github.com/zfhassaan) (`https://github.com/zfhassaan`)
- **Portfolio:** [zfhassaan.github.io](https://zfhassaan.github.io) (`https://zfhassaan.github.io`)
- **Role:** Full Stack Software Architect & Systems Engineer

[![GitHub Profile](https://img.shields.io/badge/GitHub-zfhassaan-181717?style=flat&logo=github)](https://github.com/zfhassaan)
[![Portfolio Website](https://img.shields.io/badge/Portfolio-zfhassaan.github.io-4F46E5?style=flat&logo=safari)](https://zfhassaan.github.io)

---

## License
This project is licensed under the MIT License.

