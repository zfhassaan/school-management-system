-- name: GetSchool :one
SELECT * FROM schools LIMIT 1;

-- name: InsertSchool :exec
INSERT INTO schools (
    id, name, legal_name, code, logo_document_id,
    address_line1, address_line2, city, region, postal_code, country,
    phone, email, website, timezone, currency_code, currency_minor_units,
    locale, created_at, updated_at, row_version
) VALUES (
    ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?
);

-- name: UpdateSchool :exec
UPDATE schools SET
    name = ?, legal_name = ?, code = ?, logo_document_id = ?,
    address_line1 = ?, address_line2 = ?, city = ?, region = ?, postal_code = ?, country = ?,
    phone = ?, email = ?, website = ?, timezone = ?, currency_code = ?, currency_minor_units = ?,
    locale = ?, updated_at = ?, row_version = row_version + 1
WHERE id = ? AND row_version = ?;

-- name: GetSchoolSettings :one
SELECT * FROM school_settings WHERE school_id = ?;

-- name: InsertSchoolSettings :exec
INSERT INTO school_settings (
    school_id, admission_no_format, invoice_prefix, receipt_prefix,
    rounding_unit_minor, rounding_mode, default_due_day, late_fee_grace_days,
    discount_combination, require_second_approver, date_format,
    created_at, updated_at, row_version
) VALUES (
    ?, ?, ?, ?,
    ?, ?, ?, ?,
    ?, ?, ?,
    ?, ?, ?
);

-- name: UpdateSchoolSettings :exec
UPDATE school_settings SET
    admission_no_format = ?, invoice_prefix = ?, receipt_prefix = ?,
    rounding_unit_minor = ?, rounding_mode = ?, default_due_day = ?, late_fee_grace_days = ?,
    discount_combination = ?, require_second_approver = ?, date_format = ?,
    updated_at = ?, row_version = row_version + 1
WHERE school_id = ? AND row_version = ?;

-- name: CountUsers :one
SELECT count(*) FROM users;

-- name: GetUser :one
SELECT * FROM users WHERE id = ?;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ? COLLATE NOCASE;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at ASC;

-- name: InsertUser :exec
INSERT INTO users (
    id, username, display_name, email, password_hash, password_changed_at,
    must_change_password, status, failed_login_count, lockout_count, locked_until,
    last_login_at, created_at, updated_at, row_version
) VALUES (
    ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?,
    ?, ?, ?, ?
);

-- name: UpdateUser :exec
UPDATE users SET
    username = ?, display_name = ?, email = ?, password_hash = ?, password_changed_at = ?,
    must_change_password = ?, status = ?, failed_login_count = ?, lockout_count = ?,
    locked_until = ?, last_login_at = ?, updated_at = ?, row_version = row_version + 1
WHERE id = ? AND row_version = ?;

-- name: ClearUserRoles :exec
DELETE FROM user_roles WHERE user_id = ?;

-- name: InsertUserRole :exec
INSERT INTO user_roles (user_id, role_id, assigned_by, assigned_at)
VALUES (?, ?, ?, ?);

-- name: GetUserRoles :many
SELECT r.* FROM roles r
JOIN user_roles ur ON ur.role_id = r.id
WHERE ur.user_id = ?
ORDER BY r.name ASC;

-- name: GetUserRoleCodes :many
SELECT r.code FROM roles r
JOIN user_roles ur ON ur.role_id = r.id
WHERE ur.user_id = ?
ORDER BY r.code ASC;

-- name: GetEffectivePermissions :many
SELECT DISTINCT p.code FROM permissions p
JOIN role_permissions rp ON rp.permission_code = p.code
JOIN user_roles ur ON ur.role_id = rp.role_id
WHERE ur.user_id = ?
ORDER BY p.code ASC;

-- name: CountActiveWithRole :one
SELECT count(DISTINCT u.id) FROM users u
JOIN user_roles ur ON ur.user_id = u.id
JOIN roles r ON r.id = ur.role_id
WHERE r.code = ? AND u.status = 'active';

-- name: GetRole :one
SELECT * FROM roles WHERE id = ?;

-- name: GetRoleByCode :one
SELECT * FROM roles WHERE code = ?;

-- name: ListRolesWithUserCount :many
SELECT r.*, count(ur.user_id) as user_count
FROM roles r
LEFT JOIN user_roles ur ON ur.role_id = r.id
GROUP BY r.id
ORDER BY r.is_system DESC, r.name ASC;

-- name: InsertRole :exec
INSERT INTO roles (id, code, name, description, is_system, created_at, updated_at, row_version)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateRole :exec
UPDATE roles SET
    name = ?, description = ?, updated_at = ?, row_version = row_version + 1
WHERE id = ? AND row_version = ? AND is_system = 0;

-- name: DeleteRole :exec
DELETE FROM roles WHERE id = ? AND is_system = 0;

-- name: ClearRolePermissions :exec
DELETE FROM role_permissions WHERE role_id = ?;

-- name: InsertRolePermission :exec
INSERT INTO role_permissions (role_id, permission_code)
VALUES (?, ?);

-- name: GetRolePermissions :many
SELECT permission_code FROM role_permissions
WHERE role_id = ?
ORDER BY permission_code ASC;

-- name: UpsertPermission :exec
INSERT INTO permissions (code, module, description, is_sensitive)
VALUES (?, ?, ?, ?)
ON CONFLICT(code) DO UPDATE SET
    module = excluded.module,
    description = excluded.description,
    is_sensitive = excluded.is_sensitive;

-- name: ListPermissions :many
SELECT * FROM permissions ORDER BY module ASC, code ASC;

-- name: ReplaceUserRecoveryCodes :exec
DELETE FROM recovery_codes WHERE user_id = ?;

-- name: InsertRecoveryCode :exec
INSERT INTO recovery_codes (id, user_id, code_hash, used_at, created_at)
VALUES (?, ?, ?, ?, ?);

-- name: ListUnusedRecoveryCodes :many
SELECT * FROM recovery_codes WHERE user_id = ? AND used_at IS NULL;

-- name: MarkRecoveryCodeUsed :exec
UPDATE recovery_codes SET used_at = ? WHERE id = ?;

-- name: CreateSession :exec
INSERT INTO user_sessions (
    id, user_id, token_hash, device_id, created_at, last_seen_at,
    idle_expires_at, absolute_expires_at, revoked_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NULL);

-- name: GetSession :one
SELECT * FROM user_sessions WHERE id = ? AND revoked_at IS NULL;

-- name: TouchSession :exec
UPDATE user_sessions SET
    last_seen_at = ?, idle_expires_at = ?
WHERE id = ? AND revoked_at IS NULL;

-- name: RevokeSession :exec
UPDATE user_sessions SET revoked_at = ? WHERE id = ?;

-- name: RevokeAllUserSessions :exec
UPDATE user_sessions SET revoked_at = ? WHERE user_id = ? AND revoked_at IS NULL;

-- name: AppendAuditLog :exec
INSERT INTO audit_logs (
    id, occurred_at, actor_user_id, device_id, action, entity_type,
    entity_id, correlation_id, reason, before_json, after_json, source
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: ListAuditLogs :many
SELECT * FROM audit_logs
WHERE (? = '' OR entity_type = ?)
  AND (? = '' OR entity_id = ?)
  AND (? = '' OR actor_user_id = ?)
ORDER BY occurred_at DESC
LIMIT ?;

-- name: AppendOutboxEvent :exec
INSERT INTO sync_outbox (
    operation_id, device_id, entity_type, entity_id, operation_type,
    command_name, base_revision, new_revision, payload_json,
    payload_schema_version, idempotency_key, occurred_at, status
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'pending');

-- name: GetLocalDevice :one
SELECT * FROM devices WHERE is_local = 1 LIMIT 1;

-- name: InsertLocalDevice :exec
INSERT INTO devices (
    id, device_code, name, kind, is_local, public_key, status, registered_at, last_seen_at
) VALUES (?, ?, ?, ?, 1, ?, 'active', ?, ?);

