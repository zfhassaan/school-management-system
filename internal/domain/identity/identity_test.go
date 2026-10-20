package identity

import (
	"testing"
	"time"

	"school-management/internal/domain/shared"
)

func TestLockoutEscalates(t *testing.T) {
	now := time.Date(2026, 10, 20, 7, 0, 0, 0, time.UTC)
	u := &User{}
	p := DefaultLockoutPolicy
	for i := 0; i < p.MaxAttempts-1; i++ {
		if u.RegisterFailedLogin(now, p) {
			t.Fatalf("locked too early at attempt %d", i+1)
		}
	}
	if !u.RegisterFailedLogin(now, p) {
		t.Fatal("expected lock at threshold")
	}
	if !u.IsLocked(now) || u.LockedUntil.Sub(now) != 15*time.Minute {
		t.Fatalf("unexpected first lock %v", u.LockedUntil)
	}
	for i := 0; i < p.MaxAttempts; i++ {
		u.RegisterFailedLogin(now, p)
	}
	if u.LockedUntil.Sub(now) != 30*time.Minute {
		t.Fatalf("expected doubled lock, got %v", u.LockedUntil.Sub(now))
	}
	for range 20 {
		for i := 0; i < p.MaxAttempts; i++ {
			u.RegisterFailedLogin(now, p)
		}
	}
	if u.LockedUntil.Sub(now) != p.MaxDuration {
		t.Fatalf("expected cap at %v, got %v", p.MaxDuration, u.LockedUntil.Sub(now))
	}
	if u.IsLocked(now.Add(25 * time.Hour)) {
		t.Fatal("lock must expire")
	}
	u.RegisterSuccessfulLogin(now)
	if u.LockedUntil != nil || u.LockoutCount != 0 || u.FailedLoginCount != 0 {
		t.Fatal("successful login must reset throttling")
	}
}

func TestPasswordPolicy(t *testing.T) {
	cases := map[string]bool{
		"short1":            false,
		"onlyletterslong":   false,
		"1234567890123":     false,
		"Correct horse 42!": true,
		"adminadmin1":       false, // equals username
	}
	for pw, ok := range cases {
		var v shared.Validator
		DefaultPasswordPolicy.Validate(&v, "password", pw, "adminadmin1")
		if v.Valid() != ok {
			t.Errorf("password %q valid=%v want %v (%v)", pw, v.Valid(), ok, v.Err())
		}
	}
}

func TestValidateProfileAndRole(t *testing.T) {
	var v shared.Validator
	ValidateProfile(&v, "ab", "", "not-an-email")
	err := v.Err().(*shared.Error)
	for _, f := range []string{"username", "displayName", "email"} {
		if _, ok := err.Fields[f]; !ok {
			t.Errorf("expected field error for %s", f)
		}
	}
	var rv shared.Validator
	ValidateRole(&rv, "accounts_lead", "Accounts Lead", []Permission{PermPaymentsRecord, "bogus.perm"})
	if rv.Valid() {
		t.Fatal("unknown permission must fail")
	}
}

func TestSystemRolesUseKnownPermissions(t *testing.T) {
	seen := map[string]bool{}
	for _, r := range SystemRoles() {
		if seen[r.Code] {
			t.Fatalf("duplicate role %s", r.Code)
		}
		seen[r.Code] = true
		for _, p := range r.Permissions {
			if !IsKnownPermission(p) {
				t.Fatalf("role %s has unknown permission %s", r.Code, p)
			}
		}
	}
	codes := map[Permission]bool{}
	for _, p := range Catalog {
		if codes[p.Code] {
			t.Fatalf("duplicate permission %s", p.Code)
		}
		codes[p.Code] = true
	}
}

func TestSessionExpiry(t *testing.T) {
	now := time.Date(2026, 10, 20, 7, 0, 0, 0, time.UTC)
	s := Session{IdleExpiresAt: now.Add(time.Minute), AbsoluteExpiresAt: now.Add(time.Hour)}
	if s.Expired(now) {
		t.Fatal("should be valid")
	}
	if !s.Expired(now.Add(time.Minute)) {
		t.Fatal("idle expiry")
	}
}

