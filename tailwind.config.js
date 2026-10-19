/** @type {import('tailwindcss').Config} */
const plugin = require('tailwindcss/plugin');
module.exports = {
  content: ["./**/*.{html,js,templ}"],
  theme: {
    extend: {
      colors: {
        canvas: 'var(--color-canvas)',
        surface: 'var(--color-surface)',
        'surface-alt': 'var(--color-surface-alt)',
        'surface-overlay': 'var(--color-surface-overlay)',
        field: 'var(--color-field)',
        subtle: 'var(--color-subtle)',
        'info-soft': 'var(--color-info-soft)',
        border: 'var(--color-border)',
        'border-strong': 'var(--color-border-strong)',
        ink: 'var(--color-ink)',
        'ink-muted': 'var(--color-ink-muted)',
        'ink-subtle': 'var(--color-ink-subtle)',
        'ink-disabled': 'var(--color-ink-disabled)',
        action: 'var(--color-action)',
        'action-hover': 'var(--color-action-hover)',
        'action-active': 'var(--color-action-active)',
        'action-fg': 'var(--color-action-fg)',
        'action-soft': 'var(--color-action-soft)',
        highlight: 'var(--color-highlight)',
        selected: 'var(--color-selected)',
        success: 'var(--color-success)',
        'success-soft': 'var(--color-success-soft)',
        'success-border': 'var(--color-success-border)',
        warning: 'var(--color-warning)',
        'warning-soft': 'var(--color-warning-soft)',
        'warning-border': 'var(--color-warning-border)',
        danger: 'var(--color-danger)',
        'danger-fg': 'var(--color-danger-fg)',
        'danger-soft': 'var(--color-danger-soft)',
        'danger-border': 'var(--color-danger-border)',
        focus: 'var(--color-focus)',
      },
    },
  },
  plugins: [
    plugin(function({ addBase }) {
      addBase({
         'html': { fontSize: "13px" },
       })
     }),
  ],
}
