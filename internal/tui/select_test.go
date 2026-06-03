package tui

import (
	"reflect"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

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
	got := updated.(model).query
	want := "あ"
	if got != want {
		t.Fatalf("query = %q, want %q", got, want)
	}
}

func TestModelEnterSelectsHighlightedMatch(t *testing.T) {
	m := newModel([]string{"alpha", "beta"})

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := updated.(model).selected
	want := "beta"
	if got != want {
		t.Fatalf("selected = %q, want %q", got, want)
	}
}

func TestModelEnterCreatesTypedNameWhenNoMatch(t *testing.T) {
	m := newModel([]string{"dev.shp"})
	m.query = "new-session"
	m.refilter()

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := updated.(model).selected
	want := "new-session"
	if got != want {
		t.Fatalf("selected = %q, want %q", got, want)
	}
}

func TestModelEnterCreatesDefaultSuffixWhenNoMatch(t *testing.T) {
	m := newModel([]string{"dev.shp"})
	m.defaultItem = "dev.shp"
	m.query = "-another"
	m.refilter()

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := updated.(model).selected
	want := "dev.shp-another"
	if got != want {
		t.Fatalf("selected = %q, want %q", got, want)
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

func TestModelScrollsCursorIntoVisibleRange(t *testing.T) {
	m := newModel([]string{"a", "b", "c", "d", "e"})

	updated, _ := m.Update(tea.WindowSizeMsg{Height: 7})
	m = updated.(model)
	for i := 0; i < 4; i++ {
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = updated.(model)
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
