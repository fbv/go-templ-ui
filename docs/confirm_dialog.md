# ConfirmDialog Component

Modal dialog for confirming dangerous or important actions.

## Import

```go
import "github.com/fbv/go-templ-ui/view/ui"
```

## Props

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `ModalID` | `string` | Yes | Unique identifier for the modal element |
| `Title` | `string` | Yes | Dialog title |
| `Message` | `string` | Yes | Confirmation message |
| `ConfirmLabel` | `string` | Yes | Text for the confirm button |
| `CancelLabel` | `string` | Yes | Text for the cancel button |
| `ConfirmURL` | `string` | Yes | URL to navigate to when confirmed |
| `Style` | `string` | No | Button style: `"default"` (blue) or `"danger"` (red) |

## Usage

### Basic Confirm Dialog

```go
@ui.ConfirmDialog(&ui.ConfirmDialogProps{
    ModalID:      "confirm-delete",
    Title:        "Delete Item",
    Message:      "Are you sure you want to delete this item? This action cannot be undone.",
    ConfirmLabel: "Delete",
    CancelLabel:  "Cancel",
    ConfirmURL:   "/items/123/delete",
})
```

### Danger Style (Red Button)

```go
@ui.ConfirmDialog(&ui.ConfirmDialogProps{
    ModalID:      "confirm-delete-all",
    Title:        "Delete All Items",
    Message:      "This will permanently delete all items. This action cannot be undone.",
    ConfirmLabel: "Yes, Delete All",
    CancelLabel:  "No, Keep Them",
    ConfirmURL:   "/items/delete-all",
    Style:        "danger",
})
```

### Triggering the Dialog

Use TW Elements attributes on a button to trigger the modal:

```go
<button
    data-twe-toggle="modal"
    data-twe-target="#confirm-delete"
    class="px-4 py-2 text-white bg-red-600 rounded-lg hover:bg-red-700"
>
    Delete
</button>

@ui.ConfirmDialog(&ui.ConfirmDialogProps{
    ModalID:      "confirm-delete",
    Title:        "Confirm Delete",
    Message:      "Are you sure?",
    ConfirmLabel: "Delete",
    CancelLabel:  "Cancel",
    ConfirmURL:   "/items/123/delete",
    Style:        "danger",
})
```

### Full Example with Table Row

```go
// Table with delete buttons
templ UserTable(users []User) {
    <table class="w-full">
        <thead>
            <tr>
                <th>Name</th>
                <th>Email</th>
                <th>Actions</th>
            </tr>
        </thead>
        <tbody>
            for _, user := range users {
                <tr>
                    <td>{ user.Name }</td>
                    <td>{ user.Email }</td>
                    <td>
                        <button
                            data-twe-toggle="modal"
                            data-twe-target={ "#confirm-" + user.ID }
                            class="text-red-600 hover:text-red-800"
                        >
                            Delete
                        </button>
                    </td>
                </tr>
                @ui.ConfirmDialog(&ui.ConfirmDialogProps{
                    ModalID:      "confirm-" + user.ID,
                    Title:        "Delete User",
                    Message:      "Are you sure you want to delete " + user.Name + "?",
                    ConfirmLabel: "Delete",
                    CancelLabel:  "Cancel",
                    ConfirmURL:   "/users/" + user.ID + "/delete",
                    Style:        "danger",
                })
            }
        </tbody>
    </table>
}
```

## Styles

| Style | Button Color | Use Case |
|-------|-------------|----------|
| `""` (default) | Blue | Standard confirmations |
| `"danger"` | Red | Destructive actions (delete, remove) |

## How It Works

- Uses TW Elements modal system (`data-twe-modal-init`, `data-twe-modal-dialog-ref`)
- Cancel button dismisses the modal via `data-twe-modal-dismiss`
- Confirm button navigates to `ConfirmURL` via `<a>` tag
- Modal is hidden by default and shown via `data-twe-toggle="modal"` on trigger element

## Dark Mode

Supports dark mode with appropriate background and text color adjustments.

## Accessibility

- Modal is focus-trapped when open
- Close button has `aria-label="Close"`
- ESC key dismisses the modal
- Clicking outside the modal dismisses it
- Focus returns to trigger element when modal closes
