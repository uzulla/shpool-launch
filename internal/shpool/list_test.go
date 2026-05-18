package shpool

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestListUsesTextList(t *testing.T) {
	withFakeShpool(t, `#!/bin/sh
if [ "$1" = "list" ] && [ "$#" -eq 1 ]; then
	printf 'NAME\tSTATUS\nfoo\tdisconnected\n'
	exit 0
fi
echo "unexpected args: $*" >&2
exit 2
`)

	got, err := List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	want := []string{"foo"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("List() = %v, want %v", got, want)
	}
}

func TestListReturnsListUnavailable(t *testing.T) {
	withFakeShpool(t, `#!/bin/sh
exit 1
`)

	_, err := List()
	if !errors.Is(err, ErrListUnavailable) {
		t.Fatalf("List() error = %v, want ErrListUnavailable", err)
	}
}

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
		{
			name: "uppercase session is not treated as header",
			in: `NAME    STATUS
API     disconnected
`,
			want: []string{"API"},
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

func TestListUnavailable(t *testing.T) {
	cases := []struct {
		msg  string
		want bool
	}{
		{msg: "", want: true},
		{msg: "could not connect to daemon", want: true},
		{msg: "Error: daemonizing: launched daemon, but control socket never came up", want: true},
		{msg: "connection refused", want: true},
		{msg: "permission denied", want: false},
	}
	for _, tc := range cases {
		t.Run(tc.msg, func(t *testing.T) {
			got := listUnavailable(tc.msg)
			if got != tc.want {
				t.Errorf("listUnavailable(%q) = %v, want %v", tc.msg, got, tc.want)
			}
		})
	}
}

func withFakeShpool(t *testing.T, script string) {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "shpool")
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake shpool: %v", err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}
