package shared

import (
	"encoding/json"
	"strings"
	"time"
)

// Date is a calendar date without time or timezone (date of birth, due date).
// It must never be converted through a midnight time.Time in another zone.
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

const dateLayout = "2006-01-02"

// ParseDate parses YYYY-MM-DD strictly.
func ParseDate(s string) (Date, error) {
	t, err := time.Parse(dateLayout, s)
	if err != nil || t.Format(dateLayout) != s {
		return Date{}, &Error{Code: CodeValidation, Message: "invalid date", Fields: map[string]string{"date": "must be YYYY-MM-DD"}}
	}
	return DateOf(t), nil
}

// DateOf returns the calendar date of t in t's own location.
func DateOf(t time.Time) Date {
	y, m, d := t.Date()
	return Date{Year: y, Month: m, Day: d}
}

// CivilDate constructs a Date from year, month, and day components.
func CivilDate(year int, month time.Month, day int) Date {
	return Date{Year: year, Month: month, Day: day}
}

// TodayIn returns the civil date of instant now in the given location.
func TodayIn(now time.Time, loc *time.Location) Date {
	return DateOf(now.In(loc))
}

func (d Date) IsZero() bool { return d.Year == 0 && d.Month == 0 && d.Day == 0 }

func (d Date) String() string {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC).Format(dateLayout)
}

func (d Date) t() time.Time { return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC) }

func (d Date) Before(o Date) bool { return d.t().Before(o.t()) }
func (d Date) After(o Date) bool  { return d.t().After(o.t()) }

// AddDays returns the date n days later (negative n moves backwards).
func (d Date) AddDays(n int) Date { return DateOf(d.t().AddDate(0, 0, n)) }

// DaysUntil returns the number of days from d to o.
func (d Date) DaysUntil(o Date) int { return int(o.t().Sub(d.t()).Hours() / 24) }

// MarshalJSON marshals Date to a JSON string formatted as "YYYY-MM-DD".
// A zero Date marshals as null.
func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(d.String())
}

// UnmarshalJSON unmarshals Date from a JSON string ("YYYY-MM-DD"), null,
// or a legacy struct object ({"Year":..., "Month":..., "Day":...}).
func (d *Date) UnmarshalJSON(b []byte) error {
	trimmed := strings.TrimSpace(string(b))
	if trimmed == "null" || trimmed == `""` || trimmed == "" {
		*d = Date{}
		return nil
	}

	// Try string unmarshal
	if strings.HasPrefix(trimmed, `"`) && strings.HasSuffix(trimmed, `"`) {
		s := strings.Trim(trimmed, `"`)
		parsed, err := ParseDate(s)
		if err != nil {
			return err
		}
		*d = parsed
		return nil
	}

	// Fallback for object format {Year, Month, Day} or {year, month, day}
	var aux struct {
		Year       int        `json:"Year"`
		Month      time.Month `json:"Month"`
		Day        int        `json:"Day"`
		YearLower  int        `json:"year"`
		MonthLower time.Month `json:"month"`
		DayLower   int        `json:"day"`
	}
	if err := json.Unmarshal(b, &aux); err == nil {
		y := aux.Year
		if y == 0 {
			y = aux.YearLower
		}
		m := aux.Month
		if m == 0 {
			m = aux.MonthLower
		}
		day := aux.Day
		if day == 0 {
			day = aux.DayLower
		}
		if y != 0 && m != 0 && day != 0 {
			*d = Date{Year: y, Month: m, Day: day}
			return nil
		}
	}

	return &Error{Code: CodeValidation, Message: "invalid date format", Fields: map[string]string{"date": "must be YYYY-MM-DD"}}
}

// MarshalText implements encoding.TextMarshaler.
func (d Date) MarshalText() ([]byte, error) {
	if d.IsZero() {
		return []byte(""), nil
	}
	return []byte(d.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (d *Date) UnmarshalText(text []byte) error {
	s := strings.TrimSpace(string(text))
	if s == "" || s == "null" {
		*d = Date{}
		return nil
	}
	parsed, err := ParseDate(s)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

