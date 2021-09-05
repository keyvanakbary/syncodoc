package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type FakeFileReader struct {
	Files map[string][]byte
}

func (r *FakeFileReader) Read(path string) ([]byte, error) {
	contents, ok := r.Files[path]
	if !ok {
		return nil, fmt.Errorf("no file mapping found for path %s", path)
	}

	return contents, nil
}

var fakeReader = &FakeFileReader{Files: map[string][]byte{
	"file/one": []byte(`one`),
	"file/two": []byte(`two`),
}}

type FakeSnippetReader struct {
	Snippets map[string]Snippet
}

func (r *FakeSnippetReader) ReadFull(path string) (Snippet, error) {
	snippet, ok := r.Snippets[path]
	if !ok {
		return nil, fmt.Errorf("no file mapping found for path %s", path)
	}

	return snippet, nil
}

func (r *FakeSnippetReader) ReadPart(path string, name string) (Snippet, error) {
	snippet, ok := r.Snippets[path]
	if !ok {
		return nil, fmt.Errorf("no file mapping found for path %s", path)
	}

	return snippet, nil
}

var fakeSnippetReader = &FakeSnippetReader{Snippets: map[string]Snippet{
	"/file/one.go": Snippet(`one`),
	"/file/two.go": Snippet(`two`),
	"/file/somepath.go": Snippet(`somepath`),
	"/code/relative.go": Snippet(`relative`),
}}

func TestSyncMoreThanOneBlock(t *testing.T) {
	content := []byte(`
pre one
%% sync /file/one.go
%% end

pre two
%% sync /file/two.go
%% end
`)

	expected := []byte(`
pre one
%% sync /file/one.go
{lang="go"}
~~~~~~~~
one
~~~~~~~~
%% end

pre two
%% sync /file/two.go
{lang="go"}
~~~~~~~~
two
~~~~~~~~
%% end
`)

	s := &Syncer{fakeSnippetReader, &MarkuaFormatter{}, fakeReader}

	synced := s.Sync("irrelevant", content)

	assert.Equal(t, string(expected), string(synced))
}

func TestSyncDoesNothing(t *testing.T) {
	content := []byte(`
%% sync /foo.go
%% end
`)

	expected := []byte(`
%% sync /foo.go
%% end
`)

	s := &Syncer{fakeSnippetReader, &MarkuaFormatter{}, fakeReader}

	synced := s.Sync("irrelevant", content)

	assert.Equal(t, string(expected), string(synced))
}

func TestAppendsCurrentPathToRelativePath(t *testing.T) {
	content := []byte(`
%% sync ../relative.go
%% end
`)

	expected := []byte(`
%% sync ../relative.go
{lang="go"}
~~~~~~~~
relative
~~~~~~~~
%% end
`)

	s := &Syncer{fakeSnippetReader, &MarkuaFormatter{}, fakeReader}

	synced := s.Sync("/code/src/", content)

	assert.Equal(t, string(expected), string(synced))
}



func TestExtractSyncs(t *testing.T) {
	content := []byte(`Line 1
Line 2
Line 3
%% sync ../somepath.go?lang=sql#snippet
%% end
`)
	syncs := extractSyncs(content)

	assert.Len(t, syncs, 1)
	assert.Equal(t, Sync{Lang: "sql", Path: "../somepath.go", Part: "snippet", LineStart: 4, LineEnd: 5}, syncs[0])
}



func TestReplaceSyncs(t *testing.T) {
	content := []byte(`Line 1
Line 2
Line 3
%% sync /file/somepath.go?lang=sql#snippet
%% end
`)
	expected := []byte(`Line 1
Line 2
Line 3
%% sync /file/somepath.go?lang=sql#snippet
{lang="sql"}
~~~~~~~~
somepath
~~~~~~~~
%% end
`)

	s := &Syncer{fakeSnippetReader, &MarkuaFormatter{}, fakeReader}

	synced, _ := s.replaceSync("irrelevant", content, Sync{
		Lang: "sql", Path: "/file/somepath.go", Part: "snippet", LineStart: 4, LineEnd: 5,
	})

	assert.Equal(t, string(expected), string(synced))
}

func TestReplacesCommentedCode(t *testing.T) {
	content := []byte(`
%% sync ../relative.go
This is complicated%
%% end
`)

	expected := []byte(`
%% sync ../relative.go
{lang="go"}
~~~~~~~~
relative
~~~~~~~~
%% end
`)

	s := &Syncer{fakeSnippetReader, &MarkuaFormatter{}, fakeReader}

	synced := s.Sync("/code/src/", content)

	assert.Equal(t, string(expected), string(synced))
}
