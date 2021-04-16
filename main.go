package main

import (
	"fmt"
	"github.com/stoewer/go-strcase"
	"io/ioutil"
	"os"
	"strings"
)

func appendNamespace(namespace string, snippets *Snippets) *Snippets {
	snippets.Each(func(s *Snippet) {
		s.Name = []byte(namespace + string(s.Name))
	})

	return snippets
}

func extractSnippets(dir string, namespace string) *Snippets {
	snippets := NewSnippets()
	files, _ := ioutil.ReadDir(dir)

	for _, file := range files {
		path := dir + "/"
		filepath := path + file.Name()
		if file.IsDir() {
			snippets.Merge(extractSnippets(filepath, namespace + strcase.KebabCase(file.Name()) + "/"))
		} else {
			contents, _ := ioutil.ReadFile(filepath)
			snippets.Merge(appendNamespace(namespace, Extract(contents)))
		}
	}
	return snippets
}

func replaceSnippets(dir string, snippets *Snippets) {
	files, _ := ioutil.ReadDir(dir)

	for _, file := range files {
		filepath := dir + "/" + file.Name()
		if file.IsDir() {
			replaceSnippets(filepath, snippets)
		} else {
			fmt.Println("\t" + filepath)
			contents, _ := ioutil.ReadFile(filepath)
			syncer := &Syncer{snippets}
			replace := syncer.Sync(contents)
			_ = ioutil.WriteFile(filepath, replace, 0644)
		}
	}
}

func main() {
	sources := os.Args[1]
	target := os.Args[2]

	snippets := NewSnippets()
	for _, source := range strings.Split(sources, ",") {
		snippets.Merge(extractSnippets(source, ""))
	}

	fmt.Println("Snippets found:")
	snippets.Each(func (s *Snippet) {
		fmt.Println("\t" + string(s.Name))
	})

	fmt.Println("Files synced:")
	replaceSnippets(target, snippets)
}
