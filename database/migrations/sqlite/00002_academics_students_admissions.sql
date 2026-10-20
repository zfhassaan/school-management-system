-- +goose Up
-- +goose StatementBegin

-- 1. Sections (per academic year, class, campus)
CREATE TABLE sections (
	id TEXT PRIMARY KEY,
	class_id TEXT NOT NULL REFERENCES classes(id),
	academic_year_id TEXT NOT NULL REFERENCES academic_years(id),
	campus_id TEXT REFERENCES campuses(id),
	name TEXT NOT NULL,
	capacity INTEGER CHECK (capacity IS NULL OR capacity > 0),
	class_teacher_user_id TEXT REFERENCES users(id),
	is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0,1)),
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	row_version INTEGER NOT NULL DEFAULT 1
) STRICT;
CREATE UNIQUE INDEX ux_sections_identity ON sections(class_id, academic_year_id, coalesce(campus_id,''), name);
CREATE INDEX ix_sections_year ON sections(academic_year_id, class_id);

-- 2. Students & Search
CREATE TABLE students (
	id TEXT PRIMARY KEY,
	school_id TEXT NOT NULL REFERENCES schools(id),
	admission_no TEXT NOT NULL,
	first_name TEXT NOT NULL,
	middle_name TEXT,
	last_name TEXT NOT NULL,
	search_name TEXT NOT NULL,
	date_of_birth TEXT NOT NULL CHECK (date_of_birth = date(date_of_birth)),
	gender TEXT NOT NULL CHECK (gender IN ('male','female','other','unspecified')),
	national_id TEXT,
	address_line1 TEXT,
	address_line2 TEXT,
	city TEXT,
	region TEXT,
	postal_code TEXT,
	country TEXT,
	photo_document_id TEXT,
	admission_date TEXT NOT NULL CHECK (admission_date = date(admission_date)),
	status TEXT NOT NULL CHECK (status IN ('active','inactive','graduated','withdrawn','alumni')),
	notes TEXT,
	deleted_at TEXT,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	row_version INTEGER NOT NULL DEFAULT 1,
	UNIQUE (school_id, admission_no)
) STRICT;
CREATE INDEX ix_students_search_name ON students(search_name, id) WHERE deleted_at IS NULL;
CREATE INDEX ix_students_status_name ON students(status, search_name, id) WHERE deleted_at IS NULL;

-- Trigram FTS table for typo-tolerant / substring search
CREATE VIRTUAL TABLE students_fts USING fts5(
	search_name,
	admission_no,
	content='students',
	content_rowid='rowid',
	tokenize='trigram'
);

-- Triggers to synchronize students_fts
CREATE TRIGGER trg_students_ai AFTER INSERT ON students
BEGIN
	INSERT INTO students_fts(rowid, search_name, admission_no)
	VALUES (new.rowid, new.search_name, new.admission_no);
END;

CREATE TRIGGER trg_students_ad AFTER DELETE ON students
BEGIN
	INSERT INTO students_fts(students_fts, rowid, search_name, admission_no)
	VALUES ('delete', old.rowid, old.search_name, old.admission_no);
END;

CREATE TRIGGER trg_students_au AFTER UPDATE ON students
BEGIN
	INSERT INTO students_fts(students_fts, rowid, search_name, admission_no)
	VALUES ('delete', old.rowid, old.search_name, old.admission_no);
	INSERT INTO students_fts(rowid, search_name, admission_no)
	VALUES (new.rowid, new.search_name, new.admission_no);
END;

CREATE TABLE student_status_history (
	id TEXT PRIMARY KEY,
	student_id TEXT NOT NULL REFERENCES students(id),
	from_status TEXT,
	to_status TEXT NOT NULL,
	effective_on TEXT NOT NULL CHECK (effective_on = date(effective_on)),
	reason TEXT NOT NULL,
	changed_by TEXT NOT NULL REFERENCES users(id),
	changed_at TEXT NOT NULL
) STRICT;
CREATE INDEX ix_student_status_history ON student_status_history(student_id, changed_at);

-- 3. Guardians & Relationships
CREATE TABLE guardians (
	id TEXT PRIMARY KEY,
	school_id TEXT NOT NULL REFERENCES schools(id),
	first_name TEXT NOT NULL,
	last_name TEXT,
	phone_primary TEXT NOT NULL,
	phone_primary_normalized TEXT NOT NULL,
	phone_secondary TEXT,
	email TEXT,
	occupation TEXT,
	national_id TEXT,
	address_line1 TEXT,
	address_line2 TEXT,
	city TEXT,
	region TEXT,
	postal_code TEXT,
	country TEXT,
	deleted_at TEXT,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	row_version INTEGER NOT NULL DEFAULT 1
) STRICT;
CREATE INDEX ix_guardians_phone ON guardians(phone_primary_normalized) WHERE deleted_at IS NULL;
CREATE INDEX ix_guardians_national_id ON guardians(national_id) WHERE national_id IS NOT NULL;

CREATE TABLE student_guardians (
	student_id TEXT NOT NULL REFERENCES students(id),
	guardian_id TEXT NOT NULL REFERENCES guardians(id),
	relationship TEXT NOT NULL CHECK (relationship IN ('father','mother','guardian','grandparent','sibling','relative','emergency_contact','other')),
	is_primary INTEGER NOT NULL DEFAULT 0 CHECK (is_primary IN (0,1)),
	is_emergency_contact INTEGER NOT NULL DEFAULT 0 CHECK (is_emergency_contact IN (0,1)),
	receives_billing INTEGER NOT NULL DEFAULT 0 CHECK (receives_billing IN (0,1)),
	can_pickup INTEGER NOT NULL DEFAULT 0 CHECK (can_pickup IN (0,1)),
	priority INTEGER NOT NULL DEFAULT 1,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	row_version INTEGER NOT NULL DEFAULT 1,
	PRIMARY KEY (student_id, guardian_id)
) STRICT;
CREATE UNIQUE INDEX ux_student_primary_guardian ON student_guardians(student_id) WHERE is_primary = 1;
CREATE INDEX ix_student_guardians_guardian ON student_guardians(guardian_id);

-- 4. Admissions, Documents & Form Builder
CREATE TABLE admission_form_definitions (
	id TEXT PRIMARY KEY,
	school_id TEXT NOT NULL REFERENCES schools(id),
	code TEXT NOT NULL,
	name TEXT NOT NULL,
	description TEXT,
	is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0,1)),
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	row_version INTEGER NOT NULL DEFAULT 1,
	UNIQUE (school_id, code)
) STRICT;

CREATE TABLE admission_form_versions (
	id TEXT PRIMARY KEY,
	form_definition_id TEXT NOT NULL REFERENCES admission_form_definitions(id),
	version_no INTEGER NOT NULL CHECK (version_no >= 1),
	status TEXT NOT NULL CHECK (status IN ('draft','published','retired')),
	schema_json TEXT NOT NULL CHECK (json_valid(schema_json)),
	schema_format_version INTEGER NOT NULL DEFAULT 1,
	schema_hash TEXT NOT NULL,
	published_at TEXT,
	published_by TEXT REFERENCES users(id),
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	row_version INTEGER NOT NULL DEFAULT 1,
	UNIQUE (form_definition_id, version_no),
	CHECK ((status = 'draft') = (published_at IS NULL))
) STRICT;
CREATE UNIQUE INDEX ux_form_one_draft ON admission_form_versions(form_definition_id) WHERE status = 'draft';
CREATE UNIQUE INDEX ux_form_one_published ON admission_form_versions(form_definition_id) WHERE status = 'published';

CREATE TABLE admission_applications (
	id TEXT PRIMARY KEY,
	school_id TEXT NOT NULL REFERENCES schools(id),
	application_no TEXT NOT NULL,
	form_version_id TEXT NOT NULL REFERENCES admission_form_versions(id),
	academic_year_id TEXT NOT NULL REFERENCES academic_years(id),
	applied_class_id TEXT NOT NULL REFERENCES classes(id),
	status TEXT NOT NULL CHECK (status IN ('draft','submitted','under_review','approved','rejected','enrolled','withdrawn')),
	applicant_first_name TEXT NOT NULL,
	applicant_last_name TEXT NOT NULL,
	applicant_dob TEXT CHECK (applicant_dob IS NULL OR applicant_dob = date(applicant_dob)),
	applicant_gender TEXT,
	guardian_phone_normalized TEXT,
	student_id TEXT REFERENCES students(id),
	submitted_at TEXT,
	decided_at TEXT,
	decided_by TEXT REFERENCES users(id),
	decision_reason TEXT,
	created_by TEXT NOT NULL REFERENCES users(id),
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	row_version INTEGER NOT NULL DEFAULT 1,
	UNIQUE (school_id, application_no),
	CHECK ((status = 'enrolled') = (student_id IS NOT NULL))
) STRICT;
CREATE INDEX ix_applications_status ON admission_applications(academic_year_id, status, created_at);
CREATE INDEX ix_applications_phone ON admission_applications(guardian_phone_normalized);

CREATE TABLE admission_form_responses (
	id TEXT PRIMARY KEY,
	application_id TEXT NOT NULL UNIQUE REFERENCES admission_applications(id),
	form_version_id TEXT NOT NULL REFERENCES admission_form_versions(id),
	response_json TEXT NOT NULL CHECK (json_valid(response_json)),
	response_hash TEXT NOT NULL,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	row_version INTEGER NOT NULL DEFAULT 1
) STRICT;

CREATE TABLE admission_status_history (
	id TEXT PRIMARY KEY,
	application_id TEXT NOT NULL REFERENCES admission_applications(id),
	from_status TEXT,
	to_status TEXT NOT NULL,
	reason TEXT,
	changed_by TEXT NOT NULL REFERENCES users(id),
	changed_at TEXT NOT NULL
) STRICT;
CREATE INDEX ix_admission_status_history ON admission_status_history(application_id, changed_at);

CREATE TABLE student_documents (
	id TEXT PRIMARY KEY,
	student_id TEXT REFERENCES students(id),
	application_id TEXT REFERENCES admission_applications(id),
	document_type_code TEXT NOT NULL,
	file_name TEXT NOT NULL,
	mime_type TEXT NOT NULL,
	size_bytes INTEGER NOT NULL CHECK (size_bytes > 0),
	sha256 TEXT NOT NULL CHECK (length(sha256) = 64),
	storage_key TEXT NOT NULL,
	version_no INTEGER NOT NULL DEFAULT 1,
	uploaded_by TEXT NOT NULL REFERENCES users(id),
	deleted_at TEXT,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	row_version INTEGER NOT NULL DEFAULT 1,
	CHECK (student_id IS NOT NULL OR application_id IS NOT NULL)
) STRICT;
CREATE INDEX ix_documents_student ON student_documents(student_id, document_type_code) WHERE deleted_at IS NULL;
CREATE INDEX ix_documents_application ON student_documents(application_id) WHERE deleted_at IS NULL;
CREATE INDEX ix_documents_sha ON student_documents(sha256);

-- 5. Promotions & Enrollments
CREATE TABLE promotion_batches (
	id TEXT PRIMARY KEY,
	from_academic_year_id TEXT NOT NULL REFERENCES academic_years(id),
	to_academic_year_id TEXT NOT NULL REFERENCES academic_years(id),
	status TEXT NOT NULL CHECK (status IN ('previewed','executed','rolled_back')),
	summary_json TEXT NOT NULL CHECK (json_valid(summary_json)),
	executed_by TEXT REFERENCES users(id),
	executed_at TEXT,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	row_version INTEGER NOT NULL DEFAULT 1,
	CHECK (from_academic_year_id <> to_academic_year_id)
) STRICT;

CREATE TABLE enrollments (
	id TEXT PRIMARY KEY,
	student_id TEXT NOT NULL REFERENCES students(id),
	academic_year_id TEXT NOT NULL REFERENCES academic_years(id),
	class_id TEXT NOT NULL REFERENCES classes(id),
	section_id TEXT REFERENCES sections(id),
	roll_no TEXT,
	status TEXT NOT NULL CHECK (status IN ('active','promoted','repeated','transferred','withdrawn','completed','cancelled')),
	enrolled_on TEXT NOT NULL CHECK (enrolled_on = date(enrolled_on)),
	ended_on TEXT CHECK (ended_on IS NULL OR ended_on = date(ended_on)),
	previous_enrollment_id TEXT REFERENCES enrollments(id),
	promotion_batch_id TEXT REFERENCES promotion_batches(id),
	admission_application_id TEXT REFERENCES admission_applications(id),
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	row_version INTEGER NOT NULL DEFAULT 1,
	CHECK (ended_on IS NULL OR ended_on >= enrolled_on)
) STRICT;
CREATE UNIQUE INDEX ux_enrollments_student_year ON enrollments(student_id, academic_year_id) WHERE status <> 'cancelled';
CREATE UNIQUE INDEX ux_enrollments_one_active ON enrollments(student_id) WHERE status = 'active';
CREATE INDEX ix_enrollments_roster ON enrollments(academic_year_id, class_id, section_id, status, student_id);
CREATE UNIQUE INDEX ux_enrollments_roll ON enrollments(section_id, roll_no) WHERE roll_no IS NOT NULL AND status = 'active';

-- Immutability triggers for published admission form versions
CREATE TRIGGER trg_form_version_immutable BEFORE UPDATE OF schema_json ON admission_form_versions
WHEN OLD.status <> 'draft'
BEGIN
	SELECT RAISE(ABORT, 'IMMUTABLE: published or retired admission form versions cannot be modified');
END;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS enrollments;
DROP TABLE IF EXISTS promotion_batches;
DROP TABLE IF EXISTS student_documents;
DROP TABLE IF EXISTS admission_status_history;
DROP TABLE IF EXISTS admission_form_responses;
DROP TABLE IF EXISTS admission_applications;
DROP TABLE IF EXISTS admission_form_versions;
DROP TABLE IF EXISTS admission_form_definitions;
DROP TABLE IF EXISTS student_guardians;
DROP TABLE IF EXISTS guardians;
DROP TABLE IF EXISTS student_status_history;
DROP TRIGGER IF EXISTS trg_students_au;
DROP TRIGGER IF EXISTS trg_students_ad;
DROP TRIGGER IF EXISTS trg_students_ai;
DROP TABLE IF EXISTS students_fts;
DROP TABLE IF EXISTS students;
DROP TABLE IF EXISTS sections;
-- +goose StatementEnd

