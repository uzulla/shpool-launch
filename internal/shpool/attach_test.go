package shpool

import (
	"reflect"
	"testing"
)

func TestAttachArgs(t *testing.T) {
	cases := []struct {
		test  string
		name  string
		force bool
		want  []string
	}{
		{
			test:  "normal",
			name:  "my-session",
			force: false,
			want:  []string{"shpool", "attach", "--dir", ".", "--", "my-session"},
		},
		{
			test:  "force",
			name:  "my-session",
			force: true,
			want:  []string{"shpool", "attach", "-f", "--dir", ".", "--", "my-session"},
		},
		{
			test:  "dash-leading name stays positional after --",
			name:  "-repo-another",
			force: false,
			want:  []string{"shpool", "attach", "--dir", ".", "--", "-repo-another"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.test, func(t *testing.T) {
			got := attachArgs(tc.name, tc.force)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("attachArgs(%q, %v) = %v, want %v", tc.name, tc.force, got, tc.want)
			}
		})
	}
}
