package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSyncMoreThanOneBlock(t *testing.T) {
	content := []byte(`
pre one
%% sync php file/one
%%

pre two
%% sync php file/two
%%
`)

	expected := []byte(`
pre one
%% sync php file/one
{lang="php"}
~~~~~~~~
one
~~~~~~~~
%%

pre two
%% sync php file/two
{lang="php"}
~~~~~~~~
two
~~~~~~~~
%%
`)

	s := &Syncer{
		Snippets: NewSnippets().
			Put(newSnippet("file/one", "one")).
			Put(newSnippet("file/two", "two")),
	}

	synced := s.Sync(content)

	assert.Equal(t, string(expected), string(synced))
}

func TestSyncDoesNothing(t *testing.T) {
	content := []byte(`
%% sync php file/one
%%
`)

	expected := []byte(`
%% sync php file/one
%%
`)

	s := &Syncer{
		Snippets: NewSnippets(),
	}

	synced := s.Sync(content)

	assert.Equal(t, string(expected), string(synced))
}
