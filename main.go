package main

import (
	"fmt"
	"os"
)

func main() {
	target := os.Args[1]

	fmt.Println("Files synced:")
	syncer := &Syncer{&FileSnippetReader{}, &MarkuaFormatter{}, &IOUtilFileReader{}}
	syncer.SyncAll(target)
}
