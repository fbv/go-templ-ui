package i18n

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// fallback layouts are used when a locale omits a layout, so formatting never
// produces an empty string.
const (
	fallbackDateLayout     = "2006-01-02"
	fallbackDateTimeLayout = "2006-01-02 15:04"
	fallbackTimeLayout     = "15:04"
)

const (
	defaultCurrencyDecimals = 2
	percentSign             = "%"
)

// Number formats value with the locale's decimal separator and digit grouping.
// decimals is the number of fraction digits (negative counts as zero); a
// grouped integer part of 1234567.891 with decimals 2 renders as "1,234,567.89"
// in English and "1 234 567,89" in Russian. Negatives keep a leading "-".
func Number(ctx context.Context, value float64, decimals int) string {
	return formatNumber(From(ctx), value, decimals)
}

func formatNumber(l *Locale, value float64, decimals int) string {
	if decimals < 0 {
		decimals = 0
	}
	// NaN and the infinities render verbatim: there are no digits to group and
	// a sign must stay part of the token.
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return strconv.FormatFloat(value, 'f', decimals, 64)
	}
	nf := l.Numbers
	decimal := nf.Decimal
	if decimal == "" {
		decimal = "."
	}

	s := strconv.FormatFloat(value, 'f', decimals, 64)
	sign := ""
	if strings.HasPrefix(s, "-") {
		sign = "-"
		s = s[1:]
	}
	intPart, fracPart, hasFrac := strings.Cut(s, ".")
	intPart = groupDigits(intPart, nf.Group, nf.GroupSize)

	out := sign + intPart
	if hasFrac {
		out += decimal + fracPart
	}
	return out
}

// groupDigits inserts sep every size digits of digits, counting from the right.
// A zero size means three; an empty separator means no grouping.
func groupDigits(digits, sep string, size int) string {
	if size <= 0 {
		size = 3
	}
	if sep == "" || len(digits) <= size {
		return digits
	}
	var b strings.Builder
	first := len(digits) % size
	if first == 0 {
		first = size
	}
	b.WriteString(digits[:first])
	for i := first; i < len(digits); i += size {
		b.WriteString(sep)
		b.WriteString(digits[i : i+size])
	}
	return b.String()
}

// Currency formats amount as a currency in the locale's convention. English
// renders "USD" as "$1,234.50", Russian "RUB" as "1 234,50 ₽". A code the
// locale does not describe falls back to the code itself as a prefix symbol
// ("CHF 1,234.50"). The fraction digit count comes from the code's
// CurrencyFormat.Decimals (default 2), and a negative sign is always the
// outermost character, for finite and non-finite amounts alike: "-$1,234.50",
// "-1 234,50 ₽", "-$Inf".
func Currency(ctx context.Context, amount float64, code string) string {
	l := From(ctx)
	cf, ok := l.Numbers.Currency[code]
	if !ok || cf.Symbol == "" {
		cf = CurrencyFormat{Symbol: code, Position: "before", Space: true}
	}
	decimals := cf.Decimals
	if decimals == 0 {
		decimals = defaultCurrencyDecimals
	} else if decimals < 0 {
		decimals = 0
	}

	num := formatNumber(l, amount, decimals)
	// Take the sign out of the number and put it around the whole currency
	// token, so it is outermost for every value including -Inf. NaN and +Inf
	// carry no sign and stay verbatim inside the symbol.
	negative := strings.HasPrefix(num, "-")
	if negative {
		num = num[1:]
	}
	return joinCurrency(cf, num, negative)
}

func joinCurrency(cf CurrencyFormat, num string, negative bool) string {
	gap := ""
	if cf.Space {
		gap = " "
	}
	out := cf.Symbol + gap + num
	if cf.Position == "after" {
		out = num + gap + cf.Symbol
	}
	if negative {
		out = "-" + out
	}
	return out
}

// Percent formats value as a percentage of the locale's number convention. The
// percent sign follows the number, separated by a space when the locale sets
// NumberFormat.PercentSpace.
func Percent(ctx context.Context, value float64, decimals int) string {
	gap := ""
	if From(ctx).Numbers.PercentSpace {
		gap = " "
	}
	return Number(ctx, value, decimals) + gap + percentSign
}

// Date formats t using the locale's date layout.
func Date(ctx context.Context, t time.Time) string {
	return t.Format(orDefault(From(ctx).DateLayout, fallbackDateLayout))
}

// DateTime formats t using the locale's date-and-time layout.
func DateTime(ctx context.Context, t time.Time) string {
	return t.Format(orDefault(From(ctx).DateTimeLayout, fallbackDateTimeLayout))
}

// Time formats the clock time of t using the locale's time layout.
func Time(ctx context.Context, t time.Time) string {
	return t.Format(orDefault(From(ctx).TimeLayout, fallbackTimeLayout))
}

func orDefault(layout, fallback string) string {
	if layout == "" {
		return fallback
	}
	return layout
}

// MonthName returns the locale's full name for month (January first). A name
// the locale leaves empty falls back to Default()'s name, so a partial locale
// never renders a blank. An out-of-range month returns the empty string rather
// than panicking.
func MonthName(ctx context.Context, month time.Month) string {
	return monthName(From(ctx), month)
}

func monthName(l *Locale, month time.Month) string {
	idx := int(month) - 1
	if idx < 0 || idx >= len(l.Months) {
		return ""
	}
	if s := l.Months[idx]; s != "" {
		return s
	}
	if d := Default(); d != nil && d != l {
		return d.Months[idx]
	}
	return ""
}

// WeekdayName returns the locale's short weekday name for day. Weekday arrays
// are Monday first, so a Sunday reports index 6. A name the locale leaves empty
// falls back to Default()'s name.
func WeekdayName(ctx context.Context, day time.Weekday) string {
	return weekdayName(From(ctx), day)
}

func weekdayName(l *Locale, day time.Weekday) string {
	idx := (int(day) + 6) % 7 // time.Sunday=0 -> Monday-first index 6
	if s := l.Weekdays[idx]; s != "" {
		return s
	}
	if d := Default(); d != nil && d != l {
		return d.Weekdays[idx]
	}
	return ""
}

// MonthYear formats t as the locale's month and year: MonthYearLayout when the
// locale sets one, otherwise the month name followed by the year, e.g.
// "March 2026" in English and "Март 2026" in Russian. Go reference-time layouts
// expand month names in English, so a locale that wants its own month names
// should leave MonthYearLayout empty. The year alone is used when the locale
// cannot supply a month name, so the result never begins with a stray space.
func MonthYear(ctx context.Context, t time.Time) string {
	l := From(ctx)
	if l.MonthYearLayout != "" {
		return t.Format(l.MonthYearLayout)
	}
	name := monthName(l, t.Month())
	if name == "" {
		return strconv.Itoa(t.Year())
	}
	return fmt.Sprintf("%s %d", name, t.Year())
}
