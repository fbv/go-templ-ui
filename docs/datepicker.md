# DatePicker Component

Native HTML date and time input fields with consistent styling.

## Import

```go
import "github.com/fbv/go-templ-ui/view/ui"
```

## Props

| Property | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `Name` | `string` | Yes | - | HTML name attribute for form submission |
| `Label` | `string` | No | - | Label text displayed above the input |
| `Value` | `string` | No | - | Pre-filled value |
| `Placeholder` | `string` | No | Auto | Placeholder text (auto-generated if empty) |
| `Required` | `bool` | No | `false` | Whether the field is required |
| `Min` | `string` | No | - | Minimum allowed value |
| `Max` | `string` | No | - | Maximum allowed value |
| `Step` | `string` | No | - | Step increment |
| `Type` | `string` | No | `"date"` | Input type (see variants below) |
| `Attrs` | `map[string]string` | No | - | Additional HTML attributes |

## Components

| Component | Type | Description |
|-----------|------|-------------|
| `DatePicker` | `date` | Standard date picker |
| `TimePicker` | `time` | Time-only picker |
| `DateTimePicker` | `datetime-local` | Combined date and time |
| `DateRangePicker` | Two dates | Start and end date pair |

## Usage

### Basic Date Picker

```go
@ui.DatePicker(&ui.DatePickerProps{
    Name:  "birth_date",
    Label: "Date of Birth",
})
```

### Time Picker

```go
@ui.TimePicker(&ui.DatePickerProps{
    Name:  "meeting_time",
    Label: "Meeting Time",
})
```

### DateTime Picker

```go
@ui.DateTimePicker(&ui.DatePickerProps{
    Name:  "event_datetime",
    Label: "Event Date & Time",
})
```

### Date Range Picker

```go
@ui.DateRangePicker("start_date", "end_date", "Date Range")
// Renders two date inputs with labels "Date Range (От)" and "Date Range (До)"
```

### With Constraints

```go
@ui.DatePicker(&ui.DatePickerProps{
    Name:     "delivery_date",
    Label:    "Delivery Date",
    Required: true,
    Min:      "2024-01-01",
    Max:      "2024-12-31",
})
```

### With Custom Value

```go
@ui.DatePicker(&ui.DatePickerProps{
    Name:  "created_at",
    Label: "Created At",
    Value: "2024-06-15",
})
```

### With Step (Time)

```go
@ui.TimePicker(&ui.DatePickerProps{
    Name:  "reminder_time",
    Label: "Reminder Time",
    Step:  "15", // 15-minute intervals
})
```

## Type Variants

| Type | Input Type | Placeholder | Value Format |
|------|-----------|-------------|--------------|
| `"date"` | `date` | `YYYY-MM-DD` | `2024-06-15` |
| `"time"` | `time` | `HH:MM` | `14:30` |
| `"month"` | `month` | `YYYY-MM` | `2024-06` |
| `"week"` | `week` | `YYYY-WXX` | `2024-W25` |
| `"datetime-local"` | `datetime-local` | `YYYY-MM-DD HH:MM` | `2024-06-15T14:30` |

## Form Integration

```go
@ui.Form(&ui.FormProps{
    Action: "/events",
    Method: "POST",
}) {
    @ui.DatePicker(&ui.DatePickerProps{
        Name:     "start_date",
        Label:    "Start Date",
        Required: true,
    })
    
    @ui.TimePicker(&ui.DatePickerProps{
        Name:  "start_time",
        Label: "Start Time",
    })
    
    @ui.Submit("Create Event")
}
```

## Server-Side Reading

```go
app.Post("/events", func(c *fiber.Ctx) error {
    startDate := c.FormValue("start_date") // "2024-06-15"
    startTime := c.FormValue("start_time") // "14:30"
    
    // Parse the date
    layout := "2006-01-02"
    date, err := time.Parse(layout, startDate)
    if err != nil {
        // handle error
    }
    
    return c.SendString("Event created!")
})
```

## Dark Mode

Supports dark mode with appropriate background, border, and text colors.

## Browser Support

Date/time pickers use native HTML5 `<input type="date/time">` elements. Browser support:

- **Chrome/Edge**: Full support with custom calendar UI
- **Firefox**: Full support
- **Safari**: Full support (iOS shows native picker)
- **Mobile**: Native date/time pickers on all platforms

## Accessibility

- Uses native `<input>` elements for built-in accessibility
- Label is properly associated with input via `for` attribute
- Calendar icon is decorative with `aria-hidden="true"`
- Supports keyboard navigation
- Screen readers announce the input type and constraints
