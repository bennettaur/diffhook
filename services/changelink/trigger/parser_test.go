package trigger

import (
	"github.com/bennettaur/changelink/services/changelink/models"
	"testing"
)

func Test_findOverlapOneEach(t *testing.T) {
	type args struct {
		diffLines    []models.LineRange
		watcherLines []models.LineRange
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "No Overlap",
			args: args{
				diffLines:    []models.LineRange{{50, 100}},
				watcherLines: []models.LineRange{{20, 40}},
			},
			want: false,
		},
		{
			name: "Watcher end equals diff start",
			args: args{
				diffLines:    []models.LineRange{{50, 100}},
				watcherLines: []models.LineRange{{20, 50}},
			},
			want: true,
		},
		{
			name: "Watcher end overlaps diff start",
			args: args{
				diffLines:    []models.LineRange{{50, 100}},
				watcherLines: []models.LineRange{{20, 60}},
			},
			want: true,
		},
		{
			name: "Watcher start equals diff end",
			args: args{
				diffLines:    []models.LineRange{{50, 100}},
				watcherLines: []models.LineRange{{100, 160}},
			},
			want: true,
		},
		{
			name: "Watcher start overlaps diff end",
			args: args{
				diffLines:    []models.LineRange{{50, 100}},
				watcherLines: []models.LineRange{{90, 160}},
			},
			want: true,
		},
		{
			name: "Watcher contained in diff",
			args: args{
				diffLines:    []models.LineRange{{50, 100}},
				watcherLines: []models.LineRange{{60, 90}},
			},
			want: true,
		},
		{
			name: "Watcher contains diff",
			args: args{
				diffLines:    []models.LineRange{{50, 100}},
				watcherLines: []models.LineRange{{40, 110}},
			},
			want: true,
		},
		{
			name: "Watcher contains diff (1 line)",
			args: args{
				diffLines:    []models.LineRange{{50, 50}},
				watcherLines: []models.LineRange{{40, 110}},
			},
			want: true,
		},
		{
			name: "Watcher equals diff",
			args: args{
				diffLines:    []models.LineRange{{50, 100}},
				watcherLines: []models.LineRange{{50, 100}},
			},
			want: true,
		},
		{
			name: "Watcher equals diff (1 line)",
			args: args{
				diffLines:    []models.LineRange{{50, 50}},
				watcherLines: []models.LineRange{{50, 50}},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findOverlap(tt.args.diffLines, tt.args.watcherLines); got != tt.want {
				t.Errorf("findOverlap() = %v, want %v", got, tt.want)
			}

			// The reverse should also be true
			if got := findOverlap(tt.args.watcherLines, tt.args.diffLines); got != tt.want {
				t.Errorf("findOverlap() reversed = %v, want %v", got, tt.want)
			}
		})
	}
}


func Test_findOverlapMultiple(t *testing.T) {
	type args struct {
		diffLines    []models.LineRange
		watcherLines []models.LineRange
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "No Overlap, multi diff, one watch",
			args: args{
				diffLines:    []models.LineRange{
					{10, 10},
					{50, 60},
					{80, 100},
					{150, 200},
				},
				watcherLines: []models.LineRange{{20, 40}},
			},
			want: false,
		},
		{
			name: "No Overlap, multi diff, multi watch",
			args: args{
				diffLines:    []models.LineRange{
					{25, 25},
					{50, 60},
					{80, 100},
					{150, 200},
				},
				watcherLines: []models.LineRange{
					{0, 20},
					{30, 40},
					{61, 70},
					{120, 120},
					{210, 211},
				},
			},
			want: false,
		},
		{
			name: "Watcher end equals diff start",
			args: args{
				diffLines:    []models.LineRange{
					{30, 30},
					{50, 60},
					{80, 100},
					{150, 200},
				},
				watcherLines: []models.LineRange{
					{0, 20},
					{30, 40},
					{61, 80},
					{120, 120},
					{210, 211},
				},
			},
			want: true,
		},
		{
			name: "Watcher end overlaps diff start",
			args: args{
				diffLines:    []models.LineRange{
					{30, 30},
					{50, 60},
					{80, 100},
					{150, 200},
				},
				watcherLines: []models.LineRange{
					{0, 20},
					{30, 40},
					{61, 85},
					{120, 120},
					{210, 211},
				},
			},
			want: true,
		},
		{
			name: "Watcher start equals diff end",
			args: args{
				diffLines:    []models.LineRange{
					{30, 30},
					{50, 60},
					{80, 100},
					{150, 200},
				},
				watcherLines: []models.LineRange{
					{0, 20},
					{30, 40},
					{60, 70},
					{120, 120},
					{210, 211},
				},
			},
			want: true,
		},
		{
			name: "Watcher start overlaps diff end",
			args: args{
				diffLines:    []models.LineRange{
					{30, 30},
					{50, 60},
					{80, 100},
					{150, 200},
				},
				watcherLines: []models.LineRange{
					{0, 20},
					{30, 40},
					{55, 70},
					{120, 120},
					{210, 211},
				},
			},
			want: true,
		},
		{
			name: "Watcher contained in diff",
			args: args{
				diffLines:    []models.LineRange{
					{30, 30},
					{50, 60},
					{80, 100},
					{150, 200},
				},
				watcherLines: []models.LineRange{
					{0, 20},
					{30, 40},
					{81, 90},
					{120, 120},
					{210, 211},
				},
			},
			want: true,
		},
		{
			name: "Watcher contains diff",
			args: args{
				diffLines:    []models.LineRange{
					{30, 30},
					{50, 60},
					{80, 100},
					{150, 200},
				},
				watcherLines: []models.LineRange{
					{0, 20},
					{30, 40},
					{70, 110},
					{120, 120},
					{210, 211},
				},
			},
			want: true,
		},
		{
			name: "Watcher contains diff (1 line)",
			args: args{
				diffLines:    []models.LineRange{
					{30, 30},
					{50, 60},
					{80, 80},
					{150, 200},
				},
				watcherLines: []models.LineRange{
					{0, 20},
					{30, 40},
					{70, 90},
					{120, 120},
					{210, 211},
				},
			},
			want: true,
		},
		{
			name: "Watcher equals diff",
			args: args{
				diffLines:    []models.LineRange{
					{30, 30},
					{50, 60},
					{80, 100},
					{150, 200},
				},
				watcherLines: []models.LineRange{
					{0, 20},
					{30, 40},
					{80, 100},
					{120, 120},
					{210, 211},
				},
			},
			want: true,
		},
		{
			name: "Watcher equals diff (1 line)",
			args: args{
				diffLines:    []models.LineRange{
					{30, 30},
					{50, 60},
					{80, 80},
					{150, 200},
				},
				watcherLines: []models.LineRange{
					{0, 20},
					{30, 40},
					{80, 80},
					{120, 120},
					{210, 211},
				},
			},
			want: true,
		},
		{
			name: "Watcher has overlapping segments, but not with diff",
			args: args{
				diffLines:    []models.LineRange{
					{30, 30},
					{50, 60},
					{80, 100},
					{150, 200},
				},
				watcherLines: []models.LineRange{
					{0, 20},
					{31, 40},
					{35, 45},
					{120, 120},
					{210, 211},
				},
			},
			want: false,
		},
		{
			name: "Watcher has overlapping segments, and overlaps diff",
			args: args{
				diffLines:    []models.LineRange{
					{30, 30},
					{50, 60},
					{80, 100},
					{150, 200},
				},
				watcherLines: []models.LineRange{
					{0, 20},
					{31, 40},
					{35, 55},
					{120, 120},
					{210, 211},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findOverlap(tt.args.diffLines, tt.args.watcherLines); got != tt.want {
				t.Errorf("findOverlap() = %v, want %v", got, tt.want)
			}

			// The reverse should also be true
			if got := findOverlap(tt.args.watcherLines, tt.args.diffLines); got != tt.want {
				t.Errorf("findOverlap() reversed = %v, want %v", got, tt.want)
			}
		})
	}
}