# Desktop School Management System (SMS)

> **Enterprise-grade, offline-first desktop school management software with AES-256 encrypted database at rest, executive printable admission forms, automated fee invoicing, and in-app GitHub auto-updater.**

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Wails](https://img.shields.io/badge/Wails-v2.15+-DF0000?style=flat)](https://wails.io)
[![Vue 3](https://img.shields.io/badge/Vue-3.5+-4FC08D?style=flat&logo=vue.js)](https://vuejs.org)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.9+-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org)
[![Tailwind CSS](https://img.shields.io/badge/TailwindCSS-v4.0-38B2AC?style=flat&logo=tailwind-css)](https://tailwindcss.com)
[![Encryption: AES-256-GCM](https://img.shields.io/badge/Security-AES--256--GCM-green.svg)](https://en.wikipedia.org/wiki/Galois/Counter_Mode)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

---

## Overview

The **Desktop School Management System (SMS)** is a secure, standalone Windows application engineered for primary, secondary, and higher-education institutions. Designed with an **offline-first philosophy**, it ensures complete operational autonomy without requiring an internet connection, while featuring local **AES-256-GCM encryption at rest**, automated hot-patching via GitHub Releases, and zero-headache local installations.

Built using **Go (Golang)**, **Wails v2**, **Vue 3**, and **SQLite (ModernC)** following **Clean Architecture (Hexagonal) & CQRS-lite**, the application provides lightning-fast native performance with minimal memory usage.

---

## Core Capabilities

### 🛡️ Security & Zero-Leak Encrypted Vault
- **AES-256-GCM Encryption at Rest:** The local SQLite database is encrypted on disk into `.db.enc` vaults. Unauthorized users cannot inspect plain SQL tables or extract student records.
- **Hardware-Entropy Key Generation:** Cryptographic keys are generated deterministically per installation using Windows DPAPI and machine hardware signatures.
- **Decrypted Backup & Restore Center:** Administrators can export decrypted SQL snapshots on-demand for manual auditing, off-site archiving, or migration, with safe encrypted restore mechanisms.

### 🎓 Dynamic Admissions & Executive Print Engine
- **Visual Form Builder:** Drag-and-drop form schema builder with real-time field validations, drafts, and versioning.
- **Executive Institutional Print Format:** Official admission forms formatted with institutional double borders, student photo boxes, dynamic school logos, structured parent/student data tables, and perforated tear-off verification slips.
- **Silent Print Optimization:** Automatic browser print header/footer suppression and exact CMYK print color preservation.

### 🏫 Dynamic Branding & Windows Integration
- **Live Taskbar & Window Icon Switching:** Uploading a new school logo dynamically regenerates multi-resolution `.ico` bundles and hot-updates the Windows Taskbar, titlebar, and webview favicon at runtime.
- **Local Asset Storage:** Logos, student photos, and identity documents are SHA-256 hashed and stored deduplicated inside the secure local installation directory.

### 💰 Fees Engine, Invoicing & FIFO Payments
- **Batch Fee Generation:** Idempotent mass invoice generation for entire academic classes and terms.
- **Automated FIFO Allocations:** Payments systematically offset oldest overdue balances first before settling current charges.
- **Strict Sequential Receipt Numbering:** Unbroken receipt audit sequences (e.g., `RCP-2026-000001`) enforced via atomic database sequences.
- **Advance Credit Ledger:** Automatic overpayment rollover credited to future invoices.

### 🔄 In-App GitHub Releases Auto-Updater
- **One-Click Update Checker:** Queries the GitHub Releases API directly from the app header and settings menu.
- **In-App Changelog Viewer:** Live Markdown release notes preview.
- **Safe Windows Hot-Patching:** Renames locked running binaries (`.old`), installs new executables seamlessly, and restarts without requiring administrative command-line operations.

---

## Technology Stack

| Layer | Technology | Purpose |
| :--- | :--- | :--- |
| **Backend Core** | Go 1.26+ | Clean architecture domain logic, CQRS commands/queries |
| **Desktop Bridge** | Wails v2.15+ | High-performance Go-to-Webview IPC bridge |
| **Database** | SQLite (`modernc.org/sqlite`) | Pure Go WAL-mode embedded SQL with zero CGO dependencies |
| **Data Security** | AES-256-GCM | Encrypted database vault with machine-level key derivation |
| **Frontend** | Vue 3 + TypeScript | Composition API, Pinia stores, Vue Router |
| **Styling** | Tailwind CSS v4.0 | Responsive institutional design, optimized print stylesheet |

---
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

Designed, architected, and built with precision by **Hassaan Ali**.

- **GitHub:** [@zfhassaan](https://github.com/zfhassaan) (`https://github.com/zfhassaan`)
- **Portfolio:** [zfhassaan.github.io](https://zfhassaan.github.io) (`https://zfhassaan.github.io`)
- **Role:** Full Stack Software Architect & Systems Engineer

[![GitHub Profile](https://img.shields.io/badge/GitHub-zfhassaan-181717?style=flat&logo=github)](https://github.com/zfhassaan)
[![Portfolio Website](https://img.shields.io/badge/Portfolio-zfhassaan.github.io-4F46E5?style=flat&logo=safari)](https://zfhassaan.github.io)

---

## License
This project is licensed under the MIT License.

