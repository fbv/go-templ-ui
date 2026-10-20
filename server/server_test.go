package server

import (
	"embed"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func get(t *testing.T, path string) string {
	t.Helper()
	app := New(embed.FS{}).routes()
	resp, err := app.Test(httptest.NewRequest("GET", path, nil))
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode, "GET %s must be registered", path)
	return string(body)
}

// The demo calendar's month navigation must hit a route that returns just the
// calendar element under the demo locale — never a full HTML document.
func TestI18nCalendarMonthRouteReturnsAFragment(t *testing.T) {
	body := get(t, "/showcase/i18n/calendar/month?year=2026&month=3")

	assert.NotContains(t, body, "<!doctype", "the month route must not return a document")
	assert.NotContains(t, body, "<html", "the month route must not return a document")
	assert.Contains(t, body, `id="i18n-calendar-month"`, "the route returns the calendar element")
	assert.Contains(t, body, "Март 2026", "the fragment renders under the demo locale")
}

// The demo page itself is a document, and it is Russian.
func TestI18nDemoPageRouteReturnsARussianDocument(t *testing.T) {
	body := get(t, "/showcase/i18n")

	assert.Contains(t, body, "<!doctype html>", "the demo page is a full document")
	assert.Contains(t, body, `<html lang="ru">`, "the demo page selects the Russian locale")
	assert.Contains(t, body, "1 234,50 ₽", "the ru currency sample")
}

// Every other registered demo page keeps the default English locale.
func TestOtherDemoRoutesStayEnglish(t *testing.T) {
	body := get(t, "/showcase/badge")
	assert.Contains(t, body, `<html lang="en">`)
	assert.NotContains(t, body, "Выйти")
	assert.Contains(t, body, "Log out", "row keeps the English catalog text")
}
