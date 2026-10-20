-- name: ListFeeHeads :many
SELECT * FROM fee_heads
WHERE school_id = ?
ORDER BY is_active DESC, name ASC;

-- name: GetFeeHead :one
SELECT * FROM fee_heads
WHERE id = ? LIMIT 1;

-- name: InsertFeeHead :exec
INSERT INTO fee_heads (
    id, school_id, code, name, description, is_recurring,
    default_amount_minor, is_active, created_at, updated_at, row_version
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1);

-- name: UpdateFeeHead :exec
UPDATE fee_heads
SET name = ?, description = ?, is_recurring = ?,
    default_amount_minor = ?, is_active = ?, updated_at = ?,
    row_version = row_version + 1
WHERE id = ?;

-- name: ListFeeStructures :many
SELECT * FROM fee_structures
WHERE school_id = ? AND academic_year_id = ?
ORDER BY name ASC;

-- name: GetFeeStructure :one
SELECT * FROM fee_structures
WHERE id = ? LIMIT 1;

-- name: GetFeeStructureForClass :one
SELECT * FROM fee_structures
WHERE academic_year_id = ? AND class_id = ? AND is_active = 1
LIMIT 1;

-- name: InsertFeeStructure :exec
INSERT INTO fee_structures (
    id, school_id, academic_year_id, class_id, campus_id, code,
    name, description, is_active, created_at, updated_at, row_version
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1);

-- name: ListFeeStructureVersions :many
SELECT * FROM fee_structure_versions
WHERE fee_structure_id = ?
ORDER BY version_no DESC;

-- name: GetFeeStructureVersion :one
SELECT * FROM fee_structure_versions
WHERE id = ? LIMIT 1;

-- name: GetPublishedFeeStructureVersion :one
SELECT * FROM fee_structure_versions
WHERE fee_structure_id = ? AND status = 'published'
LIMIT 1;

-- name: InsertFeeStructureVersion :exec
INSERT INTO fee_structure_versions (
    id, fee_structure_id, version_no, status, changelog,
    created_by, created_at, updated_at, row_version
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1);

-- name: PublishFeeStructureVersion :exec
UPDATE fee_structure_versions
SET status = 'published', published_by = ?, published_at = ?, updated_at = ?,
    row_version = row_version + 1
WHERE id = ?;

-- name: ListFeeStructureItems :many
SELECT fsi.*, fh.code as fee_head_code, fh.name as fee_head_name
FROM fee_structure_items fsi
JOIN fee_heads fh ON fsi.fee_head_id = fh.id
WHERE fsi.fee_structure_version_id = ?
ORDER BY fh.name ASC;

-- name: InsertFeeStructureItem :exec
INSERT INTO fee_structure_items (
    id, fee_structure_version_id, fee_head_id, amount_minor,
    frequency, due_day, is_optional, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: ListDiscountPolicies :many
SELECT * FROM discount_policies
WHERE school_id = ?
ORDER BY is_active DESC, name ASC;

-- name: GetDiscountPolicy :one
SELECT * FROM discount_policies
WHERE id = ? LIMIT 1;

-- name: InsertDiscountPolicy :exec
INSERT INTO discount_policies (
    id, school_id, code, name, description, discount_type,
    rate_basis_points, amount_minor, fee_head_id, requires_approval,
    is_active, created_at, updated_at, row_version
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1);

-- name: UpdateDiscountPolicy :exec
UPDATE discount_policies
SET name = ?, description = ?, discount_type = ?, rate_basis_points = ?,
    amount_minor = ?, fee_head_id = ?, requires_approval = ?,
    is_active = ?, updated_at = ?, row_version = row_version + 1
WHERE id = ?;

-- name: ListStudentDiscounts :many
SELECT sd.*, dp.code as policy_code, dp.name as policy_name,
       dp.discount_type, dp.rate_basis_points, dp.amount_minor as policy_amount_minor,
       dp.fee_head_id as policy_fee_head_id
FROM student_discounts sd
JOIN discount_policies dp ON sd.discount_policy_id = dp.id
WHERE sd.student_id = ? AND sd.academic_year_id = ?
ORDER BY sd.created_at DESC;

-- name: ListPendingStudentDiscounts :many
SELECT sd.*, dp.name as policy_name, s.first_name, s.last_name, s.admission_no
FROM student_discounts sd
JOIN discount_policies dp ON sd.discount_policy_id = dp.id
JOIN students s ON sd.student_id = s.id
WHERE s.school_id = ? AND sd.status = 'pending'
ORDER BY sd.created_at ASC;

-- name: GetStudentDiscount :one
SELECT * FROM student_discounts
WHERE id = ? LIMIT 1;

-- name: InsertStudentDiscount :exec
INSERT INTO student_discounts (
    id, student_id, discount_policy_id, academic_year_id, status,
    approved_by, approved_at, rejection_reason, valid_from, valid_to,
    created_by, created_at, updated_at, row_version
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1);

-- name: UpdateStudentDiscountStatus :exec
UPDATE student_discounts
SET status = ?, approved_by = ?, approved_at = ?,
    rejection_reason = ?, updated_at = ?, row_version = row_version + 1
WHERE id = ?;

-- name: ListInvoices :many
SELECT inv.*, s.first_name, s.last_name, s.admission_no
FROM invoices inv
JOIN students s ON inv.student_id = s.id
WHERE inv.school_id = ?
  AND (? = '' OR inv.status = ?)
  AND (? = '' OR inv.academic_year_id = ?)
  AND (? = '' OR inv.student_id = ?)
ORDER BY inv.created_at DESC
LIMIT ?;

-- name: GetInvoice :one
SELECT inv.*, s.first_name, s.last_name, s.admission_no
FROM invoices inv
JOIN students s ON inv.student_id = s.id
WHERE inv.id = ? LIMIT 1;

-- name: GetInvoiceByNo :one
SELECT * FROM invoices
WHERE school_id = ? AND invoice_no = ? LIMIT 1;

-- name: InsertInvoice :exec
INSERT INTO invoices (
    id, school_id, student_id, academic_year_id, invoice_no, issue_date,
    due_date, billing_period, subtotal_amount_minor, discount_amount_minor,
    total_amount_minor, paid_amount_minor, balance_minor, status,
    created_by, created_at, updated_at, row_version
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1);

-- name: UpdateInvoiceStatus :exec
UPDATE invoices
SET status = ?, paid_amount_minor = ?, balance_minor = ?, updated_at = ?,
    row_version = row_version + 1
WHERE id = ?;

-- name: VoidInvoice :exec
UPDATE invoices
SET status = 'void', void_reason = ?, voided_by = ?, voided_at = ?,
    updated_at = ?, row_version = row_version + 1
WHERE id = ?;

-- name: ListInvoiceItems :many
SELECT * FROM invoice_items
WHERE invoice_id = ?
ORDER BY created_at ASC;

-- name: InsertInvoiceItem :exec
INSERT INTO invoice_items (
    id, invoice_id, fee_head_id, fee_head_name, description,
    original_amount_minor, discount_amount_minor, net_amount_minor, created_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: CountInvoicesInPeriod :one
SELECT COUNT(*) FROM invoices
WHERE school_id = ? AND billing_period = ?;
