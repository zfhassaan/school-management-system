package shared

import (
	"time"

	"github.com/google/uuid"
)

// ID is a UUIDv7 identifier rendered in canonical 36-character form.
type ID string

func (id ID) String() string { return string(id) }

// IsZero reports whether the identifier is empty.
func (id ID) IsZero() bool { return id == "" }

// ParseID validates s as a UUID and returns it in canonical lowercase form.
func ParseID(s string) (ID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return "", &Error{Code: CodeValidation, Message: "invalid identifier", Fields: map[string]string{"id": "must be a UUID"}}
	}
	return ID(u.String()), nil
}

// IDGenerator creates new identifiers. Injected so tests can be deterministic.
type IDGenerator interface {
	NewID() ID
}

// UUIDv7Generator produces time-ordered UUIDv7 identifiers.
type UUIDv7Generator struct{}

func (UUIDv7Generator) NewID() ID {
	return ID(uuid.Must(uuid.NewV7()).String())
}

// Clock abstracts time so business rules are testable and timezone-aware.
type Clock interface {
	// Now returns the current instant in UTC.
	Now() time.Time
}

// SystemClock is the production clock.
type SystemClock struct{}

func (SystemClock) Now() time.Time { return time.Now().UTC() }

// FixedClock is a test clock that can be advanced.
type FixedClock struct{ T time.Time }

func (c *FixedClock) Now() time.Time { return c.T.UTC() }

// Advance moves the clock forward by d.
func (c *FixedClock) Advance(d time.Duration) { c.T = c.T.Add(d) }

