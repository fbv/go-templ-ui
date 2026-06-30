package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/fbv/go-templ-ui/view/ui"
)

type User struct {
	ID    string
	Name  string
	Email string
}

type ChatRoom struct {
	ID       string
	Name     string
	Members  []string
	Messages []ui.ChatMessage
}

var (
	users = map[string]*User{
		"u1": {ID: "u1", Name: "Алексей Петров", Email: "alexey@example.com"},
		"u2": {ID: "u2", Name: "Мария Иванова", Email: "maria@example.com"},
		"u3": {ID: "u3", Name: "Дмитрий Козлов", Email: "dmitry@example.com"},
		"u4": {ID: "u4", Name: "Елена Смирнова", Email: "elena@example.com"},
	}

	chatRooms = map[string]*ChatRoom{
		"chat1": {
			ID:      "chat1",
			Name:    "Мария Иванова",
			Members: []string{"u1", "u2"},
			Messages: []ui.ChatMessage{
				{ID: "m1", SenderID: "u2", SenderName: "Мария Иванова", Content: "Привет! Как дела с проектом?", Time: "10:30", Type: "text"},
				{ID: "m2", SenderID: "u1", SenderName: "Алексей Петров", Content: "Привет! Всё идёт по плану, закончил основной функционал", Time: "10:32", Type: "text", IsOwn: true},
				{ID: "m3", SenderID: "u2", SenderName: "Мария Иванова", Content: "Отлично! А когда будет готово для тестирования?", Time: "10:33", Type: "text"},
				{ID: "m4", SenderID: "u1", SenderName: "Алексей Петров", Content: "Думаю, до конца недели. Вот скриншот текущего состояния", Time: "10:35", Type: "text", IsOwn: true},
				{ID: "m5", SenderID: "u1", SenderName: "Алексей Петров", Content: "", Time: "10:35", Type: "image", FileURL: "https://placehold.co/600x400/e2e8f0/475569?text=Screenshot", FileName: "screenshot.png", IsOwn: true},
				{ID: "m6", SenderID: "u2", SenderName: "Мария Иванова", Content: "Выглядит хорошо! Ещё отправлю тебе документацию", Time: "10:38", Type: "text"},
				{ID: "m7", SenderID: "u2", SenderName: "Мария Иванова", Content: "", Time: "10:38", Type: "file", FileURL: "#", FileName: "docs.pdf", FileSize: "2.4 MB"},
				{ID: "m8", SenderID: "u1", SenderName: "Алексей Петров", Content: "Спасибо, посмотрю!", Time: "10:40", Type: "text", IsOwn: true},
			},
		},
		"chat2": {
			ID:      "chat2",
			Name:    "Дмитрий Козлов",
			Members: []string{"u1", "u3"},
			Messages: []ui.ChatMessage{
				{ID: "m9", SenderID: "u3", SenderName: "Дмитрий Козлов", Content: "Добрый день! Нужна помощь с настройкой CI/CD", Time: "14:00", Type: "text"},
				{ID: "m10", SenderID: "u1", SenderName: "Алексей Петров", Content: "Конечно! Что конкретно не работает?", Time: "14:05", Type: "text", IsOwn: true},
				{ID: "m11", SenderID: "u3", SenderName: "Дмитрий Козлов", Content: "Пайплайн падает на этапе тестов. Вот логи", Time: "14:10", Type: "text"},
				{ID: "m12", SenderID: "u3", SenderName: "Дмитрий Козлов", Content: "", Time: "14:10", Type: "file", FileURL: "#", FileName: "ci-logs.txt", FileSize: "45 KB"},
				{ID: "m13", SenderID: "u1", SenderName: "Алексей Петров", Content: "Попробуй обновить версию Go в конфиге", Time: "14:15", Type: "text", IsOwn: true},
			},
		},
		"chat3": {
			ID:      "chat3",
			Name:    "Елена Смирнова",
			Members: []string{"u1", "u4"},
			Messages: []ui.ChatMessage{
				{ID: "m14", SenderID: "system", SenderName: "", Content: "Чат создан", Time: "09:00", Type: "system"},
				{ID: "m15", SenderID: "u4", SenderName: "Елена Смирнова", Content: "Привет! Добро пожаловать в команду 🎉", Time: "09:01", Type: "text"},
				{ID: "m16", SenderID: "u1", SenderName: "Алексей Петров", Content: "Спасибо! Рад быть здесь", Time: "09:05", Type: "text", IsOwn: true},
			},
		},
	}

	currentUserID = "u1"
	storeMu       sync.RWMutex
)

func GetChatListItems() []ui.ChatListItem {
	storeMu.RLock()
	defer storeMu.RUnlock()

	items := make([]ui.ChatListItem, 0, len(chatRooms))
	for _, room := range chatRooms {
		lastMsg := ""
		lastTime := ""
		if len(room.Messages) > 0 {
			msg := room.Messages[len(room.Messages)-1]
			if msg.Type == "text" {
				lastMsg = msg.Content
			} else if msg.Type == "file" {
				lastMsg = "📎 " + msg.FileName
			} else if msg.Type == "image" {
				lastMsg = "🖼 " + msg.FileName
			} else if msg.Type == "system" {
				lastMsg = msg.Content
			}
			lastTime = msg.Time
		}
		items = append(items, ui.ChatListItem{
			ID:          room.ID,
			Name:        room.Name,
			LastMessage: lastMsg,
			Time:        lastTime,
			UnreadCount: 0,
			Active:      false,
		})
	}
	return items
}

func GetChatMessages(chatID string) []ui.ChatMessage {
	storeMu.RLock()
	defer storeMu.RUnlock()

	if room, ok := chatRooms[chatID]; ok {
		return room.Messages
	}
	return nil
}

func GetChatName(chatID string) string {
	storeMu.RLock()
	defer storeMu.RUnlock()

	if room, ok := chatRooms[chatID]; ok {
		return room.Name
	}
	return ""
}

func AddMessage(chatID, senderID, content, msgType, fileURL, fileName, fileSize string) *ui.ChatMessage {
	storeMu.Lock()
	defer storeMu.Unlock()

	room, ok := chatRooms[chatID]
	if !ok {
		return nil
	}

	sender := users[senderID]
	senderName := ""
	if sender != nil {
		senderName = sender.Name
	}

	msg := ui.ChatMessage{
		ID:           fmt.Sprintf("m%d", len(room.Messages)+1),
		SenderID:     senderID,
		SenderName:   senderName,
		Content:      content,
		Time:         time.Now().Format("15:04"),
		Type:         msgType,
		FileURL:      fileURL,
		FileName:     fileName,
		FileSize:     fileSize,
		IsOwn:        senderID == currentUserID,
	}

	room.Messages = append(room.Messages, msg)
	return &msg
}
