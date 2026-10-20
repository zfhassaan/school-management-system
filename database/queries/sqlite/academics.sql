-- name: ListAcademicYears :many
SELECT * FROM academic_years WHERE school_id = ? ORDER BY starts_on DESC;

-- name: GetAcademicYear :one
SELECT * FROM academic_years WHERE id = ?;

-- name: GetActiveAcademicYear :one
SELECT * FROM academic_years WHERE school_id = ? AND status = 'active' LIMIT 1;

-- name: InsertAcademicYear :exec
INSERT INTO academic_years (id, school_id, code, name, starts_on, ends_on, status, created_at, updated_at, row_version)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1);

-- name: UpdateAcademicYear :exec
UPDATE academic_years SET
    name = ?, starts_on = ?, ends_on = ?, status = ?, updated_at = ?, row_version = row_version + 1
WHERE id = ? AND row_version = ?;

-- name: ListAcademicTerms :many
SELECT * FROM academic_terms WHERE academic_year_id = ? ORDER BY sequence ASC;

-- name: InsertAcademicTerm :exec
INSERT INTO academic_terms (id, academic_year_id, sequence, name, starts_on, ends_on, created_at, updated_at, row_version)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1);

-- name: DeleteAcademicTermsByYear :exec
DELETE FROM academic_terms WHERE academic_year_id = ?;

-- name: ListClasses :many
SELECT * FROM classes WHERE school_id = ? ORDER BY level_order ASC;

-- name: GetClass :one
SELECT * FROM classes WHERE id = ?;

-- name: InsertClass :exec
INSERT INTO classes (id, school_id, code, name, level_order, next_class_id, is_active, created_at, updated_at, row_version)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1);

-- name: UpdateClass :exec
UPDATE classes SET
    code = ?, name = ?, level_order = ?, next_class_id = ?, is_active = ?, updated_at = ?, row_version = row_version + 1
WHERE id = ? AND row_version = ?;

-- name: ListSectionsByYear :many
SELECT s.*, c.name as class_name, c.level_order as class_level_order
FROM sections s
JOIN classes c ON c.id = s.class_id
WHERE s.academic_year_id = ?
ORDER BY c.level_order ASC, s.name ASC;

-- name: GetSection :one
SELECT s.*, c.name as class_name
FROM sections s
JOIN classes c ON c.id = s.class_id
WHERE s.id = ?;

-- name: InsertSection :exec
INSERT INTO sections (id, class_id, academic_year_id, campus_id, name, capacity, class_teacher_user_id, is_active, created_at, updated_at, row_version)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1);

-- name: UpdateSection :exec
UPDATE sections SET
    name = ?, capacity = ?, class_teacher_user_id = ?, is_active = ?, updated_at = ?, row_version = row_version + 1
WHERE id = ? AND row_version = ?;

-- name: ListCampuses :many
SELECT * FROM campuses WHERE school_id = ? ORDER BY name ASC;

-- name: InsertCampus :exec
INSERT INTO campuses (id, school_id, code, name, address, is_active, created_at, updated_at, row_version)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1);
