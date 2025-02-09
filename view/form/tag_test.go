package form

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUITag(t *testing.T) {
	tag := UITag("aaa,type=password, bbb, ddd=\"fff, \\\"ggg\\\"\"")
	p, ok := tag.Lookup("type")
	assert.True(t, ok)
	assert.Equal(t, "password", p)
	p, ok = tag.Lookup("aaa")
	assert.True(t, ok)
	assert.Empty(t, p)
	p, ok = tag.Lookup("bbb")
	assert.True(t, ok)
	assert.Empty(t, p)
	p, ok = tag.Lookup("ccc")
	assert.False(t, ok)
	assert.Empty(t, p)
	p, ok = tag.Lookup("ddd")
	assert.True(t, ok)
	assert.Equal(t, "fff, \"ggg\"", p)
}
