package server

import (
	"embed"
	"net/http"

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
	s.fiber = fiber.New(fiber.Config{})
	s.fiber.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/showcase/components")
	})
	r := s.fiber.Group("/showcase")
	r.Use("/static", filesystem.New(filesystem.Config{
		Root:   http.FS(s.content),
		Browse: true,
	}))
	r.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/showcase/components")
	})
	r.Get("/components", view.StaticPage(showcase.Showcase()))
	r.Get("/login", view.StaticPage(showcase.LoginPage()))
	r.Get("/icons", view.StaticPage(showcase.IconsPage()))
	r.Get("/panels", view.StaticPage(showcase.PanelsPage()))
	return s.fiber.Listen(":8082")
}
