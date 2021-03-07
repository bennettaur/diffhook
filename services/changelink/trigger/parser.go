package trigger

import (
	"fmt"
	"github.com/bennettaur/changelink/services/changelink/models"
	"io"
	"log"
	"os"
	"sort"

	"github.com/sourcegraph/go-diff/diff"
)

func ParseDiff() {
	diffFile := os.Stdin
	defer func() {
		err := diffFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	r := diff.NewMultiFileDiffReader(diffFile)
	for i := 0; ; i++ {
		fileIndex := fmt.Sprintf("file(%d)", i)
		fileDiff, err := r.ReadFile()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("err reading file %s: %s", fileIndex, err)
			continue
		}

		// Assumes hunks are sorted
		changedLineRanges := getDiffLineRanges(fileDiff)
		watchers, err := models.FindWatchersForFile(fileDiff.OrigName)

		if err != nil {
			log.Printf("err getting watchers for file %s: %s", fileDiff.OrigName, err)
			continue
		}

		var actionsToRun []models.Action

		for _, watcher := range watchers {
			sort.Slice(watcher.Lines, func(i, j int) bool {
				return watcher.Lines[i].StartLine < watcher.Lines[j].StartLine
			})

			if len(watcher.Lines) == 0 || findOverlap(changedLineRanges, watcher.Lines){
				actionsToRun = append(actionsToRun, watcher.Actions...)
				continue
			}
		}

		for _, action := range actionsToRun {
			fmt.Printf("Triggered %s action: %s\n", action.ActionType, action.Message)
		}
	}
}

func getDiffLineRanges(diff *diff.FileDiff) []models.LineRange {
	var ranges []models.LineRange
	for _, hunk := range diff.Hunks {
		ranges = append(
			ranges,
			models.LineRange{
				StartLine: int(hunk.OrigStartLine),
				EndLine:   int(hunk.OrigStartLine + hunk.OrigLines),
			},
		)
	}
	return ranges
}

func findOverlap(diffLines, watcherLines []models.LineRange) bool {
	var watchIndex, diffIndex int

	for {
		// Check if we've run out of diffs or watchers to check
		if watchIndex >= len(watcherLines) || diffIndex >= len(diffLines) {
			return false
		}

		// Check we overlap or are inside the diff range
		if watcherLines[watchIndex].StartLine >= diffLines[diffIndex].StartLine {
			// Check if start line is within the diff range and return
			if watcherLines[watchIndex].StartLine <= diffLines[diffIndex].EndLine {
				return true
			}

			// Current watchline is greater than diff, so advance the diff
			diffIndex += 1
			continue
		}

		// Check if our end overlaps the start
		if watcherLines[watchIndex].EndLine >= diffLines[diffIndex].StartLine {
			return true
		}

		// Watch line occurs before diff range, so move on to the next watcher LineRange forward
		watchIndex += 1
	}
}