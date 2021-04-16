package main

import (
	"bytes"
	"regexp"
)

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

func trim(snippet []byte) []byte {
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

func Extract(contents []byte) *Snippets {
	snippets := NewSnippets()

	beginSnippetRe := regexp.MustCompile(`(?://|{#|<!--|--|#)\s*snippet ([a-zA-Z0-9-]+)`)
	endSnippetRe := regexp.MustCompile(`(?://|{#|<!--|--|#)\s*end-snippet`)

	beginIgnoreRe := regexp.MustCompile(`(?://|{#|<!--|--|#)\s*ignore`)
	endIgnoreRe := regexp.MustCompile(`(?://|{#|<!--|--|#)\s*end-ignore`)

	var sName []byte
	var sLines [][]byte
	ignore := false
	for _, line := range lines(contents) {
		if endSnippetRe.Match(line) {
			snippets.Put(&Snippet{sName, trim(bytes.Join(sLines, []byte("\n")))})
			sName = nil
			sLines = nil
			ignore = false
		}

		if beginIgnoreRe.Match(line) {
			ignore = true
			nSpaces := numSpaces(line)
			sLines = append(sLines, append(line[:nSpaces], "// ..."...))
		}

		if sLines != nil && ignore == false {
			sLines = append(sLines, line)
		}

		if beginSnippetRe.Match(line) {
			matches := beginSnippetRe.FindSubmatch(line)
			sName = matches[1]
			sLines = [][]byte{}
		}

		if endIgnoreRe.Match(line) {
			ignore = false
		}

	}

	return snippets
}
