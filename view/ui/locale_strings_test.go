package ui

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/fbv/go-templ-ui/view/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These assertions pin the locale surface of the migrated components: every
// user-facing string — including the ones the allowlist guard cannot see (Go
// helpers, {{ }} statements, attribute builders) — must come from the catalog
// and the format helpers. Each string is checked in both built-in locales.

func renderIn(t *testing.T, l *i18n.Locale, c templ.Component) string {
	t.Helper()
	ctx := context.Background()
	if l != nil {
		ctx = i18n.WithLocale(ctx, l)
	}
	var buf strings.Builder
	require.NoError(t, c.Render(ctx, &buf))
	return buf.String()
}

func TestSidebarStringsFollowTheLocale(t *testing.T) {
	p := &SidebarProps{
		Items:       []*SidebarItem{{Name: "Home", URL: "/"}},
		User:        &UserInfo{Name: "Anna", Email: "anna@example.com", Avatar: "/a.png"},
		LogoutURL:   "/logout",
		Collapsible: true,
		ID:          "side",
	}

	ru := []string{"Выйти", "Свернуть", "Переключить боковую панель", "Аватар пользователя"}
	en := []string{"Log out", "Collapse", "Toggle sidebar", "User avatar"}

	htmlEn := renderIn(t, nil, Sidebar(p))
	htmlRu := renderIn(t, i18n.Russian, Sidebar(p))
	for i, s := range en {
		assert.Contains(t, htmlEn, s, "English sidebar must render %q", s)
		assert.NotContains(t, htmlEn, ru[i], "English sidebar must not render %q", ru[i])
	}
	for i, s := range ru {
		assert.Contains(t, htmlRu, s, "Russian sidebar must render %q", s)
		assert.NotContains(t, htmlRu, en[i], "Russian sidebar must not render %q", en[i])
	}

	htmlEn = renderIn(t, nil, SidebarNew())
	htmlRu = renderIn(t, i18n.Russian, SidebarNew())
	assert.Contains(t, htmlEn, ">New<")
	assert.Contains(t, htmlRu, ">Новый<")
	assert.NotContains(t, htmlEn, "Новый")
}

func TestCalendarMonthStringsAndFormatsFollowTheLocale(t *testing.T) {
	day := time.Date(2026, time.March, 9, 9, 0, 0, 0, time.UTC)
	p := &CalendarMonthProps{
		ID:      "cal",
		Year:    2026,
		Month:   time.March,
		LoadURL: "/showcase/i18n/calendar/month",
		Events: []CalendarEvent{
			{ID: "1", Title: "One", Start: day, End: day.Add(time.Hour), Color: "blue"},
			{ID: "2", Title: "Two", Start: day, End: day.Add(time.Hour), Color: "green"},
			{ID: "3", Title: "Three", Start: day, End: day.Add(time.Hour), Color: "red"},
			{ID: "4", Title: "Four", Start: day, End: day.Add(time.Hour), Color: "gray"},
		},
	}
	htmlEn := renderIn(t, nil, CalendarMonth(p))
	htmlRu := renderIn(t, i18n.Russian, CalendarMonth(p))

	// Month title and weekday headers come from calendar.go's helpers and must
	// resolve through i18n.MonthYear / i18n.WeekdayName.
	assert.Contains(t, htmlEn, "March 2026", "English month title")
	assert.NotContains(t, htmlEn, "Март 2026")
	assert.Contains(t, htmlRu, "Март 2026", "Russian month title")
	assert.NotContains(t, htmlRu, "March 2026")

	assert.Contains(t, htmlEn, ">Mon<", "English weekday header")
	assert.NotContains(t, htmlEn, ">Пн<")
	assert.Contains(t, htmlRu, ">Пн<", "Russian weekday header")
	assert.NotContains(t, htmlRu, ">Mon<")

	// Button and panel labels.
	assert.Contains(t, htmlEn, ">Today<")
	assert.Contains(t, htmlRu, ">Сегодня<")
	assert.Contains(t, htmlEn, ">Choose a day<")
	assert.Contains(t, htmlRu, ">Выберите день<")

	// The "+N more" overflow hint.
	assert.Contains(t, htmlEn, "+1 more")
	assert.NotContains(t, htmlEn, "ещё")
	assert.Contains(t, htmlRu, "+1 ещё")
	assert.NotContains(t, htmlRu, "more")

	// The empty-day label travels as an attribute the day-click handler reads.
	assert.Contains(t, htmlEn, `data-empty-label="No events"`, "English empty-day label")
	assert.NotContains(t, htmlEn, `data-empty-label="Нет событий"`)
	assert.Contains(t, htmlRu, `data-empty-label="Нет событий"`, "Russian empty-day label")
	assert.NotContains(t, htmlRu, `data-empty-label="No events"`)

	// The handler must read the attribute, never a hardcoded label.
	script := renderIn(t, nil, CalendarMonthScript())
	assert.Contains(t, script, "data-empty-label")
	assert.NotContains(t, script, "Нет событий")
	assert.NotContains(t, script, "No events")
}

func TestCalendarTimelineStringsAndFormatsFollowTheLocale(t *testing.T) {
	day := time.Date(2026, time.March, 9, 10, 0, 0, 0, time.UTC)
	p := &CalendarTimelineProps{
		ID: "tl",
		Events: []CalendarEvent{
			{ID: "e1", Title: "Demo", Start: day, End: day.Add(time.Hour), Color: "blue"},
			{ID: "e2", Title: "All day demo", Start: day, End: day, Color: "green", AllDay: true},
		},
	}
	htmlEn := renderIn(t, nil, CalendarTimeline(p))
	htmlRu := renderIn(t, i18n.Russian, CalendarTimeline(p))

	assert.Contains(t, htmlEn, ">Events calendar<")
	assert.Contains(t, htmlRu, ">Календарь событий<")

	// Day-group weekday name and date resolve through i18n.WeekdayName/Date.
	assert.Contains(t, htmlEn, "Mon, 03/09/2026", "English day-group heading")
	assert.NotContains(t, htmlEn, "Пн, 09.03.2026")
	assert.Contains(t, htmlRu, "Пн, 09.03.2026", "Russian day-group heading")
	assert.NotContains(t, htmlRu, "Mon, 03/09/2026")

	assert.Contains(t, htmlEn, ">All day<")
	assert.Contains(t, htmlRu, ">Весь день<")

	// Empty state.
	empty := &CalendarTimelineProps{ID: "tl-empty"}
	htmlEn = renderIn(t, nil, CalendarTimeline(empty))
	htmlRu = renderIn(t, i18n.Russian, CalendarTimeline(empty))
	assert.Contains(t, htmlEn, ">No events<")
	assert.Contains(t, htmlRu, ">Нет событий<")

	// Loading sentinel.
	loading := &CalendarTimelineProps{ID: "tl-more", HasMore: true, NextCursor: "c"}
	htmlEn = renderIn(t, nil, CalendarTimeline(loading))
	htmlRu = renderIn(t, i18n.Russian, CalendarTimeline(loading))
	assert.Contains(t, htmlEn, "Loading…")
	assert.NotContains(t, htmlEn, "Загрузка")
	assert.Contains(t, htmlRu, "Загрузка…")
	assert.NotContains(t, htmlRu, "Loading")
}

func TestCalendarDayAndWeekStringsAndFormatsFollowTheLocale(t *testing.T) {
	day := time.Date(2026, time.March, 9, 10, 0, 0, 0, time.UTC)
	allDay := []CalendarEvent{{ID: "e1", Title: "All day demo", Start: day, End: day, Color: "green", AllDay: true}}

	dayP := &CalendarDayProps{ID: "day", Date: day, Events: allDay}
	htmlEn := renderIn(t, nil, CalendarDay(dayP))
	htmlRu := renderIn(t, i18n.Russian, CalendarDay(dayP))

	// The heading pairs a weekday name and a date; both halves follow the
	// locale rather than mixing a hardcoded name with a formatted date.
	assert.Contains(t, htmlEn, "Mon, 03/09/2026", "English day heading")
	assert.NotContains(t, htmlEn, "Пн, 09.03.2026")
	assert.NotContains(t, htmlEn, "Понедельник", "the day heading must not keep a hardcoded weekday name")
	assert.Contains(t, htmlRu, "Пн, 09.03.2026", "Russian day heading")
	assert.NotContains(t, htmlRu, "Mon, 03/09/2026")

	assert.Contains(t, htmlEn, ">Today<", "English today button")
	assert.NotContains(t, htmlEn, "Сегодня")
	assert.Contains(t, htmlRu, ">Сегодня<", "Russian today button")
	assert.NotContains(t, htmlRu, ">Today<")

	assert.Contains(t, htmlEn, ">All day<", "English all-day label")
	assert.NotContains(t, htmlEn, "Весь день")
	assert.Contains(t, htmlRu, ">Весь день<", "Russian all-day label")
	assert.NotContains(t, htmlRu, ">All day<")

	weekP := &CalendarWeekProps{ID: "week", StartDate: day, Events: allDay}
	htmlEn = renderIn(t, nil, CalendarWeek(weekP))
	htmlRu = renderIn(t, i18n.Russian, CalendarWeek(weekP))

	assert.Contains(t, htmlEn, "03/09/2026 — 03/15/2026", "English week heading")
	assert.NotContains(t, htmlEn, "09.03.2026")
	assert.Contains(t, htmlRu, "09.03.2026 — 15.03.2026", "Russian week heading")
	assert.NotContains(t, htmlRu, "03/09/2026")

	assert.Contains(t, htmlEn, ">Today<", "English today button")
	assert.NotContains(t, htmlEn, "Сегодня")
	assert.Contains(t, htmlRu, ">Сегодня<", "Russian today button")
	assert.NotContains(t, htmlRu, ">Today<")

	assert.Contains(t, htmlEn, ">All day<", "English all-day label")
	assert.NotContains(t, htmlEn, "Весь день")
	assert.Contains(t, htmlRu, ">Весь день<", "Russian all-day label")
	assert.NotContains(t, htmlRu, ">All day<")
}

func TestStatStringsFollowTheLocale(t *testing.T) {
	p := &StatProps{Title: "Revenue", Value: "1 234", Change: "+12%", ChangeType: "up"}
	htmlEn := renderIn(t, nil, StatCard(p))
	htmlRu := renderIn(t, i18n.Russian, StatCard(p))
	assert.Contains(t, htmlEn, "vs previous period")
	assert.NotContains(t, htmlEn, "vs прошлый период")
	assert.Contains(t, htmlRu, "vs прошлый период")
	assert.NotContains(t, htmlRu, "vs previous period")
}

func TestThemeToggleStringsFollowTheLocale(t *testing.T) {
	htmlEn := renderIn(t, nil, ThemeToggle())
	htmlRu := renderIn(t, i18n.Russian, ThemeToggle())
	assert.Contains(t, htmlEn, ">Theme<")
	assert.Contains(t, htmlEn, `title="Toggle light/dark theme"`)
	assert.Contains(t, htmlEn, `aria-label="Toggle light/dark theme"`)
	assert.NotContains(t, htmlEn, "Тема")
	assert.Contains(t, htmlRu, ">Тема<")
	assert.Contains(t, htmlRu, `title="Переключить светлую/тёмную тему"`)
	assert.Contains(t, htmlRu, `aria-label="Переключить светлую/тёмную тему"`)
	assert.NotContains(t, htmlRu, "Toggle light/dark theme")
}

func TestAlertStringsFollowTheLocale(t *testing.T) {
	type want struct{ en, ru string }
	titles := map[string]want{
		"error":   {"Error", "Ошибка"},
		"warning": {"Warning", "Предупреждение"},
		"success": {"Success", "Успех"},
		"info":    {"Info", "Информация"},
	}
	for style, w := range titles {
		p := &AlertProps{Title: "T", Style: style, Dismissible: true}
		htmlEn := renderIn(t, nil, Alert(p))
		htmlRu := renderIn(t, i18n.Russian, Alert(p))
		assert.Contains(t, htmlEn, ">"+w.en+"<", "English %s sr-only title", style)
		assert.NotContains(t, htmlEn, ">"+w.ru+"<")
		assert.Contains(t, htmlRu, ">"+w.ru+"<", "Russian %s sr-only title", style)
		assert.NotContains(t, htmlRu, ">"+w.en+"<")
	}

	p := &AlertProps{Style: "success", Dismissible: true}
	htmlEn := renderIn(t, nil, Alert(p))
	htmlRu := renderIn(t, i18n.Russian, Alert(p))
	assert.Contains(t, htmlEn, `aria-label="Close"`)
	assert.Contains(t, htmlEn, ">Close<")
	assert.NotContains(t, htmlEn, "Закрыть")
	assert.Contains(t, htmlRu, `aria-label="Закрыть"`)
	assert.Contains(t, htmlRu, ">Закрыть<")
	assert.NotContains(t, htmlRu, ">Close<")
}

func TestToastStringsFollowTheLocale(t *testing.T) {
	p := &ToastProps{ID: "t1", Title: "Saved", Style: "success"}
	htmlEn := renderIn(t, nil, Toast(p))
	htmlRu := renderIn(t, i18n.Russian, Toast(p))
	assert.Contains(t, htmlEn, `aria-label="Close"`)
	assert.Contains(t, htmlEn, ">Close<")
	assert.NotContains(t, htmlEn, "Закрыть")
	assert.Contains(t, htmlRu, `aria-label="Закрыть"`)
	assert.Contains(t, htmlRu, ">Закрыть<")
	assert.NotContains(t, htmlRu, ">Close<")
}

func TestMigratedComponentsNoLongerCarryCatalogLiterals(t *testing.T) {
	// The Russian catalog text must only ever appear through a Russian context;
	// a default-context render of every migrated component is English-only.
	bad := []string{
		"Выберите день", "Календарь событий", "Выйти", "Свернуть",
		"Переключить", "Аватар", "vs прошлый период", "Закрыть", "Нет событий",
		"Загрузка", "Сегодня",
	}
	comps := map[string]templ.Component{
		"sidebar":      Sidebar(&SidebarProps{Items: []*SidebarItem{{Name: "H", URL: "/"}}, User: &UserInfo{Name: "A", Email: "e"}, LogoutURL: "/l", Collapsible: true, ID: "s"}),
		"sidebar-new":  SidebarNew(),
		"calendar":     CalendarMonth(&CalendarMonthProps{ID: "c", Year: 2026, Month: time.March}),
		"timeline":     CalendarTimeline(&CalendarTimelineProps{ID: "t"}),
		"stat":         StatCard(&StatProps{Title: "x", Value: "1", Change: "+1%", ChangeType: "up"}),
		"theme-toggle": ThemeToggle(),
		"alert":        Alert(&AlertProps{Style: "error", Dismissible: true}),
		"toast":        Toast(&ToastProps{ID: "t", Style: "error"}),
	}
	for name, c := range comps {
		html := renderIn(t, nil, c)
		for _, s := range bad {
			assert.NotContains(t, html, s, "default-context %s must not render %q", name, s)
		}
	}
}
