package main

import (
	"log"

	"github.com/laenzlinger/setlist/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(cmd.Instance(), "./docs")
	if err != nil {
		log.Fatal(err)
	}
}
