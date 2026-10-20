# State — the uniform five data states

`view/ui` carries one shared value for the five states every data-bearing
component must be able to render (CAP-5, EXPERIENCE.md § State Patterns):
a `DataState` enum, a `StateProps` struct, and the `StateView` surface that
renders each state's treatment. A component with data to show renders its
content; a component whose data is loading, missing, failed, forbidden, or
only partially available sets `State` and renders the matching surface —
never a blank area.

`Panel` is the first carrier: an additive `State *StateProps` field that
replaces its children (header kept) when set. Ticket 1.7 applies the same
value to table, list and select without redesigning the API.

Every user-facing string is a prop. The component ships no hardcoded literal,
so it is locale-neutral: a consumer supplies `Message`, `ActionLabel` and
`RetryLabel` from its own catalog (see [i18n.md](i18n.md)).

## Import

```go
import "github.com/fbv/go-templ-ui/view/ui"
```

## Props

`DataState` is a string enum:

| Value | Treatment |
|-------|-----------|
| `StateLoading` (`"loading"`) | `Skeleton` rows matching the expected layout and count. Never a spinner over a blank; the host's height is preserved by the row-shaped skeleton. |
| `StateEmpty` (`"empty"`) | Neutral, illustration-free message in the subheading typography role, plus one primary action. |
| `StateError` (`"error"`) | Neutral surface with a `{colors.danger}` accent and a retry action. |
| `StatePermissionDenied` (`"permission-denied"`) | Neutral surface in `{colors.ink-muted}`. No retry control is rendered — retrying cannot help. |
| `StatePartial` (`"partial"`) | Renders the host's children plus an inline notice naming what is unavailable. |

`StateProps`:

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `State` | `ui.DataState` | No | Which of the five treatments to render. |
| `Message` | `string` | No | The consumer-supplied message: what the empty state offers, why the error happened, who to ask, or which fields are unavailable. |
| `ActionLabel` | `string` | No | Empty state's primary action label. With `ActionURL`, renders exactly one action control; alone it renders nothing. |
| `ActionURL` | `string` | No | Empty state's primary action target. Both label and URL must be set for the control to render. |
| `RetryLabel` | `string` | No | Error state's retry label. Never rendered by permission-denied. |
| `RetryURL` | `string` | No | Error state's retry target. Both label and URL must be set for the control to render. |
| `SkeletonRows` | `int` | No | Loading: skeleton row count matching the expected layout. Zero renders a single-line skeleton. |
| `SkeletonCols` | `int` | No | Loading: cells per skeleton row. Zero or one renders one column. |

Fallbacks, in full:

- `StateView(nil)` renders its children untouched (and nothing when there are
  none) — the pre-state behaviour.
- An unrecognised `State` value (including `""`) renders the empty treatment:
  a visible surface, never a blank area and never a panic.
- A state other than `partial` replaces any children with its surface; only
  `partial` emits them.

## Usage

### Carrying a state on a Panel

`PanelProps.State` is additive: with it unset, `Panel` renders exactly as
before. When set, the panel keeps its header and renders the state surface in
place of its children.

```go
@ui.Panel(&ui.PanelProps{Header: "Accounts", State: &ui.StateProps{
	State:        ui.StateLoading,
	SkeletonRows: 5,
	SkeletonCols: 4,
}}) {
	@accountTable(accounts)
}
```

### Loading

```go
&ui.StateProps{State: ui.StateLoading, SkeletonRows: 5, SkeletonCols: 4}
```

Renders `SkeletonTable(rows, cols)` — row-shaped placeholders that match the
expected layout and count, so nothing shifts when the data arrives. A spinner
is never rendered over a blank host.

### Empty

```go
&ui.StateProps{
	State:       ui.StateEmpty,
	Message:     "No accounts yet.",
	ActionLabel: "Create account",
	ActionURL:   "/accounts/new",
}
```

The message renders in the subheading typography role; the action renders as
exactly one primary control (action-blue link on `--color-action`).

### Error

```go
&ui.StateProps{
	State:      ui.StateError,
	Message:    "Accounts could not be loaded. Retry, or contact support if it persists.",
	RetryLabel: "Retry",
	RetryURL:   "/accounts",
}
```

The surface is the neutral subtle fill with a `--color-danger-border` accent;
the message is danger-coloured text that names the cause and the next step;
the retry control is the same primary-action treatment as the empty state's
action. The component renders the control — re-issuing the request is the
consumer's handler behind `RetryURL`.

### Permission denied

```go
&ui.StateProps{
	State:   ui.StatePermissionDenied,
	Message: "You do not have access to this record. Ask an administrator for access.",
}
```

The surface is neutral with muted ink. No retry control and no action control
render, whatever `RetryLabel`/`RetryURL`/`ActionLabel`/`ActionURL` carry —
the message alone explains why and who to ask, and the surface is never
styled as a system failure.

### Partial data — the composition rule

Partial data is the one state that is not a full-surface replacement: it
renders the host's content *plus* a notice. `StateView` accepts children, and
only `StatePartial` emits them:

```go
@ui.Panel(&ui.PanelProps{Header: "Accounts", State: &ui.StateProps{
	State:   ui.StatePartial,
	Message: "Owner and ARR are unavailable.",
}}) {
	@accountTable(accounts)
}
```

The children render as usual; the notice follows them as an inline row naming
what is unavailable. A caller who passes children to any other state gets the
state surface — the "never blank" guarantee takes precedence.

### Using StateView directly

`StateView` is the shared value, not a per-component variant, so any component
can host it:

```go
@ui.StateView(&ui.StateProps{State: ui.StateError, Message: msg, RetryLabel: retry, RetryURL: u}) {
	@myContent()
}
```

### The `/showcase/state` demo

`/showcase/state` (`showcase/state.templ`) is this document demonstrated end
to end under the `DemoPage` scaffold's variants section: one `Panel` per
state — loading skeleton rows, empty with a primary action, error with a
retry, permission-denied with no retry, and partial data with an inline
notice — each with consumer-supplied text.

## Styles

`StateView`'s own markup takes every colour and type role from the L2 token
layers (see `view/ui/tokens.templ`); there are no hardcoded palette classes
and no `dark:` variants in it — light and dark recolour by swapping token
values alone. Loading is the exception: it reuses the pre-existing
`SkeletonTable`/`Skeleton` components, whose palette classes
(`bg-gray-200`, `dark:bg-gray-700`, `dark:border-gray-700`) are known debt
owned by the refactor sweep, not `StateView` markup.

| Treatment | Surface | Message | Control |
|-----------|---------|---------|---------|
| Loading | `SkeletonTable` / `Skeleton` (pre-existing components) | — | — |
| Empty | `--color-subtle` on `--color-border`, rounded-lg | subheading role, `--color-ink` | primary action on `--color-action` / `--color-action-fg` |
| Error | `--color-subtle` on `--color-danger-border` | body role, `--color-danger` | primary action on `--color-action` |
| Permission denied | `--color-subtle` on `--color-border` | body role, `--color-ink-muted` | none |
| Partial notice | `--color-subtle` on `--color-border` | body-sm role, `--color-ink-muted` | — |

The focus ring on the action control is the house ring: 2px
`--color-focus` with a 2px offset, matching the rest of the library.

## Dark Mode

Theme switching flips L1 primitive values (`--g-*`) on `[data-theme="dark"]`
only. `StateView`'s own markup and classes are identical in both themes; the
`--color-*` tokens resolve to the active L1 values. The loading skeleton's
`dark:` variants live on the pre-existing `Skeleton` components (known debt,
see Styles), not on `StateView` markup.

See `view/ui/tokens.templ` for the three token layers and
`view/ui/theme_toggle.templ` for the switch.

## Accessibility

- **Never a blank surface.** With a state set, the component renders that
  state's surface; an unrecognised value renders the empty treatment. The
  only blank-adjacent path is `StateView(nil)`, which is the explicit
  "no state" passthrough and renders the caller's children unchanged.
- **The message is the surface's accessible content for empty, error,
  permission-denied and partial.** In those four states the consumer's
  `Message` is visible text in the surface, read by assistive technology as
  the content of the region; the library adds no `aria-label` over it and
  hardcodes no string — the consumer's wording is what gets announced.
  Loading is different: it ignores `Message` and emits only decorative
  `div`s via `SkeletonTable`/`Skeleton` — no text and no status role.
- **Permission denied is not a failure.** It uses muted ink on a neutral
  surface, renders no retry, and its message (supplied by the consumer)
  explains why access is denied and who to ask.
- **Error is identified in text**, not colour alone (WCAG 1.4.1): the danger
  accent accompanies message text that names the cause and the next step.
- **Loading does not shift the layout.** The skeleton matches the expected
  layout and count, so the surface the data arrives into is already the right
  shape.
- **Controls meet the target-size floor** (≥ 24×24 CSS px) and take the
  library's single focus ring.
