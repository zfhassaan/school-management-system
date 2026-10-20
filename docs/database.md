# 💾 Database Architecture & Schema Specification

## 1. Design Principles & Conventions

The database schema is engineered for high-throughput relational integrity, deterministic financial auditing, and multi-device cloud synchronization.

### 1.1 Identifiers: UUIDv7
- All table primary keys use **UUIDv7** (RFC 9562).
- **SQLite Storage:** Stored as `TEXT(36)` strings.
- **PostgreSQL Storage:** Stored as native `uuid` types.
- **Why UUIDv7?** Unlike random UUIDv4, UUIDv7 embeds a 48-bit millisecond Unix timestamp in its most significant bits. B-Tree index inserts remain sequential and localized, preventing index fragmentation while allowing disconnected clients to generate collision-free IDs offline.

### 1.2 Monetary Accuracy: Zero Floating-Point Numbers
- All monetary amounts are stored as **64-bit integers (`int64`) in minor currency units** (e.g., cents, pence, paise).
  - `$100.00` is stored as `10000`.
  - `$0.50` is stored as `50`.
- Percentages (discounts, taxes, late fees) are stored as **integer basis points (`bps`)**:
  - `1%` = `100 bps`
  - `20%` = `2000 bps`
  - `100%` = `10000 bps`
- Floating-point data types (`REAL`, `FLOAT`, `DOUBLE`) are strictly forbidden in domain and persistence logic.

### 1.3 Time & Date Representation
- **Instants (Audit timestamps, transaction times):** Stored in **UTC ISO-8601** format with millisecond precision: `YYYY-MM-DDTHH:MM:SS.sssZ`.
- **Civil Dates (Birthdays, due dates, billing periods):** Stored as strict **`YYYY-MM-DD`** strings (length 10).
- **Timezone Invariance:** Civil dates are immune to daylight savings time and timezone shifts. An invoice due on `2026-09-30` remains due on that exact date regardless of whether viewed in London, New York, or Tokyo.

---

## 2. Entity-Relationship Model

```mermaid
erDiagram
    SCHOOL ||--o{ CAMPUS : has
    SCHOOL ||--o{ ACADEMIC_YEAR : configures
    ACADEMIC_YEAR ||--o{ ACADEMIC_TERM : divides_into
    ACADEMIC_YEAR ||--o{ CLASS : offers
    CLASS ||--o{ SECTION : contains
    
    STUDENT ||--o{ STUDENT_GUARDIAN : links
    GUARDIAN ||--o{ STUDENT_GUARDIAN : links
    STUDENT ||--o{ ENROLLMENT : has
    SECTION ||--o{ ENROLLMENT : assigns
    
    FEE_HEAD ||--o{ FEE_STRUCTURE_ITEM : categorized_in
    FEE_STRUCTURE ||--o{ FEE_STRUCTURE_VERSION : versions
    FEE_STRUCTURE_VERSION ||--o{ FEE_STRUCTURE_ITEM : defines
    
    STUDENT ||--o{ INVOICE : billed_to
    INVOICE ||--o{ INVOICE_ITEM : itemizes
    
    STUDENT ||--o{ PAYMENT : makes
    PAYMENT ||--o{ PAYMENT_ALLOCATION : allocates_to
    INVOICE ||--o{ PAYMENT_ALLOCATION : receives
    STUDENT ||--o{ CREDIT_BALANCE : maintains
    STUDENT ||--o{ FINANCIAL_ADJUSTMENT : granted
```

---

## 3. Schema Manifest by Module

### 3.1 Core & Identity
- **`schools`:** School root entity, configuration JSON, default currency code, and timezone.
- **`campuses`:** Campus locations, contact info, and status.
- **`users`:** Staff and admin accounts, email, Argon2id password hash, and active status.
- **`roles`:** RBAC role definitions (`superadmin`, `school_admin`, `accountant`, `registrar`).
- **`permissions`:** System permissions (e.g., `students:create`, `fees:publish`, `payments:reverse`).
- **`role_permissions` & `user_roles`:** Join tables establishing role assignments and permission sets.
- **`audit_logs`:** Append-only log recording actor ID, IP, resource type, action, diff JSON, and timestamp.

### 3.2 Academics
- **`academic_years`:** Calendar/fiscal schooling cycles (e.g., `2026-2027`) with start and end dates.
- **`academic_terms`:** Sub-divisions (e.g., Fall Term, Spring Term) with distinct deadlines.
- **`classes`:** Grades/standards offered (e.g., `Grade 1`, `Grade 2`, `Grade 10`).
- **`sections`:** Classroom sub-divisions (e.g., `Section A`, `Section B`) with student capacity limits.

### 3.3 Students, Guardians & Admissions
- **`students`:** Student demographic records, admission numbers, medical notes, and status (`applied`, `enrolled`, `promoted`, `withdrawn`, `graduated`, `suspended`).
- **`guardians`:** Parent/guardian profiles with national IDs, emails, and contact numbers.
- **`student_guardians`:** Association table detailing relationship (`father`, `mother`, `guardian`), pickup authorization, primary billing contact flag, and emergency priority.
- **`enrollments`:** Historical mapping of students to classes and sections for specific academic terms.
- **`admission_form_definitions`:** Metadata for admission forms.
- **`admission_form_versions`:** Versioned JSON schemas defining fields, sections, rules, and document requirements.
- **`admission_applications`:** Dynamic responses submitted by applicants against specific form versions.
- **`student_documents`:** Content-addressed file metadata (`sha256` checksum, MIME type, byte size, file path).

### 3.4 Fees & Invoicing
- **`fee_heads`:** Categorical definitions of school fees (`Tuition`, `Lab Fee`, `Sports`, `Transportation`).
- **`fee_structures`:** Class-level fee containers.
- **`fee_structure_versions`:** Immutable published versions of fee schedules with effective date ranges.
- **`fee_structure_items`:** Individual fee heads and their corresponding minor amounts within a structure version.
- **`discount_policies`:** Automated discount rule criteria (sibling counts, merit grades, staff children).
- **`student_discounts`:** Assigned discounts linked to specific students.
- **`invoices`:** Master billing records containing total amount, discount amount, paid amount, adjusted amount, balance, status (`draft`, `issued`, `partially_paid`, `paid`, `void`, `cancelled`), and due date.
- **`invoice_items`:** Line-item breakdown of invoice charges tied back to original fee heads.

### 3.5 Payments, Receipts & Credits
- **`payments`:** Immutable records of payments processed via cash, bank transfer, cheque, or POS.
- **`receipt_sequences`:** Atomic counter table guaranteeing unbroken calendar-year sequential receipt numbering (`RCP-YYYY-XXXXXX`).
- **`payment_allocations`:** Detailed ledger tracing which exact invoice(s) a payment was applied to.
- **`credit_balances`:** Student advance balance ledger holding funds received in excess of outstanding dues.
- **`payment_reversals`:** Formal reversal tracking for bounced cheques or erroneous entries.
- **`financial_adjustments`:** Auditable waivers, discounts, or write-offs applied against invoices.

### 3.6 Offline Synchronization & Outbox
- **`sync_outbox`:** Transactional outbox capturing all database mutations for asynchronous sync streaming.
- **`sync_checkpoints`:** High-water marks recording the last synced transaction ID between node and server.
- **`sync_conflicts`:** Detailed ledger for recording two-way conflicts requiring human or deterministic resolution.
- **`sync_tombstones`:** Soft-delete markers ensuring deletions propagate across offline clients.

---

## 4. Migration Strategy

Database migrations are managed via **Goose** and embedded directly into the Go executable via `embed.FS`.

### Migration Sequence:
```text
database/migrations/sqlite/
├── 00001_initial_schema.sql                  (Core, Identity, Settings, Outbox)
├── 00002_academics_students_admissions.sql   (Academics, Students, Forms, Docs)
├── 00003_fees_and_invoicing.sql              (Fee Heads, Structures, Invoices)
└── 00004_payments_allocations_adjustments.sql(Payments, Receipts, Credit Ledger)
```

On desktop launch, `sqlite.Open` automatically checks the current version and applies pending migrations inside a transaction before starting the application window.

---

## 5. Performance Indexing & Query Optimization

To guarantee sub-10ms response times on datasets with tens of thousands of student records:

1. **Student Search Index:**
   ```sql
   CREATE INDEX idx_students_search ON students(school_id, status, last_name, first_name);
   CREATE INDEX idx_students_admission_no ON students(school_id, admission_number);
   ```
2. **Guardian Phone Lookup Index:**
   ```sql
   CREATE INDEX idx_guardians_phone ON guardians(school_id, phone_primary);
   ```
3. **Calculated Overdue Query Index:**
   ```sql
   CREATE INDEX idx_invoices_due ON invoices(school_id, status, due_date);
   ```
4. **Receipt Sequence Lock:**
   The `receipt_sequences` table enforces atomic increment per year using:
   ```sql
   INSERT INTO receipt_sequences (school_id, calendar_year, last_sequence)
   VALUES (?, ?, 1)
   ON CONFLICT(school_id, calendar_year)
   DO UPDATE SET last_sequence = receipt_sequences.last_sequence + 1
   RETURNING last_sequence;
   ```
