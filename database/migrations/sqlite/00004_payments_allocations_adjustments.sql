-- +goose Up
-- +goose StatementBegin

-- 1. Receipt Sequences for atomic sequential receipt generation per year
CREATE TABLE receipt_sequences (
	school_id TEXT NOT NULL REFERENCES schools(id),
	year INTEGER NOT NULL,
	current_val INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY (school_id, year)
) STRICT;

-- 2. Student Credit Balances (Advance payments, overpayments, credit concessions)
CREATE TABLE credit_balances (
	id TEXT PRIMARY KEY,
	school_id TEXT NOT NULL REFERENCES schools(id),
	student_id TEXT NOT NULL REFERENCES students(id) ON DELETE RESTRICT,
	amount_minor INTEGER NOT NULL DEFAULT 0 CHECK (amount_minor >= 0),
	updated_at TEXT NOT NULL,
	created_at TEXT NOT NULL,
	UNIQUE (school_id, student_id)
) STRICT;

CREATE INDEX ix_credit_balances_student ON credit_balances(student_id);

-- 3. Student Credit Balance Ledger (Audit trail for every credit addition / drawdown)
CREATE TABLE credit_balance_ledger (
	id TEXT PRIMARY KEY,
	credit_balance_id TEXT NOT NULL REFERENCES credit_balances(id) ON DELETE RESTRICT,
	student_id TEXT NOT NULL REFERENCES students(id) ON DELETE RESTRICT,
	amount_minor INTEGER NOT NULL,
	balance_after_minor INTEGER NOT NULL CHECK (balance_after_minor >= 0),
	reference_type TEXT NOT NULL CHECK (reference_type IN ('overpayment', 'advance', 'invoice_application', 'refund', 'adjustment')),
	reference_id TEXT,
	notes TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL
) STRICT;

CREATE INDEX ix_credit_ledger_student ON credit_balance_ledger(student_id, created_at);

-- 4. Payments
CREATE TABLE payments (
	id TEXT PRIMARY KEY,
	school_id TEXT NOT NULL REFERENCES schools(id),
	student_id TEXT NOT NULL REFERENCES students(id) ON DELETE RESTRICT,
	receipt_no TEXT NOT NULL UNIQUE,
	payment_date TEXT NOT NULL,
	payment_method TEXT NOT NULL CHECK (payment_method IN ('cash', 'bank_transfer', 'cheque', 'pos_card', 'mobile_money', 'credit_balance')),
	amount_minor INTEGER NOT NULL CHECK (amount_minor > 0),
	allocated_amount_minor INTEGER NOT NULL DEFAULT 0 CHECK (allocated_amount_minor >= 0 AND allocated_amount_minor <= amount_minor),
	unallocated_amount_minor INTEGER NOT NULL DEFAULT 0 CHECK (unallocated_amount_minor >= 0),
	reference_no TEXT NOT NULL DEFAULT '',
	payer_name TEXT NOT NULL DEFAULT '',
	payer_phone TEXT NOT NULL DEFAULT '',
	status TEXT NOT NULL DEFAULT 'posted' CHECK (status IN ('posted', 'reversed')),
	recorded_by_id TEXT NOT NULL REFERENCES users(id),
	recorded_at TEXT NOT NULL,
	notes TEXT NOT NULL DEFAULT ''
) STRICT;

CREATE INDEX ix_payments_student ON payments(student_id, payment_date DESC);
CREATE INDEX ix_payments_receipt ON payments(receipt_no);
CREATE INDEX ix_payments_school_date ON payments(school_id, payment_date DESC);

-- 5. Payment Allocations (Allocates payment amount to specific invoices and invoice line items)
CREATE TABLE payment_allocations (
	id TEXT PRIMARY KEY,
	payment_id TEXT NOT NULL REFERENCES payments(id) ON DELETE RESTRICT,
	invoice_id TEXT NOT NULL REFERENCES invoices(id) ON DELETE RESTRICT,
	invoice_item_id TEXT REFERENCES invoice_items(id) ON DELETE RESTRICT,
	allocated_amount_minor INTEGER NOT NULL CHECK (allocated_amount_minor > 0),
	allocated_at TEXT NOT NULL,
	allocated_by_id TEXT NOT NULL REFERENCES users(id)
) STRICT;

CREATE INDEX ix_payment_allocations_payment ON payment_allocations(payment_id);
CREATE INDEX ix_payment_allocations_invoice ON payment_allocations(invoice_id);

-- 6. Payment Reversals
CREATE TABLE payment_reversals (
	id TEXT PRIMARY KEY,
	payment_id TEXT NOT NULL UNIQUE REFERENCES payments(id) ON DELETE RESTRICT,
	reversal_date TEXT NOT NULL,
	reason TEXT NOT NULL,
	reversed_by_id TEXT NOT NULL REFERENCES users(id),
	reversed_at TEXT NOT NULL
) STRICT;

-- 7. Financial Adjustments (Waivers, write-offs, corrections)
CREATE TABLE financial_adjustments (
	id TEXT PRIMARY KEY,
	school_id TEXT NOT NULL REFERENCES schools(id),
	student_id TEXT NOT NULL REFERENCES students(id) ON DELETE RESTRICT,
	invoice_id TEXT REFERENCES invoices(id) ON DELETE RESTRICT,
	invoice_item_id TEXT REFERENCES invoice_items(id) ON DELETE RESTRICT,
	adjustment_type TEXT NOT NULL CHECK (adjustment_type IN ('waiver', 'write_off', 'debit_correction', 'credit_correction')),
	amount_minor INTEGER NOT NULL CHECK (amount_minor > 0),
	reason TEXT NOT NULL,
	status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
	requested_by_id TEXT NOT NULL REFERENCES users(id),
	requested_at TEXT NOT NULL,
	decided_by_id TEXT REFERENCES users(id),
	decided_at TEXT,
	decision_notes TEXT NOT NULL DEFAULT ''
) STRICT;

CREATE INDEX ix_financial_adjustments_student ON financial_adjustments(student_id);
CREATE INDEX ix_financial_adjustments_invoice ON financial_adjustments(invoice_id);
CREATE INDEX ix_financial_adjustments_status ON financial_adjustments(status);

-- 8. Immutability Triggers
CREATE TRIGGER trg_payments_prevent_delete
BEFORE DELETE ON payments
BEGIN
	SELECT RAISE(ABORT, 'payments are financial records and cannot be deleted; use reversal workflow instead');
END;

CREATE TRIGGER trg_payment_allocations_prevent_update
BEFORE UPDATE ON payment_allocations
BEGIN
	SELECT RAISE(ABORT, 'payment allocations are immutable financial records and cannot be modified');
END;

CREATE TRIGGER trg_payment_allocations_prevent_delete
BEFORE DELETE ON payment_allocations
BEGIN
	SELECT RAISE(ABORT, 'payment allocations are immutable financial records and cannot be deleted');
END;

CREATE TRIGGER trg_payment_reversals_prevent_update
BEFORE UPDATE ON payment_reversals
BEGIN
	SELECT RAISE(ABORT, 'payment reversal records are permanent and immutable and cannot be modified');
END;

CREATE TRIGGER trg_payment_reversals_prevent_delete
BEFORE DELETE ON payment_reversals
BEGIN
	SELECT RAISE(ABORT, 'payment reversal records are permanent and immutable and cannot be deleted');
END;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_payment_reversals_prevent_delete;
DROP TRIGGER IF EXISTS trg_payment_reversals_prevent_update;
DROP TRIGGER IF EXISTS trg_payment_allocations_prevent_delete;
DROP TRIGGER IF EXISTS trg_payment_allocations_prevent_update;
DROP TRIGGER IF EXISTS trg_payments_prevent_delete;
DROP TABLE IF EXISTS financial_adjustments;
DROP TABLE IF EXISTS payment_reversals;
DROP TABLE IF EXISTS payment_allocations;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS credit_balance_ledger;
DROP TABLE IF EXISTS credit_balances;
DROP TABLE IF EXISTS receipt_sequences;
-- +goose StatementEnd
