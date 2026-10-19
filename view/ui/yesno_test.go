package ui

import (
	"context"
	"strings"
	"testing"

	"github.com/fbv/go-templ-ui/view/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func renderYesNo(t *testing.T, ctx context.Context, yes bool) string {
	t.Helper()
	var buf strings.Builder
	require.NoError(t, YesNo(yes).Render(ctx, &buf))
	return buf.String()
}

func TestYesNoUsesDefaultEnglishCatalog(t *testing.T) {
	yes := renderYesNo(t, context.Background(), true)
	no := renderYesNo(t, context.Background(), false)

	assert.Contains(t, yes, ">Yes<")
	assert.NotContains(t, yes, "Да")
	assert.Contains(t, no, ">No<")
	assert.NotContains(t, no, "Нет")
}

func TestYesNoFollowsTheContextLocale(t *testing.T) {
	ctx := i18n.WithLocale(context.Background(), i18n.Russian)

	yes := renderYesNo(t, ctx, true)
	no := renderYesNo(t, ctx, false)

	assert.Contains(t, yes, ">Да<", "a Russian context must render Да")
	assert.Contains(t, no, ">Нет<", "a Russian context must render Нет")
}
