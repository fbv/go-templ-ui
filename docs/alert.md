# Alert Component

Display important messages to users with different severity levels.

## Import

```go
import "github.com/fbv/go-templ-ui/view/ui"
```

## Props

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `Title` | `string` | No | Bold title displayed before the message |
| `Style` | `string` | No | Alert type: `"info"` (default), `"success"`, `"warning"`, `"error"` |
| `Dismissible` | `bool` | No | Show a close button (requires additional JS for hiding) |
| `Attrs` | `map[string]string` | No | Additional HTML attributes |

## Usage

### Basic Alert

```go
@ui.Alert(&ui.AlertProps{
    Title: "Info",
    Style: "info",
}) {
    This is an informational message.
}
```

### Success Alert

```go
@ui.Alert(&ui.AlertProps{
    Title: "Success!",
    Style: "success",
}) {
    Your changes have been saved successfully.
}
```

### Warning Alert

```go
@ui.Alert(&ui.AlertProps{
    Title: "Warning",
    Style: "warning",
}) {
    Your session is about to expire. Please save your work.
}
```

### Error Alert

```go
@ui.Alert(&ui.AlertProps{
    Title: "Error",
    Style: "error",
}) {
    Something went wrong. Please try again later.
}
```

### With Children

```go
@ui.Alert(&ui.AlertProps{Style: "info"}) {
    <p>Line 1 of the message.</p>
    <p>Line 2 of the message.</p>
}
```

### With Custom Attributes

```go
@ui.Alert(&ui.AlertProps{
    Title: "Notice",
    Attrs: map[string]string{
        "id":    "my-alert",
        "class": "mb-8",
    },
}) {
    Customized alert content.
}
```

## Styles

| Style | Color | Use Case |
|-------|-------|----------|
| `"info"` | Blue | General information |
| `"success"` | Green | Successful operations |
| `"warning"` | Yellow | Warnings and cautions |
| `"error"` | Red | Errors and failures |

## Dark Mode

The Alert component supports dark mode automatically with appropriate color adjustments.

## Accessibility

- Uses `role="alert"` for screen reader announcements
- Dismissible alerts include `aria-label="Close"` on the close button
- SVG icons are hidden from screen readers with `aria-hidden="true"`
