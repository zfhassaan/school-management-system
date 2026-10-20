// Package school models the single school profile and its operational settings.
package school

import (
	"context"
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"school-management/internal/domain/shared"
)

// Address is a postal address alias for shared.Address.
type Address = shared.Address

// School is the tenant root. The system manages exactly one school.
type School struct {
	ID                 shared.ID
	Name               string
	LegalName          string
	Code               string
	LogoDocumentID     *shared.ID
	Address            Address
	Phone              string
	Email              string
	Website            string
	Timezone           string
	CurrencyCode       string
	CurrencyMinorUnits int
	Locale             string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Revision           int64
}

var (
	schoolCodePattern = regexp.MustCompile(`^[A-Z0-9][A-Z0-9-]{1,19}$`)
	currencyPattern   = regexp.MustCompile(`^[A-Z]{3}$`)
)

// Validate checks the school profile.
func (s *School) Validate() error {
	var v shared.Validator
	name := strings.TrimSpace(s.Name)
	v.Check(name != "" && utf8.RuneCountInString(name) <= 200, "name", "required, at most 200 characters")
	v.Check(schoolCodePattern.MatchString(s.Code), "code", "2-20 characters: uppercase letters, digits or '-'")
	_, err := time.LoadLocation(s.Timezone)
	v.Check(s.Timezone != "" && err == nil, "timezone", "must be a valid IANA timezone")
	v.Check(currencyPattern.MatchString(s.CurrencyCode), "currencyCode", "must be a 3-letter ISO code")
	v.Check(s.CurrencyMinorUnits >= 0 && s.CurrencyMinorUnits <= 3, "currencyMinorUnits", "must be between 0 and 3")
	if s.Email != "" {
		_, err := mail.ParseAddress(s.Email)
		v.Check(err == nil, "email", "invalid email address")
	}
	v.Check(utf8.RuneCountInString(s.Phone) <= 30, "phone", "at most 30 characters")
	return v.Err()
}

// Location returns the school's timezone location.
func (s *School) Location() *time.Location {
	loc, err := time.LoadLocation(s.Timezone)
	if err != nil {
		return time.UTC
	}
	return loc
}

// DiscountCombination selects how multiple discounts combine (fee engine).
type DiscountCombination string

const (
	ExclusiveOverrides DiscountCombination = "exclusive_overrides"
	BestForStudent     DiscountCombination = "best_for_student"
)

// Settings holds operational configuration for the school.
type Settings struct {
	SchoolID              shared.ID
	AdmissionNoFormat     string
	InvoicePrefix         string
	ReceiptPrefix         string
	RoundingUnitMinor     int64
	RoundingMode          shared.RoundingMode
	DefaultDueDay         int
	LateFeeGraceDays      int
	DiscountCombination   DiscountCombination
	RequireSecondApprover bool
	DateFormat            string
	UpdatedAt             time.Time
	Revision              int64
}

// DefaultSettings returns sensible defaults for a new school.
func DefaultSettings(schoolID shared.ID) Settings {
	return Settings{
		SchoolID:            schoolID,
		AdmissionNoFormat:   "{YYYY}-{SEQ:5}",
		InvoicePrefix:       "INV",
		ReceiptPrefix:       "RCP",
		RoundingUnitMinor:   1,
		RoundingMode:        shared.RoundHalfUp,
		DefaultDueDay:       10,
		DiscountCombination: ExclusiveOverrides,
		DateFormat:          "dd/MM/yyyy",
	}
}

var (
	prefixPattern       = regexp.MustCompile(`^[A-Z0-9-]{1,10}$`)
	admissionFmtPattern = regexp.MustCompile(`\{SEQ(:[1-9])?\}`)
)

// Validate checks settings consistency.
func (s *Settings) Validate() error {
	var v shared.Validator
	v.Check(admissionFmtPattern.MatchString(s.AdmissionNoFormat) && len(s.AdmissionNoFormat) <= 40,
		"admissionNoFormat", "must contain {SEQ} or {SEQ:n}, at most 40 characters")
	v.Check(prefixPattern.MatchString(s.InvoicePrefix), "invoicePrefix", "1-10 characters: A-Z, 0-9 or '-'")
	v.Check(prefixPattern.MatchString(s.ReceiptPrefix), "receiptPrefix", "1-10 characters: A-Z, 0-9 or '-'")
	v.Check(s.RoundingUnitMinor >= 1 && s.RoundingUnitMinor <= 100_000, "roundingUnitMinor", "must be between 1 and 100000")
	switch s.RoundingMode {
	case shared.RoundHalfUp, shared.RoundHalfEven, shared.RoundDown, shared.RoundUp:
	default:
		v.Add("roundingMode", "unsupported rounding mode")
	}
	v.Check(s.DefaultDueDay >= 1 && s.DefaultDueDay <= 28, "defaultDueDay", "must be between 1 and 28")
	v.Check(s.LateFeeGraceDays >= 0 && s.LateFeeGraceDays <= 365, "lateFeeGraceDays", "must be between 0 and 365")
	v.Check(s.DiscountCombination == ExclusiveOverrides || s.DiscountCombination == BestForStudent,
		"discountCombination", "unsupported combination policy")
	switch s.DateFormat {
	case "dd/MM/yyyy", "MM/dd/yyyy", "yyyy-MM-dd", "dd-MMM-yyyy":
	default:
		v.Add("dateFormat", "unsupported date format")
	}
	return v.Err()
}

// Repository persists the school profile and settings.
type Repository interface {
	// Get returns the single school or shared.ErrNotFound before setup.
	Get(ctx context.Context) (*School, error)
	Insert(ctx context.Context, s *School) error
	Update(ctx context.Context, s *School) error
	GetSettings(ctx context.Context, schoolID shared.ID) (*Settings, error)
	InsertSettings(ctx context.Context, s *Settings, at time.Time) error
	UpdateSettings(ctx context.Context, s *Settings, at time.Time) error
}

