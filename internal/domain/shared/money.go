package shared

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Money is an exact monetary amount in minor currency units (e.g. paisa, cents).
// Floating point is never used for currency.
type Money struct {
	Minor    int64
	Currency string
}

// RoundingMode controls how fractional minor units are resolved.
type RoundingMode string

const (
	RoundHalfUp   RoundingMode = "half_up"
	RoundHalfEven RoundingMode = "half_even"
	RoundDown     RoundingMode = "down"
	RoundUp       RoundingMode = "up"
)

// BasisPointsWhole is 100% expressed in basis points.
const BasisPointsWhole int64 = 10_000

var errOverflow = Invariant("monetary overflow")

func NewMoney(minor int64, currency string) Money {
	return Money{Minor: minor, Currency: currency}
}

func (m Money) IsZero() bool     { return m.Minor == 0 }
func (m Money) IsNegative() bool { return m.Minor < 0 }

func (m Money) sameCurrency(o Money) error {
	if m.Currency != o.Currency {
		return Invariant(fmt.Sprintf("currency mismatch: %s vs %s", m.Currency, o.Currency))
	}
	return nil
}

// Add returns m + o, failing on currency mismatch or overflow.
func (m Money) Add(o Money) (Money, error) {
	if err := m.sameCurrency(o); err != nil {
		return Money{}, err
	}
	r := m.Minor + o.Minor
	if (o.Minor > 0 && r < m.Minor) || (o.Minor < 0 && r > m.Minor) {
		return Money{}, errOverflow
	}
	return Money{Minor: r, Currency: m.Currency}, nil
}

// Sub returns m - o, failing on currency mismatch or overflow.
func (m Money) Sub(o Money) (Money, error) {
	if o.Minor == math.MinInt64 {
		return Money{}, errOverflow
	}
	return m.Add(Money{Minor: -o.Minor, Currency: o.Currency})
}

// MulInt returns m * n, failing on overflow.
func (m Money) MulInt(n int64) (Money, error) {
	if m.Minor == 0 || n == 0 {
		return Money{Currency: m.Currency}, nil
	}
	r := m.Minor * n
	if r/n != m.Minor {
		return Money{}, errOverflow
	}
	return Money{Minor: r, Currency: m.Currency}, nil
}

// PercentBps returns m * bps / 10000 rounded with mode. Uses integer math only.
func (m Money) PercentBps(bps int64, mode RoundingMode) (Money, error) {
	num := m.Minor * bps
	if bps != 0 && num/bps != m.Minor {
		return Money{}, errOverflow
	}
	return Money{Minor: DivRound(num, BasisPointsWhole, mode), Currency: m.Currency}, nil
}

// RoundTo rounds m to a multiple of unit minor units (e.g. 100 = whole rupees).
func (m Money) RoundTo(unit int64, mode RoundingMode) Money {
	if unit <= 1 {
		return m
	}
	return Money{Minor: DivRound(m.Minor, unit, mode) * unit, Currency: m.Currency}
}

// DivRound divides a by b (b > 0) applying the rounding mode symmetrically
// around zero.
func DivRound(a, b int64, mode RoundingMode) int64 {
	if b <= 0 {
		panic("DivRound: non-positive divisor")
	}
	neg := a < 0
	if neg {
		a = -a
	}
	q, r := a/b, a%b
	if r != 0 {
		switch mode {
		case RoundUp:
			q++
		case RoundDown:
		case RoundHalfEven:
			if r*2 > b || (r*2 == b && q%2 == 1) {
				q++
			}
		default: // half up
			if r*2 >= b {
				q++
			}
		}
	}
	if neg {
		return -q
	}
	return q
}

// Format renders the amount with the given number of minor digits, e.g.
// 350000 with 2 digits -> "3,500.00". Currency symbols are a presentation concern.
func (m Money) Format(minorDigits int) string {
	neg := m.Minor < 0
	v := m.Minor
	if neg {
		v = -v
	}
	s := strconv.FormatInt(v, 10)
	var whole, frac string
	if minorDigits > 0 {
		for len(s) <= minorDigits {
			s = "0" + s
		}
		whole, frac = s[:len(s)-minorDigits], s[len(s)-minorDigits:]
	} else {
		whole = s
	}
	var b strings.Builder
	if neg {
		b.WriteByte('-')
	}
	for i, c := range whole {
		if i > 0 && (len(whole)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(c)
	}
	if frac != "" {
		b.WriteByte('.')
		b.WriteString(frac)
	}
	return b.String()
}

