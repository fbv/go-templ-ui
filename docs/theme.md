# Theme — Custom Semantic Colour Themes

`view/theme` lets a consumer author a brand theme against the library's
semantic colour tokens. A theme is a plain value: a name and the values of the
31 `--color-*` roles it overrides. Themes live in a package registry and render
as one `[data-theme="name"]` block over the aliases `ui.TokenStyles()` declares,
so a consumer rebrands without touching component markup and without forking
the built-in light/dark palettes.

A theme is the colour roles only. Radius, spacing and type scales stay
structural: rebranding is a palette change, not a layout change.

## Import

```go
import "github.com/fbv/go-templ-ui/view/theme"
```

## Props

None. This is a package, not a component: a registered theme is selected
through `data-theme`, so there is no props struct to change and no component
takes a theme argument.

## Usage

### API reference

| Function | Description |
|----------|-------------|
| `Register(t *Theme) error` | Adds a theme to the package registry, keyed by `t.Name`. Callable from a consumer's `init`. A `nil` theme, a name that does not match `[A-Za-z][A-Za-z0-9_-]*`, the reserved names `light`/`dark`, or a name already taken is an error and leaves the registry untouched. |
| `Get(name string) (*Theme, bool)` | The registered theme for a name; `false` for an unknown or empty name. The built-in light/dark palettes are not registry entries. |
| `Names() []string` | The registered names in sorted order — the list `layout.Page`'s pre-paint init accepts beside the built-in `light`/`dark`. |
| `Check(t *Theme) []Violation` | Every documented pair below 4.5:1, naming the pair, both colours and the ratio. Nil (empty) when the theme holds. |
| `Styles()` | A templ component emitting one `[data-theme="name"]` block per registered theme. `layout.Page` emits it beside `ui.TokenStyles()`. |

Types and methods:

| Identifier | Description |
|------------|-------------|
| `Theme` | `Name` plus the 31 semantic colour roles below. An empty field means "keep the `:root` alias" and is omitted from the emitted stylesheet. |
| `Violation` | `Pair`, `FGName`, `BGName`, `FG`, `BG`, `Ratio`. `String()` renders e.g. `ink (#EEEEEE) on canvas (#FFFFFF) = 1.16`. |

### The Theme value and key set

A `Theme` is a plain value; every field is one semantic role, spelled as the
role name without the `--color-` prefix (`Canvas` is `--color-canvas`). These
31 colour roles are the whole authoring surface — the component tokens
(`--rounded-*`, `--space-*`, `--font-*`, `--text-*`) are not theme values.

| Field | Token | Meaning |
|-------|-------|---------|
| `Canvas` | `--color-canvas` | Page background |
| `Surface` | `--color-surface` | Cards, panels, sidebar |
| `SurfaceAlt` | `--color-surface-alt` | Alternate panel fill |
| `SurfaceOverlay` | `--color-surface-overlay` | Overlays (modals, dropdowns) |
| `Field` | `--color-field` | Form field fill |
| `Subtle` | `--color-subtle` | Muted fill (neutral badge, washes) |
| `InfoSoft` | `--color-info-soft` | Soft informational wash |
| `Border` | `--color-border` | Hairline border |
| `BorderStrong` | `--color-border-strong` | Stronger border |
| `Ink` | `--color-ink` | Primary text |
| `InkMuted` | `--color-ink-muted` | Secondary text |
| `InkSubtle` | `--color-ink-subtle` | Tertiary text |
| `InkDisabled` | `--color-ink-disabled` | Disabled text |
| `Action` | `--color-action` | Links, primary affordances |
| `ActionHover` | `--color-action-hover` | Affordance hover |
| `ActionActive` | `--color-action-active` | Affordance pressed |
| `ActionFg` | `--color-action-fg` | Text on an action fill |
| `ActionSoft` | `--color-action-soft` | Soft action wash |
| `Highlight` | `--color-highlight` | Hover wash |
| `Selected` | `--color-selected` | Selected fill |
| `Success` | `--color-success` | Positive status |
| `SuccessSoft` | `--color-success-soft` | Soft positive wash |
| `SuccessBorder` | `--color-success-border` | Positive border |
| `Warning` | `--color-warning` | Caution status |
| `WarningSoft` | `--color-warning-soft` | Soft caution wash |
| `WarningBorder` | `--color-warning-border` | Caution border |
| `Danger` | `--color-danger` | Destructive/error status |
| `DangerFg` | `--color-danger-fg` | Text on a danger fill |
| `DangerSoft` | `--color-danger-soft` | Soft danger wash |
| `DangerBorder` | `--color-danger-border` | Danger border |
| `Focus` | `--color-focus` | Focus ring |

Values are written into the stylesheet as given; they are developer-supplied
data, like any CSS the consumer owns. The contrast check reads 6-digit hex
(`#1A2B4A`, with or without the `#`).

### Selection

Selection reuses `data-theme`, and the value that reaches that attribute comes
from `localStorage`. `layout.Page`'s pre-paint init runs before paint on every
load: it reads `localStorage`'s `theme` key, accepts any registered name beside
the built-in `light`/`dark`, falls back to `light` for anything else (an empty
or unavailable `localStorage`, an unknown name), and writes the result to
`data-theme` on `<html>` — overwriting whatever attribute the server rendered.
The `dark` class is added only for `dark`. `ui.ThemeToggle`'s `toggleTheme`
writes only `light` or `dark` to the same key.

So the selection path a consumer can follow for a custom theme is `localStorage`:
write the registered name to the `theme` key yourself — your own theme picker,
or a first-visit script — then reload; the pre-paint init applies it before
paint:

```js
localStorage.setItem("theme", "brand"); // then reload
```

Setting `data-theme="brand"` as a static attribute is not enough: the pre-paint
init overwrites the attribute from `localStorage` on every load.

## Custom theme

A `Theme` is plain data. Register one in your own package's `init`:

```go
package app

import "github.com/fbv/go-templ-ui/view/theme"

func init() {
	_ = theme.Register(&theme.Theme{
		Name:           "brand",
		Canvas:         "#FAF7F2",
		Surface:        "#FFFFFF",
		SurfaceAlt:     "#F0EBE3",
		Field:          "#FFFFFF",
		Subtle:         "#F0EBE3",
		InfoSoft:       "#E8E0D4",
		Border:         "#D9D0C3",
		BorderStrong:   "#BFB3A2",
		Ink:            "#1A2B4A",
		InkMuted:       "#3D4A63",
		Action:         "#0B5FA5",
		ActionFg:       "#FFFFFF",
		Success:        "#1F6B3A",
		Warning:        "#8A5A00",
		Danger:         "#A32B22",
		Focus:          "#0B5FA5",
	})
}
```

Select it through `localStorage`, the way the pre-paint init reads it (see
Selection above) — the library ships no picker for registered names:

```js
localStorage.setItem("theme", "brand"); // then reload
```

Roles the theme leaves out resolve through the `:root` aliases, so a partial
theme is usable: rebrand what matters first and fill the rest in over time.
Registering a name that is already taken is an error and leaves the registry
unchanged — a name is taken once, so two themes cannot own the same selector.
The names `light` and `dark` are reserved for the built-in palettes. A name
must match `[A-Za-z][A-Za-z0-9_-]*` so it is usable in the emitted CSS and in
the pre-paint init's name list.

## Contrast check

`Check` runs the same floor the library holds its built-in palettes to: every
documented text/background pair at ≥ 4.5:1 (WCAG AA body text). Run it on a
theme before shipping it.

| Foreground | Background | Why |
|------------|------------|-----|
| `ink` | `canvas`, `surface`, `surface-alt`, `field`, `subtle`, `info-soft` | Body text on every surface it can sit on |
| `ink-muted` | `canvas`, `surface`, `surface-alt`, `field`, `subtle`, `info-soft` | Secondary text on every surface it can sit on |
| `action` | `surface` | Links and primary affordances on cards/panels |
| `success` | `surface` | Positive status on cards/panels |
| `warning` | `surface` | Caution status on cards/panels |
| `danger` | `surface` | Destructive status on cards/panels |

A pair that fails reports its name, both colours and the ratio; a theme that
holds every pair reports nothing. Roles a theme leaves unset are checked as
the page paints them — through the `:root` alias, which is the built-in light
palette — so a partial theme is proved against what it will actually render.
A colour the check cannot read (not 6-digit hex) reports a ratio of `0` and
therefore fails, instead of panicking. The built-in light and dark palettes
report zero violations under the same check.

```go
if v := theme.Check(myTheme); len(v) > 0 {
	for _, pair := range v {
		log.Println(pair) // ink (#EEEEEE) on canvas (#FFFFFF) = 1.16
	}
}
```

## Styles

`Styles()` emits one `[data-theme="name"]` block per registered theme,
declaring only the roles that theme supplies. `layout.Page` places it beside
`ui.TokenStyles()`, so the blocks sit after the `:root` aliases they override
and every omitted role keeps its `:root` fallback. Components consume the same
semantic tokens either way — no component markup or prop changes.

## Dark Mode

A custom theme is orthogonal to `dark`. Selecting a custom theme sets
`data-theme="brand"` instead of `light`/`dark`, and `layout.Page` adds
Tailwind's `dark` class only when the selected name is `dark` — a custom theme
never gains the class. `ui.ThemeToggle` continues to flip `data-theme` between
`light` and `dark`.

A custom theme recolours only the elements that read the `--color-*` semantic
tokens. The component templates are not limited to those tokens: many still
carry fixed Tailwind palette classes (`bg-white`, `text-gray-500`,
`bg-blue-600`, …) and `dark:` variant utilities. Those elements keep their
built-in Tailwind colours and ignore the theme. While a custom theme is
selected the `dark` class is absent, so their `dark:` variants are inactive and
their light Tailwind values paint beside the theme's palette; selecting `dark`
turns those variants on again.

## Accessibility

- **Contrast (WCAG 1.4.3).** `Check` enforces ≥ 4.5:1 on the documented
  text/background pairs above and names every violation, so a theme can be
  proved before it ships. The built-in light and dark palettes hold the same
  floor.
- **Unreadable values fail loudly.** A value that is not 6-digit hex reports a
  `0` ratio and appears in the violation list rather than being skipped.
- **Meaning never rests on colour alone.** A palette recolours roles; it does
  not add or remove cues. Keep the icon or text cue that accompanies a status
  in your own templates.
