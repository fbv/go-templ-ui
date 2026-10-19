package i18n

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	enCtx = context.Background()
	ruCtx = WithLocale(context.Background(), Russian)
)

func TestNumberMatrix(t *testing.T) {
	cases := []struct {
		name     string
		ctx      context.Context
		value    float64
		decimals int
		want     string
	}{
		{"en grouped", enCtx, 1234567.891, 2, "1,234,567.89"},
		{"ru grouped", ruCtx, 1234567.891, 2, "1 234 567,89"},
		{"en no decimals", enCtx, 1234567.891, 0, "1,234,568"},
		{"ru no decimals", ruCtx, 1234567.891, 0, "1 234 568"},
		{"en negative", enCtx, -1234.5, 2, "-1,234.50"},
		{"ru negative", ruCtx, -1234.5, 2, "-1 234,50"},
		{"en below group size", enCtx, 999, 0, "999"},
		{"en exact group", enCtx, 1000, 0, "1,000"},
		{"ru zero", ruCtx, 0, 2, "0,00"},
		{"negative decimals clamp", enCtx, 12.6, -1, "13"},
		{"NaN verbatim", enCtx, math.NaN(), 2, "NaN"},
		{"+Inf verbatim", enCtx, math.Inf(1), 2, "+Inf"},
		{"+Inf in ru, ungrouped", ruCtx, math.Inf(1), 2, "+Inf"},
		{"-Inf verbatim", enCtx, math.Inf(-1), 2, "-Inf"},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, Number(tc.ctx, tc.value, tc.decimals), tc.name)
	}
}

func TestCurrencyMatrix(t *testing.T) {
	cases := []struct {
		name   string
		ctx    context.Context
		amount float64
		code   string
		want   string
	}{
		{"en USD", enCtx, 1234.5, "USD", "$1,234.50"},
		{"ru RUB", ruCtx, 1234.5, "RUB", "1 234,50 ₽"},
		{"en EUR", enCtx, 99, "EUR", "€99.00"},
		{"ru EUR", ruCtx, 99, "EUR", "€99,00"},
		{"en unknown code", enCtx, 1234.5, "CHF", "CHF 1,234.50"},
		{"ru unknown code", ruCtx, 1234.5, "CHF", "CHF 1 234,50"},
		{"en negative sign outermost", enCtx, -12.345, "USD", "-$12.35"},
		{"ru negative sign outermost", ruCtx, -12.345, "RUB", "-12,35 ₽"},
		{"en -Inf sign outermost", enCtx, math.Inf(-1), "USD", "-$Inf"},
		{"ru -Inf sign outermost", ruCtx, math.Inf(-1), "RUB", "-Inf ₽"},
		{"en +Inf verbatim inside symbol", enCtx, math.Inf(1), "USD", "$+Inf"},
		{"en NaN verbatim inside symbol", enCtx, math.NaN(), "USD", "$NaN"},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, Currency(tc.ctx, tc.amount, tc.code), tc.name)
	}
}

func TestCurrencyHonoursPerLocaleDecimals(t *testing.T) {
	three := &Locale{
		Tag: "decimals-test",
		Numbers: NumberFormat{
			Decimal: ".",
			Group:   ",",
			Currency: map[string]CurrencyFormat{
				"USD": {Symbol: "$", Position: "before", Decimals: 3},
				"JPY": {Symbol: "¥", Position: "before", Decimals: -1},
			},
		},
	}
	require.NoError(t, Register(three))
	ctx := WithLocale(context.Background(), three)

	assert.Equal(t, "$12.346", Currency(ctx, 12.3456, "USD"), "CurrencyFormat.Decimals drives the fraction digits")
	assert.Equal(t, "¥12", Currency(ctx, 12.3456, "JPY"), "a negative Decimals counts as zero")
	// Currencies absent from the table keep the default of two digits.
	assert.Equal(t, "EUR 12.35", Currency(ctx, 12.3456, "EUR"))
}

func TestPercentMatrix(t *testing.T) {
	assert.Equal(t, "12.5%", Percent(enCtx, 12.5, 1))
	assert.Equal(t, "12,5%", Percent(ruCtx, 12.5, 1))
}

func TestPercentSpaceIsLocaleDriven(t *testing.T) {
	spaced := &Locale{
		Tag: "percent-test",
		Numbers: NumberFormat{
			Decimal:      ".",
			PercentSpace: true,
		},
	}
	require.NoError(t, Register(spaced))
	ctx := WithLocale(context.Background(), spaced)

	assert.Equal(t, "12.5 %", Percent(ctx, 12.5, 1))
	assert.Equal(t, "12.5%", Percent(enCtx, 12.5, 1), "English keeps no space")
}

func TestDateMatrix(t *testing.T) {
	d := time.Date(2026, time.March, 9, 14, 5, 6, 0, time.UTC)
	assert.Equal(t, "03/09/2026", Date(enCtx, d))
	assert.Equal(t, "09.03.2026", Date(ruCtx, d))
}

func TestDateTimeAndTimeMatrix(t *testing.T) {
	d := time.Date(2026, time.March, 9, 14, 5, 6, 0, time.UTC)
	assert.Equal(t, "03/09/2026 14:05", DateTime(enCtx, d))
	assert.Equal(t, "09.03.2026 14:05", DateTime(ruCtx, d))
	assert.Equal(t, "14:05", Time(enCtx, d))
	assert.Equal(t, "14:05", Time(ruCtx, d))
}

func TestZeroTimeFormatsAsDate(t *testing.T) {
	var zero time.Time
	assert.Equal(t, "01/01/0001", Date(enCtx, zero))
	assert.Equal(t, "01.01.0001", Date(ruCtx, zero))
	assert.NotPanics(t, func() {
		_ = DateTime(enCtx, zero)
		_ = Time(enCtx, zero)
		_ = MonthYear(enCtx, zero)
	})
}

func TestMonthYearMatrix(t *testing.T) {
	d := time.Date(2026, time.March, 9, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, "March 2026", MonthYear(enCtx, d))
	assert.Equal(t, "Март 2026", MonthYear(ruCtx, d))
}

func TestMonthYearLayoutIsLocaleDriven(t *testing.T) {
	reordered := &Locale{
		Tag:             "monthyear-test",
		MonthYearLayout: "2006, January",
	}
	require.NoError(t, Register(reordered))
	ctx := WithLocale(context.Background(), reordered)

	d := time.Date(2026, time.March, 9, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, "2026, March", MonthYear(ctx, d))
}

func TestMonthNameAndWeekdayName(t *testing.T) {
	assert.Equal(t, "January", MonthName(enCtx, time.January))
	assert.Equal(t, "Январь", MonthName(ruCtx, time.January))
	assert.Equal(t, "December", MonthName(enCtx, time.December))
	assert.Equal(t, "Декабрь", MonthName(ruCtx, time.December))

	// Monday first in both locales.
	assert.Equal(t, "Mon", WeekdayName(enCtx, time.Monday))
	assert.Equal(t, "Пн", WeekdayName(ruCtx, time.Monday))
	assert.Equal(t, "Sun", WeekdayName(enCtx, time.Sunday))
	assert.Equal(t, "Вс", WeekdayName(ruCtx, time.Sunday))
}

func TestOutOfRangeValuesDoNotPanic(t *testing.T) {
	assert.NotPanics(t, func() {
		assert.Equal(t, "", MonthName(enCtx, time.Month(0)))
		assert.Equal(t, "", MonthName(enCtx, time.Month(13)))
	})
}

func TestFormatsFallBackWhenLocaleOmitsThem(t *testing.T) {
	bare := &Locale{Tag: "bare-test"}
	ctx := WithLocale(context.Background(), bare)

	// No NumberFormat: "." separates the fraction and nothing groups digits.
	assert.Equal(t, "1234.50", Number(ctx, 1234.5, 2))
	// No currency table: the code stands in as the symbol.
	assert.Equal(t, "USD 1234.50", Currency(ctx, 1234.5, "USD"))
	// No layouts: stable fallbacks rather than an empty string.
	d := time.Date(2026, time.March, 9, 14, 5, 0, 0, time.UTC)
	assert.Equal(t, "2026-03-09", Date(ctx, d))
	assert.Equal(t, "2026-03-09 14:05", DateTime(ctx, d))
	assert.Equal(t, "14:05", Time(ctx, d))
}

func TestMissingNameArraysFallBackToDefault(t *testing.T) {
	bare := &Locale{Tag: "bare-names-test"}
	ctx := WithLocale(context.Background(), bare)

	// A locale that omits Months/Weekdays borrows Default()'s names; it never
	// renders an empty name or a month-year with a stray leading space.
	assert.Equal(t, "March", MonthName(ctx, time.March))
	assert.Equal(t, "Mon", WeekdayName(ctx, time.Monday))
	assert.Equal(t, "March 2026", MonthYear(ctx, time.Date(2026, time.March, 9, 0, 0, 0, 0, time.UTC)))

	partial := &Locale{Tag: "partial-names-test"}
	partial.Months[2] = "März" // March only
	partial.Weekdays[0] = "Mo" // Monday only
	ctx = WithLocale(context.Background(), partial)
	assert.Equal(t, "März", MonthName(ctx, time.March), "a present name wins")
	assert.Equal(t, "January", MonthName(ctx, time.January), "an empty entry falls back")
	assert.Equal(t, "Mo", WeekdayName(ctx, time.Monday))
	assert.Equal(t, "Tue", WeekdayName(ctx, time.Tuesday))
}

func TestGroupDigits(t *testing.T) {
	assert.Equal(t, "1", groupDigits("1", ",", 3))
	assert.Equal(t, "12", groupDigits("12", ",", 3))
	assert.Equal(t, "123", groupDigits("123", ",", 3))
	assert.Equal(t, "1,234", groupDigits("1234", ",", 3))
	assert.Equal(t, "1,234,567", groupDigits("1234567", ",", 3))
	assert.Equal(t, "1234567", groupDigits("1234567", "", 3))
	assert.Equal(t, "12 34 56", groupDigits("123456", " ", 2))
	assert.Equal(t, "1,234,560", groupDigits("1234560", ",", 0), "zero group size means three")
}
