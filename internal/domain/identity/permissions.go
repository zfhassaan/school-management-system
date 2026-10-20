package identity

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

	"school-management/internal/domain/shared"
)

// Permission is a fine-grained capability code such as "payments.record".
type Permission string

// PermissionInfo describes a permission in the catalog.
type PermissionInfo struct {
	Code        Permission
	Module      string
	Description string
	// Sensitive permissions require a reason (and optionally re-authentication).
	Sensitive bool
}

// Permission catalog. The Go catalog is the single source of truth and is
// synchronized into the database at startup.
const (
	PermDashboardView Permission = "dashboard.view"

	PermSchoolView     Permission = "school.view"
	PermSettingsManage Permission = "settings.manage"

	PermUsersView   Permission = "users.view"
	PermUsersManage Permission = "users.manage"
	PermRolesManage Permission = "roles.manage"

	PermAcademicsView    Permission = "academics.view"
	PermAcademicsManage  Permission = "academics.manage"
	PermAcademicsPromote Permission = "academics.promote"

	PermStudentsView   Permission = "students.view"
	PermStudentsCreate Permission = "students.create"
	PermStudentsEdit   Permission = "students.edit"
	PermStudentsStatus Permission = "students.status"

	PermGuardiansView   Permission = "guardians.view"
	PermGuardiansManage Permission = "guardians.manage"

	PermDocumentsView   Permission = "documents.view"
	PermDocumentsManage Permission = "documents.manage"

	PermAdmissionsView    Permission = "admissions.view"
	PermAdmissionsCreate  Permission = "admissions.create"
	PermAdmissionsReview  Permission = "admissions.review"
	PermAdmissionsDecide  Permission = "admissions.decide"
	PermAdmissionsEnroll  Permission = "admissions.enroll"
	PermAdmissionFormsMgt Permission = "admissions.forms.manage"

	PermFeesView    Permission = "fees.view"
	PermFeesEdit    Permission = "fees.edit"
	PermFeesPublish Permission = "fees.publish"
	PermFeesAssign  Permission = "fees.assign"
	PermFeesApprove Permission = "fees.approve"

	PermInvoicesView     Permission = "invoices.view"
	PermInvoicesGenerate Permission = "invoices.generate"
	PermInvoicesVoid     Permission = "invoices.void"

	PermPaymentsView    Permission = "payments.view"
	PermPaymentsRecord  Permission = "payments.record"
	PermPaymentsReverse Permission = "payments.reverse"
	PermPaymentsRefund  Permission = "payments.refund"

	PermAdjustmentsRequest Permission = "adjustments.request"
	PermAdjustmentsApprove Permission = "adjustments.approve"

	PermReportsStudents  Permission = "reports.students"
	PermReportsFinance   Permission = "reports.finance"
	PermReportsAdmission Permission = "reports.admissions"

	PermBackupManage Permission = "backup.manage"
	PermSyncManage   Permission = "sync.manage"
	PermAuditView    Permission = "audit.view"
)

// Catalog lists every permission known to the application.
var Catalog = []PermissionInfo{
	{PermDashboardView, "dashboard", "View dashboard", false},
	{PermSchoolView, "settings", "View school profile and settings", false},
	{PermSettingsManage, "settings", "Manage school profile and settings", true},
	{PermUsersView, "users", "View users", false},
	{PermUsersManage, "users", "Create, edit, disable users and reset passwords", true},
	{PermRolesManage, "users", "Create and edit roles and permissions", true},
	{PermAcademicsView, "academics", "View academic years, classes and sections", false},
	{PermAcademicsManage, "academics", "Manage academic years, classes and sections", false},
	{PermAcademicsPromote, "academics", "Promote students between academic years", true},
	{PermStudentsView, "students", "View students", false},
	{PermStudentsCreate, "students", "Create students", false},
	{PermStudentsEdit, "students", "Edit student profiles", false},
	{PermStudentsStatus, "students", "Change student status", true},
	{PermGuardiansView, "guardians", "View guardians", false},
	{PermGuardiansManage, "guardians", "Manage guardians", false},
	{PermDocumentsView, "documents", "View student documents", false},
	{PermDocumentsManage, "documents", "Upload and remove student documents", false},
	{PermAdmissionsView, "admissions", "View admission applications", false},
	{PermAdmissionsCreate, "admissions", "Create admission applications", false},
	{PermAdmissionsReview, "admissions", "Review admission applications", false},
	{PermAdmissionsDecide, "admissions", "Approve or reject admissions", true},
	{PermAdmissionsEnroll, "admissions", "Enroll approved applicants", false},
	{PermAdmissionFormsMgt, "admissions", "Design and publish admission forms", true},
	{PermFeesView, "fees", "View fee structures", false},
	{PermFeesEdit, "fees", "Edit draft fee structures, rules and discounts", false},
	{PermFeesPublish, "fees", "Publish and activate fee structures", true},
	{PermFeesAssign, "fees", "Assign fees and overrides to students", false},
	{PermFeesApprove, "fees", "Approve overrides and concessions", true},
	{PermInvoicesView, "billing", "View invoices", false},
	{PermInvoicesGenerate, "billing", "Generate and issue invoices", false},
	{PermInvoicesVoid, "billing", "Void invoices", true},
	{PermPaymentsView, "billing", "View payments and receipts", false},
	{PermPaymentsRecord, "billing", "Record payments", false},
	{PermPaymentsReverse, "billing", "Reverse payments", true},
	{PermPaymentsRefund, "billing", "Refund payments and credit balances", true},
	{PermAdjustmentsRequest, "billing", "Request financial adjustments", false},
	{PermAdjustmentsApprove, "billing", "Approve financial adjustments", true},
	{PermReportsStudents, "reports", "Student reports", false},
	{PermReportsFinance, "reports", "Financial reports", false},
	{PermReportsAdmission, "reports", "Admission reports", false},
	{PermBackupManage, "system", "Create backups and restore data", true},
	{PermSyncManage, "system", "Configure cloud synchronization", true},
	{PermAuditView, "system", "View audit logs", false},
}

// IsKnownPermission reports whether p is in the catalog.
func IsKnownPermission(p Permission) bool {
	return slices.ContainsFunc(Catalog, func(i PermissionInfo) bool { return i.Code == p })
}

// System role codes.
const (
	RoleSuperAdmin       = "super_admin"
	RoleSchoolAdmin      = "school_admin"
	RoleAdmissionOfficer = "admission_officer"
	RoleAccountant       = "accountant"
	RoleTeacher          = "teacher"
	RoleReadOnly         = "read_only"
)

// SystemRoleTemplate seeds a system role. Templates are applied only when the
// role does not exist yet, so administrators can customize them later
// (except super_admin, which always has every permission).
type SystemRoleTemplate struct {
	Code        string
	Name        string
	Description string
	Permissions []Permission
}

func allPermissions() []Permission {
	out := make([]Permission, len(Catalog))
	for i, p := range Catalog {
		out[i] = p.Code
	}
	return out
}

func except(excluded ...Permission) []Permission {
	var out []Permission
	for _, p := range allPermissions() {
		if !slices.Contains(excluded, p) {
			out = append(out, p)
		}
	}
	return out
}

// SystemRoles returns the seeded roles.
func SystemRoles() []SystemRoleTemplate {
	return []SystemRoleTemplate{
		{RoleSuperAdmin, "Super Administrator", "Unrestricted access, including users, roles and data restore", allPermissions()},
		{RoleSchoolAdmin, "School Administrator", "Manages school operations except role definitions and sync", except(PermRolesManage, PermSyncManage)},
		{RoleAdmissionOfficer, "Admission Officer", "Handles admissions, students and guardians", []Permission{
			PermDashboardView, PermSchoolView, PermAcademicsView, PermStudentsView, PermStudentsCreate, PermStudentsEdit,
			PermGuardiansView, PermGuardiansManage, PermDocumentsView, PermDocumentsManage, PermAdmissionsView,
			PermAdmissionsCreate, PermAdmissionsReview, PermAdmissionsEnroll, PermReportsAdmission, PermReportsStudents,
		}},
		{RoleAccountant, "Accountant", "Manages fees, invoices and payments", []Permission{
			PermDashboardView, PermSchoolView, PermAcademicsView, PermStudentsView, PermGuardiansView, PermFeesView,
			PermFeesEdit, PermFeesAssign, PermInvoicesView, PermInvoicesGenerate, PermPaymentsView, PermPaymentsRecord,
			PermAdjustmentsRequest, PermReportsFinance, PermReportsStudents,
		}},
		{RoleTeacher, "Teacher", "Read-only access to students and academics", []Permission{
			PermDashboardView, PermSchoolView, PermAcademicsView, PermStudentsView, PermGuardiansView,
		}},
		{RoleReadOnly, "Read-only Staff", "Read-only access to operational data", []Permission{
			PermDashboardView, PermSchoolView, PermAcademicsView, PermStudentsView, PermGuardiansView, PermDocumentsView,
			PermAdmissionsView, PermFeesView, PermInvoicesView, PermPaymentsView,
		}},
	}
}

// Role groups permissions and is assigned to users.
type Role struct {
	ID          shared.ID
	Code        string
	Name        string
	Description string
	IsSystem    bool
	Permissions []Permission
	UserCount   int
	Revision    int64
}

var roleCodePattern = regexp.MustCompile(`^[a-z][a-z0-9_]{1,49}$`)

// ValidateRole validates editable role fields.
func ValidateRole(v *shared.Validator, code, name string, perms []Permission) {
	v.Check(roleCodePattern.MatchString(code), "code", "2-50 characters: lowercase letters, digits or '_'")
	n := strings.TrimSpace(name)
	v.Check(n != "" && utf8.RuneCountInString(n) <= 100, "name", "required, at most 100 characters")
	for _, p := range perms {
		if !IsKnownPermission(p) {
			v.Add("permissions", "unknown permission "+string(p))
			break
		}
	}
}

// PermissionSet is an immutable set of permissions for fast checks.
type PermissionSet map[Permission]struct{}

func NewPermissionSet(perms []Permission) PermissionSet {
	s := make(PermissionSet, len(perms))
	for _, p := range perms {
		s[p] = struct{}{}
	}
	return s
}

func (s PermissionSet) Has(p Permission) bool {
	_, ok := s[p]
	return ok
}

// Sorted returns the permissions in deterministic order.
func (s PermissionSet) Sorted() []Permission {
	out := make([]Permission, 0, len(s))
	for p := range s {
		out = append(out, p)
	}
	slices.Sort(out)
	return out
}

