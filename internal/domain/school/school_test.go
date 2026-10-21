package school_test

import (
	"testing"
	"time"

	"school-management/internal/domain/school"
	"school-management/internal/domain/shared"
)

func TestSchoolTimezoneValidation(t *testing.T) {
	timezones := []string{
		"Asia/Karachi",
		"UTC",
		"Asia/Dubai",
		"Asia/Riyadh",
		"Asia/Dhaka",
		"Asia/Kolkata",
		"America/New_York",
		"Europe/London",
	}

	for _, tz := range timezones {
		s := &school.School{
			ID:                 shared.ID("018f0000-0000-7000-8000-000000000001"),
			Name:               "Test School",
			Code:               "TEST-01",
			Timezone:           tz,
			CurrencyCode:       "PKR",
			CurrencyMinorUnits: 2,
		}
		if err := s.Validate(); err != nil {
			t.Errorf("expected timezone %q to validate successfully, got err: %v", tz, err)
		}
	}
}

func TestSchoolInvalidTimezone(t *testing.T) {
	s := &school.School{
		ID:                 shared.ID("018f0000-0000-7000-8000-000000000001"),
		Name:               "Test School",
		Code:               "TEST-01",
		Timezone:           "Invalid/NonExistent_Zone",
		CurrencyCode:       "PKR",
		CurrencyMinorUnits: 2,
	}
	err := s.Validate()
	if err == nil {
		t.Errorf("expected error for invalid timezone, got nil")
	}
}

func TestTimeLoadLocationDirect(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Karachi")
	if err != nil {
		t.Fatalf("failed to load embedded Asia/Karachi location: %v", err)
	}
	if loc.String() != "Asia/Karachi" {
		t.Errorf("expected Asia/Karachi, got %v", loc.String())
	}
}
