-- +goose Up
-- +goose StatementBegin

-- 1. Schools & Configuration
CREATE TABLE schools (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL CHECK (length(name) BETWEEN 1 AND 200),
	legal_name TEXT,
	code TEXT NOT NULL UNIQUE,
	logo_document_id TEXT,
	address_line1 TEXT,
	address_line2 TEXT,
	city TEXT,
	region TEXT,
	postal_code TEXT,
	country TEXT,
	phone TEXT,
	email TEXT,
	website TEXT,
	timezone TEXT NOT NULL,
	currency_code TEXT NOT NULL CHECK (length(currency_code) = 3),
	currency_minor_units INTEGER NOT NULL CHECK (currency_minor_units BETWEEN 0 AND 3),
	locale TEXT NOT NULL DEFAULT 'en',
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	row_version INTEGER NOT NULL DEFAULT 1
) STRICT;

CREATE TABLE school_settings (
	school_id TEXT PRIMARY KEY REFERENCES schools(id),
	admission_no_format TEXT NOT NULL DEFAULT '{YYYY}-{SEQ:5}',
	invoice_prefix TEXT NOT NULL DEFAULT 'INV',
	receipt_prefix TEXT NOT NULL DEFAULT 'RCP',
	rounding_unit_minor INTEGER NOT NULL DEFAULT 1 CHECK (rounding_unit_minor >= 1),
	rounding_mode TEXT NOT NULL DEFAULT 'half_up' CHECK (rounding_mode IN ('half_up','half_even','down','up')),
	default_due_day INTEGER NOT NULL DEFAULT 10 CHECK (default_due_day BETWEEN 1 AND 28),
	late_fee_grace_days INTEGER NOT NULL DEFAULT 0 CHECK (late_fee_grace_days >= 0),
	discount_combination TEXT NOT NULL DEFAULT 'exclusive_overrides' CHECK (discount_combination IN ('exclusive_overrides','best_for_student')),
	require_second_approver INTEGER NOT NULL DEFAULT 0 CHECK (require_second_approver IN (0,1)),
	date_format TEXT NOT NULL DEFAULT 'dd/MM/yyyy',
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	row_version INTEGER NOT NULL DEFAULT 1
) STRICT;

CREATE TABLE campuses (
	id TEXT PRIMARY KEY,
	school_id TEXT NOT NULL REFERENCES schools(id),
	code TEXT NOT NULL,
	name TEXT NOT NULL,
	address TEXT,
	is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0,1)),
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	row_version INTEGER NOT NULL DEFAULT 1,
	UNIQUE (school_id, code)
) STRICT;

CREATE TABLE academic_years (
	id TEXT PRIMARY KEY,
	school_id TEXT NOT NULL REFERENCES schools(id),
	code TEXT NOT NULL,
	name TEXT NOT NULL,
	starts_on TEXT NOT NULL CHECK (starts_on = date(starts_on)),
	ends_on TEXT NOT NULL CHECK (ends_on = date(ends_on)),
	status TEXT NOT NULL CHECK (status IN ('planned','active','closed')),
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	row_version INTEGER NOT NULL DEFAULT 1,
	UNIQUE (school_id, code),
	CHECK (ends_on > starts_on)
) STRICT;
CREATE UNIQUE INDEX ux_academic_years_one_active ON academic_years(school_id) WHERE status = 'active';

CREATE TABLE academic_terms (
	id TEXT PRIMARY KEY,
	academic_year_id TEXT NOT NULL REFERENCES academic_years(id),
	sequence INTEGER NOT NULL CHECK (sequence >= 1),
	name TEXT NOT NULL,
	starts_on TEXT NOT NULL CHECK (starts_on = date(starts_on)),
	ends_on TEXT NOT NULL CHECK (ends_on = date(ends_on)),
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	row_version INTEGER NOT NULL DEFAULT 1,
	UNIQUE (academic_year_id, sequence),
	CHECK (ends_on > starts_on)
) STRICT;

CREATE TABLE classes (
	id TEXT PRIMARY KEY,
	school_id TEXT NOT NULL REFERENCES schools(id),
	code TEXT NOT NULL,
	name TEXT NOT NULL,
	level_order INTEGER NOT NULL,
	next_class_id TEXT REFERENCES classes(id),
	is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0,1)),
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	row_version INTEGER NOT NULL DEFAULT 1,
	UNIQUE (school_id, code),
	UNIQUE (school_id, level_order)
) STRICT;

-- 2. Administration & Security
CREATE TABLE users (
	id TEXT PRIMARY KEY,
	username TEXT NOT NULL UNIQUE COLLATE NOCASE,
	display_name TEXT NOT NULL,
	email TEXT,
	password_hash TEXT NOT NULL,
	password_changed_at TEXT NOT NULL,
	must_change_password INTEGER NOT NULL DEFAULT 0 CHECK (must_change_password IN (0,1)),
	status TEXT NOT NULL CHECK (status IN ('active','disabled')),
	failed_login_count INTEGER NOT NULL DEFAULT 0,
	lockout_count INTEGER NOT NULL DEFAULT 0,
	locked_until TEXT,
	last_login_at TEXT,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	row_version INTEGER NOT NULL DEFAULT 1
) STRICT;

CREATE TABLE roles (
	id TEXT PRIMARY KEY,
	code TEXT NOT NULL UNIQUE,
	name TEXT NOT NULL,
	description TEXT,
	is_system INTEGER NOT NULL DEFAULT 0 CHECK (is_system IN (0,1)),
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	row_version INTEGER NOT NULL DEFAULT 1
) STRICT;

CREATE TABLE permissions (
	code TEXT PRIMARY KEY,
	module TEXT NOT NULL,
	description TEXT NOT NULL,
	is_sensitive INTEGER NOT NULL DEFAULT 0 CHECK (is_sensitive IN (0,1))
) STRICT;

CREATE TABLE role_permissions (
	role_id TEXT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
	permission_code TEXT NOT NULL REFERENCES permissions(code),
	PRIMARY KEY (role_id, permission_code)
) STRICT;

CREATE TABLE user_roles (
	user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	role_id TEXT NOT NULL REFERENCES roles(id),
	assigned_by TEXT REFERENCES users(id),
	assigned_at TEXT NOT NULL,
	PRIMARY KEY (user_id, role_id)
) STRICT;

CREATE TABLE user_sessions (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES users(id),
	token_hash TEXT NOT NULL UNIQUE,
	device_id TEXT,
	created_at TEXT NOT NULL,
	last_seen_at TEXT NOT NULL,
	idle_expires_at TEXT NOT NULL,
	absolute_expires_at TEXT NOT NULL,
	revoked_at TEXT
) STRICT;
CREATE INDEX ix_sessions_user ON user_sessions(user_id) WHERE revoked_at IS NULL;

CREATE TABLE recovery_codes (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES users(id),
	code_hash TEXT NOT NULL,
	used_at TEXT,
	created_at TEXT NOT NULL
) STRICT;

CREATE TABLE audit_logs (
	id TEXT PRIMARY KEY,
	occurred_at TEXT NOT NULL,
	actor_user_id TEXT REFERENCES users(id),
	device_id TEXT,
	action TEXT NOT NULL,
	entity_type TEXT NOT NULL,
	entity_id TEXT,
	correlation_id TEXT NOT NULL,
	reason TEXT,
	before_json TEXT CHECK (before_json IS NULL OR json_valid(before_json)),
	after_json TEXT CHECK (after_json IS NULL OR json_valid(after_json)),
	source TEXT NOT NULL CHECK (source IN ('desktop','http','system','sync'))
) STRICT;
CREATE INDEX ix_audit_entity ON audit_logs(entity_type, entity_id, occurred_at);
CREATE INDEX ix_audit_actor ON audit_logs(actor_user_id, occurred_at);
CREATE INDEX ix_audit_time ON audit_logs(occurred_at);

CREATE TABLE application_settings (
	key TEXT PRIMARY KEY,
	value_json TEXT NOT NULL CHECK (json_valid(value_json)),
	updated_by TEXT REFERENCES users(id),
	updated_at TEXT NOT NULL
) STRICT;

CREATE TABLE backup_history (
	id TEXT PRIMARY KEY,
	kind TEXT NOT NULL CHECK (kind IN ('manual','scheduled','pre_migration','pre_restore')),
	destination TEXT NOT NULL,
	file_name TEXT NOT NULL,
	size_bytes INTEGER,
	sha256 TEXT,
	schema_version INTEGER NOT NULL,
	status TEXT NOT NULL CHECK (status IN ('running','succeeded','failed','verified')),
	error TEXT,
	started_by TEXT REFERENCES users(id),
	started_at TEXT NOT NULL,
	completed_at TEXT
) STRICT;
CREATE INDEX ix_backup_history_time ON backup_history(started_at);

-- 3. Synchronization & Devices
CREATE TABLE devices (
	id TEXT PRIMARY KEY,
	device_code TEXT NOT NULL UNIQUE CHECK (length(device_code) BETWEEN 1 AND 4),
	name TEXT NOT NULL,
	kind TEXT NOT NULL CHECK (kind IN ('desktop','server','lan_client')),
	is_local INTEGER NOT NULL DEFAULT 0 CHECK (is_local IN (0,1)),
	public_key TEXT,
	status TEXT NOT NULL CHECK (status IN ('active','revoked')),
	registered_at TEXT NOT NULL,
	last_seen_at TEXT
) STRICT;
CREATE UNIQUE INDEX ux_devices_one_local ON devices(is_local) WHERE is_local = 1;

CREATE TABLE sync_outbox (
	seq INTEGER PRIMARY KEY AUTOINCREMENT,
	operation_id TEXT NOT NULL UNIQUE,
	device_id TEXT NOT NULL REFERENCES devices(id),
	entity_type TEXT NOT NULL,
	entity_id TEXT NOT NULL,
	operation_type TEXT NOT NULL CHECK (operation_type IN ('create','update','delete','command')),
	command_name TEXT,
	base_revision INTEGER,
	new_revision INTEGER,
	payload_json TEXT NOT NULL CHECK (json_valid(payload_json)),
	payload_schema_version INTEGER NOT NULL,
	idempotency_key TEXT NOT NULL,
	occurred_at TEXT NOT NULL,
	status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','sent','acked','failed','dead')),
	attempts INTEGER NOT NULL DEFAULT 0,
	next_attempt_at TEXT,
	last_error TEXT
);
CREATE INDEX ix_outbox_pending ON sync_outbox(status, next_attempt_at, seq);

CREATE TABLE sync_checkpoints (
	id TEXT PRIMARY KEY,
	peer TEXT NOT NULL,
	direction TEXT NOT NULL CHECK (direction IN ('push','pull')),
	cursor TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	UNIQUE (peer, direction)
) STRICT;

CREATE TABLE sync_conflicts (
	id TEXT PRIMARY KEY,
	entity_type TEXT NOT NULL,
	entity_id TEXT NOT NULL,
	local_operation_id TEXT,
	remote_operation_id TEXT,
	base_revision INTEGER,
	local_payload_json TEXT,
	remote_payload_json TEXT,
	policy TEXT NOT NULL,
	status TEXT NOT NULL CHECK (status IN ('open','resolved_local','resolved_remote','resolved_merged','dismissed')),
	resolution_json TEXT,
	detected_at TEXT NOT NULL,
	resolved_by TEXT REFERENCES users(id),
	resolved_at TEXT
) STRICT;
CREATE INDEX ix_conflicts_open ON sync_conflicts(status, detected_at);

CREATE TABLE sync_tombstones (
	entity_type TEXT NOT NULL,
	entity_id TEXT NOT NULL,
	revision INTEGER NOT NULL,
	deleted_at TEXT NOT NULL,
	deleted_by TEXT,
	device_id TEXT,
	PRIMARY KEY (entity_type, entity_id)
) STRICT;

CREATE TABLE sync_schema_versions (
	entity_type TEXT NOT NULL,
	payload_schema_version INTEGER NOT NULL,
	introduced_in_app_version TEXT NOT NULL,
	PRIMARY KEY (entity_type, payload_schema_version)
) STRICT;

-- 4. Audit immutability triggers
CREATE TRIGGER trg_audit_logs_no_delete BEFORE DELETE ON audit_logs
BEGIN
	SELECT RAISE(ABORT, 'IMMUTABLE: audit logs cannot be deleted');
END;

CREATE TRIGGER trg_audit_logs_no_update BEFORE UPDATE ON audit_logs
BEGIN
	SELECT RAISE(ABORT, 'IMMUTABLE: audit logs cannot be updated');
END;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sync_schema_versions;
DROP TABLE IF EXISTS sync_tombstones;
DROP TABLE IF EXISTS sync_conflicts;
DROP TABLE IF EXISTS sync_checkpoints;
DROP TABLE IF EXISTS sync_outbox;
DROP TABLE IF EXISTS devices;
DROP TABLE IF EXISTS backup_history;
DROP TABLE IF EXISTS application_settings;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS recovery_codes;
DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS classes;
DROP TABLE IF EXISTS academic_terms;
DROP TABLE IF EXISTS academic_years;
DROP TABLE IF EXISTS campuses;
DROP TABLE IF EXISTS school_settings;
DROP TABLE IF EXISTS schools;
-- +goose StatementEnd

