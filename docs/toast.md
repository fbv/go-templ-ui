# Toast Component

Floating notification messages that appear and auto-dismiss after a duration.

## Import

```go
import "github.com/fbv/go-templ-ui/view/ui"
```

## Setup

Add `ToastScript()` to your page layout to enable toast functionality:

```go
templ Page() {
    // ... head content ...
    <body>
        { children... }
        @ui.ToastScript()
    </body>
}
```

## Props

| Property | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `ID` | `string` | Yes | - | Unique identifier for the toast element |
| `Title` | `string` | No | - | Bold title displayed at the top |
| `Style` | `string` | No | `"info"` | Toast type: `"info"`, `"success"`, `"warning"`, `"error"` |
| `Position` | `string` | No | `"bottom-right"` | Screen position (see positions below) |
| `Duration` | `int` | No | `3000` | Auto-hide duration in milliseconds |
| `Attrs` | `map[string]string` | No | - | Additional HTML attributes |

## Positions

| Position | Description |
|----------|-------------|
| `"top-left"` | Top left corner |
| `"top-center"` | Top center |
| `"top-right"` | Top right corner |
| `"bottom-left"` | Bottom left corner |
| `"bottom-center"` | Bottom center |
| `"bottom-right"` | Bottom right corner (default) |

## Usage

### Basic Toast

```go
// Place the toast in your template (hidden by default)
@ui.Toast(&ui.ToastProps{
    ID:    "my-toast",
    Title: "Info",
    Style: "info",
}) {
    This is an informational message.
}

// Show it with JavaScript
<button onclick="showToast('my-toast')">Show Toast</button>
```

### Different Styles

```go
// Success toast
@ui.Toast(&ui.ToastProps{
    ID:        "success-toast",
    Title:     "Success",
    Style:     "success",
    Position:  "top-right",
    Duration:  5000,
}) {
    Operation completed successfully!
}

// Error toast
@ui.Toast(&ui.ToastProps{
    ID:       "error-toast",
    Title:    "Error",
    Style:    "error",
    Position: "top-left",
    Duration: 10000,
}) {
    An error occurred. Please try again.
}
```

### Programmatic Show/Hide

```go
// Show a toast
<button onclick="showToast('my-toast')">Show</button>

// Hide a toast
<button onclick="hideToast('my-toast')">Hide</button>
```

## JavaScript API

### `showToast(id)`

Shows the toast with the specified ID. Starts the auto-hide timer.

```javascript
showToast('my-toast');
```

### `hideToast(id)`

Immediately hides the toast with the specified ID.

```javascript
hideToast('my-toast');
```

## Styling

Each toast includes:
- An icon matching the style (info/success/warning/error)
- A progress bar showing remaining time
- A close button

## Dark Mode

Toasts support dark mode automatically with appropriate background and text colors.

## Accessibility

- Uses `role="alert"` for screen reader announcements
- Close button includes `aria-label="Close"`
- SVG icons are decorative with `aria-hidden="true"`
