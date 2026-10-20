-- name: GetAndIncrementReceiptSeq :one
INSERT INTO receipt_sequences (school_id, year, current_val)
VALUES (?, ?, 1)
ON CONFLICT (school_id, year)
DO UPDATE SET current_val = receipt_sequences.current_val + 1
RETURNING current_val;

-- name: GetCreditBalanceByStudent :one
SELECT * FROM credit_balances
WHERE school_id = ? AND student_id = ? LIMIT 1;

-- name: UpsertCreditBalance :exec
INSERT INTO credit_balances (
    id, school_id, student_id, amount_minor, updated_at, created_at
) VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT (school_id, student_id)
DO UPDATE SET amount_minor = excluded.amount_minor, updated_at = excluded.updated_at;

-- name: InsertCreditLedgerEntry :exec
INSERT INTO credit_balance_ledger (
    id, credit_balance_id, student_id, amount_minor, balance_after_minor,
    reference_type, reference_id, notes, created_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: ListCreditLedger :many
SELECT * FROM credit_balance_ledger
WHERE student_id = ?
ORDER BY created_at DESC;

-- name: InsertPayment :exec
INSERT INTO payments (
    id, school_id, student_id, receipt_no, payment_date, payment_method,
    amount_minor, allocated_amount_minor, unallocated_amount_minor, reference_no,
    payer_name, payer_phone, status, recorded_by_id, recorded_at, notes
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdatePaymentAllocation :exec
UPDATE payments
SET allocated_amount_minor = ?, unallocated_amount_minor = ?
WHERE id = ?;

-- name: UpdatePaymentStatus :exec
UPDATE payments
SET status = ?
WHERE id = ?;

-- name: GetPaymentByID :one
SELECT p.*, s.first_name AS student_first_name, s.last_name AS student_last_name, s.admission_no,
       u.display_name AS recorded_by_name
FROM payments p
JOIN students s ON p.student_id = s.id
JOIN users u ON p.recorded_by_id = u.id
WHERE p.id = ? LIMIT 1;

-- name: ListPayments :many
SELECT p.*, s.first_name AS student_first_name, s.last_name AS student_last_name, s.admission_no,
       u.display_name AS recorded_by_name
FROM payments p
JOIN students s ON p.student_id = s.id
JOIN users u ON p.recorded_by_id = u.id
WHERE p.school_id = ?
  AND (? = '' OR p.student_id = ?)
  AND (? = '' OR p.status = ?)
  AND (? = '' OR p.payment_method = ?)
ORDER BY p.recorded_at DESC
LIMIT ?;

-- name: InsertPaymentAllocation :exec
INSERT INTO payment_allocations (
    id, payment_id, invoice_id, invoice_item_id, allocated_amount_minor,
    allocated_at, allocated_by_id
) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: ListAllocationsByPayment :many
SELECT pa.*, inv.invoice_no, inv.billing_period
FROM payment_allocations pa
JOIN invoices inv ON pa.invoice_id = inv.id
WHERE pa.payment_id = ?
ORDER BY pa.allocated_at ASC;

-- name: ListAllocationsByInvoice :many
SELECT pa.*, p.receipt_no, p.payment_date, p.payment_method
FROM payment_allocations pa
JOIN payments p ON pa.payment_id = p.id
WHERE pa.invoice_id = ?
ORDER BY pa.allocated_at ASC;

-- name: InsertPaymentReversal :exec
INSERT INTO payment_reversals (
    id, payment_id, reversal_date, reason, reversed_by_id, reversed_at
) VALUES (?, ?, ?, ?, ?, ?);

-- name: GetPaymentReversal :one
SELECT pr.*, u.display_name AS reversed_by_name
FROM payment_reversals pr
JOIN users u ON pr.reversed_by_id = u.id
WHERE pr.payment_id = ? LIMIT 1;

-- name: InsertFinancialAdjustment :exec
INSERT INTO financial_adjustments (
    id, school_id, student_id, invoice_id, invoice_item_id, adjustment_type,
    amount_minor, reason, status, requested_by_id, requested_at, decision_notes
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetFinancialAdjustmentByID :one
SELECT fa.*, s.first_name AS student_first_name, s.last_name AS student_last_name, s.admission_no,
       req.display_name AS requested_by_name,
       dec.display_name AS decided_by_name,
       inv.invoice_no
FROM financial_adjustments fa
JOIN students s ON fa.student_id = s.id
JOIN users req ON fa.requested_by_id = req.id
LEFT JOIN users dec ON fa.decided_by_id = dec.id
LEFT JOIN invoices inv ON fa.invoice_id = inv.id
WHERE fa.id = ? LIMIT 1;

-- name: UpdateFinancialAdjustmentStatus :exec
UPDATE financial_adjustments
SET status = ?, decided_by_id = ?, decided_at = ?, decision_notes = ?
WHERE id = ?;

-- name: ListFinancialAdjustments :many
SELECT fa.*, s.first_name AS student_first_name, s.last_name AS student_last_name, s.admission_no,
       req.display_name AS requested_by_name,
       dec.display_name AS decided_by_name,
       inv.invoice_no
FROM financial_adjustments fa
JOIN students s ON fa.student_id = s.id
JOIN users req ON fa.requested_by_id = req.id
LEFT JOIN users dec ON fa.decided_by_id = dec.id
LEFT JOIN invoices inv ON fa.invoice_id = inv.id
WHERE fa.school_id = ?
  AND (? = '' OR fa.student_id = ?)
  AND (? = '' OR fa.status = ?)
ORDER BY fa.requested_at DESC
LIMIT ?;

-- name: GetUnpaidInvoicesForStudent :many
SELECT * FROM invoices
WHERE school_id = ? AND student_id = ? AND status IN ('issued', 'partially_paid')
ORDER BY issue_date ASC, created_at ASC;

