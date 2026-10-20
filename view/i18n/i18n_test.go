package i18n

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultIsEnglish(t *testing.T) {
	require.NotNil(t, Default())
	assert.Equal(t, "en", Default().Tag)
	assert.Same(t, English, Default())
}

func TestGetBuiltInLocales(t *testing.T) {
	en, ok := Get("en")
	require.True(t, ok, "en must be registered")
	assert.Same(t, English, en)

	ru, ok := Get("ru")
	require.True(t, ok, "ru must be registered")
	assert.Same(t, Russian, ru)
	assert.Equal(t, "Русский", ru.Name)
}

func TestGetUnknownLocaleReportsFalse(t *testing.T) {
	for _, tag := range []string{"xx", "", "EN", "en-US"} {
		l, ok := Get(tag)
		assert.False(t, ok, "tag %q must not resolve", tag)
		assert.Nil(t, l, "tag %q must not return a locale", tag)
	}
}

func TestRegisterRejectsEmptyTag(t *testing.T) {
	before, _ := Get("zz")
	require.Nil(t, before)

	err := Register(&Locale{Name: "No tag"})
	require.Error(t, err)
	_, ok := Get("")
	assert.False(t, ok, "an empty tag must never enter the registry")

	require.Error(t, Register(nil), "a nil locale must be rejected")
}

func TestRegisterCustomLocale(t *testing.T) {
	custom := &Locale{
		Tag:  "xx-test",
		Name: "Testish",
		Strings: map[string]string{
			KeyYes:   "Yep",
			"custom": "Custom %s",
		},
		Months:   [12]string{"M1", "M2", "M3", "M4", "M5", "M6", "M7", "M8", "M9", "M10", "M11", "M12"},
		Weekdays: [7]string{"D1", "D2", "D3", "D4", "D5", "D6", "D7"},
	}
	require.NoError(t, Register(custom))

	got, ok := Get("xx-test")
	require.True(t, ok)
	assert.Same(t, custom, got)

	ctx := WithLocale(context.Background(), custom)
	assert.Equal(t, "Yep", T(ctx, KeyYes))
	assert.Equal(t, "Custom value", T(ctx, "custom", "value"))
	// A key the custom locale lacks falls back to English.
	assert.Equal(t, English.Strings[KeyCancel], T(ctx, KeyCancel))
	// Keys the custom locale provides are used for formats too.
	assert.Equal(t, "M3", monthName(custom, time.March))
	assert.Equal(t, "D7", weekdayName(custom, 0)) // Sunday is last in a Monday-first set
}

func TestWithLocaleNilLeavesContextUnchanged(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, ctx, WithLocale(ctx, nil))
	assert.Same(t, Default(), From(WithLocale(ctx, nil)))
}

func TestFromEmptyContextFallsBackToEnglish(t *testing.T) {
	assert.Same(t, English, From(context.Background()))
	assert.Same(t, English, From(nil), "a nil context must not panic")
	// A context carrying an unregistered locale reports that locale; its
	// missing keys still resolve through Default() (see TestTNonEmptyKeysNeverReturnEmptyString).
	carried := &Locale{Tag: "xx"}
	assert.Same(t, carried, From(WithLocale(context.Background(), carried)))
}

func TestFromIgnoresATaglessLocale(t *testing.T) {
	// layout.Page renders <html lang={ From(ctx).Tag }>, so a carried locale
	// with no tag must resolve to Default() rather than produce lang="".
	tagless := &Locale{Name: "No tag", Strings: map[string]string{KeyYes: "Yep"}}
	assert.Same(t, Default(), From(WithLocale(context.Background(), tagless)))
	assert.Equal(t, "Yes", T(WithLocale(context.Background(), tagless), KeyYes))
}

func TestTResolvesPerLocale(t *testing.T) {
	en := context.Background()
	ru := WithLocale(context.Background(), Russian)

	assert.Equal(t, "Yes", T(en, KeyYes))
	assert.Equal(t, "Да", T(ru, KeyYes))
	assert.Equal(t, "No events", T(en, KeyNoEvents))
	assert.Equal(t, "Нет событий", T(ru, KeyNoEvents))
	assert.Equal(t, "3 selected", T(en, KeySelectedCount, 3))
	assert.Equal(t, "Выбрано: 3", T(ru, KeySelectedCount, 3))
	assert.Equal(t, "Page 2 of 7", T(en, KeyPageOf, 2, 7))
	assert.Equal(t, "Страница 2 из 7", T(ru, KeyPageOf, 2, 7))
}

func TestTMissingKeyReturnsTheKey(t *testing.T) {
	assert.Equal(t, "no-such-key", T(context.Background(), "no-such-key"))
	assert.Equal(t, "no-such-key", T(WithLocale(context.Background(), Russian), "no-such-key"))
}

func TestTNonEmptyKeysNeverReturnEmptyString(t *testing.T) {
	sparse := &Locale{Tag: "sparse", Strings: map[string]string{KeyYes: "", "": ""}}
	require.NoError(t, Register(sparse))
	ctx := WithLocale(context.Background(), sparse)

	assert.Equal(t, "Yes", T(ctx, KeyYes), "an empty entry falls back to English")
	assert.Equal(t, "", T(ctx, ""), "an empty key returns the empty key")
	for _, key := range []string{KeyYes, KeyNo, KeyClose, KeyLoading, KeyNoData, KeyPageOf, "unknown"} {
		assert.NotEmpty(t, T(ctx, key), "key %q must not resolve empty", key)
	}
}

func TestRegisteringTheFallbackTagReplacesDefault(t *testing.T) {
	replacement := &Locale{
		Tag:     "en",
		Name:    "Replacement English",
		Strings: map[string]string{KeyYes: "Yeah"},
	}
	require.NoError(t, Register(replacement))
	// Restore the built-in English for the rest of the package's tests.
	t.Cleanup(func() { require.NoError(t, Register(English)) })

	assert.Same(t, replacement, Default(), "Default must follow the locale registered under the fallback tag")
	assert.Equal(t, "Yeah", T(context.Background(), KeyYes))
	// A locale that lacks the key falls through to the replacement Default.
	sparse := &Locale{Tag: "fallback-probe"}
	assert.Equal(t, "Yeah", T(WithLocale(context.Background(), sparse), KeyYes))
	// A key the replacement does not carry has no Default entry either, so it
	// surfaces as the key itself.
	assert.Equal(t, KeyCancel, T(context.Background(), KeyCancel))
}

func TestEveryBuiltInCatalogKeyIsPresent(t *testing.T) {
	keys := []string{
		KeyYes, KeyNo, KeyClose, KeyLoading, KeyToday, KeyAllDay, KeyNoEvents,
		KeyNoData, KeyCancel, KeyConfirm, KeySubmit, KeySearch, KeyTheme,
		KeyCollapse, KeyNew, KeyLogout, KeySelectedCount, KeyPageOf,
		KeyChooseDay, KeyMore, KeyEventsCalendar, KeyToggleSidebar, KeyUserAvatar,
		KeyToggleTheme, KeyVsPreviousPeriod,
		KeyError, KeyWarning, KeySuccess, KeyInfo,
	}
	for _, l := range []*Locale{English, Russian} {
		for _, key := range keys {
			s, ok := l.Lookup(key)
			assert.True(t, ok, "%s: missing key %q", l.Tag, key)
			assert.NotEmpty(t, s, "%s: empty entry for %q", l.Tag, key)
		}
		assert.Len(t, l.Months, 12, "%s: month names", l.Tag)
		assert.Len(t, l.Weekdays, 7, "%s: weekday names", l.Tag)
		assert.NotEmpty(t, l.DateLayout, "%s: date layout", l.Tag)
		assert.NotEmpty(t, l.DateTimeLayout, "%s: date-time layout", l.Tag)
		assert.NotEmpty(t, l.TimeLayout, "%s: time layout", l.Tag)
	}
}
