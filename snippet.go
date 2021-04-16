package main

import "fmt"

type Snippet struct {
	Name []byte
	Contents []byte
}

type Snippets struct {
	Map map[string]*Snippet
}

func NewSnippets() *Snippets {
	return &Snippets{map[string]*Snippet{}}
}

func (ss *Snippets) ByName(name []byte) (*Snippet, error) {
	snippet := ss.Map[string(name)]
	if snippet == nil {
		return nil, fmt.Errorf("no contents found for %ss", name)
	}

	return snippet, nil
}

func (ss *Snippets) Contains(name []byte) bool {
	_, err := ss.ByName(name)

	return err == nil
}

func (ss *Snippets) Empty() bool {
	return len(ss.Map) == 0
}

func (ss *Snippets) Each(fn func(*Snippet)) {
	for _, value := range ss.Map {
		fn(value)
	}
}

func (ss *Snippets) Put(s *Snippet) *Snippets {
	ss.Map[string(s.Name)] = s

	return ss
}

func (ss *Snippets) Merge(s *Snippets) {
	s.Each(func(snippet *Snippet) {
		ss.Put(snippet)
	})
}
