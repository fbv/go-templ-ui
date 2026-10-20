package server

import (
	"embed"
	"net/http"
	"strings"

	"github.com/fbv/go-templ-ui/docs"
	"github.com/fbv/go-templ-ui/showcase"
	"github.com/fbv/go-templ-ui/view"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

type Server struct {
	content embed.FS
	fiber   *fiber.App
}

func New(content embed.FS) *Server {
	return &Server{
		content: content,
	}
}

func (s *Server) Run() error {
	return s.routes().Listen(":8082")
}

// routes builds the app and registers every route, so tests can exercise the
// real routing table without listening on a socket.
func (s *Server) routes() *fiber.App {
	s.fiber = fiber.New(fiber.Config{})
	s.fiber.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/showcase/components")
	})
	r := s.fiber.Group("/showcase")
	r.Use("/static", filesystem.New(filesystem.Config{
		Root:   http.FS(s.content),
		Browse: true,
	}))
	// Demo pages link to docs/<name>.md relative to /showcase/<name>, which
	// resolves to /showcase/docs/<name>.md. Mount the embedded docs tree so
	// that href is a live URL instead of a 404.
	r.Use("/docs", filesystem.New(filesystem.Config{
		Root:   http.FS(docs.Content),
		Browse: true,
	}))
	r.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/showcase/components")
	})
	r.Get("/components", view.StaticPage(showcase.Showcase()))
	r.Get("/login", view.StaticPage(showcase.LoginPage()))
	r.Get("/icons", view.StaticPage(showcase.IconsPage()))
	r.Get("/panels", view.StaticPage(showcase.PanelsPage()))
	r.Get("/calendar", view.StaticPage(showcase.CalendarPage()))
	r.Post("/api/selected", func(c *fiber.Ctx) error {
		ids := c.FormValue("ids")
		return c.SendString("Выбраны элементы: " + ids)
	})
	r.Get("/tree", view.StaticPage(showcase.TreePage()))
	// The i18n demo calendar's month navigation swaps the calendar element
	// alone; the demo page's LoadURL comes from the same constant.
	r.Get(strings.TrimPrefix(showcase.I18nCalendarMonthURL, "/showcase"), showcase.I18nCalendarMonthHandler())
	for _, d := range showcase.Demos() {
		r.Get(strings.TrimPrefix(d.Route, "/showcase"), view.StaticPage(d.Render(d)))
	}
	return s.fiber
}
