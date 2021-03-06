package trigger

import (
	"os"
	"testing"

	"github.com/bennettaur/diffhook/services/diffhook/models"
	"github.com/bennettaur/diffhook/services/diffhook/models/actions"
	"github.com/sourcegraph/go-diff/diff"
	"github.com/stretchr/testify/assert"
)

func Test_findOverlapOneEach(t *testing.T) {
	type args struct {
		diffLines    []actions.LineRange
		watcherLines []actions.LineRange
	}
	tests := []struct {
		name      string
		args      args
		wantFound bool
	}{
		{
			name: "No Overlap",
			args: args{
				diffLines:    []actions.LineRange{{50, 100}},
				watcherLines: []actions.LineRange{{20, 40}},
			},
			wantFound: false,
		},
		{
			name: "Watcher end equals diff start",
			args: args{
				diffLines:    []actions.LineRange{{50, 100}},
				watcherLines: []actions.LineRange{{20, 50}},
			},
			wantFound: true,
		},
		{
			name: "Watcher end overlaps diff start",
			args: args{
				diffLines:    []actions.LineRange{{50, 100}},
				watcherLines: []actions.LineRange{{20, 60}},
			},
			wantFound: true,
		},
		{
			name: "Watcher start equals diff end",
			args: args{
				diffLines:    []actions.LineRange{{50, 100}},
				watcherLines: []actions.LineRange{{100, 160}},
			},
			wantFound: true,
		},
		{
			name: "Watcher start overlaps diff end",
			args: args{
				diffLines:    []actions.LineRange{{50, 100}},
				watcherLines: []actions.LineRange{{90, 160}},
			},
			wantFound: true,
		},
		{
			name: "Watcher contained in diff",
			args: args{
				diffLines:    []actions.LineRange{{50, 100}},
				watcherLines: []actions.LineRange{{60, 90}},
			},
			wantFound: true,
		},
		{
			name: "Watcher contains diff",
			args: args{
				diffLines:    []actions.LineRange{{50, 100}},
				watcherLines: []actions.LineRange{{40, 110}},
			},
			wantFound: true,
		},
		{
			name: "Watcher contains diff (1 line)",
			args: args{
				diffLines:    []actions.LineRange{{50, 50}},
				watcherLines: []actions.LineRange{{40, 110}},
			},
			wantFound: true,
		},
		{
			name: "Watcher equals diff",
			args: args{
				diffLines:    []actions.LineRange{{50, 100}},
				watcherLines: []actions.LineRange{{50, 100}},
			},
			wantFound: true,
		},
		{
			name: "Watcher equals diff (1 line)",
			args: args{
				diffLines:    []actions.LineRange{{50, 50}},
				watcherLines: []actions.LineRange{{50, 50}},
			},
			wantFound: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFileDiff := makeMockFileDiff(tt.args.diffLines)
			got := findOverlap(tt.args.diffLines, tt.args.watcherLines, mockFileDiff)

			expected := &actions.TriggeredLines{DiffLines: tt.args.diffLines[0], WatchedLines: tt.args.watcherLines[0]}

			if tt.wantFound && !equalTriggeredLines(got, expected) ||
				!tt.wantFound && got != nil {
				t.Errorf("findOverlap() = %v, want %v", got, expected)
			}

			// The reverse should also be true
			mockFileDiff = makeMockFileDiff(tt.args.watcherLines)
			expected = &actions.TriggeredLines{DiffLines: tt.args.watcherLines[0], WatchedLines: tt.args.diffLines[0]}
			got = findOverlap(tt.args.watcherLines, tt.args.diffLines, mockFileDiff)
			if tt.wantFound && !equalTriggeredLines(got, expected) ||
				!tt.wantFound && got != nil {
				t.Errorf("findOverlap() reverse = %v, want %v", got, expected)
			}
		})
	}
}

func Test_findOverlapMultiple(t *testing.T) {
	type args struct {
		diffLines    []actions.LineRange
		watcherLines []actions.LineRange
	}
	tests := []struct {
		name string
		args args
		want *actions.TriggeredLines
	}{
		{
			name: "No Overlap, multi diff, one watch",
			args: args{
				diffLines: []actions.LineRange{
					{10, 10},
					{50, 60},
					{80, 100},
					{150, 200},
				},
				watcherLines: []actions.LineRange{{20, 40}},
			},
			want: nil,
		},
		{
			name: "No Overlap, multi diff, multi watch",
			args: args{
				diffLines: []actions.LineRange{
					{25, 25},
					{50, 60},
					{80, 100},
					{150, 200},
				},
				watcherLines: []actions.LineRange{
					{0, 20},
					{30, 40},
					{61, 70},
					{120, 120},
					{210, 211},
				},
			},
			want: nil,
		},
		{
			name: "Watcher end equals diff start",
			args: args{
				diffLines: []actions.LineRange{
					{30, 30},
					{50, 60},
					{80, 100}, // Trigger
					{150, 200},
				},
				watcherLines: []actions.LineRange{
					{0, 20},
					{35, 40},
					{61, 80}, // Trigger
					{120, 120},
					{210, 211},
				},
			},
			want: &actions.TriggeredLines{
				DiffLines:    actions.LineRange{StartLine: 80, EndLine: 100},
				WatchedLines: actions.LineRange{StartLine: 61, EndLine: 80},
			},
		},
		{
			name: "Watcher end overlaps diff start",
			args: args{
				diffLines: []actions.LineRange{
					{30, 30},
					{50, 60},
					{80, 100}, // Trigger
					{150, 200},
				},
				watcherLines: []actions.LineRange{
					{0, 20},
					{31, 40},
					{61, 85}, // Trigger
					{120, 120},
					{210, 211},
				},
			},
			want: &actions.TriggeredLines{
				DiffLines:    actions.LineRange{StartLine: 80, EndLine: 100},
				WatchedLines: actions.LineRange{StartLine: 61, EndLine: 85},
			},
		},
		{
			name: "Watcher start equals diff end",
			args: args{
				diffLines: []actions.LineRange{
					{30, 30},
					{50, 60}, // Trigger
					{80, 100},
					{150, 200},
				},
				watcherLines: []actions.LineRange{
					{0, 20},
					{35, 40},
					{60, 70}, // Trigger
					{120, 120},
					{210, 211},
				},
			},
			want: &actions.TriggeredLines{
				DiffLines:    actions.LineRange{StartLine: 50, EndLine: 60},
				WatchedLines: actions.LineRange{StartLine: 60, EndLine: 70},
			},
		},
		{
			name: "Watcher start overlaps diff end",
			args: args{
				diffLines: []actions.LineRange{
					{30, 30},
					{50, 60}, // Trigger
					{80, 100},
					{150, 200},
				},
				watcherLines: []actions.LineRange{
					{0, 20},
					{35, 40},
					{55, 70}, // Trigger
					{120, 120},
					{210, 211},
				},
			},
			want: &actions.TriggeredLines{
				DiffLines:    actions.LineRange{StartLine: 50, EndLine: 60},
				WatchedLines: actions.LineRange{StartLine: 55, EndLine: 70},
			},
		},
		{
			name: "Watcher contained in diff",
			args: args{
				diffLines: []actions.LineRange{
					{30, 30},
					{50, 60},
					{80, 100}, // Trigger
					{150, 200},
				},
				watcherLines: []actions.LineRange{
					{0, 20},
					{35, 40},
					{81, 90}, // Trigger
					{120, 120},
					{210, 211},
				},
			},
			want: &actions.TriggeredLines{
				DiffLines:    actions.LineRange{StartLine: 80, EndLine: 100},
				WatchedLines: actions.LineRange{StartLine: 81, EndLine: 90},
			},
		},
		{
			name: "Watcher contains diff",
			args: args{
				diffLines: []actions.LineRange{
					{30, 30},
					{50, 60},
					{80, 100}, // Trigger
					{150, 200},
				},
				watcherLines: []actions.LineRange{
					{0, 20},
					{31, 40},
					{70, 110}, // Trigger
					{120, 120},
					{210, 211},
				},
			},
			want: &actions.TriggeredLines{
				DiffLines:    actions.LineRange{StartLine: 80, EndLine: 100},
				WatchedLines: actions.LineRange{StartLine: 70, EndLine: 110},
			},
		},
		{
			name: "Watcher contains diff (1 line)",
			args: args{
				diffLines: []actions.LineRange{
					{30, 30},
					{50, 60},
					{80, 80}, // Trigger
					{150, 200},
				},
				watcherLines: []actions.LineRange{
					{0, 20},
					{35, 40},
					{70, 90}, // Trigger
					{120, 120},
					{210, 211},
				},
			},
			want: &actions.TriggeredLines{
				DiffLines:    actions.LineRange{StartLine: 80, EndLine: 80},
				WatchedLines: actions.LineRange{StartLine: 70, EndLine: 90},
			},
		},
		{
			name: "Watcher equals diff",
			args: args{
				diffLines: []actions.LineRange{
					{30, 30},
					{50, 60},
					{80, 100}, // Trigger
					{150, 200},
				},
				watcherLines: []actions.LineRange{
					{0, 20},
					{35, 40},
					{80, 100}, // Trigger
					{120, 120},
					{210, 211},
				},
			},
			want: &actions.TriggeredLines{
				DiffLines:    actions.LineRange{StartLine: 80, EndLine: 100},
				WatchedLines: actions.LineRange{StartLine: 80, EndLine: 100},
			},
		},
		{
			name: "Watcher equals diff (1 line)",
			args: args{
				diffLines: []actions.LineRange{
					{30, 30},
					{50, 60},
					{80, 80}, // Trigger
					{150, 200},
				},
				watcherLines: []actions.LineRange{
					{0, 20},
					{35, 40},
					{80, 80}, // Trigger
					{120, 120},
					{210, 211},
				},
			},
			want: &actions.TriggeredLines{
				DiffLines:    actions.LineRange{StartLine: 80, EndLine: 80},
				WatchedLines: actions.LineRange{StartLine: 80, EndLine: 80},
			},
		},
		{
			name: "Watcher has overlapping segments, but not with diff",
			args: args{
				diffLines: []actions.LineRange{
					{30, 30},
					{50, 60},
					{80, 100},
					{150, 200},
				},
				watcherLines: []actions.LineRange{
					{0, 20},
					{31, 40},
					{35, 45},
					{120, 120},
					{210, 211},
				},
			},
			want: nil,
		},
		{
			name: "Watcher has overlapping segments, and overlaps diff",
			args: args{
				diffLines: []actions.LineRange{
					{30, 30},
					{50, 60}, // Trigger
					{80, 100},
					{150, 200},
				},
				watcherLines: []actions.LineRange{
					{0, 20},
					{31, 40},
					{35, 55}, // Trigger
					{120, 120},
					{210, 211},
				},
			},
			want: &actions.TriggeredLines{
				DiffLines:    actions.LineRange{StartLine: 50, EndLine: 60},
				WatchedLines: actions.LineRange{StartLine: 35, EndLine: 55},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFileDiff := makeMockFileDiff(tt.args.diffLines)
			if got := findOverlap(tt.args.diffLines, tt.args.watcherLines, mockFileDiff); !equalTriggeredLines(got, tt.want) {
				t.Errorf("findOverlap() = %v, want %v", got, tt.want)
			}

			// The reverse should also be true
			if tt.want != nil {
				tt.want.WatchedLines, tt.want.DiffLines = tt.want.DiffLines, tt.want.WatchedLines
			}
			mockFileDiff = makeMockFileDiff(tt.args.watcherLines)
			if got := findOverlap(tt.args.watcherLines, tt.args.diffLines, mockFileDiff); !equalTriggeredLines(got, tt.want) {
				t.Errorf("findOverlap() reversed = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetActions(t *testing.T) {
	tests := []struct {
		name             string
		watcherFixture   string
		storeFixture     string
		wantWatcherNames []string
	}{
		{
			name:           "one line changed",
			watcherFixture: "../../../test/one_line.diff",
			storeFixture:   "../../../test/.diffhook.yml",
			wantWatcherNames: []string{
				"Slack Watcher",
				"Any Line Log Watch",
				"Any Log Watch",
			},
		},
		{
			name:           "multiple lines and chunks changed",
			watcherFixture: "../../../test/multiple.diff",
			storeFixture:   "../../../test/.diffhook.yml",
			wantWatcherNames: []string{
				"Slack Watcher",
				"Multiple Line Log Watch",
				"Any Line Log Watch",
				"Any Log Watch",
			},
		},
		{
			name:           "file renamed",
			watcherFixture: "../../../test/rename.diff",
			storeFixture:   "../../../test/.diffhook.yml",
			wantWatcherNames: []string{"Rename Log Watch", "Any Log Watch"},
		},
		{
			name:           "file moved",
			watcherFixture: "../../../test/move.diff",
			storeFixture:   "../../../test/.diffhook.yml",
			wantWatcherNames: []string{"Move Log Watch", "Any Log Watch"},
		},
		{
			name:           "file perms changed",
			watcherFixture: "../../../test/mode.diff",
			storeFixture:   "../../../test/.diffhook.yml",
			wantWatcherNames: []string{"Permission Log Watch", "Any Log Watch"},
		},
		{
			name:           "file deleted",
			watcherFixture: "../../../test/delete.diff",
			storeFixture:   "../../../test/.diffhook.yml",
			wantWatcherNames: []string{
				"Delete Log Watch",
				"Any Log Watch",
				"Slack Watcher",
				"Multiple Line Log Watch",
				"Any Line Log Watch",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			models.SetLocalStore(tt.storeFixture)
			f, err := os.Open(tt.watcherFixture)
			if err != nil {
				t.Errorf("Error opening file: %s", err)
				return
			}
			defer f.Close()
			mr := diff.NewMultiFileDiffReader(f)
			watchers := TriggerWatchers(mr)
			var watcherNames []string
			for _, w := range watchers {
				watcherNames = append(watcherNames, w.Watcher.Name)
			}

			assert.ElementsMatch(t, tt.wantWatcherNames, watcherNames)
		})
	}
}

func equalTriggeredLines(x, y *actions.TriggeredLines) bool {
	if x == nil && y == nil {
		return true
	} else if x == nil || y == nil {
		return false
	}

	if x.DiffLines.StartLine == y.DiffLines.StartLine &&
		x.DiffLines.EndLine == y.DiffLines.EndLine &&
		x.WatchedLines.StartLine == y.WatchedLines.StartLine &&
		x.WatchedLines.EndLine == y.WatchedLines.EndLine {
		return true
	}
	return false
}

func makeMockFileDiff(diffLines []actions.LineRange) *diff.FileDiff {
	return &diff.FileDiff{
		Hunks: make([]*diff.Hunk, len(diffLines)),
	}
}