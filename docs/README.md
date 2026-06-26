# Go Templ UI - Component Documentation

## Components

### Navigation
- [Sidebar](sidebar.md) - Navigation sidebar with collapsible mode

### Feedback & Notifications
- [Alert](alert.md) - Banner notifications (info, success, warning, error)
- [Toast](toast.md) - Floating notifications with auto-dismiss
- [Spinner](spinner.md) - Loading indicators

### Data Display
- [Skeleton](skeleton.md) - Placeholder content while loading
- [Pagination](pagination.md) - Page navigation for lists/tables

### Forms
- [DatePicker](datepicker.md) - Date and time input fields

### Dialogs
- [ConfirmDialog](confirm_dialog.md) - Confirmation modal for dangerous actions

## Quick Reference

| Component | Purpose | Key Props |
|-----------|---------|-----------|
| `Sidebar` | Navigation sidebar | `Items`, `User`, `Collapsible`, `ID` |
| `Alert` | Static banner notifications | `Title`, `Style`, `Dismissible` |
| `Toast` | Auto-hiding floating alerts | `ID`, `Title`, `Style`, `Duration` |
| `Spinner` | Loading indicator | `Size`, `Style`, `Label` |
| `Skeleton` | Content placeholder | `Lines`, `Style` |
| `Pagination` | Page navigation | `CurrentPage`, `TotalPages`, `BaseURL` |
| `DatePicker` | Date/time input | `Name`, `Label`, `Type` |
| `ConfirmDialog` | Confirmation modal | `ModalID`, `Title`, `Message`, `ConfirmURL` |
