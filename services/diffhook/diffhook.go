package main

import (
	"github.com/bennettaur/diffhook/services/diffhook/trigger"
	"github.com/sourcegraph/go-diff/diff"
	"log"
	"os"
)

func main() {
	diffFile := os.Stdin
	defer func() {
		err := diffFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	r := diff.NewMultiFileDiffReader(diffFile)
	trigger.TriggerWatchers(r)
}
