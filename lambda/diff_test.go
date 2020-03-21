package main

import (
	"reflect"
	"sort"
	"testing"
)

func Test_diff_ToString(t *testing.T) {
	type fields struct {
		removed []string
		added   []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "no_changes",
			fields: fields{
				removed: []string{},
				added:   []string{},
			},
			want: "no changes",
		},
		{
			name: "removed",
			fields: fields{
				removed: []string{"1"},
				added:   []string{},
			},
			want: "removed: 1\nnone were added",
		},
		{
			name: "added",
			fields: fields{
				removed: []string{},
				added:   []string{"1"},
			},
			want: "none were removed\nadded: 1",
		},
		{
			name: "complex",
			fields: fields{
				removed: []string{"1", "2"},
				added:   []string{"3"},
			},
			want: "removed: 1, 2\nadded: 3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := diff{
				removed: tt.fields.removed,
				added:   tt.fields.added,
			}
			if got := d.ToString(); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateDiff(t *testing.T) {
	type args struct {
		old []string
		new []string
	}
	tests := []struct {
		name string
		args args
		want diff
	}{
		{
			name: "empty",
			args: args{
				old: []string{},
				new: []string{},
			},
			want: diff{
				removed: []string{},
				added:   []string{},
			},
		},
		{
			name: "no_change",
			args: args{
				old: []string{"1"},
				new: []string{"1"},
			},
			want: diff{
				removed: []string{},
				added:   []string{},
			},
		},
		{
			name: "add",
			args: args{
				old: []string{},
				new: []string{"1"},
			},
			want: diff{
				removed: []string{},
				added:   []string{"1"},
			},
		},
		{
			name: "remove",
			args: args{
				old: []string{"1"},
				new: []string{},
			},
			want: diff{
				removed: []string{"1"},
				added:   []string{},
			},
		},
		{
			name: "complex",
			args: args{
				old: []string{"1", "2", "3", "4"},
				new: []string{"1", "2", "5", "6"},
			},
			want: diff{
				removed: []string{"3", "4"},
				added:   []string{"5", "6"},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateDiff(tt.args.old, tt.args.new); !reflect.DeepEqual(got, tt.want) {
				sort.Strings(got.removed)
				sort.Strings(got.added)
				sort.Strings(tt.want.added)
				sort.Strings(tt.want.removed)
				t.Errorf("calculateDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_calculateDiff(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculateDiff([]string{"1", "2", "3", "4"}, []string{"1", "2", "5", "6"})
	}
}
