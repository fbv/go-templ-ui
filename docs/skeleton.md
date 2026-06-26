# Skeleton Component

Placeholder content that displays while data is loading.

## Import

```go
import "github.com/fbv/go-templ-ui/view/ui"
```

## Props (Skeleton)

| Property | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `Lines` | `int` | No | `0` | Number of text lines to show (0 = single block) |
| `Style` | `string` | No | `"default"` | Preset: `"default"`, `"avatar"`, `"card"`, `"button"`, `"title"` |
| `Width` | `string` | No | Auto | Custom width (CSS value) |
| `Height` | `string` | No | Auto | Custom height (CSS value) |
| `Attrs` | `map[string]string` | No | - | Additional HTML attributes |

## Preset Components

| Component | Description |
|-----------|-------------|
| `SkeletonCard()` | Card-shaped skeleton with title and text lines |
| `SkeletonTable(rows, cols)` | Table with specified rows and columns |
| `SkeletonAvatar()` | User avatar with name/email lines |

## Usage

### Basic Skeleton

```go
// Single block
@ui.Skeleton(&ui.SkeletonProps{})

// Multiple lines
@ui.Skeleton(&ui.SkeletonProps{Lines: 3})
```

### Custom Size

```go
// Custom width and height
@ui.Skeleton(&ui.SkeletonProps{
    Width:  "200px",
    Height: "24px",
})
```

### Title Skeleton

```go
@ui.Skeleton(&ui.SkeletonProps{
    Style: "title",
    Width: "50%",
})
```

### Avatar Skeleton

```go
// Preset avatar skeleton
@ui.SkeletonAvatar()

// Or custom
@ui.Skeleton(&ui.SkeletonProps{
    Style:  "avatar",
    Width:  "3rem",
    Height: "3rem",
})
```

### Card Skeleton

```go
// Preset card skeleton
@ui.SkeletonCard()
```

### Table Skeleton

```go
// 5 rows, 4 columns
@ui.SkeletonTable(5, 4)
```

### Combined Example

```go
templ LoadingPage() {
    <div class="space-y-6">
        // Title area
        @ui.Skeleton(&ui.SkeletonProps{
            Style: "title",
            Width: "40%",
        })
        
        // Avatar + info
        @ui.SkeletonAvatar()
        
        // Content lines
        @ui.Skeleton(&ui.SkeletonProps{Lines: 5})
        
        // Card content
        @ui.SkeletonCard()
        
        // Table
        @ui.SkeletonTable(8, 5)
    </div>
}
```

### With Conditional Rendering

```go
templ UserList(users []User, loading bool) {
    if loading {
        @ui.SkeletonTable(len(users), 3)
    } else {
        <table>
            for _, user := range users {
                <tr>
                    <td>{ user.Name }</td>
                    <td>{ user.Email }</td>
                </tr>
            }
        </table>
    }
}
```

## Preset Styles

### SkeletonCard

Renders a card-shaped skeleton with:
- Title line (3/4 width)
- 3 text lines of varying widths

```go
@ui.SkeletonCard()
```

### SkeletonTable

Renders a table skeleton with specified dimensions:

```go
// rows=5, cols=4
@ui.SkeletonTable(5, 4)
```

### SkeletonAvatar

Renders a user profile skeleton with:
- Circular avatar placeholder
- Two text lines (name and email)

```go
@ui.SkeletonAvatar()
```

## Style Presets

| Style | Default Width | Default Height | Description |
|-------|--------------|----------------|-------------|
| `"default"` | 100% | 1rem | Basic text line |
| `"avatar"` | 3rem | 3rem | Circular avatar |
| `"card"` | 100% | 10rem | Card container |
| `"button"` | 6rem | 2.5rem | Button shape |
| `"title"` | 50% | 1rem | Section title |

## Animation

All skeletons use Tailwind's `animate-pulse` class for a subtle pulsing animation that indicates loading.

## Dark Mode

Supports dark mode with `dark:bg-gray-700` for skeleton backgrounds.

## Accessibility

- Skeletons are decorative; consider adding `aria-hidden="true"` if used alongside loading text
- Provide a visible loading indicator for screen readers
- Use `role="status"` with `aria-live="polite"` for dynamic loading states
