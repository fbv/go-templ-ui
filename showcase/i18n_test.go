package showcase

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func renderDemo(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf strings.Builder
	require.NoError(t, c.Render(context.Background(), &buf))
	return buf.String()
}

func i18nRow(t *testing.T) Demo {
	t.Helper()
	for _, d := range Demos() {
		if d.Name == "i18n" {
			return d
		}
	}
	t.Fatal("no i18n row in the registry")
	return Demo{}
}

// The demo page must render under Russian: document language, catalog strings
// and the ru format rows — while every other page keeps English.
func TestI18nDemoPageRendersUnderRussian(t *testing.T) {
	html := renderDemo(t, i18nRow(t).Render(i18nRow(t)))

	assert.Contains(t, html, `<html lang="ru">`, "the demo page selects the Russian locale")

	assert.Contains(t, html, "Свернуть", "sidebar collapse label")
	assert.Contains(t, html, "Выйти", "sidebar logout label")
	assert.Contains(t, html, ">Тема<", "theme toggle label")
	assert.Contains(t, html, "Переключить светлую/тёмную тему", "theme toggle title")

	assert.Contains(t, html, "vs прошлый период", "stat comparison label")
	assert.Contains(t, html, "Закрыть", "alert/toast close label")
	assert.Contains(t, html, ">Успех<", "alert sr-only title")

	assert.Contains(t, html, "Март 2026", "calendar month title")
	assert.Contains(t, html, ">Пн<", "calendar weekday header")
	assert.Contains(t, html, ">Сегодня<", "calendar today button")
	assert.Contains(t, html, ">Выберите день<", "calendar day panel")
	assert.Contains(t, html, `data-empty-label="Нет событий"`, "calendar empty-day label")
	assert.Contains(t, html, ">Календарь событий<", "timeline title")
	assert.Contains(t, html, "Пн, 09.03.2026", "timeline day-group heading")
}

func TestI18nDemoPageFormatsReadTheRussianRows(t *testing.T) {
	html := renderDemo(t, i18nRow(t).Render(i18nRow(t)))

	assert.Contains(t, html, "1 234,50", "number sample")
	assert.Contains(t, html, "1 234,50 ₽", "currency sample")
	assert.Contains(t, html, "09.03.2026", "date sample")
	assert.Contains(t, html, "Март 2026", "month-year sample")

	// The English forms must not leak in. English RUB is "₽1,234.50" (symbol
	// before, no space) — not "$1,234.50".
	for _, s := range []string{"1,234.50", "₽1,234.50", "$1,234.50", "03/09/2026", "March 2026"} {
		assert.NotContains(t, html, s, "English format form %q must not appear", s)
	}
}

func TestI18nDemoPageCalendarNavigationIsFragmentWired(t *testing.T) {
	html := renderDemo(t, i18nRow(t).Render(i18nRow(t)))

	assert.Contains(t, html, `id="i18n-calendar-month"`, "the demo calendar carries the id the fragment swaps")
	// Prev/next/Today must point at the registered fragment endpoint, never at
	// the page's own document route.
	assert.Contains(t, html, I18nCalendarMonthURL, "month navigation must use the fragment route")
	assert.NotContains(t, html, `hx-get="/showcase/i18n?`, "month navigation must not hit the document route")
	// Every hx-get on the page targets the fragment endpoint.
	targets := strings.Split(html, `hx-get="`)
	require.NotEmpty(t, targets[1:], "the month navigation must render hx-get buttons")
	for _, target := range targets[1:] {
		url, _, quoted := strings.Cut(target, `"`)
		require.True(t, quoted, "an hx-get value must be quoted")
		assert.True(t, strings.HasPrefix(url, I18nCalendarMonthURL), "hx-get %q must point at %s", url, I18nCalendarMonthURL)
	}
}

func TestI18nDemoPageTimelineHasNoDeadLoadMore(t *testing.T) {
	html := renderDemo(t, i18nRow(t).Render(i18nRow(t)))

	assert.NotContains(t, html, `class="calendar-timeline-sentinel"`, "a sentinel with no endpoint must not render")
	assert.NotContains(t, html, `class="calendar-timeline-loading`, "a stuck loading spinner must not render")
	assert.NotContains(t, html, `data-load-more="`, "no load-more control")
	assert.NotContains(t, html, "timeline/more", "no unregistered load-more URL")
}

func TestI18nCalendarMonthFragmentIsNotADocument(t *testing.T) {
	html := renderDemo(t, I18nCalendarMonth(2026, time.March))

	assert.NotContains(t, html, "<!doctype", "the fragment must not carry a document shell")
	assert.NotContains(t, html, "<html", "the fragment must not carry a document shell")
	assert.Contains(t, html, `id="i18n-calendar-month"`, "the fragment is the calendar element itself")
	assert.Contains(t, html, "Март 2026", "the fragment renders under the demo locale")
}

func TestOtherDemoPagesStayEnglish(t *testing.T) {
	for _, d := range Demos() {
		if d.Name == "i18n" {
			continue
		}
		html := renderDemo(t, d.Render(d))
		assert.Contains(t, html, `<html lang="en">`, "row %q must stay on the default locale", d.Name)
		assert.NotContains(t, html, "Выйти", "row %q must not render Russian strings", d.Name)
		assert.Contains(t, html, "Log out", "row %q keeps the English catalog text", d.Name)
	}
}
