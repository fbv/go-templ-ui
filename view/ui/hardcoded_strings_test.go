package ui

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	parser "github.com/a-h/templ/parser/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This guard enforces CAP-4's "no new component hardcodes a user-facing string".
//
// A user-facing string is a literal text node in a templ body, or the literal
// value of a user-facing attribute (aria-label, title, placeholder, alt). Every
// one of them must go through i18n.T. The allowlist records the strings that
// already existed at the baseline revision; no entry may be added (fix the
// component instead), and stale entries must be removed.
//
// Extraction is delegated to templ's own parser (github.com/a-h/templ/parser/v2)
// and a parser.Visitor, so the AST decides what is text and what is Go code.
// Text inside if/for/switch bodies is reached because the visitor descends into
// every child node; text produced by an expression, a component call, a
// <script> or a <style> is never a *parser.Text node and is never reported.
//
// Extracted strings are compared as a multiset ("file: text" entries, in sorted
// order, duplicates included): a second copy of an allowlisted string still
// counts as a new hardcoded string.

const allowlistPath = "i18n_allowlist.txt"

// userFacingAttrs are the attributes whose literal value is read aloud or shown
// to the user, so they are part of the locale surface.
var userFacingAttrs = map[string]bool{
	"aria-label":  true,
	"title":       true,
	"placeholder": true,
	"alt":         true,
}

// extractor fixtures. Each shape the guard must get right has one, with its
// expected set written down in the test below.
const (
	// A bare literal is the thing the guard exists to report; the same text
	// routed through i18n.T is not.
	fixtureBare = `package ui

templ Fixture() {
	<div class="p-2">
		<span class="text-green-600">Да</span>
	</div>
}
`
	fixtureLocalized = `package ui

import "github.com/fbv/go-templ-ui/view/i18n"

templ Fixture() {
	<div class="p-2">
		<span class="text-green-600">{ i18n.T(ctx, i18n.KeyYes) }</span>
	</div>
}
`
	// Go statements, <script>/<style> bodies and expressions are not text
	// nodes, while markup inside for/if bodies still is.
	fixtureNoise = `package ui

import "github.com/fbv/go-templ-ui/view/i18n"

templ Fixture(p *Props) {
	{{ label := "not a text node" }}
	<script>
		var s = "Нет событий";
	</script>
	<style>
		.fixture::after { content: 'Ещё'; }
	</style>
	for _, item := range p.Items {
		if item.Active {
			<span>{ i18n.T(ctx, i18n.KeyYes) }</span>
		}
	}
}
`
	// A literal whose rendered text begins with a Go keyword or contains := is
	// still text: the parser decides, not the words. templ itself refuses a
	// bare line beginning with "if"/"for"/"switch" (see the rejection fixture
	// below), so the parseable spelling of such text is entity-escaped; the
	// guard reports the node rather than dropping it.
	fixtureKeywordLiterals = `package ui

templ Fixture() {
	<p>&#105;f you continue, data will be lost</p>
	<p>&#102;or the record, this stays</p>
	<p>Assign with x := 1 in code</p>
}
`
	// templ's parser rejects an unescaped keyword at the start of a text line,
	// so such a literal can never be silently dropped: the file does not parse.
	fixtureUnescapedKeyword = `package ui

templ Fixture() {
	<p>if you continue, data will be lost</p>
}
`
	// A brace inside a text run is an inline expression: the text on both sides
	// is kept and the Go identifier between the braces is never reported.
	fixtureInlineExpression = `package ui

templ Fixture(n int) {
	<p>Total { n } results</p>
}
`
	// A label that ends in ':' is text, not a Go switch label.
	fixtureTrailingColon = `package ui

templ Fixture() {
	<p>Warning:</p>
	<p>data will be lost</p>
}
`
	// An '@' inside a text run is not a component call.
	fixtureEmail = `package ui

templ Fixture() {
	<p>Contact support@example.com</p>
}
`
	// Text before an if block is its own text node and must survive.
	fixtureTextBeforeIf = `package ui

templ Fixture(x bool) {
	<div>Total
		if x {
			<span>yes</span>
		}
	</div>
}
`
	// A <style> block's contents are raw text, not text nodes, and markup after
	// it is still reached (the parser owns the brace matching).
	fixtureStyleThenMarkup = `package ui

templ Fixture() {
	<style>
.x {
	content: 'Ещё';
}
	</style>
	<p>After the style block</p>
}
`
	// A component call at statement position after a text line is a call, not
	// part of the literal.
	fixtureComponentAfterText = `package ui

templ Fixture() {
	<div>Total
		@Price()
	</div>
}
`
	// Unlocalized user-facing attributes are collected, so they can be reported.
	fixtureAttributes = `package ui

templ Fixture() {
	<div>
		<button aria-label="Search">ok</button>
		<input placeholder="Filter rows" title="Row filter" alt="Filter"/>
	</div>
}
`
	// Two copies of the same literal are two entries, never one.
	fixtureDuplicate = `package ui

templ Fixture() {
	<div>
		<span>Close</span>
		<span>Close</span>
	</div>
}
`
	// Go control flow around markup is not text: only the rendered strings
	// survive, including switch bodies.
	fixtureControlFlow = `package ui

templ Fixture(b bool, n int) {
	if b {
		<span>ok</span>
	} else if n > 0 {
		<span>no</span>
	} else {
		<span>maybe</span>
	}
	for i := 0; i < n; i++ {
		<span>row</span>
	}
	switch n {
		case 1:
			<span>one</span>
		default:
			<span>many</span>
	}
}
`
)

func TestHardcodedStringsGuard(t *testing.T) {
	t.Run("allows only allowlisted strings", func(t *testing.T) {
		got := collectHardcodedStrings(t, "*.templ")

		want, err := readAllowlist(allowlistPath)
		require.NoError(t, err, "read %s", allowlistPath)

		missing, extra := diffMultisets(want, got)
		assert.Empty(t, extra, "hardcoded user-facing string(s) must go through i18n.T: %s", strings.Join(extra, "; "))
		assert.Empty(t, missing, "stale entries in %s (string no longer present): %s", allowlistPath, strings.Join(missing, "; "))
	})

	t.Run("reports a bare literal and accepts a localized one", func(t *testing.T) {
		assert.Equal(t, []string{"fixture.templ: Да"}, hardcodedEntries(t, "fixture.templ", fixtureBare),
			"a bare literal must be reported with its file")
		assert.Empty(t, hardcodedEntries(t, "fixture.templ", fixtureLocalized),
			"text routed through i18n.T must not be reported")
	})

	t.Run("ignores Go code, scripts and styles", func(t *testing.T) {
		assert.Empty(t, hardcodedEntries(t, "fixture.templ", fixtureNoise),
			"Go statements, <script>/<style> bodies and expressions are not text nodes")
	})

	t.Run("keeps literals whose rendered text begins with a Go keyword or contains :=", func(t *testing.T) {
		assert.Equal(t, []string{
			"fixture.templ: &#102;or the record, this stays",
			"fixture.templ: &#105;f you continue, data will be lost",
			"fixture.templ: Assign with x := 1 in code",
		}, hardcodedEntries(t, "fixture.templ", fixtureKeywordLiterals))
	})

	t.Run("rejects an unescaped keyword-leading literal instead of dropping it", func(t *testing.T) {
		_, err := userFacingStrings(fixtureUnescapedKeyword)
		require.Error(t, err, "templ cannot parse keyword-leading text; the guard must fail loudly, never skip it")
	})

	t.Run("keeps text around an inline expression and never reports its identifier", func(t *testing.T) {
		got := hardcodedEntries(t, "fixture.templ", fixtureInlineExpression)
		assert.Equal(t, []string{"fixture.templ: Total", "fixture.templ: results"}, got)
		assert.NotContains(t, got, "fixture.templ: n")
	})

	t.Run("reports a label that ends in a colon", func(t *testing.T) {
		assert.Equal(t, []string{"fixture.templ: Warning:", "fixture.templ: data will be lost"},
			hardcodedEntries(t, "fixture.templ", fixtureTrailingColon))
	})

	t.Run("keeps an email address as one text node", func(t *testing.T) {
		assert.Equal(t, []string{"fixture.templ: Contact support@example.com"},
			hardcodedEntries(t, "fixture.templ", fixtureEmail))
	})

	t.Run("keeps text immediately before an if block", func(t *testing.T) {
		assert.Equal(t, []string{"fixture.templ: Total", "fixture.templ: yes"},
			hardcodedEntries(t, "fixture.templ", fixtureTextBeforeIf))
	})

	t.Run("reaches markup after a style block", func(t *testing.T) {
		assert.Equal(t, []string{"fixture.templ: After the style block"},
			hardcodedEntries(t, "fixture.templ", fixtureStyleThenMarkup))
	})

	t.Run("does not swallow a component call after a text line", func(t *testing.T) {
		assert.Equal(t, []string{"fixture.templ: Total"},
			hardcodedEntries(t, "fixture.templ", fixtureComponentAfterText))
	})

	t.Run("reports an unlocalized user-facing attribute", func(t *testing.T) {
		assert.Equal(t, []string{
			"fixture.templ: Filter",      // alt
			"fixture.templ: Filter rows", // placeholder
			"fixture.templ: Row filter",  // title
			"fixture.templ: Search",      // aria-label
			"fixture.templ: ok",          // text node
		}, hardcodedEntries(t, "fixture.templ", fixtureAttributes))
	})

	t.Run("counts a duplicated literal twice", func(t *testing.T) {
		got := hardcodedEntries(t, "fixture.templ", fixtureDuplicate)
		assert.Equal(t, []string{"fixture.templ: Close", "fixture.templ: Close"}, got)

		// Against a single allowlist entry, the second copy is a new string.
		missing, extra := diffMultisets([]string{"fixture.templ: Close"}, got)
		assert.Empty(t, missing)
		assert.Equal(t, []string{"fixture.templ: Close"}, extra, "a duplicate must not hide behind an allowlisted copy")
	})

	t.Run("reaches text inside control-flow bodies only", func(t *testing.T) {
		assert.Equal(t, []string{
			"fixture.templ: many",
			"fixture.templ: maybe",
			"fixture.templ: no",
			"fixture.templ: ok",
			"fixture.templ: one",
			"fixture.templ: row",
		}, hardcodedEntries(t, "fixture.templ", fixtureControlFlow))
	})

	t.Run("fails in both mismatch directions", func(t *testing.T) {
		missing, extra := diffMultisets([]string{"a.templ: Close"}, []string{"a.templ: Close", "a.templ: Search"})
		assert.Empty(t, missing)
		assert.Equal(t, []string{"a.templ: Search"}, extra, "a new hardcoded string must be reported")

		missing, extra = diffMultisets([]string{"a.templ: Close", "a.templ: Search"}, []string{"a.templ: Close"})
		assert.Equal(t, []string{"a.templ: Search"}, missing, "a stale allowlist entry must be reported")
		assert.Empty(t, extra)
	})
}

// collectHardcodedStrings renders every matching templ file's user-facing
// strings as sorted "file: text" entries, duplicates included.
func collectHardcodedStrings(t *testing.T, pattern string) []string {
	t.Helper()
	files, err := filepath.Glob(pattern)
	require.NoError(t, err)
	require.NotEmpty(t, files, "no templ files matched %s", pattern)

	var out []string
	for _, file := range files {
		src, err := os.ReadFile(file)
		require.NoError(t, err)
		out = append(out, hardcodedEntries(t, filepath.Base(file), string(src))...)
	}
	sort.Strings(out)
	return out
}

// hardcodedEntries renders a source's user-facing strings as sorted "file: text"
// entries, duplicates included.
func hardcodedEntries(t *testing.T, name, src string) []string {
	t.Helper()
	found, err := userFacingStrings(src)
	require.NoError(t, err, "parse %s", name)

	out := make([]string, 0, len(found))
	for _, text := range found {
		out = append(out, name+": "+text)
	}
	sort.Strings(out)
	return out
}

// userFacingStrings parses src with templ's parser and collects every literal
// text node and every literal value of a user-facing attribute.
func userFacingStrings(src string) ([]string, error) {
	tf, err := parser.ParseString(src)
	if err != nil {
		return nil, err
	}
	e := &stringExtractor{}
	if err := tf.Visit(e); err != nil {
		return nil, err
	}
	return e.found, nil
}

// stringExtractor walks the templ AST. The parser's Visitor only dispatches the
// node it is given, so each composite node recurses into its children here.
type stringExtractor struct {
	found []string
}

func (e *stringExtractor) add(s string) {
	s = strings.Join(strings.Fields(s), " ")
	if s == "" {
		return
	}
	e.found = append(e.found, s)
}

func visitNodes(e *stringExtractor, nodes []parser.Node) error {
	for _, n := range nodes {
		if err := n.Visit(e); err != nil {
			return err
		}
	}
	return nil
}

func visitAttributes(e *stringExtractor, attrs []parser.Attribute) error {
	for _, a := range attrs {
		if err := a.Visit(e); err != nil {
			return err
		}
	}
	return nil
}

func (e *stringExtractor) VisitTemplateFile(tf *parser.TemplateFile) error {
	for _, n := range tf.Nodes {
		if err := n.Visit(e); err != nil {
			return err
		}
	}
	return nil
}

func (e *stringExtractor) VisitHTMLTemplate(t *parser.HTMLTemplate) error {
	return visitNodes(e, t.Children)
}

func (e *stringExtractor) VisitElement(el *parser.Element) error {
	if err := visitAttributes(e, el.Attributes); err != nil {
		return err
	}
	return visitNodes(e, el.Children)
}

func (e *stringExtractor) VisitText(t *parser.Text) error {
	e.add(t.Value)
	return nil
}

func (e *stringExtractor) VisitConstantAttribute(a *parser.ConstantAttribute) error {
	if userFacingAttrs[strings.ToLower(a.Key.String())] {
		e.add(a.Value)
	}
	return nil
}

func (e *stringExtractor) VisitConditionalAttribute(a *parser.ConditionalAttribute) error {
	if err := visitAttributes(e, a.Then); err != nil {
		return err
	}
	return visitAttributes(e, a.Else)
}

func (e *stringExtractor) VisitIfExpression(n *parser.IfExpression) error {
	return visitNodes(e, n.ChildNodes())
}

func (e *stringExtractor) VisitSwitchExpression(n *parser.SwitchExpression) error {
	return visitNodes(e, n.ChildNodes())
}

func (e *stringExtractor) VisitForExpression(n *parser.ForExpression) error {
	return visitNodes(e, n.ChildNodes())
}

func (e *stringExtractor) VisitTemplElementExpression(n *parser.TemplElementExpression) error {
	return visitNodes(e, n.Children)
}

// The remaining nodes cannot contain a user-facing literal text node: Go code,
// expressions, comments, script bodies, CSS and template calls.

func (e *stringExtractor) VisitTemplateFileGoExpression(*parser.TemplateFileGoExpression) error {
	return nil
}
func (e *stringExtractor) VisitPackage(*parser.Package) error { return nil }
func (e *stringExtractor) VisitWhitespace(*parser.Whitespace) error {
	return nil
}
func (e *stringExtractor) VisitCSSTemplate(*parser.CSSTemplate) error { return nil }
func (e *stringExtractor) VisitConstantCSSProperty(*parser.ConstantCSSProperty) error {
	return nil
}
func (e *stringExtractor) VisitExpressionCSSProperty(*parser.ExpressionCSSProperty) error {
	return nil
}
func (e *stringExtractor) VisitDocType(*parser.DocType) error { return nil }
func (e *stringExtractor) VisitScriptElement(*parser.ScriptElement) error {
	return nil
}
func (e *stringExtractor) VisitRawElement(*parser.RawElement) error { return nil }
func (e *stringExtractor) VisitBoolConstantAttribute(*parser.BoolConstantAttribute) error {
	return nil
}
func (e *stringExtractor) VisitBoolExpressionAttribute(*parser.BoolExpressionAttribute) error {
	return nil
}
func (e *stringExtractor) VisitExpressionAttribute(*parser.ExpressionAttribute) error {
	return nil
}
func (e *stringExtractor) VisitSpreadAttributes(*parser.SpreadAttributes) error { return nil }
func (e *stringExtractor) VisitGoComment(*parser.GoComment) error               { return nil }
func (e *stringExtractor) VisitHTMLComment(*parser.HTMLComment) error           { return nil }
func (e *stringExtractor) VisitCallTemplateExpression(*parser.CallTemplateExpression) error {
	return nil
}
func (e *stringExtractor) VisitChildrenExpression(*parser.ChildrenExpression) error {
	return nil
}
func (e *stringExtractor) VisitGoCode(*parser.GoCode) error { return nil }
func (e *stringExtractor) VisitStringExpression(*parser.StringExpression) error {
	return nil
}
func (e *stringExtractor) VisitScriptTemplate(*parser.ScriptTemplate) error { return nil }
func (e *stringExtractor) VisitFallthrough(*parser.Fallthrough) error       { return nil }

func readAllowlist(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var entries []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimRight(line, "\r")
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		entries = append(entries, line)
	}
	sort.Strings(entries)
	return entries, nil
}

// diffMultisets returns entries in want but not got, and entries in got but not
// want, honouring multiplicity.
func diffMultisets(want, got []string) (missing, extra []string) {
	wantCount := countStrings(want)
	gotCount := countStrings(got)
	for _, s := range sortedKeys(wantCount) {
		for i := gotCount[s]; i < wantCount[s]; i++ {
			missing = append(missing, s)
		}
	}
	for _, s := range sortedKeys(gotCount) {
		for i := wantCount[s]; i < gotCount[s]; i++ {
			extra = append(extra, s)
		}
	}
	sort.Strings(missing)
	sort.Strings(extra)
	return missing, extra
}

func countStrings(in []string) map[string]int {
	counts := make(map[string]int, len(in))
	for _, s := range in {
		counts[s]++
	}
	return counts
}

func sortedKeys(counts map[string]int) []string {
	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
