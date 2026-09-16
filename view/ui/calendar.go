package ui

import (
	"fmt"
	"time"
)

type CalendarEvent struct {
	ID      string
	Title   string
	Start   time.Time
	End     time.Time
	Color   string
	AllDay  bool
}

func calendarEventColor(color string) string {
	switch color {
	case "gray":
		return "bg-gray-100 text-gray-800 border-gray-300 dark:bg-gray-700 dark:text-gray-300 dark:border-gray-600"
	case "red":
		return "bg-red-100 text-red-800 border-red-300 dark:bg-red-900 dark:text-red-300 dark:border-red-700"
	case "green":
		return "bg-green-100 text-green-800 border-green-300 dark:bg-green-900 dark:text-green-300 dark:border-green-700"
	case "yellow":
		return "bg-yellow-100 text-yellow-800 border-yellow-300 dark:bg-yellow-900 dark:text-yellow-300 dark:border-yellow-700"
	case "purple":
		return "bg-purple-100 text-purple-800 border-purple-300 dark:bg-purple-900 dark:text-purple-300 dark:border-purple-700"
	default:
		return "bg-blue-100 text-blue-800 border-blue-300 dark:bg-blue-900 dark:text-blue-300 dark:border-blue-700"
	}
}

func calendarEventDotColor(color string) string {
	switch color {
	case "gray":
		return "bg-gray-400"
	case "red":
		return "bg-red-500"
	case "green":
		return "bg-green-500"
	case "yellow":
		return "bg-yellow-500"
	case "purple":
		return "bg-purple-500"
	default:
		return "bg-blue-500"
	}
}

func calendarEventPosition(start, end time.Time, hourHeight int) string {
	startMinutes := start.Hour()*60 + start.Minute()
	endMinutes := end.Hour()*60 + end.Minute()
	if endMinutes <= startMinutes {
		endMinutes = startMinutes + 30
	}
	top := startMinutes * hourHeight / 60
	height := (endMinutes - startMinutes) * hourHeight / 60
	if height < 20 {
		height = 20
	}
	return fmt.Sprintf("top: %dpx; height: %dpx;", top, height)
}

func isSameDay(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}

func formatTime(t time.Time) string {
	return fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
}

func formatDate(t time.Time) string {
	return fmt.Sprintf("%02d.%02d.%d", t.Day(), t.Month(), t.Year())
}

func formatMonthYear(t time.Time) string {
	months := []string{
		"Январь", "Февраль", "Март", "Апрель", "Май", "Июнь",
		"Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь",
	}
	return fmt.Sprintf("%s %d", months[t.Month()-1], t.Year())
}

func dayOfWeekHeaders() []string {
	return []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}
}

func daysInMonth(year int, month time.Month) []time.Time {
	first := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	last := first.AddDate(0, 1, -1)
	days := make([]time.Time, 0, last.Day())
	for d := first; !d.After(last); d = d.AddDate(0, 0, 1) {
		days = append(days, d)
	}
	return days
}

func calendarGrid(year int, month time.Month) [][]time.Time {
	first := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	weekday := first.Weekday()
	// Monday=0 offset
	offset := (int(weekday) + 6) % 7

	last := first.AddDate(0, 1, -1)
	totalCells := offset + last.Day()
	rows := (totalCells + 6) / 7

	grid := make([]time.Time, rows*7)
	empty := time.Time{}
	for i := range grid {
		grid[i] = empty
	}
	for d := 0; d < last.Day(); d++ {
		grid[offset+d] = time.Date(year, month, d+1, 0, 0, 0, 0, time.UTC)
	}

	result := make([][]time.Time, rows)
	for r := 0; r < rows; r++ {
		result[r] = grid[r*7 : (r+1)*7]
	}
	return result
}

func weekStart(date time.Time) time.Time {
	weekday := date.Weekday()
	offset := (int(weekday) + 6) % 7
	return date.AddDate(0, 0, -offset)
}

func hourSlots() []int {
	slots := make([]int, 24)
	for i := 0; i < 24; i++ {
		slots[i] = i
	}
	return slots
}

func eventsForDay(events []CalendarEvent, day time.Time) []CalendarEvent {
	var result []CalendarEvent
	for _, e := range events {
		if isSameDay(e.Start, day) || isSameDay(e.End, day) ||
			(e.Start.Before(day) && e.End.After(day)) {
			result = append(result, e)
		}
	}
	return result
}

func allDayEvents(events []CalendarEvent) []CalendarEvent {
	var result []CalendarEvent
	for _, e := range events {
		if e.AllDay {
			result = append(result, e)
		}
	}
	return result
}

func timedEvents(events []CalendarEvent) []CalendarEvent {
	var result []CalendarEvent
	for _, e := range events {
		if !e.AllDay {
			result = append(result, e)
		}
	}
	return result
}
