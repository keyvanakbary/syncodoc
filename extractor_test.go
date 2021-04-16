package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNumSpaces(t *testing.T) {
	cases := map[string]int{
		"zero": 0,
		"   3 spaces": 3,
		"    4 spaces": 4,
	}

	for line, expected := range cases {
		assert.Equal(t, expected, numSpaces([]byte(line)))
	}
}

func TestTrimsAllSpaces(t *testing.T) {
	snippet := []byte(
`  First line
    Second line
      Third line`)

	expected := []byte(
`First line
  Second line
    Third line`)

	assert.Equal(t, string(expected), string(trim(snippet)))
}

func TestTrimsSpacesOnly(t *testing.T) {
	snippet := []byte(
		`  First line
Second line
      Third line`)

	expected := []byte(
		`First line
Second line
    Third line`)

	assert.Equal(t, string(expected), string(trim(snippet)))
}

func newSnippet(name string, contents string) *Snippet {
	return &Snippet{[]byte(name), []byte(contents)}
}

func assertSnippetExists(t *testing.T, ss *Snippets, s *Snippet) {
	if !ss.Contains(s.Name) {
		assert.Failf(t, "No snippet found with name %s", string(s.Name))
		return
	}
	snippet, _ := ss.ByName(s.Name)
	assert.Equal(t, string(s.Contents), string(snippet.Contents))
}

func AssertSnippetsExists(t *testing.T, ss1 *Snippets, ss2 *Snippets) {
	ss2.Each(func(snippet *Snippet) {
		assertSnippetExists(t, ss1, snippet)
	})
}

func TestExtractsNothing(t *testing.T) {
	snippets := Extract([]byte(`Some text`))

	assert.True(t, snippets.Empty(), "There should be no snippets")
}

func TestExtractsSingleSnippet(t *testing.T) {
	snippets := Extract([]byte(`
Some text
// snippet foo
Some foo
// end-snippet
`))

	assertSnippetExists(t, snippets, newSnippet("foo", "Some foo"))
}

func TestExtractsSnippedMultiLine(t *testing.T) {
	snippets := Extract([]byte(`
Some text
// snippet foo
First line
Second line
// end-snippet
`))

	assertSnippetExists(t, snippets, newSnippet("foo", `First line
Second line`))
}

func TestIgnoreBlock(t *testing.T) {
	snippets := Extract([]byte(`
Some text
// snippet foo
// ignore
First line
Second line
// end-ignore
Third line
// end-snippet
`))

	assertSnippetExists(t, snippets, newSnippet("foo", `// ...
Third line`))
}


func TestExtractsMultipleSnippets(t *testing.T) {
	snippets := Extract([]byte(`
Some text
// snippet foo
Some foo
// end-snippet

More text
// snippet bar
Some bar
// end-snippet
`))

	AssertSnippetsExists(t, snippets, NewSnippets().
		Put(newSnippet("foo", "Some foo")).
		Put(newSnippet("bar", "Some bar")))
}

func TestExtractTrimsLeftSpace(t *testing.T) {
	snippets := Extract([]byte(`
	Some text
    // snippet foo
    Line one
        Line two
            Line three
    // end-snippet
`))

	assertSnippetExists(t, snippets, newSnippet("foo", `Line one
    Line two
        Line three`))
}

func TestExtractIgnoreRespectsSpacing(t *testing.T) {
	snippets := Extract([]byte(`
Some text
// snippet foo
Line one
    // ignore
        Line two
            Line three
	// end-ignore
    Line four
// end-snippet
`))

	assertSnippetExists(t, snippets, newSnippet("foo", `Line one
    // ...
    Line four`))
}