# Go Templ UI - Component Documentation

## Showcase & docs convention

Every component in `view/ui` ships with both of the following. Later epics must
follow this convention — the Badge demo is the reference proof.

1. **Demo page** at `/showcase/<name>`, built on `showcase.DemoPage`
   (`showcase/demo.templ`). The scaffold renders: title, description, props
   summary table, variant/state matrix (each variant has a live preview box and
   a `ui.Code` snippet), and a link to `docs/<name>.md`.
2. **Docs page** at `docs/<name>.md`, following the template used by
   [alert.md](alert.md) and [badge.md](badge.md): Import / Props / Usage /
   Styles / Dark Mode / Accessibility.

The demo page is registered in `showcase/registry.go` (`Demo` struct: name,
title, icon, route, docs path, render); the docs page is a markdown file in
`docs/` that the row's `DocsPath` names. `server/server.go` iterates the
registry to register routes and mounts the embedded docs tree;
`showcase/showcase.templ` builds sidebar links from the registry. The
demo page receives its title and docs path from its registry row —
`Render` has signature `func(d Demo) templ.Component` and is called as
`d.Render(d)` — so a demo never restates them and its sidebar link, page
heading and docs URL cannot diverge. Add a category link and a Quick Reference
row below when you add a component.

Components consume semantic/component tokens only (see `view/ui/tokens.templ`).
Never hardcode palette classes in a component or a demo. The theme toggle ships
on demo pages (`showcase.DemoPage`), never on the not-yet-migrated mega
showcase page.

## Components

### Navigation
- [Sidebar](sidebar.md) - Navigation sidebar with collapsible mode

### Feedback & Notifications
- [Alert](alert.md) - Banner notifications (info, success, warning, error)
- [Toast](toast.md) - Floating notifications with auto-dismiss
- [Spinner](spinner.md) - Loading indicators

### Data Display
- [Badge](badge.md) - Status pill labels (token-based, light/dark)
- [Skeleton](skeleton.md) - Placeholder content while loading
- [Pagination](pagination.md) - Page navigation for lists/tables

### Forms
- [DatePicker](datepicker.md) - Date and time input fields

### Dialogs
- [ConfirmDialog](confirm_dialog.md) - Confirmation modal for dangerous actions

### Theming & i18n
- [i18n](i18n.md) - Locale, strings and formats carried in `context.Context`
  (package, not a component: no demo page)

## Quick Reference

| Component | Purpose | Key Props |
|-----------|---------|-----------|
| `Sidebar` | Navigation sidebar | `Items`, `User`, `Collapsible`, `ID` |
| `Alert` | Static banner notifications | `Title`, `Style`, `Dismissible` |
| `Toast` | Auto-hiding floating alerts | `ID`, `Title`, `Style`, `Duration` |
| `Spinner` | Loading indicator | `Size`, `Style`, `Label` |
| `Badge` | Status pill label | `Name`, `Style` |
| `Skeleton` | Content placeholder | `Lines`, `Style` |
| `Pagination` | Page navigation | `CurrentPage`, `TotalPages`, `BaseURL` |
| `DatePicker` | Date/time input | `Name`, `Label`, `Type` |
| `ConfirmDialog` | Confirmation modal | `ModalID`, `Title`, `Message`, `ConfirmURL` |
| `i18n` | Locale, strings and formats | `WithLocale`, `T`, `Number`, `Currency`, `Date`, `MonthYear` |
