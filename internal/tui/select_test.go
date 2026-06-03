package tui

import (
	"reflect"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func asModel(t *testing.T, updated tea.Model) model {
	t.Helper()
	m, ok := updated.(model)
	if !ok {
		t.Fatalf("Update returned %T, want model", updated)
	}
	return m
}

func TestOrderWithDefault(t *testing.T) {
	cases := []struct {
		name        string
		items       []string
		defaultItem string
		want        []string
	}{
		{
			name:        "default not in items is prepended",
			items:       []string{"a", "b"},
			defaultItem: "x",
			want:        []string{"x", "a", "b"},
		},
		{
			name:        "default already present is moved to front",
			items:       []string{"a", "b", "x", "c"},
			defaultItem: "x",
			want:        []string{"x", "a", "b", "c"},
		},
		{
			name:        "empty default keeps order",
			items:       []string{"a", "b"},
			defaultItem: "",
			want:        []string{"a", "b"},
		},
		{
			name:        "no items but default",
			items:       nil,
			defaultItem: "x",
			want:        []string{"x"},
		},
		{
			name:        "all empty",
			items:       nil,
			defaultItem: "",
			want:        []string{},
		},
		{
			name:        "duplicate items deduped",
			items:       []string{"a", "a", "b"},
			defaultItem: "",
			want:        []string{"a", "b"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := orderWithDefault(tc.items, tc.defaultItem)
			if len(got) == 0 && len(tc.want) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("orderWithDefault(%v, %q) = %v, want %v", tc.items, tc.defaultItem, got, tc.want)
			}
		})
	}
}

func TestModelBackspaceRemovesWholeRune(t *testing.T) {
	m := newModel([]string{"alpha"})
	m.query = "あい"

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	got := asModel(t, updated).query
	want := "あ"
	if got != want {
		t.Fatalf("query = %q, want %q", got, want)
	}
}

func TestModelEnterSelectsHighlightedMatch(t *testing.T) {
	m := newModel([]string{"alpha", "beta"})

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = asModel(t, updated)
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := asModel(t, updated).selected
	want := "beta"
	if got != want {
		t.Fatalf("selected = %q, want %q", got, want)
	}
	if cmd == nil {
		t.Fatalf("cmd = nil, want quit command")
	}
}

func TestModelEnterCreatesTypedNameWhenNoMatch(t *testing.T) {
	m := newModel([]string{"dev.shp"})
	m.query = "new-session"
	m.refilter()

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := asModel(t, updated).selected
	want := "new-session"
	if got != want {
		t.Fatalf("selected = %q, want %q", got, want)
	}
	if cmd == nil {
		t.Fatalf("cmd = nil, want quit command")
	}
}

func TestModelEnterCreatesDefaultSuffixWhenNoMatch(t *testing.T) {
	m := newModel([]string{"dev.shp"})
	m.defaultItem = "dev.shp"
	m.query = "-another"
	m.refilter()

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := asModel(t, updated).selected
	want := "dev.shp-another"
	if got != want {
		t.Fatalf("selected = %q, want %q", got, want)
	}
	if cmd == nil {
		t.Fatalf("cmd = nil, want quit command")
	}
}

func TestModelEnterRejectsDefaultSuffixWithWhitespace(t *testing.T) {
	m := newModel([]string{"dev.shp"})
	m.defaultItem = "dev.shp"
	m.query = "-another name"
	m.refilter()

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := asModel(t, updated).selected
	if got != "" {
		t.Fatalf("selected = %q, want empty", got)
	}
	if cmd != nil {
		t.Fatalf("cmd = %v, want nil", cmd)
	}
}

func TestModelViewShowsNewDefaultSuffixCandidate(t *testing.T) {
	m := newModel([]string{"dev.shp"})
	m.defaultItem = "dev.shp"
	m.query = "-another"
	m.refilter()

	got := m.View()
	if !strings.Contains(got, "dev.shp-another") {
		t.Fatalf("View() = %q, want new session candidate", got)
	}
}

func TestModelEnterRejectsNewNameWithWhitespace(t *testing.T) {
	m := newModel([]string{"dev.shp"})
	m.query = "new session"
	m.refilter()

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := asModel(t, updated).selected
	if got != "" {
		t.Fatalf("selected = %q, want empty", got)
	}
	if cmd != nil {
		t.Fatalf("cmd = %v, want nil", cmd)
	}
}

func TestModelEnterRejectsNewNameWithFullWidthSpace(t *testing.T) {
	m := newModel([]string{"dev.shp"})
	m.query = "新規　セッション" // contains a full-width space (U+3000)
	m.refilter()

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := asModel(t, updated).selected
	if got != "" {
		t.Fatalf("selected = %q, want empty", got)
	}
	if cmd != nil {
		t.Fatalf("cmd = %v, want nil", cmd)
	}
}

func TestModelEnterRejectsDashOnlyNewName(t *testing.T) {
	m := newModel([]string{"dev.shp"})
	m.defaultItem = "dev.shp"
	m.query = "--"
	m.refilter()

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := asModel(t, updated).selected
	if got != "" {
		t.Fatalf("selected = %q, want empty", got)
	}
	if cmd != nil {
		t.Fatalf("cmd = %v, want nil", cmd)
	}
}

func TestModelEnterRejectsSuffixWhenDefaultLeadsWithDash(t *testing.T) {
	m := newModel([]string{"-repo"})
	m.defaultItem = "-repo" // e.g. derived from a "~/-repo" cwd
	m.query = "-another"
	m.refilter()

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := asModel(t, updated).selected
	if got != "" {
		t.Fatalf("selected = %q, want empty (dash-leading name must not be created)", got)
	}
	if cmd != nil {
		t.Fatalf("cmd = %v, want nil", cmd)
	}
}

func TestModelEnterNormalizesNewNameToSafeCharset(t *testing.T) {
	m := newModel([]string{"dev.shp"})
	m.query = "foo/bar"
	m.refilter()

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := asModel(t, updated).selected
	want := "foo_bar"
	if got != want {
		t.Fatalf("selected = %q, want %q", got, want)
	}
	if cmd == nil {
		t.Fatalf("cmd = nil, want quit command")
	}
}

func TestModelViewLabelsSanitizedCollisionAsExisting(t *testing.T) {
	m := newModel([]string{"foo_bar", "dev.shp"})
	m.sessions = []string{"foo_bar", "dev.shp"} // real `shpool list` sessions
	m.query = "foo/bar"                         // sanitizes to the existing "foo_bar"
	m.refilter()

	got := m.View()
	if !strings.Contains(got, "foo_bar") {
		t.Fatalf("View() = %q, want candidate foo_bar", got)
	}
	if strings.Contains(got, "(new)") {
		t.Fatalf("View() = %q, must not label an existing session as (new)", got)
	}
	if !strings.Contains(got, "(existing)") {
		t.Fatalf("View() = %q, want (existing) label", got)
	}
}

func TestModelViewLabelsSyntheticDefaultCollisionAsNew(t *testing.T) {
	// defaultItem is the cwd-derived name but is NOT a running session, so a
	// query that sanitizes to it must be labeled (new): Enter creates it.
	m := newModel([]string{"foo_bar", "dev.shp"})
	m.defaultItem = "foo_bar"
	m.sessions = []string{"dev.shp"} // foo_bar is the injected default, not real
	m.query = "foo/bar"
	m.refilter()

	got := m.View()
	if strings.Contains(got, "(existing)") {
		t.Fatalf("View() = %q, must not label a non-running default as (existing)", got)
	}
	if !strings.Contains(got, "(new)") {
		t.Fatalf("View() = %q, want (new) label", got)
	}
}

func TestModelEnterRejectsLeadingDashWithoutDefault(t *testing.T) {
	m := newModel([]string{"dev.shp"})
	m.query = "-scratch"
	m.refilter()

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := asModel(t, updated).selected
	if got != "" {
		t.Fatalf("selected = %q, want empty", got)
	}
	if cmd != nil {
		t.Fatalf("cmd = %v, want nil", cmd)
	}
}

func TestModelViewShowsInvalidNewSessionName(t *testing.T) {
	m := newModel([]string{"dev.shp"})
	m.query = "new session"
	m.refilter()

	got := m.View()
	if !strings.Contains(got, "(invalid session name)") {
		t.Fatalf("View() = %q, want invalid session message", got)
	}
}

func TestModelScrollsCursorIntoVisibleRange(t *testing.T) {
	m := newModel([]string{"a", "b", "c", "d", "e"})

	updated, _ := m.Update(tea.WindowSizeMsg{Height: 7})
	m = asModel(t, updated)
	for i := 0; i < 4; i++ {
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = asModel(t, updated)
	}

	if m.cursor != 4 {
		t.Fatalf("cursor = %d, want 4", m.cursor)
	}
	if m.offset != 2 {
		t.Fatalf("offset = %d, want 2", m.offset)
	}
	start, end := m.visibleRange()
	if start != 2 || end != 5 {
		t.Fatalf("visibleRange = (%d, %d), want (2, 5)", start, end)
	}
}

func TestFilter(t *testing.T) {
	items := []string{
		"work__company-a__api",
		"work__company-b__api",
		"sandbox__api-test",
		"misc__notes",
	}
	cases := []struct {
		name  string
		query string
		want  []string
	}{
		{
			name:  "empty query returns all",
			query: "",
			want: []string{
				"work__company-a__api",
				"work__company-b__api",
				"sandbox__api-test",
				"misc__notes",
			},
		},
		{
			name:  "single token substring",
			query: "api",
			want: []string{
				"work__company-a__api",
				"work__company-b__api",
				"sandbox__api-test",
			},
		},
		{
			name:  "AND of two tokens",
			query: "company api",
			want: []string{
				"work__company-a__api",
				"work__company-b__api",
			},
		},
		{
			name:  "case insensitive",
			query: "API",
			want: []string{
				"work__company-a__api",
				"work__company-b__api",
				"sandbox__api-test",
			},
		},
		{
			name:  "no match",
			query: "zzz",
			want:  []string{},
		},
		{
			name:  "extra whitespace ignored",
			query: "  company    api  ",
			want: []string{
				"work__company-a__api",
				"work__company-b__api",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Filter(items, tc.query)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Filter(%q) = %v, want %v", tc.query, got, tc.want)
			}
		})
	}
}
