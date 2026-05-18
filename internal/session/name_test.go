package session

import "testing"

func TestFromPath(t *testing.T) {
	cases := []struct {
		name string
		path string
		home string
		want string
	}{
		{
			name: "home-relative simple",
			path: "/home/uzulla/work/company-a/api",
			home: "/home/uzulla",
			want: "work.company-a.api",
		},
		{
			name: "spaces replaced",
			path: "/Users/uzulla/src/foo bar/api",
			home: "/Users/uzulla",
			want: "src.foo_bar.api",
		},
		{
			name: "outside home",
			path: "/tmp/test/api",
			home: "/home/uzulla",
			want: "tmp.test.api",
		},
		{
			name: "no home given",
			path: "/tmp/test/api",
			home: "",
			want: "tmp.test.api",
		},
		{
			name: "home itself",
			path: "/home/uzulla",
			home: "/home/uzulla",
			want: "",
		},
		{
			name: "trailing slash",
			path: "/home/uzulla/work/",
			home: "/home/uzulla",
			want: "work",
		},
		{
			name: "japanese chars replaced",
			path: "/Users/uzulla/プロジェクト/api",
			home: "/Users/uzulla",
			// 6 runes (プロジェクト) → 6 "_" + "." separator + "api"
			want: "______.api",
		},
		{
			name: "dots in dirname preserved",
			path: "/srv/app-v1.2.3/api",
			home: "",
			want: "srv.app-v1.2.3.api",
		},
		{
			name: "underscores preserved",
			path: "/srv/my_app/web",
			home: "",
			want: "srv.my_app.web",
		},
		{
			name: "home prefix-only match should not match other dir",
			path: "/home/uzulla-other/work",
			home: "/home/uzulla",
			want: "home.uzulla-other.work",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := FromPath(tc.path, tc.home)
			if got != tc.want {
				t.Errorf("FromPath(%q, %q) = %q, want %q", tc.path, tc.home, got, tc.want)
			}
		})
	}
}

func TestFromPath_JapaneseLen(t *testing.T) {
	got := FromPath("/Users/u/プロジェクト", "/Users/u")
	want := "______"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
