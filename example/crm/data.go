package main

import (
	"github.com/a-h/templ"
	"github.com/fbv/go-templ-ui/view/ui"
)

// ===== Clients =====

type Client struct {
	ID        int
	Name      string
	Email     string
	Phone     string
	Company   string
	Status    string // active, lead, inactive
	CreatedAt string
}

var clients = []Client{
	{1, "Иван Петров", "ivan@example.com", "+7 900 123-45-67", "ООО Ромашка", "active", "2024-01-15"},
	{2, "Мария Сидорова", "maria@example.com", "+7 900 234-56-78", "ЗАО Альфа", "active", "2024-02-20"},
	{3, "Алексей Козлов", "alexey@example.com", "+7 900 345-67-89", "ИП Козлов", "lead", "2024-03-10"},
	{4, "Елена Новикова", "elena@example.com", "+7 900 456-78-90", "ООО Бета", "active", "2024-04-05"},
	{5, "Дмитрий Волков", "dmitry@example.com", "+7 900 567-89-01", "ПАО Гамма", "inactive", "2024-05-12"},
	{6, "Ольга Морозова", "olga@example.com", "+7 900 678-90-12", "ООО Дельта", "active", "2024-06-18"},
	{7, "Сергей Лебедев", "sergey@example.com", "+7 900 789-01-23", "ЗАО Эпсилон", "lead", "2024-07-22"},
	{8, "Наталья Кузнецова", "natalya@example.com", "+7 900 890-12-34", "ООО Дзета", "active", "2024-08-30"},
	{9, "Павел Соколов", "pavel@example.com", "+7 900 901-23-45", "ИП Соколов", "active", "2024-09-14"},
	{10, "Анна Попова", "anna@example.com", "+7 900 012-34-56", "ООО Эта", "lead", "2024-10-01"},
	{11, "Михаил Федоров", "mikhail@example.com", "+7 900 111-22-33", "ПАО Тета", "active", "2024-10-15"},
	{12, "Татьяна Егорова", "tatiana@example.com", "+7 900 222-33-44", "ЗАО Йота", "inactive", "2024-11-02"},
}

type ClientDataSource []Client

func (ds ClientDataSource) GetColumnCount() int { return 5 }
func (ds ClientDataSource) GetColumnName(col int) string {
	names := []string{"Имя", "Email", "Компания", "Статус", "Дата"}
	return names[col]
}
func (ds ClientDataSource) GetColumnStyle(col int) string { return "" }
func (ds ClientDataSource) GetRowCount() int              { return len(ds) }
func (ds ClientDataSource) GetData(row, col int) templ.Component {
	c := ds[row]
	switch col {
	case 0:
		return ui.Text(c.Name)
	case 1:
		return ui.Text(c.Email)
	case 2:
		return ui.Text(c.Company)
	case 3:
		return ui.Badge(&ui.BadgeProps{Name: clientStatusName(c.Status), Style: clientStatusStyle(c.Status)})
	case 4:
		return ui.Text(c.CreatedAt)
	}
	return ui.Text("")
}
func (ds ClientDataSource) GetRowMeta(row int) *ui.RowMeta {
	return &ui.RowMeta{HRef: "/crm/clients/" + itoa(ds[row].ID)}
}

func clientStatusName(s string) string {
	switch s {
	case "active":
		return "Активный"
	case "lead":
		return "Лид"
	default:
		return "Неактивный"
	}
}

func clientStatusStyle(s string) string {
	switch s {
	case "active":
		return "green"
	case "lead":
		return "blue"
	default:
		return "gray"
	}
}

// ===== Deals =====

type Deal struct {
	ID       int
	Title    string
	ClientID int
	Client   string
	Amount   string
	Stage    string // new, negotiation, proposal, closed_won, closed_lost
	Date     string
}

var deals = []Deal{
	{1, "Внедрение CRM", 1, "ООО Ромашка", "1 500 000 ₽", "negotiation", "2024-11-01"},
	{2, "Разработка сайта", 2, "ЗАО Альфа", "850 000 ₽", "proposal", "2024-11-05"},
	{3, "Техподдержка", 4, "ООО Бета", "360 000 ₽", "closed_won", "2024-10-20"},
	{4, "Миграция в облако", 6, "ООО Дельта", "2 200 000 ₽", "new", "2024-11-10"},
	{5, "Аудит безопасности", 8, "ООО Дзета", "450 000 ₽", "negotiation", "2024-11-08"},
	{6, "Интеграция 1С", 9, "ИП Соколов", "720 000 ₽", "proposal", "2024-11-12"},
	{7, "Разработка мобильного приложения", 11, "ПАО Тета", "3 100 000 ₽", "new", "2024-11-15"},
	{8, "Обновление ERP", 1, "ООО Ромашка", "980 000 ₽", "closed_lost", "2024-09-15"},
	{9, "Настройка серверов", 4, "ООО Бета", "280 000 ₽", "closed_won", "2024-10-01"},
	{10, "Консалтинг", 2, "ЗАО Альфа", "150 000 ₽", "closed_won", "2024-10-10"},
}

type DealDataSource []Deal

func (ds DealDataSource) GetColumnCount() int { return 5 }
func (ds DealDataSource) GetColumnName(col int) string {
	names := []string{"Сделка", "Клиент", "Сумма", "Стадия", "Дата"}
	return names[col]
}
func (ds DealDataSource) GetColumnStyle(col int) string { return "" }
func (ds DealDataSource) GetRowCount() int              { return len(ds) }
func (ds DealDataSource) GetData(row, col int) templ.Component {
	d := ds[row]
	switch col {
	case 0:
		return ui.Text(d.Title)
	case 1:
		return ui.Text(d.Client)
	case 2:
		return ui.Text(d.Amount)
	case 3:
		return ui.Badge(&ui.BadgeProps{Name: dealStageName(d.Stage), Style: dealStageStyle(d.Stage)})
	case 4:
		return ui.Text(d.Date)
	}
	return ui.Text("")
}
func (ds DealDataSource) GetRowMeta(row int) *ui.RowMeta {
	return &ui.RowMeta{HRef: "#"}
}

func dealStageName(s string) string {
	switch s {
	case "new":
		return "Новая"
	case "negotiation":
		return "Переговоры"
	case "proposal":
		return "КП отправлен"
	case "closed_won":
		return "Выиграна"
	case "closed_lost":
		return "Проиграна"
	default:
		return s
	}
}

func dealStageStyle(s string) string {
	switch s {
	case "new":
		return "blue"
	case "negotiation":
		return "yellow"
	case "proposal":
		return "indigo"
	case "closed_won":
		return "green"
	case "closed_lost":
		return "red"
	default:
		return "gray"
	}
}

// ===== Tasks =====

type Task struct {
	ID       int
	Title    string
	ClientID int
	Client   string
	Priority string // high, medium, low
	Status   string // todo, in_progress, done
	DueDate  string
}

var tasks = []Task{
	{1, "Подготовить КП для ООО Ромашка", 1, "ООО Ромашка", "high", "todo", "2024-11-20"},
	{2, "Звонок клиенту ЗАО Альфа", 2, "ЗАО Альфа", "high", "in_progress", "2024-11-18"},
	{3, "Отправить договор ИП Соколов", 9, "ИП Соколов", "medium", "todo", "2024-11-22"},
	{4, "Презентация демо ООО Бета", 4, "ООО Бета", "medium", "done", "2024-11-15"},
	{5, "Согласование ТЗ ПАО Тета", 11, "ПАО Тета", "high", "todo", "2024-11-25"},
	{6, "Завершить аудит ООО Дзета", 8, "ООО Дзета", "low", "in_progress", "2024-11-30"},
	{7, "Обновить презентацию компании", 0, "", "low", "todo", "2024-12-01"},
	{8, "Подготовить отчёт за Q4", 0, "", "medium", "todo", "2024-12-10"},
}

// ===== Timeline =====

type TimelineActivity struct {
	Title   string
	Date    string
	Content string
	Color   string
}

var clientTimeline = []TimelineActivity{
	{Title: "Созвон по проекту", Date: "15.11.2024", Content: "Обсудили технические требования к внедрению CRM", Color: "blue"},
	{Title: "Отправлено КП", Date: "10.11.2024", Content: "Коммерческое предложение на сумму 1 500 000 ₽", Color: "green"},
	{Title: "Первичный контакт", Date: "01.11.2024", Content: "Заполнена заявка на сайте, обратный звонок", Color: "yellow"},
	{Title: "Заключён договор", Date: "20.10.2024", Content: "Подписан договор на техподдержку 360 000 ₽/год", Color: "green"},
}

// ===== Helpers =====

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	s := ""
	n := i
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
