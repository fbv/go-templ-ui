package showcase

import (
	"context"
	"io"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/fbv/go-templ-ui/view"
	"github.com/fbv/go-templ-ui/view/i18n"
	"github.com/fbv/go-templ-ui/view/ui"
	"github.com/gofiber/fiber/v2"
)

// I18nCalendarMonthURL is the fragment endpoint the demo calendar's month
// navigation points at. It returns just the CalendarMonth element under the
// demo's Russian locale, so the prev/next/Today buttons' hx-swap="outerHTML"
// replaces the calendar element instead of injecting a full HTML document.
const I18nCalendarMonthURL = "/showcase/i18n/calendar/month"

// i18nDemoDate pins the demo's date samples to 09.03.2026 (ru) / 03/09/2026
// (en) so the page — and the tests reading it — show one fixed snapshot.
var i18nDemoDate = time.Date(2026, time.March, 9, 14, 5, 0, 0, time.UTC)

// I18nDemoPage is the Render entry of the i18n registry row. It wraps the page
// in the demo locale with the same templ.ComponentFunc + view.Render form
// docs/i18n.md documents for consumers, so layout.Page emits <html lang="ru">
// and every migrated string and format below resolves through the Russian
// catalog. The rest of the showcase stays on the default English locale.
func I18nDemoPage(d Demo) templ.Component {
	return inDemoLocale(i18nDemoPage(d))
}

// I18nCalendarMonth returns the demo calendar's month fragment under the same
// Russian locale as the page it belongs to.
func I18nCalendarMonth(year int, month time.Month) templ.Component {
	return inDemoLocale(i18nCalendarMonthFragment(year, month))
}

// I18nCalendarMonthHandler serves I18nCalendarMonthURL. server/server.go
// mounts it; the demo page's LoadURL comes from the same constant, so a control
// on the page cannot point at a route that is not registered.
func I18nCalendarMonthHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		year := i18nDemoDate.Year()
		month := i18nDemoDate.Month()
		if v, err := strconv.Atoi(c.Query("year")); err == nil && v > 0 {
			year = v
		}
		if v, err := strconv.Atoi(c.Query("month")); err == nil && v >= 1 && v <= 12 {
			month = time.Month(v)
		}
		return view.Render(c, I18nCalendarMonth(year, month))
	}
}

func inDemoLocale(c templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return c.Render(i18n.WithLocale(ctx, i18n.Russian), w)
	})
}

// i18nDemoEvents returns a small fixed event set on the 3rd, 5th and 9th of
// the given month, so month navigation keeps showing sample events and the
// 9th lines up with i18nDemoDate.
func i18nDemoEvents(year int, month time.Month) []ui.CalendarEvent {
	at := func(day, hour int) time.Time {
		return time.Date(year, month, day, hour, 0, 0, 0, time.UTC)
	}
	return []ui.CalendarEvent{
		{ID: "demo-1", Title: "Планирование спринта", Start: at(3, 10), End: at(3, 11), Color: "blue"},
		{ID: "demo-2", Title: "Код-ревью", Start: at(5, 14), End: at(5, 15), Color: "purple"},
		{ID: "demo-3", Title: "Встреча с заказчиком", Start: at(9, 10), End: at(9, 11), Color: "green"},
		{ID: "demo-4", Title: "Демонстрация релиза", Start: at(9, 14), End: at(9, 15), Color: "red"},
		{ID: "demo-5", Title: "Ежегодная конференция", Start: at(9, 0), End: at(9, 23), Color: "yellow", AllDay: true},
	}
}
