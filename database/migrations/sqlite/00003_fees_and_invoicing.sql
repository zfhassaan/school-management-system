-- +goose Up
-- +goose StatementBegin

-- ============================================================================
-- Phase 4: Fee Engine & Invoicing Architecture
-- Single-school offline-first schema with integer minor units & basis points.
-- ============================================================================

-- Fee Heads: Master categories of fees (Tuition, Admission, Transport, Lab, Exam, Library)
CREATE TABLE fee_heads (
    id TEXT PRIMARY KEY NOT NULL,
    school_id TEXT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    is_recurring INTEGER NOT NULL DEFAULT 1,
    default_amount_minor INTEGER NOT NULL DEFAULT 0,
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    row_version INTEGER NOT NULL DEFAULT 1,
    UNIQUE(school_id, code)
);

CREATE INDEX idx_fee_heads_school_active ON fee_heads(school_id, is_active);

-- Fee Structures: Grouping of fee configurations for specific class, academic year, campus
CREATE TABLE fee_structures (
    id TEXT PRIMARY KEY NOT NULL,
    school_id TEXT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    academic_year_id TEXT NOT NULL REFERENCES academic_years(id) ON DELETE CASCADE,
    class_id TEXT NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    campus_id TEXT REFERENCES campuses(id) ON DELETE SET NULL,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    row_version INTEGER NOT NULL DEFAULT 1,
    UNIQUE(academic_year_id, class_id, campus_id, code)
);

CREATE INDEX idx_fee_structures_lookup ON fee_structures(school_id, academic_year_id, class_id);

-- Fee Structure Versions: Versioned snapshot of fees with draft/published state
CREATE TABLE fee_structure_versions (
    id TEXT PRIMARY KEY NOT NULL,
    fee_structure_id TEXT NOT NULL REFERENCES fee_structures(id) ON DELETE CASCADE,
    version_no INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'draft' CHECK(status IN ('draft', 'published', 'retired')),
    changelog TEXT NOT NULL DEFAULT '',
    created_by TEXT NOT NULL REFERENCES users(id),
    published_by TEXT REFERENCES users(id),
    published_at TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    row_version INTEGER NOT NULL DEFAULT 1,
    UNIQUE(fee_structure_id, version_no)
);

CREATE INDEX idx_fee_struct_ver_status ON fee_structure_versions(fee_structure_id, status);

-- Trigger: Prevent modification of published fee_structure_versions
CREATE TRIGGER trg_fee_structure_versions_prevent_update_published
BEFORE UPDATE ON fee_structure_versions
WHEN OLD.status = 'published' AND NEW.status != 'retired'
BEGIN
    SELECT RAISE(ABORT, 'Published fee structure versions are immutable');
END;

-- Fee Structure Items: Line items specifying fee head, amount, billing frequency
CREATE TABLE fee_structure_items (
    id TEXT PRIMARY KEY NOT NULL,
    fee_structure_version_id TEXT NOT NULL REFERENCES fee_structure_versions(id) ON DELETE CASCADE,
    fee_head_id TEXT NOT NULL REFERENCES fee_heads(id) ON DELETE RESTRICT,
    amount_minor INTEGER NOT NULL CHECK(amount_minor >= 0),
    frequency TEXT NOT NULL CHECK(frequency IN ('monthly', 'termly', 'annual', 'one_time')),
    due_day INTEGER NOT NULL DEFAULT 10 CHECK(due_day BETWEEN 1 AND 31),
    is_optional INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    UNIQUE(fee_structure_version_id, fee_head_id)
);

CREATE INDEX idx_fee_struct_items_ver ON fee_structure_items(fee_structure_version_id);

-- Trigger: Prevent modification of items on published versions
CREATE TRIGGER trg_fee_structure_items_prevent_mod_published
BEFORE INSERT ON fee_structure_items
WHEN (SELECT status FROM fee_structure_versions WHERE id = NEW.fee_structure_version_id) = 'published'
BEGIN
    SELECT RAISE(ABORT, 'Cannot add items to a published fee structure version');
END;

-- Discount Policies: Global rules for concessions (Sibling, Merit Scholarship, Staff Child)
CREATE TABLE discount_policies (
    id TEXT PRIMARY KEY NOT NULL,
    school_id TEXT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    discount_type TEXT NOT NULL CHECK(discount_type IN ('percentage', 'fixed')),
    rate_basis_points INTEGER NOT NULL DEFAULT 0 CHECK(rate_basis_points BETWEEN 0 AND 10000),
    amount_minor INTEGER NOT NULL DEFAULT 0 CHECK(amount_minor >= 0),
    fee_head_id TEXT REFERENCES fee_heads(id) ON DELETE SET NULL,
    requires_approval INTEGER NOT NULL DEFAULT 1,
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    row_version INTEGER NOT NULL DEFAULT 1,
    UNIQUE(school_id, code)
);

CREATE INDEX idx_discount_policies_school ON discount_policies(school_id, is_active);

-- Student Discounts: Concessions assigned to a student
CREATE TABLE student_discounts (
    id TEXT PRIMARY KEY NOT NULL,
    student_id TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    discount_policy_id TEXT NOT NULL REFERENCES discount_policies(id) ON DELETE RESTRICT,
    academic_year_id TEXT NOT NULL REFERENCES academic_years(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'approved', 'rejected', 'revoked')),
    approved_by TEXT REFERENCES users(id),
    approved_at TEXT,
    rejection_reason TEXT NOT NULL DEFAULT '',
    valid_from TEXT NOT NULL,
    valid_to TEXT,
    created_by TEXT NOT NULL REFERENCES users(id),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    row_version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_student_discounts_lookup ON student_discounts(student_id, academic_year_id, status);

-- Invoices: Financial billing bills issued to students
CREATE TABLE invoices (
    id TEXT PRIMARY KEY NOT NULL,
    school_id TEXT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id TEXT NOT NULL REFERENCES students(id) ON DELETE RESTRICT,
    academic_year_id TEXT NOT NULL REFERENCES academic_years(id) ON DELETE RESTRICT,
    invoice_no TEXT NOT NULL,
    issue_date TEXT NOT NULL,
    due_date TEXT NOT NULL,
    billing_period TEXT NOT NULL,
    subtotal_amount_minor INTEGER NOT NULL CHECK(subtotal_amount_minor >= 0),
    discount_amount_minor INTEGER NOT NULL DEFAULT 0 CHECK(discount_amount_minor >= 0),
    total_amount_minor INTEGER NOT NULL CHECK(total_amount_minor >= 0),
    paid_amount_minor INTEGER NOT NULL DEFAULT 0 CHECK(paid_amount_minor >= 0),
    balance_minor INTEGER NOT NULL CHECK(balance_minor >= 0),
    status TEXT NOT NULL DEFAULT 'draft' CHECK(status IN ('draft', 'issued', 'partially_paid', 'paid', 'overdue', 'void')),
    void_reason TEXT NOT NULL DEFAULT '',
    voided_by TEXT REFERENCES users(id),
    voided_at TEXT,
    created_by TEXT NOT NULL REFERENCES users(id),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    row_version INTEGER NOT NULL DEFAULT 1,
    UNIQUE(school_id, invoice_no)
);

CREATE INDEX idx_invoices_student_year ON invoices(student_id, academic_year_id, status);
CREATE INDEX idx_invoices_school_status ON invoices(school_id, status, due_date);
CREATE INDEX idx_invoices_period ON invoices(school_id, billing_period, student_id);

-- Invoice Items: Granular item breakdown
CREATE TABLE invoice_items (
    id TEXT PRIMARY KEY NOT NULL,
    invoice_id TEXT NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    fee_head_id TEXT NOT NULL REFERENCES fee_heads(id) ON DELETE RESTRICT,
    fee_head_name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    original_amount_minor INTEGER NOT NULL CHECK(original_amount_minor >= 0),
    discount_amount_minor INTEGER NOT NULL DEFAULT 0 CHECK(discount_amount_minor >= 0),
    net_amount_minor INTEGER NOT NULL CHECK(net_amount_minor >= 0),
    created_at TEXT NOT NULL
);

CREATE INDEX idx_invoice_items_invoice ON invoice_items(invoice_id);

-- Trigger: Prevent modifying invoice items once invoice is issued/paid
CREATE TRIGGER trg_invoice_items_prevent_mod_issued
BEFORE INSERT ON invoice_items
WHEN (SELECT status FROM invoices WHERE id = NEW.invoice_id) NOT IN ('draft')
BEGIN
    SELECT RAISE(ABORT, 'Cannot add line items to an invoice that is already issued or finalized');
END;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_invoice_items_prevent_mod_issued;
DROP TABLE IF EXISTS invoice_items;
DROP TABLE IF EXISTS invoices;
DROP TABLE IF EXISTS student_discounts;
DROP TABLE IF EXISTS discount_policies;
DROP TRIGGER IF EXISTS trg_fee_structure_items_prevent_mod_published;
DROP TABLE IF EXISTS fee_structure_items;
DROP TRIGGER IF EXISTS trg_fee_structure_versions_prevent_update_published;
DROP TABLE IF EXISTS fee_structure_versions;
DROP TABLE IF EXISTS fee_structures;
DROP TABLE IF EXISTS fee_heads;
-- +goose StatementEnd

