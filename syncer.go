package main

import (
	"fmt"
	"regexp"
)

type Syncer struct {
	Snippets *Snippets
}

func MarkuaFormat(language []byte, s *Snippet) []byte {
	return []byte(fmt.Sprintf(`%%%% sync %s %s
{lang="%s"}
~~~~~~~~
%s
~~~~~~~~
%%%%`, language, s.Name, language, s.Contents))
}

func (s *Syncer) Sync(content []byte) []byte {
	re := regexp.MustCompile(`%% sync (\w+) ([a-zA-Z0-9/-]+)\n(?:[^%]*|%[^%]*)%%`)
	return re.ReplaceAllFunc(content, func(match []byte) []byte {
		matches := re.FindSubmatch(match)

		snippet, err := s.Snippets.ByName(matches[2])

		if err != nil {
			fmt.Println("\t\tERROR: " + string(matches[2]) + " not found")
			return matches[0]
		}

		fmt.Println("\t\t" + string(matches[2]))

		return MarkuaFormat(matches[1], snippet)
	})
}
