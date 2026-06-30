package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fbv/go-templ-ui/view"
	"github.com/fbv/go-templ-ui/view/ui"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"net/http"
)

func main() {
	app := fiber.New(fiber.Config{})

	app.Use("/static", filesystem.New(filesystem.Config{
		Root:   http.FS(view.Content),
		Browse: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/chat")
	})

	app.Get("/chat", func(c *fiber.Ctx) error {
		return view.Render(c, ChatPage())
	})

	app.Get("/chat/:id", func(c *fiber.Ctx) error {
		return view.Render(c, ChatPage())
	})

	app.Post("/chat/:id/messages", handleMessagePost)
	app.Post("/chat/:id/messages/file", handleFilePost)

	app.Static("/uploads", "./tmp/uploads")

	log.Fatal(app.Listen(":8083"))
}

func handleMessagePost(c *fiber.Ctx) error {
	chatID := c.Params("id")
	content := c.FormValue("content")
	senderID := c.FormValue("sender_id")
	msgType := "text"

	var fileURL, fileName, fileSize string
	file, err := c.FormFile("file")
	if err == nil && file != nil {
		msgType = "file"
		fileName = file.Filename
		fileSize = formatFileSize(file.Size)

		uploadDir := "./tmp/uploads"
		os.MkdirAll(uploadDir, 0755)
		savePath := filepath.Join(uploadDir, fileName)
		if saveErr := c.SaveFile(file, savePath); saveErr == nil {
			fileURL = "/uploads/" + fileName
		}
	}

	if content == "" && msgType == "text" {
		return c.SendString("")
	}

	msg := AddMessage(chatID, senderID, content, msgType, fileURL, fileName, fileSize)
	if msg == nil {
		return c.SendString("")
	}

	return view.Render(c, ui.ChatMessageBubble(msg))
}

func handleFilePost(c *fiber.Ctx) error {
	chatID := c.Params("id")
	senderID := c.FormValue("sender_id")

	file, err := c.FormFile("file")
	if err != nil {
		return c.SendString("")
	}

	fileName := file.Filename
	fileSize := formatFileSize(file.Size)
	msgType := "file"

	uploadDir := "./tmp/uploads"
	os.MkdirAll(uploadDir, 0755)
	savePath := filepath.Join(uploadDir, fileName)
	var fileURL string
	if saveErr := c.SaveFile(file, savePath); saveErr == nil {
		fileURL = "/uploads/" + fileName
	}

	msg := AddMessage(chatID, senderID, "", msgType, fileURL, fileName, fileSize)
	if msg == nil {
		return c.SendString("")
	}

	return view.Render(c, ui.ChatMessageBubble(msg))
}

func formatFileSize(bytes int64) string {
	if bytes < 1024 {
		return "1 B"
	}
	if bytes < 1048576 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	}
	return fmt.Sprintf("%.1f MB", float64(bytes)/1048576)
}
