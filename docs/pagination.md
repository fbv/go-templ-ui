# Pagination Component

Navigate through paginated data with smart page number rendering.

## Import

```go
import "github.com/fbv/go-templ-ui/view/ui"
```

## Props

| Property | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `CurrentPage` | `int` | Yes | - | The currently active page (1-indexed) |
| `TotalPages` | `int` | Yes | - | Total number of pages |
| `BaseURL` | `string` | Yes | - | Base URL for page links (page number appended as `?page=N`) |
| `Style` | `string` | No | - | Reserved for future style variants |
| `ShowFirst` | `bool` | No | `false` | Show "first page" button (`<<`) |
| `ShowLast` | `bool` | No | `false` | Show "last page" button (`>>`) |
| `Attrs` | `map[string]string` | No | - | Additional HTML attributes |

## Usage

### Basic Pagination

```go
@ui.Pagination(&ui.PaginationProps{
    CurrentPage: 1,
    TotalPages:  10,
    BaseURL:     "/users",
})
// Renders: 1 2 3 4 5 6 7 8 9 10
```

### With First/Last Buttons

```go
@ui.Pagination(&ui.PaginationProps{
    CurrentPage: 5,
    TotalPages:  20,
    BaseURL:     "/orders",
    ShowFirst:   true,
    ShowLast:    true,
})
// Renders: << « 4 5 6 » >>
```

### In a List Context

```go
@ui.Pagination(&ui.PaginationProps{
    CurrentPage: 3,
    TotalPages:  15,
    BaseURL:     "/products?page",
    ShowFirst:   true,
    ShowLast:    true,
})
```

## How It Works

### Smart Page Range

When there are more than 7 pages, the pagination shows:

```
« 1 ... 4 5 6 ... 15 »
```

- Always shows first and last page
- Shows current page and neighbors
- Uses `...` for omitted pages

### URL Format

Page links are generated as `{BaseURL}?page={N}`:

| BaseURL | Page | Result |
|---------|------|--------|
| `/users` | 1 | `/users?page=1` |
| `/orders` | 5 | `/orders?page=5` |
| `/products?page` | 3 | `/products?page?page=3` |

### When Not to Render

If `TotalPages` is 1 or less, the component renders nothing.

## Server-Side Integration

```go
app.Get("/users", func(c *fiber.Ctx) error {
    page := c.QueryInt("page", 1)
    pageSize := 20
    totalItems := getUserCount()
    totalPages := (totalItems + pageSize - 1) / pageSize
    
    users := getUsers(page, pageSize)
    
    pagination := ui.Pagination(&ui.PaginationProps{
        CurrentPage: page,
        TotalPages:  totalPages,
        BaseURL:     "/users",
        ShowFirst:   true,
        ShowLast:    true,
    })
    
    // Render page with users and pagination
    return render(c, usersPage(users, pagination))
})
```

## Dark Mode

Pagination supports dark mode with appropriate border, background, and text colors for both active and inactive states.

## Accessibility

- Uses `<nav>` element for navigation landmark
- Active page is visually distinguished but not marked with `aria-current` (could be added as enhancement)
- Links are keyboard accessible
