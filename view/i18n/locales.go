package i18n

// The built-in locales. English is Default() and its wording is the
// EXPERIENCE.md default set (sentence case, terse). Russian is the second
// built-in locale — the one the showcase wires up — and doubles as the template
// for a consumer's own locale.
var (
	// English is the built-in default locale.
	English = &Locale{
		Tag:  "en",
		Name: "English",
		Strings: map[string]string{
			KeyYes:           "Yes",
			KeyNo:            "No",
			KeyClose:         "Close",
			KeyLoading:       "Loading…",
			KeyToday:         "Today",
			KeyAllDay:        "All day",
			KeyNoEvents:      "No events",
			KeyNoData:        "No data",
			KeyCancel:        "Cancel",
			KeyConfirm:       "Confirm",
			KeySubmit:        "Submit",
			KeySearch:        "Search",
			KeyTheme:         "Theme",
			KeyCollapse:      "Collapse",
			KeyNew:           "New",
			KeyLogout:        "Log out",
			KeySelectedCount: "%d selected",
			KeyPageOf:        "Page %d of %d",
		},
		Months: [12]string{
			"January", "February", "March", "April", "May", "June",
			"July", "August", "September", "October", "November", "December",
		},
		Weekdays: [7]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
		Numbers: NumberFormat{
			Decimal:   ".",
			Group:     ",",
			GroupSize: 3,
			Currency: map[string]CurrencyFormat{
				"USD": {Symbol: "$", Position: "before"},
				"EUR": {Symbol: "€", Position: "before"},
				"RUB": {Symbol: "₽", Position: "before"},
			},
		},
		DateLayout:     "01/02/2006",
		DateTimeLayout: "01/02/2006 15:04",
		TimeLayout:     "15:04",
	}

	// Russian is the built-in non-default locale.
	Russian = &Locale{
		Tag:  "ru",
		Name: "Русский",
		Strings: map[string]string{
			KeyYes:           "Да",
			KeyNo:            "Нет",
			KeyClose:         "Закрыть",
			KeyLoading:       "Загрузка…",
			KeyToday:         "Сегодня",
			KeyAllDay:        "Весь день",
			KeyNoEvents:      "Нет событий",
			KeyNoData:        "Нет данных",
			KeyCancel:        "Отмена",
			KeyConfirm:       "Подтвердить",
			KeySubmit:        "Отправить",
			KeySearch:        "Поиск",
			KeyTheme:         "Тема",
			KeyCollapse:      "Свернуть",
			KeyNew:           "Новый",
			KeyLogout:        "Выйти",
			KeySelectedCount: "Выбрано: %d",
			KeyPageOf:        "Страница %d из %d",
		},
		Months: [12]string{
			"Январь", "Февраль", "Март", "Апрель", "Май", "Июнь",
			"Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь",
		},
		Weekdays: [7]string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"},
		Numbers: NumberFormat{
			Decimal:   ",",
			Group:     " ",
			GroupSize: 3,
			Currency: map[string]CurrencyFormat{
				"USD": {Symbol: "$", Position: "before"},
				"EUR": {Symbol: "€", Position: "before"},
				"RUB": {Symbol: "₽", Position: "after", Space: true},
			},
		},
		DateLayout:     "02.01.2006",
		DateTimeLayout: "02.01.2006 15:04",
		TimeLayout:     "15:04",
	}
)

func init() {
	mustRegister(English)
	mustRegister(Russian)
}
