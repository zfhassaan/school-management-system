-- name: GetStudent :one
SELECT * FROM students WHERE id = ? AND deleted_at IS NULL;

-- name: GetStudentByAdmissionNo :one
SELECT * FROM students WHERE school_id = ? AND admission_no = ? AND deleted_at IS NULL;

-- name: InsertStudent :exec
INSERT INTO students (
    id, school_id, admission_no, first_name, middle_name, last_name, search_name,
    date_of_birth, gender, national_id, address_line1, address_line2, city, region,
    postal_code, country, photo_document_id, admission_date, status, notes,
    created_at, updated_at, row_version
) VALUES (
    ?, ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?, ?,
    ?, ?, 1
);

-- name: UpdateStudent :exec
UPDATE students SET
    first_name = ?, middle_name = ?, last_name = ?, search_name = ?,
    date_of_birth = ?, gender = ?, national_id = ?, address_line1 = ?, address_line2 = ?,
    city = ?, region = ?, postal_code = ?, country = ?, photo_document_id = ?,
    status = ?, notes = ?, updated_at = ?, row_version = row_version + 1
WHERE id = ? AND row_version = ?;

-- name: AppendStudentStatusHistory :exec
INSERT INTO student_status_history (id, student_id, from_status, to_status, effective_on, reason, changed_by, changed_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: ListStudentStatusHistory :many
SELECT * FROM student_status_history WHERE student_id = ? ORDER BY changed_at DESC;

-- name: ListStudentsCursor :many
SELECT s.*, e.class_id, c.name as class_name, e.section_id, sec.name as section_name, e.roll_no
FROM students s
LEFT JOIN enrollments e ON e.student_id = s.id AND e.status = 'active'
LEFT JOIN classes c ON c.id = e.class_id
LEFT JOIN sections sec ON sec.id = e.section_id
WHERE s.school_id = ?
  AND s.deleted_at IS NULL
  AND (? = '' OR s.status = ?)
  AND (? = '' OR e.class_id = ?)
  AND (? = '' OR e.section_id = ?)
  AND (s.search_name > ? OR (s.search_name = ? AND s.id > ?))
ORDER BY s.search_name ASC, s.id ASC
LIMIT ?;

-- name: SearchStudentsFTS :many
SELECT s.*, e.class_id, c.name as class_name, e.section_id, sec.name as section_name, e.roll_no
FROM students_fts
JOIN students s ON s.rowid = students_fts.rowid
LEFT JOIN enrollments e ON e.student_id = s.id AND e.status = 'active'
LEFT JOIN classes c ON c.id = e.class_id
LEFT JOIN sections sec ON sec.id = e.section_id
WHERE students_fts.search_name MATCH ?
  AND s.school_id = ?
  AND s.deleted_at IS NULL
LIMIT ?;

-- name: CountStudents :one
SELECT count(*) FROM students WHERE school_id = ? AND deleted_at IS NULL;

-- name: GetGuardian :one
SELECT * FROM guardians WHERE id = ? AND deleted_at IS NULL;

-- name: GetGuardianByPhone :one
SELECT * FROM guardians WHERE school_id = ? AND phone_primary_normalized = ? AND deleted_at IS NULL LIMIT 1;

-- name: InsertGuardian :exec
INSERT INTO guardians (
    id, school_id, first_name, last_name, phone_primary, phone_primary_normalized,
    phone_secondary, email, occupation, national_id, address_line1, address_line2,
    city, region, postal_code, country, created_at, updated_at, row_version
) VALUES (
    ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?, ?, 1
);

-- name: UpdateGuardian :exec
UPDATE guardians SET
    first_name = ?, last_name = ?, phone_primary = ?, phone_primary_normalized = ?,
    phone_secondary = ?, email = ?, occupation = ?, national_id = ?,
    address_line1 = ?, address_line2 = ?, city = ?, region = ?, postal_code = ?,
    country = ?, updated_at = ?, row_version = row_version + 1
WHERE id = ? AND row_version = ?;

-- name: InsertStudentGuardian :exec
INSERT INTO student_guardians (
    student_id, guardian_id, relationship, is_primary, is_emergency_contact,
    receives_billing, can_pickup, priority, created_at, updated_at, row_version
) VALUES (
    ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?, 1
);

-- name: DeleteStudentGuardians :exec
DELETE FROM student_guardians WHERE student_id = ?;

-- name: ListStudentGuardians :many
SELECT sg.*, g.first_name, g.last_name, g.phone_primary, g.phone_secondary, g.email, g.occupation
FROM student_guardians sg
JOIN guardians g ON g.id = sg.guardian_id
WHERE sg.student_id = ?
ORDER BY sg.is_primary DESC, sg.priority ASC;

-- name: FindSiblingsByGuardian :many
SELECT s.*, c.name as class_name, sec.name as section_name
FROM student_guardians sg
JOIN students s ON s.id = sg.student_id
LEFT JOIN enrollments e ON e.student_id = s.id AND e.status = 'active'
LEFT JOIN classes c ON c.id = e.class_id
LEFT JOIN sections sec ON sec.id = e.section_id
WHERE sg.guardian_id IN (
    SELECT inner_sg.guardian_id FROM student_guardians inner_sg WHERE inner_sg.student_id = ?
)
AND s.id <> ?
AND s.deleted_at IS NULL;

-- name: GetActiveEnrollment :one
SELECT e.*, c.name as class_name, sec.name as section_name, ay.code as year_code
FROM enrollments e
JOIN classes c ON c.id = e.class_id
LEFT JOIN sections sec ON sec.id = e.section_id
JOIN academic_years ay ON ay.id = e.academic_year_id
WHERE e.student_id = ? AND e.status = 'active'
LIMIT 1;

-- name: ListEnrollmentHistory :many
SELECT e.*, c.name as class_name, sec.name as section_name, ay.code as year_code
FROM enrollments e
JOIN classes c ON c.id = e.class_id
LEFT JOIN sections sec ON sec.id = e.section_id
JOIN academic_years ay ON ay.id = e.academic_year_id
WHERE e.student_id = ?
ORDER BY e.enrolled_on DESC;

-- name: InsertEnrollment :exec
INSERT INTO enrollments (
    id, student_id, academic_year_id, class_id, section_id, roll_no,
    status, enrolled_on, ended_on, previous_enrollment_id, promotion_batch_id,
    admission_application_id, created_at, updated_at, row_version
) VALUES (
    ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?,
    ?, ?, ?, 1
);

-- name: CloseEnrollment :exec
UPDATE enrollments SET
    status = ?, ended_on = ?, updated_at = ?, row_version = row_version + 1
WHERE id = ? AND row_version = ?;

-- name: ListSectionRoster :many
SELECT s.id, s.admission_no, s.first_name, s.last_name, e.roll_no, e.id as enrollment_id, e.status as enrollment_status
FROM enrollments e
JOIN students s ON s.id = e.student_id
WHERE e.section_id = ? AND e.status = 'active'
ORDER BY e.roll_no ASC, s.last_name ASC, s.first_name ASC;

-- name: InsertPromotionBatch :exec
INSERT INTO promotion_batches (
    id, from_academic_year_id, to_academic_year_id, status, summary_json,
    executed_by, executed_at, created_at, updated_at, row_version
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1);

-- name: GetPromotionBatch :one
SELECT * FROM promotion_batches WHERE id = ?;

-- name: ListAdmissionFormDefinitions :many
SELECT * FROM admission_form_definitions WHERE school_id = ? ORDER BY name ASC;

-- name: GetAdmissionFormDefinition :one
SELECT * FROM admission_form_definitions WHERE id = ?;

-- name: InsertAdmissionFormDefinition :exec
INSERT INTO admission_form_definitions (id, school_id, code, name, description, is_active, created_at, updated_at, row_version)
VALUES (?, ?, ?, ?, ?, 1, ?, ?, 1);

-- name: ListAdmissionFormVersions :many
SELECT * FROM admission_form_versions WHERE form_definition_id = ? ORDER BY version_no DESC;

-- name: GetAdmissionFormVersion :one
SELECT * FROM admission_form_versions WHERE id = ?;

-- name: GetPublishedFormVersion :one
SELECT * FROM admission_form_versions WHERE form_definition_id = ? AND status = 'published' LIMIT 1;

-- name: InsertAdmissionFormVersion :exec
INSERT INTO admission_form_versions (
    id, form_definition_id, version_no, status, schema_json, schema_format_version,
    schema_hash, published_at, published_by, created_at, updated_at, row_version
) VALUES (
    ?, ?, ?, 'draft', ?, 1,
    ?, NULL, NULL, ?, ?, 1
);

-- name: PublishAdmissionFormVersion :exec
UPDATE admission_form_versions SET
    status = 'published', published_at = ?, published_by = ?, schema_hash = ?, updated_at = ?, row_version = row_version + 1
WHERE id = ? AND status = 'draft';

-- name: RetireAdmissionFormVersion :exec
UPDATE admission_form_versions SET
    status = 'retired', updated_at = ?, row_version = row_version + 1
WHERE id = ? AND status = 'published';

-- name: ListAdmissionApplications :many
SELECT a.*, c.name as class_name, ay.code as year_code
FROM admission_applications a
JOIN classes c ON c.id = a.applied_class_id
JOIN academic_years ay ON ay.id = a.academic_year_id
WHERE a.school_id = ?
  AND (? = '' OR a.status = ?)
  AND (? = '' OR a.academic_year_id = ?)
ORDER BY a.created_at DESC
LIMIT ?;

-- name: GetAdmissionApplication :one
SELECT a.*, c.name as class_name, ay.code as year_code, r.response_json
FROM admission_applications a
JOIN classes c ON c.id = a.applied_class_id
JOIN academic_years ay ON ay.id = a.academic_year_id
LEFT JOIN admission_form_responses r ON r.application_id = a.id
WHERE a.id = ?;

-- name: InsertAdmissionApplication :exec
INSERT INTO admission_applications (
    id, school_id, application_no, form_version_id, academic_year_id, applied_class_id,
    status, applicant_first_name, applicant_last_name, applicant_dob, applicant_gender,
    guardian_phone_normalized, student_id, submitted_at, created_by, created_at, updated_at, row_version
) VALUES (
    ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?, ?, 1
);

-- name: UpdateAdmissionApplicationStatus :exec
UPDATE admission_applications SET
    status = ?, decided_at = ?, decided_by = ?, decision_reason = ?,
    student_id = ?, updated_at = ?, row_version = row_version + 1
WHERE id = ? AND row_version = ?;

-- name: InsertAdmissionFormResponse :exec
INSERT INTO admission_form_responses (id, application_id, form_version_id, response_json, response_hash, created_at, updated_at, row_version)
VALUES (?, ?, ?, ?, ?, ?, ?, 1);

-- name: AppendAdmissionStatusHistory :exec
INSERT INTO admission_status_history (id, application_id, from_status, to_status, reason, changed_by, changed_at)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: ListAdmissionStatusHistory :many
SELECT * FROM admission_status_history WHERE application_id = ? ORDER BY changed_at ASC;

-- name: InsertStudentDocument :exec
INSERT INTO student_documents (
    id, student_id, application_id, document_type_code, file_name, mime_type,
    size_bytes, sha256, storage_key, version_no, uploaded_by, created_at, updated_at, row_version
) VALUES (
    ?, ?, ?, ?, ?, ?,
    ?, ?, ?, 1, ?, ?, ?, 1
);

-- name: ListStudentDocuments :many
SELECT * FROM student_documents WHERE student_id = ? AND deleted_at IS NULL ORDER BY created_at DESC;

-- name: ListApplicationDocuments :many
SELECT * FROM student_documents WHERE application_id = ? AND deleted_at IS NULL ORDER BY created_at DESC;

-- name: GetDocument :one
SELECT * FROM student_documents WHERE id = ? AND deleted_at IS NULL;
