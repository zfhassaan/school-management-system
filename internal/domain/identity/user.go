// Package identity models local users, roles, permissions and sessions.
// Authentication works fully offline.
package identity

import (
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"school-management/internal/domain/shared"
)

type UserStatus string

const (
	UserActive   UserStatus = "active"
	UserDisabled UserStatus = "disabled"
)

func (s UserStatus) Valid() bool { return s == UserActive || s == UserDisabled }

// User is a local application account.
type User struct {
	ID                 shared.ID
	Username           string
	DisplayName        string
	Email              string
	PasswordHash       string
	PasswordChangedAt  time.Time
	MustChangePassword bool
	Status             UserStatus
	FailedLoginCount   int
	LockoutCount       int
	LockedUntil        *time.Time
	LastLoginAt        *time.Time
	RoleIDs            []shared.ID
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Revision           int64
}

var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{2,49}$`)

// NormalizeUsername trims and lowercases a username for storage and lookup.
func NormalizeUsername(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

// ValidateProfile validates the user-editable profile fields.
func ValidateProfile(v *shared.Validator, username, displayName, email string) {
	v.Check(usernamePattern.MatchString(username), "username", "3-50 characters: letters, digits, '.', '_' or '-'")
	dn := strings.TrimSpace(displayName)
	v.Check(dn != "" && utf8.RuneCountInString(dn) <= 100, "displayName", "required, at most 100 characters")
	if email != "" {
		_, err := mail.ParseAddress(email)
		v.Check(err == nil && len(email) <= 200, "email", "invalid email address")
	}
}

// LockoutPolicy controls throttling of failed login attempts.
type LockoutPolicy struct {
	MaxAttempts  int
	BaseDuration time.Duration
	MaxDuration  time.Duration
}

// DefaultLockoutPolicy locks for 15 minutes after 5 failures, doubling up to 24h.
var DefaultLockoutPolicy = LockoutPolicy{MaxAttempts: 5, BaseDuration: 15 * time.Minute, MaxDuration: 24 * time.Hour}

// IsLocked reports whether the account is currently locked.
func (u *User) IsLocked(now time.Time) bool {
	return u.LockedUntil != nil && now.Before(*u.LockedUntil)
}

// RegisterFailedLogin records a failed attempt and locks the account when the
// policy threshold is reached. It returns true when a new lock was applied.
func (u *User) RegisterFailedLogin(now time.Time, p LockoutPolicy) bool {
	u.FailedLoginCount++
	if u.FailedLoginCount < p.MaxAttempts {
		return false
	}
	d := p.BaseDuration
	for i := 0; i < u.LockoutCount && d < p.MaxDuration; i++ {
		d *= 2
	}
	if d > p.MaxDuration {
		d = p.MaxDuration
	}
	until := now.Add(d)
	u.LockedUntil = &until
	u.LockoutCount++
	u.FailedLoginCount = 0
	return true
}

// RegisterSuccessfulLogin clears throttling state.
func (u *User) RegisterSuccessfulLogin(now time.Time) {
	u.FailedLoginCount = 0
	u.LockoutCount = 0
	u.LockedUntil = nil
	u.LastLoginAt = &now
}

// Unlock clears a lock applied by the lockout policy.
func (u *User) Unlock() {
	u.FailedLoginCount = 0
	u.LockoutCount = 0
	u.LockedUntil = nil
}

// PasswordPolicy defines minimum password requirements.
type PasswordPolicy struct {
	MinLength int
}

var DefaultPasswordPolicy = PasswordPolicy{MinLength: 10}

// Validate checks pw against the policy. field names the input for errors.
func (p PasswordPolicy) Validate(v *shared.Validator, field, pw, username string) {
	n := utf8.RuneCountInString(pw)
	if n < p.MinLength {
		v.Add(field, "must be at least "+itoa(p.MinLength)+" characters")
		return
	}
	if n > 128 {
		v.Add(field, "must be at most 128 characters")
		return
	}
	var letter, digit bool
	for _, r := range pw {
		switch {
		case unicode.IsLetter(r):
			letter = true
		case unicode.IsDigit(r):
			digit = true
		}
	}
	v.Check(letter && digit, field, "must contain letters and digits")
	if username != "" {
		v.Check(!strings.EqualFold(pw, username), field, "must not equal the username")
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}

