package main

import (
	"log"
	"net/http"

	"github.com/fbv/go-templ-ui/view"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

func main() {
	app := fiber.New()

	app.Use("/static", filesystem.New(filesystem.Config{
		Root:   http.FS(view.Content),
		Browse: true,
	}))

	app.Get("/crm", func(c *fiber.Ctx) error {
		return view.Render(c, DashboardPage())
	})

	app.Get("/crm/clients", func(c *fiber.Ctx) error {
		return view.Render(c, ClientsPage())
	})

	app.Get("/crm/clients/:id", func(c *fiber.Ctx) error {
		return view.Render(c, ClientPage())
	})

	app.Get("/crm/deals", func(c *fiber.Ctx) error {
		return view.Render(c, DealsPage())
	})

	app.Get("/crm/tasks", func(c *fiber.Ctx) error {
		return view.Render(c, TasksPage())
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/crm")
	})

	log.Fatal(app.Listen(":8081"))
}
