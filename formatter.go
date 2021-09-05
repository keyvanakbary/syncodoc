package main

import "fmt"

type Formatter interface {
	Format(language string, s Snippet) []byte
}

type MarkuaFormatter struct {}

func (f* MarkuaFormatter) Format(language string, s Snippet) []byte {
	return []byte(fmt.Sprintf(`{lang="%s"}
~~~~~~~~
%s
~~~~~~~~`, language, s))
}

type MarkdownFormatter struct {}

func (f* MarkdownFormatter) Format(language string, s Snippet) []byte {
	return []byte(fmt.Sprintf("```%s\n%s\n```", language, s))
}