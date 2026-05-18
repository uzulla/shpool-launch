package shpool

import (
	"reflect"
	"testing"
)

func TestParseList(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{
			name: "two simple rows",
			in: `work__company-a__api   attached
work__company-b__api   detached
`,
			want: []string{"work__company-a__api", "work__company-b__api"},
		},
		{
			name: "with header",
			in: `NAME                    STATUS
work__company-a__api    attached
sandbox__api-test       detached
`,
			want: []string{"work__company-a__api", "sandbox__api-test"},
		},
		{
			name: "blank lines skipped",
			in: `

work__company-a__api   attached

work__company-b__api   detached

`,
			want: []string{"work__company-a__api", "work__company-b__api"},
		},
		{
			name: "empty input",
			in:   ``,
			want: nil,
		},
		{
			name: "only header",
			in:   `NAME    STATUS`,
			want: nil,
		},
		{
			name: "single column",
			in: `foo
bar
baz`,
			want: []string{"foo", "bar", "baz"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ParseList([]byte(tc.in))
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ParseList = %v, want %v", got, tc.want)
			}
		})
	}
}
