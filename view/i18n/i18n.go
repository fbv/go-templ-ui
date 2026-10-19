// Package i18n carries the locale used to render user-facing strings and
// formatted values.
//
// A Locale is a plain value: a tag, a string catalog, month/weekday names and
// number/date formats. Locales live in a package registry and travel in a
// context.Context — the one channel every templ component already receives — so
// a consumer selects a locale once and every component below it renders in that
// locale without gaining a locale prop:
//
//	ctx = i18n.WithLocale(ctx, i18n.Russian)
//	ui.YesNo(true)                          // Да
//	i18n.T(ctx, i18n.KeyNoEvents)           // Нет событий
//	i18n.Currency(ctx, 1234.5, "RUB")       // 1 234,50 ₽
//
// English is the built-in Default(); Russian ships as a second built-in locale
// and as the template for a consumer's own. Resolution never panics: a missing
// catalog entry falls back to Default()'s entry, and a non-empty key neither
// catalog knows resolves to the key itself — a visible gap rather than a silent
// blank.
package i18n

import (
	"context"
	"fmt"
	"sync"
)

// Catalog keys for the starter strings every component shares. A locale's
// catalog is a map[string]string keyed by these constants (custom keys are
// allowed for consumer strings).
const (
	KeyYes           = "yes"
	KeyNo            = "no"
	KeyClose         = "close"
	KeyLoading       = "loading"
	KeyToday         = "today"
	KeyAllDay        = "all-day"
	KeyNoEvents      = "no-events"
	KeyNoData        = "no-data"
	KeyCancel        = "cancel"
	KeyConfirm       = "confirm"
	KeySubmit        = "submit"
	KeySearch        = "search"
	KeyTheme         = "theme"
	KeyCollapse      = "collapse"
	KeyNew           = "new"
	KeyLogout        = "logout"
	KeySelectedCount = "selected-count"
	KeyPageOf        = "page-of"
)

// CurrencyFormat describes how one currency code is written in a locale.
type CurrencyFormat struct {
	// Symbol is what replaces the code in output, e.g. "$", "€", "₽".
	Symbol string
	// Position is "before" (prefix) or "after" (suffix). Anything else is
	// treated as "before".
	Position string
	// Space puts a space between the symbol and the number.
	Space bool
	// Decimals is the number of fraction digits. Zero means the default of 2;
	// a negative value counts as zero.
	Decimals int
}

// NumberFormat describes a locale's number conventions.
type NumberFormat struct {
	// Decimal separates the integer and fraction parts, e.g. "." or ",".
	Decimal string
	// Group separates digit groups in the integer part, e.g. "," or " ".
	Group string
	// GroupSize is the number of digits per group. Zero means 3.
	GroupSize int
	// Currency maps an ISO code to its local presentation. A code that is
	// absent falls back to the code used as a symbol.
	Currency map[string]CurrencyFormat
	// PercentSpace puts a space between the number and the percent sign.
	PercentSpace bool
}

// Locale is an immutable, shareable locale value. Once registered it may be
// used from any goroutine.
type Locale struct {
	// Tag is the BCP 47 language tag, e.g. "en" or "ru". It is the registry
	// key and the value components expose as the document language.
	Tag string
	// Name is a human-readable name, e.g. "English".
	Name string

	// Strings is the message catalog, keyed by the Key* constants (or any
	// consumer-defined key). Entries may be format strings for T's args.
	Strings map[string]string

	// Months holds the twelve full month names, January first.
	Months [12]string
	// Weekdays holds seven short weekday names, Monday first.
	Weekdays [7]string

	// Numbers carries the locale's number and currency conventions.
	Numbers NumberFormat

	// DateLayout, DateTimeLayout and TimeLayout are Go reference-time
	// layouts, e.g. "01/02/2006" or "02.01.2006".
	DateLayout     string
	DateTimeLayout string
	TimeLayout     string
	// MonthYearLayout is the reference-time layout MonthYear uses. When it is
	// empty, MonthYear renders the month name followed by the year.
	MonthYearLayout string
}

// Lookup returns the locale's entry for key, or false when the catalog has no
// (or an empty) entry.
func (l *Locale) Lookup(key string) (string, bool) {
	if l == nil {
		return "", false
	}
	s, ok := l.Strings[key]
	if !ok || s == "" {
		return "", false
	}
	return s, true
}

// fallbackTag is the tag Default resolves through. A locale registered under
// it replaces the built-in English everywhere, including T's last fallback.
const fallbackTag = "en"

var (
	registryMu sync.RWMutex
	registry   = map[string]*Locale{}
)

// Register adds a locale to the package registry, keyed by its Tag. It may be
// called from a consumer package's init. Re-registering an existing tag
// replaces it, so a consumer can override a built-in locale.
//
// Register returns an error when l is nil or l.Tag is empty; the registry is
// left untouched in that case.
func Register(l *Locale) error {
	if l == nil {
		return fmt.Errorf("i18n: cannot register a nil locale")
	}
	if l.Tag == "" {
		return fmt.Errorf("i18n: cannot register a locale with an empty Tag")
	}
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[l.Tag] = l
	return nil
}

func mustRegister(l *Locale) {
	if err := Register(l); err != nil {
		panic(err)
	}
}

// Get returns the registered locale for tag. An unknown or empty tag reports
// false; callers that need a usable locale should use From.
func Get(tag string) (*Locale, bool) {
	if tag == "" {
		return nil, false
	}
	registryMu.RLock()
	defer registryMu.RUnlock()
	l, ok := registry[tag]
	return l, ok
}

// Default returns the locale registered under the fallback tag ("en"), the
// fallback for every lookup. Registering a locale under that tag replaces the
// built-in English here too. It is never nil once the built-in locales are
// registered.
func Default() *Locale {
	registryMu.RLock()
	defer registryMu.RUnlock()
	if l, ok := registry[fallbackTag]; ok && l != nil {
		return l
	}
	return English
}

type localeKeyType struct{}

var localeKey = localeKeyType{}

// WithLocale returns a context carrying l. A nil locale leaves ctx unchanged.
func WithLocale(ctx context.Context, l *Locale) context.Context {
	if l == nil {
		return ctx
	}
	return context.WithValue(ctx, localeKey, l)
}

// From returns the locale carried by ctx. A context with no locale — including
// a nil context — yields Default(), so every call site has a usable locale. A
// carried locale with an empty Tag is also ignored: components expose the
// locale's Tag as the document language, and it must never be empty.
func From(ctx context.Context) *Locale {
	if ctx != nil {
		if l, ok := ctx.Value(localeKey).(*Locale); ok && l != nil && l.Tag != "" {
			return l
		}
	}
	return Default()
}

// T resolves key through the locale carried by ctx, formats it with args when
// any are given, and never returns an empty string for a non-empty key: a key
// missing from the selected locale falls back to Default()'s entry, and a key
// neither knows resolves to the key itself. An empty key returns the empty key.
func T(ctx context.Context, key string, args ...any) string {
	if key == "" {
		return ""
	}
	if s, ok := From(ctx).Lookup(key); ok {
		if len(args) == 0 {
			return s
		}
		return fmt.Sprintf(s, args...)
	}
	if s, ok := Default().Lookup(key); ok {
		if len(args) == 0 {
			return s
		}
		return fmt.Sprintf(s, args...)
	}
	return key
}
