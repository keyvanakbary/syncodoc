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

func assertSnippetExists(t *testing.T, snippets Snippets, expected Snippets) {
	for expectedName, expectedSnippet := range expected {
		assert.Contains(t, snippets, expectedName, "Snippet %s does not exist", expectedName)
		assert.Equal(t, string(expectedSnippet), string(snippets[expectedName]))
	}
}

func TestExtractsNothing(t *testing.T) {
	snippets := extract([]byte(`Some text`))

	assert.Empty(t, snippets, "There should be no snippets")
}

func TestExtractsSingleSnippet(t *testing.T) {
	snippets := extract([]byte(`
Some text
// snippet foo
Some foo
// end-snippet
`))

	assertSnippetExists(t, snippets, Snippets {"foo": Snippet("Some foo")})
}
//
func TestExtractsSnippedMultiLine(t *testing.T) {
	snippets := extract([]byte(`
Some text
// snippet foo
First line
Second line
// end-snippet
`))

	assertSnippetExists(t, snippets, Snippets {"foo": Snippet(`First line
Second line`)})
}

func TestIgnoreBlock(t *testing.T) {
	snippets := extract([]byte(`
Some text
// snippet foo
// ignore
First line
Second line
// end-ignore
Third line
// end-snippet
`))

	assertSnippetExists(t, snippets, Snippets {"foo": Snippet(`// ...
Third line`)})
}

func TestExtractsMultipleSnippets(t *testing.T) {
	snippets := extract([]byte(`
Some text
// snippet foo
Some foo
// end-snippet

More text
// snippet bar
Some bar
// end-snippet
`))

	assertSnippetExists(t, snippets, Snippets {
		"foo": Snippet("Some foo"),
		"bar": Snippet("Some bar"),
	})
}

func TestExtractTrimsLeftSpace(t *testing.T) {
	snippets := extract([]byte(`
	Some text
    // snippet foo
    Line one
        Line two
            Line three
    // end-snippet
`))

	assertSnippetExists(t, snippets, Snippets{ "foo": Snippet(`Line one
    Line two
        Line three`)})
}

func TestExtractIgnoreRespectsSpacing(t *testing.T) {
	snippets := extract([]byte(`
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

	assertSnippetExists(t, snippets, Snippets{ "foo": Snippet(`Line one
    // ...
    Line four`)})
}

func TestExtracts(t *testing.T) {
	snippets := extract([]byte(`
# snippet foo
Some text
# ignore
Line 1
Line 2
# end-ignore
# end-snippet
`))

	assertSnippetExists(t, snippets, Snippets{ "foo": Snippet(`Some text
// ...`)})
}