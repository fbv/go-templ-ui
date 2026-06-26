# Go Templ UI

A comprehensive collection of reusable UI components built with [Templ](https://templ.guide/) and styled with [Tailwind CSS](https://tailwindcss.com/). This library provides a modern, accessible, and customizable component system for Go Web applications.

## ✨ Features

- 🎨 **Modern Design**: Clean, professional components styled with Tailwind CSS
- ⚡ **High Performance**: Built with Templ for optimal server-side rendering
- 🔧 **Highly Customizable**: Flexible props and styling options
- 📱 **Responsive**: Mobile-first design approach
- ♿ **Accessible**: Built with accessibility best practices
- 🧪 **Well Tested**: Comprehensive test coverage
- 📦 **Easy Integration**: Simple installation and usage

## 🚀 Quick Start

### Installation

```bash
go get github.com/fbv/go-templ-ui
```

### Basic Usage

```go
package main

import (
    "log"
    "net/http"
    
    "github.com/fbv/go-templ-ui/view/ui"
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()
    
    // Simple page with a form
    app.Get("/", func(c *fiber.Ctx) error {
        form := ui.Form(&ui.FormProps{
            Action: "/submit",
            Method: "POST",
        }, func() templ.Component {
            return ui.Input(&ui.InputProps{
                Label: "Email",
                Name:  "email",
                Type:  "email",
                Required: true,
            })
        })
        
        return form.Render(c.Context(), c.Response().BodyWriter())
    })
    
    // Handle form submission
    app.Post("/submit", func(c *fiber.Ctx) error {
        email := c.FormValue("email")
        log.Printf("Received email: %s", email)
        return c.SendString("Form submitted successfully!")
    })
    
    log.Fatal(app.Listen(":3000"))
}
```

## 🎯 Available Components

### Layout Components
- **Card** - Container with shadow and padding
- **Panel** - Content panel with optional header
- **MasterDetails** - Two-column layout for master-detail views
- **Columns** - Flexible column layout system

### Form Components
- **Form** - Complete form wrapper with validation
- **GridForm** - Grid-based form layout
- **Input** - Text, email, password, and other input types
- **Select** - Dropdown selection component
- **Switch** - Toggle switch component
- **Checkbox** - Checkbox input
- **DatePicker** - Date and time picker with multiple variants

### Navigation Components
- **Sidebar** - Navigation sidebar with user info
- **Tabs** - Tab navigation component
- **Nav** - Navigation menu
- **Breadcrumb** - Breadcrumb navigation
- **Pagination** - Page navigation for lists and tables

### Data Display
- **Table** - Data table with sorting and pagination
- **List** - Simple and advanced list components
- **Badge** - Status and category badges
- **Collapsible** - Expandable content sections
- **Skeleton** - Placeholder content while loading

### Feedback Components
- **Alert** - Banner notifications (info, success, warning, error)
- **Toast** - Floating notifications with auto-dismiss
- **Spinner** - Loading indicators with multiple sizes and styles

### Interactive Components
- **Modal** - Popup modal dialogs
- **ConfirmDialog** - Confirmation dialog for dangerous actions
- **Button** - Various button styles and states
- **ButtonGroup** - Grouped button controls
- **Link** - Styled link component

### Utility Components
- **Headers** - H1, H2, H3 styled headers
- **YesNo** - Boolean display component
- **Block** - Content block wrapper

## 🎨 Showcase Application

The project includes a comprehensive showcase application that demonstrates all components in action.

### Running the Showcase

```bash
# Clone the repository
git clone https://github.com/fbv/go-templ-ui.git
cd go-templ-ui

# Install dependencies
go mod download

# Generate Templ components
make templ

# Build and run the showcase
make showcase
```

The showcase will be available at `http://localhost:8080`

## 🛠️ Development

### Prerequisites

- Go 1.24 or later
- [Templ](https://templ.guide/) for component generation
- [Tailwind CSS](https://tailwindcss.com/) for styling

### Development Setup

```bash
# Install Templ
go install github.com/a-h/templ/cmd/templ@latest

# Install Tailwind CSS (if not already installed)
make install_tailwindcss

# Generate components
make templ

# Watch for changes and rebuild
make tailwind-watch
```

### Build Commands

```bash
# Generate Templ components
make templ

# Build all projects
make build

# Run tests
make test

# Clean build artifacts
make clean

# Watch Tailwind CSS changes
make tailwind-watch
```

## 📖 Documentation

Detailed documentation for each component is available in the [docs](docs/) directory:

| Component | Documentation |
|-----------|---------------|
| Alert | [docs/alert.md](docs/alert.md) |
| Toast | [docs/toast.md](docs/toast.md) |
| Spinner | [docs/spinner.md](docs/spinner.md) |
| Skeleton | [docs/skeleton.md](docs/skeleton.md) |
| Pagination | [docs/pagination.md](docs/pagination.md) |
| DatePicker | [docs/datepicker.md](docs/datepicker.md) |
| ConfirmDialog | [docs/confirm_dialog.md](docs/confirm_dialog.md) |

## 📖 Component Examples

### Form with Validation

```go
// In your route handler
app.Get("/form", func(c *fiber.Ctx) error {
    form := ui.Form(&ui.FormProps{
        Action: "/submit",
        Method: "POST",
    }, func() templ.Component {
        return ui.Input(&ui.InputProps{
            Label: "Email Address",
            Name: "email",
            Type: "email",
            Placeholder: "Enter your email",
            Required: true,
        })
    })
    
    return form.Render(c.Context(), c.Response().BodyWriter())
})
```

### Sidebar Navigation

```go
// In your route handler
app.Get("/dashboard", func(c *fiber.Ctx) error {
    sidebar := ui.Sidebar(&ui.SidebarProps{
        Items: []*ui.SidebarItem{
            {
                Name: "Dashboard",
                Icon: icon.Dashboard,
                URL: "/dashboard",
                Active: true,
            },
            {
                Name: "Settings",
                Icon: icon.Settings,
                URL: "/settings",
            },
        },
        User: &ui.UserInfo{
            Name: "John Doe",
            Email: "john@example.com",
        },
        LogoutURL: "/logout",
    })
    
    return sidebar.Render(c.Context(), c.Response().BodyWriter())
})
```

### Data Table

```go
// In your route handler
app.Get("/users", func(c *fiber.Ctx) error {
    table := ui.Table(&ui.TableProps{
        Headers: []string{"Name", "Email", "Status"},
        Rows: [][]string{
            {"John Doe", "john@example.com", "Active"},
            {"Jane Smith", "jane@example.com", "Inactive"},
        },
    })
    
    return table.Render(c.Context(), c.Response().BodyWriter())
})
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- [Templ Documentation](https://templ.guide/)
- [Tailwind CSS Documentation](https://tailwindcss.com/)
- [Go Fiber Framework](https://gofiber.io/)

## 📞 Support

If you have any questions or need help, please open an issue on GitHub.
