package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
)

type Snippet []byte
type Snippets map[string]Snippet

type SnippetReader interface {
	ReadFull(path string) (Snippet, error)
	ReadPart(path string, name string) (Snippet, error)
}

type FileSnippetReader struct {

}

func (r *FileSnippetReader) ReadFull(path string) (Snippet, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snippet file not found %s", path)
	}

	return ignore(contents), nil
}

func (r *FileSnippetReader) ReadPart(path string, name string) (Snippet, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snippets file not found %s", path)
	}

	snippets := extract(contents)
	snippet, ok := snippets[name]
	if !ok {
		return nil, fmt.Errorf("no snippet %s found in %s", name, path)
	}

	return snippet, nil
}

func numSpaces(line []byte) int {
	num := 0
	for _, c := range line {
		if c == byte(' ') || c == byte('\t') {
			num++
		} else {
			break
		}
	}

	return num
}

func trimSpaces(line []byte, maxSpaces int) []byte {
	lineSpaces := numSpaces(line)

	if maxSpaces > lineSpaces {
		return line[lineSpaces:]
	} else {
		return line[maxSpaces:]
	}
}

func lines(contents []byte) [][]byte {
	return bytes.Split(contents, []byte("\n"))
}

func trim(snippet Snippet) Snippet {
	ls := lines(snippet)
	if len(ls) == 0 {
		return snippet
	}

	length := numSpaces(ls[0])
	var trimmed [][]byte
	for _, line := range ls {
		trimmed = append(trimmed, trimSpaces(line, length))
	}
	return bytes.Join(trimmed, []byte("\n"))
}

func ignore(contents []byte) []byte {
	beginIgnoreRe := regexp.MustCompile(`(?://|{#|<!--|--|#)\s*ignore`)
	endIgnoreRe := regexp.MustCompile(`(?://|{#|<!--|--|#)\s*end-ignore`)

	var snippetLines [][]byte
	ignore := false
	for _, line := range lines(contents) {

		if beginIgnoreRe.Match(line) {
			ignore = true
			nSpaces := numSpaces(line)
			snippetLines = append(snippetLines, append(line[:nSpaces], "// ..."...))
		}

		if ignore == false {
			snippetLines = append(snippetLines, line)
		}

		if endIgnoreRe.Match(line) {
			ignore = false
		}
	}

	return trim(bytes.Join(snippetLines, []byte("\n")))
}

func extract(contents []byte) Snippets {
	contents = ignore(contents)
	snippets := make(Snippets)

	beginSnippetRe := regexp.MustCompile(`(?://|{#|<!--|--|#)\s*snippet ([a-zA-Z0-9-]+)`)
	endSnippetRe := regexp.MustCompile(`(?://|{#|<!--|--|#)\s*end-snippet`)

	var snippetName string
	var snippetLines [][]byte
	for _, line := range lines(contents) {
		if endSnippetRe.Match(line) {
			snippets[snippetName] = trim(bytes.Join(snippetLines, []byte("\n")))
			snippetName = ""
			snippetLines = nil
		}

		if snippetLines != nil {
			snippetLines = append(snippetLines, line)
		}

		if beginSnippetRe.Match(line) {
			matches := beginSnippetRe.FindSubmatch(line)
			snippetName = string(matches[1])
			snippetLines = [][]byte{}
		}
	}

	return snippets
}
