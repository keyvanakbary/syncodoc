package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
)

type Syncer struct {
	SnippetReader SnippetReader
	Formatter Formatter
	FileReader FileReader
}

type FileReader interface {
	Read(path string) ([]byte, error)
}

type IOUtilFileReader struct {
}

func (r *IOUtilFileReader) Read(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func language(u *url.URL) string {
	lang := u.Query().Get("lang")

	if lang != "" {
		return lang
	}

	ext := filepath.Ext(u.Path)
	if ext != "" {
		return ext[1:]//remove the dot
	}

	return ""
}

type Sync struct {
	Path string
	Lang string
	Part string
	LineStart int
	LineEnd int
}

func extractSyncs(content []byte) []Sync {
	var syncs []Sync

	beginSync := regexp.MustCompile(`^%% sync (.*)$`)
	endSync := regexp.MustCompile(`^%% end$`)

	var sync = Sync{}
	for num, line := range bytes.Split(content, []byte("\n")) {
		numLine := num + 1

		if beginSync.Match(line) {
			matches := beginSync.FindSubmatch(line)
			path := string(matches[1])
			u, err := url.Parse(path)
			if err != nil {
				fmt.Println("\t\tERROR: Invalid snippet path format " + path)
			}

			sync = Sync{Path: u.Path, Lang: language(u), Part: u.Fragment, LineStart: numLine, LineEnd: numLine}
		}

		if endSync.Match(line) {
			if sync == (Sync{}) {
				fmt.Println("\t\tERROR: No previous sync expression")
			} else {
				sync.LineEnd = numLine

				syncs = append(syncs, sync)

				sync = Sync{}
			}
		}
	}
	return syncs
}


func (s *Syncer) readSnippet(path string, sync Sync) (Snippet, error) {
	if sync.Part == "" {
		return s.SnippetReader.ReadFull(path)
	} else {
		return s.SnippetReader.ReadPart(path, sync.Part)
	}
}

func (s *Syncer) replaceSync(currentPath string, content []byte, sync Sync) ([]byte, error) {
	var syncedLines [][]byte

	snippetPath := sync.Path
	if !filepath.IsAbs(snippetPath) {
		dir := filepath.Dir(currentPath)
		snippetPath = filepath.Join(dir, snippetPath)
	}

	snippet, err := s.readSnippet(snippetPath, sync)
	if err != nil {
		return content, err
	}
	snippetLines := lines(s.Formatter.Format(sync.Lang, snippet))

	contentLines := lines(content)
	for i := 0; i < sync.LineStart; i++ {
		syncedLines = append(syncedLines, contentLines[i])
	}

	for _, line := range snippetLines {
		syncedLines = append(syncedLines, line)
	}

	for i := sync.LineEnd - 1; i < len(contentLines); i++ {
		syncedLines = append(syncedLines, contentLines[i])
	}

	return bytes.Join(syncedLines, []byte("\n")), nil
}

func (s *Syncer) Sync(currentPath string, content []byte) []byte {
	synced := content
	diff := 0
	for _, sync := range extractSyncs(content) {
		newSync := Sync{
			Lang: sync.Lang,
			Path: sync.Path,
			Part: sync.Part,
			LineStart: sync.LineStart + diff,
			LineEnd: sync.LineEnd + diff,
		}

		newS, err := s.replaceSync(currentPath, synced, newSync)
		if err != nil {
			fmt.Printf("\t\tERROR: Line %d, %s\n", newSync.LineStart, err.Error())
		} else {
			fmt.Printf("\t\tSynced snippet %s %s\n", newSync.Path, newSync.Part)
		}

		diff = len(lines(newS)) - len(lines(content))
		synced = newS
	}

	return synced
}

func (s *Syncer) SyncAll(root string) {
	_ = filepath.Walk(root, func(path string, _ os.FileInfo, err error) error {
		fmt.Println("\t" + path)
		content, _ := s.FileReader.Read(path)
		synced := s.Sync(path, content)
		_ = ioutil.WriteFile(path, synced, 0644)
		return nil
	})
}
