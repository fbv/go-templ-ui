# Sidebar Component

Navigation sidebar with support for collapsible mode showing only icons.

## Import

```go
import "github.com/fbv/go-templ-ui/view/ui"
```

## Setup

Add `SidebarStyles()` and `SidebarScript()` to your page layout:

```go
templ Page() {
    <head>
        @ui.SidebarStyles()
    </head>
    <body>
        { children... }
        @ui.SidebarScript()
    </body>
}
```

## Props

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `Items` | `[]*SidebarItem` | Yes | Navigation items |
| `User` | `*UserInfo` | No | User info for footer |
| `LogoutURL` | `string` | No | Logout button URL |
| `Collapsible` | `bool` | No | Enable collapse/expand toggle |
| `ID` | `string` | If collapsible | Unique ID for the sidebar (required for localStorage) |

### SidebarItem Props

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `Name` | `string` | Yes | Display text |
| `Icon` | `string` | No | SVG path data for Heroicons-style icons |
| `IconRaw` | `string` | No | Raw SVG tag (alternative to Icon) |
| `URL` | `string` | Yes | Link destination |
| `Active` | `bool` | No | Highlight as active item |
| `Modifiers` | `[]templ.Component` | No | Badge components (SidebarNew, SidebarRed, etc.) |
| `Tooltip` | `string` | No | Custom tooltip (defaults to Name) |

### UserInfo Props

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `Name` | `string` | Yes | User display name |
| `Email` | `string` | Yes | User email |
| `Avatar` | `string` | No | Avatar image URL |

## Usage

### Basic Sidebar

```go
@ui.Sidebar(&ui.SidebarProps{
    Items: []*ui.SidebarItem{
        {Name: "Dashboard", Icon: icon.Dashboard, URL: "/dashboard", Active: true},
        {Name: "Users", Icon: icon.Clients, URL: "/users"},
        {Name: "Settings", Icon: icon.Settings, URL: "/settings"},
    },
    User: &ui.UserInfo{
        Name:  "John Doe",
        Email: "john@example.com",
    },
    LogoutURL: "/logout",
})
```

### Collapsible Sidebar

```go
@ui.Sidebar(&ui.SidebarProps{
    Items: []*ui.SidebarItem{
        {Name: "Dashboard", Icon: icon.Dashboard, URL: "/dashboard"},
        {Name: "Users", Icon: icon.Clients, URL: "/users"},
        {Name: "Settings", Icon: icon.Settings, URL: "/settings"},
    },
    User: &ui.UserInfo{
        Name:  "John Doe",
        Email: "john@example.com",
    },
    LogoutURL:   "/logout",
    Collapsible: true,
    ID:          "main-sidebar",
})
```

### With Modifiers/Badges

```go
@ui.Sidebar(&ui.SidebarProps{
    Items: []*ui.SidebarItem{
        {Name: "Dashboard", Icon: icon.Dashboard, URL: "/dashboard"},
        {Name: "Messages", Icon: icon.Messages, URL: "/messages",
            Modifiers: []templ.Component{ui.SidebarNew()}},
        {Name: "Tasks", Icon: icon.Tasks, URL: "/tasks",
            Modifiers: []templ.Component{ui SidebarRed("5")}},
    },
})
```

## Collapsible Mode

When `Collapsible: true`, the sidebar features:

- **Toggle button** at the top to expand/collapse
- **Compact view** showing only icons (width shrinks to 4rem)
- **Tooltips** on hover showing item names
- **localStorage persistence** - remembers state across page reloads
- **Smooth animation** with CSS transitions

### Behavior

- Text labels hide in collapsed mode
- User info hides in collapsed mode
- Logout button shows only icon
- Modifiers/badges hide in collapsed mode
- Arrow icon indicates expand/collapse direction

### Toggle Functionality

The sidebar uses JavaScript to toggle classes:

```javascript
// Toggle sidebar
toggleSidebar('main-sidebar');

// Initialize from localStorage
initSidebar('main-sidebar');
```

## Modifiers

### SidebarNew

```go
@ui.SidebarNew()
// Renders: "New" badge
```

### SidebarRed

```go
@ui.SidebarRed("5")
// Renders: Red badge with "5"
```

### SidebarGreen

```go
@ui.sidebarGreen("New")
// Renders: Green badge with "New"
```

## Icons

Use Heroicons SVG path data:

```go
import "github.com/fbv/go-templ-ui/view/icon"

&ui.SidebarItem{
    Name: "Dashboard",
    Icon: icon.Dashboard,  // SVG path data
    URL:  "/dashboard",
}
```

Or use raw SVG:

```go
&ui.SidebarItem{
    Name:    "Custom",
    IconRaw: `<svg class="w-5 h-5" ...>...</svg>`,
    URL:     "/custom",
}
```

## Dark Mode

The sidebar supports dark mode with appropriate color adjustments when used with Tailwind's dark mode classes.

## Accessibility

- Uses semantic `<nav>` element (consider wrapping in `<nav>` in your layout)
- Active state is visually distinguished
- Toggle button has `aria-label="Toggle sidebar"`
- Links are keyboard accessible
- Focus states are visible

## Layout Integration

```go
templ Dashboard() {
    <div class="flex h-screen">
        @ui.Sidebar(&ui.SidebarProps{
            Items:       items,
            User:        user,
            LogoutURL:   "/logout",
            Collapsible: true,
            ID:          "dashboard-sidebar",
        })
        <main class="flex-1 overflow-auto p-6">
            { children... }
        </main>
    </div>
}
```

## CSS Classes

| Class | Description |
|-------|-------------|
| `sidebar-root` | Base sidebar container |
| `sidebar-collapsed` | Collapsed state modifier |
| `sidebar-item` | Navigation item link |
| `sidebar-item-icon` | Icon container |
| `sidebar-item-text` | Text label |
| `sidebar-item-modifier` | Badge container |
| `sidebar-footer` | User info section |
| `sidebar-user-info` | User name/email |
| `sidebar-logout` | Logout button |
| `sidebar-logout-text` | Logout text |
| `sidebar-icon-expand` | Expand arrow icon |
| `sidebar-icon-collapse` | Collapse arrow icon |
