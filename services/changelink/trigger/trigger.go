package trigger

import (
	"fmt"
	"io"
	"log"
	"path/filepath"
	"sort"

	"github.com/bennettaur/changelink/services/changelink/models"
	"github.com/bennettaur/changelink/services/changelink/models/actions"
	"github.com/sourcegraph/go-diff/diff"
)

type TriggeredWatcher struct {
	FileDiff       *diff.FileDiff
	TriggeredLines *actions.TriggeredLines
	Watcher        models.Watcher
	Reason         string
}

func TriggerWatchers(diffReader *diff.MultiFileDiffReader) []TriggeredWatcher {
	log.Println("Starting")

	var triggeredWatchers []TriggeredWatcher
	for i := 0; ; i++ {
		fileIndex := fmt.Sprintf("file(%d)", i)
		log.Printf("Reading %s", fileIndex)
		fileDiff, err := diffReader.ReadFile()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("err reading file %s: %s", fileIndex, err)
			continue
		}

		// Assumes hunks are sorted
		changedLineRanges := getDiffLineRanges(fileDiff)
		log.Printf("Found the following line changes in %s: %v", fileIndex, changedLineRanges)
		watchers, err := models.FindWatchersForFile(fileDiff.OrigName)

		if err != nil {
			log.Printf("err getting watchers for file %s: %s", fileDiff.OrigName, err)
			continue
		}

		for _, watcher := range watchers {
			log.Printf("Checking watcher %s", watcher.Name)
			sort.Slice(watcher.Lines, func(i, j int) bool {
				return watcher.Lines[i].StartLine < watcher.Lines[j].StartLine
			})

			var triggeredWatcher *TriggeredWatcher

			if triggered, reason := specialTrigger(watcher, changedLineRanges, fileDiff); triggered {
				triggeredWatcher = &TriggeredWatcher{
					FileDiff:       fileDiff,
					TriggeredLines: nil,
					Watcher:        watcher,
					Reason:         reason,
				}
			} else {
				triggeredLines := findOverlap(changedLineRanges, watcher.Lines, fileDiff)
				if triggeredLines != nil {
					triggeredWatcher = &TriggeredWatcher{
						FileDiff:       fileDiff,
						TriggeredLines: triggeredLines,
						Watcher:        watcher,
						Reason:         "Watched lines changed",
					}
				}
			}

			if triggeredWatcher != nil {
				triggeredWatchers = append(triggeredWatchers, *triggeredWatcher)
			}
		}
	}
	return triggeredWatchers
}

func getDiffLineRanges(fileDiff *diff.FileDiff) []actions.LineRange {
	var ranges []actions.LineRange
	for _, hunk := range fileDiff.Hunks {
		// Git unified diffs include 3 lines before and after the actual hunk changes
		changeStartLine := int(hunk.OrigStartLine) + 3
		changeEndLine := changeStartLine + int(hunk.OrigLines) - 6

		ranges = append(
			ranges,
			actions.LineRange{
				StartLine: changeStartLine,
				EndLine:   changeEndLine,
			},
		)
	}
	return ranges
}

// Compares the two sets of line ranges for any overlaps and returns the indices of the overlapping line ranges
// or -1 if no overlap is found
func findOverlap(diffLines, watchedLines []actions.LineRange, fileDiff *diff.FileDiff) *actions.TriggeredLines {
	var diffIndex, watchIndex int

	for {
		// Check if we've run out of diffs or watchers to check
		if watchIndex >= len(watchedLines) || diffIndex >= len(diffLines) {
			return nil
		}

		// Check we overlap or are inside the diff range
		if watchedLines[watchIndex].StartLine >= diffLines[diffIndex].StartLine {
			// Check if start line is within the diff range and return
			if watchedLines[watchIndex].StartLine <= diffLines[diffIndex].EndLine {
				return &actions.TriggeredLines{
					DiffLines: diffLines[diffIndex],
					WatchedLines: watchedLines[watchIndex],
					Hunk: fileDiff.Hunks[diffIndex],
				}
			}

			// Current watchline is greater than diff, so advance the diff
			diffIndex += 1
			continue
		}

		// Check if our end overlaps the start
		if watchedLines[watchIndex].EndLine >= diffLines[diffIndex].StartLine {
			return &actions.TriggeredLines{
				DiffLines: diffLines[diffIndex],
				WatchedLines: watchedLines[watchIndex],
				Hunk: fileDiff.Hunks[diffIndex],
			}
		}

		// Watch line occurs before diff range, so move on to the next watcher LineRange forward
		watchIndex += 1
	}
}

func renamed(fileDiff *diff.FileDiff) bool {
	origName := filepath.Base(fileDiff.OrigName)
	newName := filepath.Base(fileDiff.NewName)
	return origName != newName && !deleted(fileDiff)
}

func moved(fileDiff *diff.FileDiff) bool {
	// Need to drop the first letter from each path because they'll always differ (i.e. a, b)
	origDir := filepath.Dir(fileDiff.OrigName[1:])
	newDir := filepath.Dir(fileDiff.NewName[1:])
	return origDir != newDir && !deleted(fileDiff)
}

func deleted(fileDiff *diff.FileDiff) bool {
	return fileDiff.NewName == "/dev/null"
}

func modeChanged(fileDiff *diff.FileDiff) bool {
	return false
}

func specialTrigger(w models.Watcher, changedLines []actions.LineRange, fileDiff *diff.FileDiff) (bool, string) {
	if w.TriggerAny {
		return true, "Any Change"
	}
	if w.TriggerAnyLine && len(changedLines) > 0 {
		return true, "Any Line"
	}
	if w.TriggerOnRename && renamed(fileDiff) {
		return true, "File Renamed"
	}
	if w.TriggerOnMove && moved(fileDiff) {
		return true, "File Moved"
	}
	if w.TriggerOnDelete && deleted(fileDiff) {
		return true, "File Deleted"
	}
	if w.TriggerOnMode && modeChanged(fileDiff) {
		return true, "File Mode Changed"
	}
	return false, ""
}