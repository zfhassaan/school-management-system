package guardians

import (
	"context"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"school-management/internal/domain/shared"
)

type Relationship string

const (
	RelFather           Relationship = "father"
	RelMother           Relationship = "mother"
	RelGuardian         Relationship = "guardian"
	RelGrandparent      Relationship = "grandparent"
	RelSibling          Relationship = "sibling"
	RelRelative         Relationship = "relative"
	RelEmergencyContact Relationship = "emergency_contact"
	RelOther            Relationship = "other"
)

type Guardian struct {
	ID                     shared.ID
	SchoolID               shared.ID
	FirstName              string
	LastName               string
	PhonePrimary           string
	PhonePrimaryNormalized string
	PhoneSecondary         string
	Email                  string
	Occupation             string
	NationalID             string
	Address                shared.Address
	CreatedAt              time.Time
	UpdatedAt              time.Time
	Revision               int64
}

var nonDigitPattern = regexp.MustCompile(`\D`)

// NormalizePhone strips all non-digit characters for robust lookup and sibling detection.
func NormalizePhone(s string) string {
	digits := nonDigitPattern.ReplaceAllString(s, "")
	// If starts with 00, trim it
	if strings.HasPrefix(digits, "00") {
		digits = digits[2:]
	}
	return digits
}

func (g *Guardian) Validate() error {
	var v shared.Validator
	v.Check(strings.TrimSpace(g.FirstName) != "", "firstName", "first name is required")
	v.Check(strings.TrimSpace(g.PhonePrimary) != "", "phonePrimary", "primary phone is required")
	norm := NormalizePhone(g.PhonePrimary)
	v.Check(len(norm) >= 7, "phonePrimary", "valid phone number with at least 7 digits is required")
	if g.Email != "" {
		_, err := mail.ParseAddress(g.Email)
		v.Check(err == nil, "email", "invalid email address")
	}
	g.PhonePrimaryNormalized = norm
	return v.Err()
}

type StudentGuardianLink struct {
	StudentID          shared.ID
	GuardianID         shared.ID
	Relationship       Relationship
	IsPrimary          bool
	IsEmergencyContact bool
	ReceivesBilling    bool
	CanPickup          bool
	Priority           int
	// Guardian expanded view
	FirstName      string
	LastName       string
	PhonePrimary   string
	PhoneSecondary string
	Email          string
	Occupation     string
}

type Repository interface {
	Get(ctx context.Context, id shared.ID) (*Guardian, error)
	GetByPhone(ctx context.Context, schoolID shared.ID, phoneNormalized string) (*Guardian, error)
	Insert(ctx context.Context, g *Guardian) error
	Update(ctx context.Context, g *Guardian) error
	ListForStudent(ctx context.Context, studentID shared.ID) ([]StudentGuardianLink, error)
	SetStudentGuardians(ctx context.Context, studentID shared.ID, links []StudentGuardianLink) error
}
