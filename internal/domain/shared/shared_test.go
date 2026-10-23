package shared

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestDivRound(t *testing.T) {
	cases := []struct {
		a, b int64
		mode RoundingMode
		want int64
	}{
		{15, 10, RoundHalfUp, 2},
		{14, 10, RoundHalfUp, 1},
		{-15, 10, RoundHalfUp, -2},
		{25, 10, RoundHalfEven, 2},
		{35, 10, RoundHalfEven, 4},
		{19, 10, RoundDown, 1},
		{11, 10, RoundUp, 2},
		{-11, 10, RoundUp, -2},
		{20, 10, RoundUp, 2},
	}
	for _, c := range cases {
		if got := DivRound(c.a, c.b, c.mode); got != c.want {
			t.Errorf("DivRound(%d,%d,%s)=%d want %d", c.a, c.b, c.mode, got, c.want)
		}
	}
}

func TestMoneyArithmetic(t *testing.T) {
	a := NewMoney(350000, "PKR")
	b := NewMoney(12345, "PKR")
	sum, err := a.Add(b)
	if err != nil || sum.Minor != 362345 {
		t.Fatalf("Add = %v, %v", sum, err)
	}
	diff, err := a.Sub(b)
	if err != nil || diff.Minor != 337655 {
		t.Fatalf("Sub = %v, %v", diff, err)
	}
	if _, err := a.Add(NewMoney(1, "USD")); !errors.Is(err, ErrInvariant) {
		t.Fatalf("expected currency mismatch invariant, got %v", err)
	}
	if _, err := NewMoney(1<<62, "PKR").Add(NewMoney(1<<62, "PKR")); err == nil {
		t.Fatal("expected overflow")
	}
	if _, err := NewMoney(1<<62, "PKR").MulInt(4); err == nil {
		t.Fatal("expected overflow on MulInt")
	}
}

func TestMoneyPercentAndRounding(t *testing.T) {
	// 10% of 3,333.33 = 333.333 -> 333.33
	p, err := NewMoney(333333, "PKR").PercentBps(1000, RoundHalfUp)
	if err != nil || p.Minor != 33333 {
		t.Fatalf("PercentBps = %v, %v", p, err)
	}
	// 12.5% of 0.04 = 0.005 -> half up 0.01
	p, _ = NewMoney(4, "PKR").PercentBps(1250, RoundHalfUp)
	if p.Minor != 1 {
		t.Fatalf("half up expected 1 got %d", p.Minor)
	}
	if r := NewMoney(345050, "PKR").RoundTo(100, RoundHalfUp); r.Minor != 345100 {
		t.Fatalf("RoundTo = %d", r.Minor)
	}
}

func TestMoneyFormat(t *testing.T) {
	cases := map[int64]string{0: "0.00", 5: "0.05", 350000: "3,500.00", -123456789: "-1,234,567.89"}
	for in, want := range cases {
		if got := NewMoney(in, "PKR").Format(2); got != want {
			t.Errorf("Format(%d) = %q want %q", in, got, want)
		}
	}
	if got := NewMoney(1234567, "JPY").Format(0); got != "1,234,567" {
		t.Errorf("Format zero digits = %q", got)
	}
}

func TestDates(t *testing.T) {
	d, err := ParseDate("2026-02-28")
	if err != nil {
		t.Fatal(err)
	}
	if d.AddDays(1).String() != "2026-03-01" {
		t.Fatalf("AddDays: %s", d.AddDays(1))
	}
	if _, err := ParseDate("2026-02-30"); err == nil {
		t.Fatal("expected invalid date")
	}
	karachi := time.FixedZone("PKT", 5*3600)
	// 20:00 UTC on Oct 20 is already Oct 21 in Karachi.
	now := time.Date(2026, 10, 20, 20, 0, 0, 0, time.UTC)
	if got := TodayIn(now, karachi).String(); got != "2026-10-21" {
		t.Fatalf("TodayIn = %s", got)
	}
}

func TestDateJSON(t *testing.T) {
	// 1. Unmarshal from string "2026-10-23"
	var d1 Date
	if err := json.Unmarshal([]byte(`"2026-10-23"`), &d1); err != nil {
		t.Fatalf("Unmarshal string: %v", err)
	}
	if d1.Year != 2026 || d1.Month != 10 || d1.Day != 23 {
		t.Fatalf("Expected 2026-10-23, got %s", d1)
	}

	// 2. Marshal to JSON string
	b, err := json.Marshal(d1)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(b) != `"2026-10-23"` {
		t.Fatalf("Expected %q, got %s", `"2026-10-23"`, string(b))
	}

	// 3. Unmarshal from object {"Year": 2026, "Month": 10, "Day": 23}
	var d2 Date
	if err := json.Unmarshal([]byte(`{"Year":2026,"Month":10,"Day":23}`), &d2); err != nil {
		t.Fatalf("Unmarshal legacy object: %v", err)
	}
	if d2.Year != 2026 || d2.Month != 10 || d2.Day != 23 {
		t.Fatalf("Expected 2026-10-23, got %s", d2)
	}

	// 4. Unmarshal null
	var d3 Date
	if err := json.Unmarshal([]byte(`null`), &d3); err != nil {
		t.Fatalf("Unmarshal null: %v", err)
	}
	if !d3.IsZero() {
		t.Fatalf("Expected zero date, got %s", d3)
	}

	// 5. Marshal zero date
	bZero, err := json.Marshal(d3)
	if err != nil {
		t.Fatalf("Marshal zero: %v", err)
	}
	if string(bZero) != "null" {
		t.Fatalf("Expected null, got %s", string(bZero))
	}

	// 6. Unmarshal struct with Date fields from client JSON
	type BatchInput struct {
		AcademicYearID string `json:"academicYearId"`
		ClassID        string `json:"classId"`
		BillingPeriod  string `json:"billingPeriod"`
		IssueDate      Date   `json:"issueDate"`
		DueDate        Date   `json:"dueDate"`
	}
	rawClientJSON := `{"academicYearId":"01a15d6e-0792-756b-b363-6638341721e8","classId":"01a15d6e-b63c-7396-b9a6-2d1aeb48dc59","billingPeriod":"2026-10","issueDate":"2026-10-23","dueDate":"2026-11-07"}`
	var batch BatchInput
	if err := json.Unmarshal([]byte(rawClientJSON), &batch); err != nil {
		t.Fatalf("Unmarshal client JSON: %v", err)
	}
	if batch.IssueDate.String() != "2026-10-23" || batch.DueDate.String() != "2026-11-07" {
		t.Fatalf("Dates mismatch: issue=%s due=%s", batch.IssueDate, batch.DueDate)
	}
}

func TestValidatorAndErrors(t *testing.T) {
	var v Validator
	v.Check(false, "name", "required")
	v.Check(false, "name", "second message ignored")
	v.Check(true, "email", "ok")
	err := v.Err()
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
	var e *Error
	if !errors.As(err, &e) || e.Fields["name"] != "required" || len(e.Fields) != 1 {
		t.Fatalf("unexpected fields %#v", e)
	}
	if CodeOf(errors.New("x")) != CodeInternal {
		t.Fatal("untyped errors must map to internal")
	}
}

