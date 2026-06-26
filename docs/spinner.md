# Spinner Component

Animated loading indicator for async operations.

## Import

```go
import "github.com/fbv/go-templ-ui/view/ui"
```

## Props

| Property | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `Size` | `string` | No | `"md"` | Spinner size: `"sm"`, `"md"`, `"lg"`, `"xl"` |
| `Style` | `string` | No | `"blue"` | Color: `"blue"`, `"white"`, `"gray"`, `"success"`, `"error"` |
| `Label` | `string` | No | - | Optional text label next to the spinner |
| `Attrs` | `map[string]string` | No | - | Additional HTML attributes |

## Sizes

| Size | Dimensions | Use Case |
|------|-----------|----------|
| `"sm"` | 16x16 px | Inline buttons, small spaces |
| `"md"` | 32x32 px | Default, general use |
| `"lg"` | 40x40 px | Page sections |
| `"xl"` | 64x64 px | Full-page loading |

## Usage

### Basic Spinner

```go
@ui.Spinner(&ui.SpinnerProps{})
```

### Different Sizes

```go
@ui.Spinner(&ui.SpinnerProps{Size: "sm"})
@ui.Spinner(&ui.SpinnerProps{Size: "md"})
@ui.Spinner(&ui.SpinnerProps{Size: "lg"})
@ui.Spinner(&ui.SpinnerProps{Size: "xl"})
```

### With Label

```go
@ui.Spinner(&ui.SpinnerProps{
    Label: "Loading...",
})
```

### White Spinner (on dark backgrounds)

```go
<div class="bg-gray-900 p-4">
    @ui.Spinner(&ui.SpinnerProps{
        Size:  "lg",
        Style: "white",
    })
</div>
```

### Success/Error Spinners

```go
@ui.Spinner(&ui.SpinnerProps{
    Size:  "md",
    Style: "success",
    Label: "Processing...",
})

@ui.Spinner(&ui.SpinnerProps{
    Size:  "md",
    Style: "error",
    Label: "Failed",
})
```

### Full-Page Overlay

```go
@ui.SpinnerOverlay()
```

This renders a centered spinner with a semi-transparent white background covering the entire viewport.

### In a Button

```go
<button class="px-4 py-2 bg-blue-600 text-white rounded" disabled>
    @ui.Spinner(&ui.SpinnerProps{Size: "sm", Style: "white"})
    <span class="ms-2">Saving...</span>
</button>
```

## SpinnerOverlay

A full-page loading overlay that centers a large spinner:

```go
@ui.SpinnerOverlay()
```

**Features:**
- Fixed positioning covering entire viewport
- Semi-transparent white background (`bg-white/75`)
- Dark mode support (`dark:bg-gray-900/75`)
- Large spinner (xl size)
- z-index of 50

## Styles

| Style | Color | Use Case |
|-------|-------|----------|
| `"blue"` | Blue | Default, general loading |
| `"white"` | White | Dark backgrounds |
| `"gray"` | Light gray | Subtle loading |
| `"success"` | Green | Processing/success states |
| `"error"` | Red | Error states |

## Animation

The spinner uses Tailwind's `animate-spin` class for a smooth rotation animation.

## Dark Mode

Supports dark mode with appropriate color adjustments for all styles.

## Accessibility

- Spinner is purely visual; add a label for screen readers
- Use `aria-label` or visible `Label` prop for accessibility
- Consider adding `role="status"` and `aria-live="polite"` for dynamic content
