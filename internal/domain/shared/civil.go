package shared

import (
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

