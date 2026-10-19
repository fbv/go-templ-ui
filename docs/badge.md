# Badge Component

Status pill for labels and state. Renders a rounded-full caption badge whose
colors come from semantic/component tokens — light and dark recolor by swapping
token values alone.

## Import

```go
import "github.com/fbv/go-templ-ui/view/ui"
```

## Props

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `Name` | `string` | Yes | Badge label text |
| `Style` | `string` | No | Visual variant: `"red"`, `"green"`, `"yellow"`, `"ping"`, `"gray"`, `"indigo"`, `"purple"`, `"blue"`. Empty or unknown `Style` falls back to the neutral default variant. |

There is no `"gray"` case in `badgeStyle`: `"gray"`, `"indigo"`, `"purple"`
and `"blue"` take the same neutral default as an empty or unknown value. That
is the "neutral" variant below — DESIGN.md reserves blue for affordances and a
badge is a non-interactive display pill, so there is no action-blue badge.
There is no panic and no hardcoded palette class on any path.

## Usage

### Basic Badge

```go
@ui.Badge(&ui.BadgeProps{Name: "Draft", Style: "gray"})
```

### Status Badges

```go
@ui.Badge(&ui.BadgeProps{Name: "Active", Style: "green"})
@ui.Badge(&ui.BadgeProps{Name: "Pending", Style: "yellow"})
@ui.Badge(&ui.BadgeProps{Name: "Overdue", Style: "red"})
```

### Default (no Style)

```go
@ui.Badge(&ui.BadgeProps{Name: "Default"})
```

## Styles

| Style | Variant | Fill | Foreground | Border |
|-------|---------|------|------------|--------|
| `"gray"`, `"indigo"`, `"purple"`, `"blue"`, `""`, unknown | neutral | `--badge-neutral-bg` | `--badge-neutral-fg` | `--badge-neutral-border` |
| `"green"` | success | `--badge-success-bg` | `--badge-success-fg` | `--badge-success-border` |
| `"yellow"` | warning | `--badge-warning-bg` | `--badge-warning-fg` | `--badge-warning-border` |
| `"red"`, `"ping"` | danger | `--badge-danger-bg` | `--badge-danger-fg` | `--badge-danger-border` |

Anatomy follows the design system: a caption-type (11px) pill whose radius is
the `--badge-radius` component token (`{rounded.full}`), a status soft fill and a
1px matching border so the pill holds its shape against the surface. Every
color class is a `var(--badge-*)` token — there are no hardcoded palette
utilities.

The **warning** pill is a tinted-neutral pill, not a status-coloured label: its
label uses `--color-ink-muted` for contrast (DESIGN.md's `warning` on
`warning-soft` measures 4.39:1 in light, below the 4.5:1 floor), so its only
status cue is the `--badge-warning-border` border. Do not read the fill as the
status signal.

Badges are non-interactive display pills, so they never take the action blue
(blue is reserved for affordances).

## Dark Mode

Theme switching flips L1 primitive values (`--g-*`) on `[data-theme="dark"]`
only. Badge markup and classes are identical in both themes; the L3
`--badge-*` tokens resolve through L2 (`--color-*`) to the active L1 values.

See `view/ui/tokens.templ` for the three token layers and `view/ui/theme_toggle.templ`
for the switch.

## Accessibility

- Badge meaning is carried by its text label, not by color alone (WCAG 1.4.1)
- Text/background pairs meet ≥ 4.5:1 contrast in both themes
- The warning pill's label is `ink-muted` for contrast (status cue = `warning-border`, not the fill)
- The pill is a `<span>` — decorative status chrome around already-visible text
