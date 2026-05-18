// Package tui provides a peco-style incremental filter UI for picking a
// shpool session name.
package tui

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ErrCancelled is returned when the user dismisses the UI (Esc / Ctrl-C).
var ErrCancelled = errors.New("selection cancelled")

// Select shows the picker for the given items and returns the chosen entry.
// Returns ErrCancelled if the user aborts.
func Select(items []string) (string, error) {
	return SelectWithDefault(items, "")
}

// SelectWithDefault shows the picker with defaultItem placed first and marked
// "(cwd)". If defaultItem is non-empty and not already in items, it is added
// to the front of the list. If it is already present, it is moved to the front.
// The cursor starts on the default item when provided.
func SelectWithDefault(items []string, defaultItem string) (string, error) {
	all := orderWithDefault(items, defaultItem)
	if len(all) == 0 {
		return "", errors.New("no items to select")
	}
	m := newModel(all)
	m.defaultItem = defaultItem
	p := tea.NewProgram(m, tea.WithAltScreen())
	final, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("tui: %w", err)
	}
	fm, ok := final.(model)
	if !ok {
		return "", errors.New("tui: unexpected model type")
	}
	if fm.cancelled {
		return "", ErrCancelled
	}
	if fm.selected == "" {
		return "", ErrCancelled
	}
	return fm.selected, nil
}

// orderWithDefault returns items with defaultItem (if non-empty) first and
// duplicates removed.
func orderWithDefault(items []string, defaultItem string) []string {
	out := make([]string, 0, len(items)+1)
	seen := make(map[string]bool, len(items)+1)
	if defaultItem != "" {
		out = append(out, defaultItem)
		seen[defaultItem] = true
	}
	for _, it := range items {
		if seen[it] {
			continue
		}
		seen[it] = true
		out = append(out, it)
	}
	return out
}

// Filter returns items matching all whitespace-separated tokens in query
// (case-insensitive substring match). An empty query returns items unchanged.
func Filter(items []string, query string) []string {
	tokens := strings.Fields(strings.ToLower(query))
	if len(tokens) == 0 {
		// Return a copy so callers can't mutate the original.
		out := make([]string, len(items))
		copy(out, items)
		return out
	}
	out := make([]string, 0, len(items))
	for _, it := range items {
		lower := strings.ToLower(it)
		ok := true
		for _, tok := range tokens {
			if !strings.Contains(lower, tok) {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, it)
		}
	}
	return out
}

type model struct {
	all         []string
	filtered    []string
	query       string
	cursor      int
	selected    string
	cancelled   bool
	defaultItem string
}

func newModel(items []string) model {
	cp := make([]string, len(items))
	copy(cp, items)
	return model{all: cp, filtered: cp}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch keyMsg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		m.cancelled = true
		return m, tea.Quit
	case tea.KeyEnter:
		if len(m.filtered) > 0 && m.cursor >= 0 && m.cursor < len(m.filtered) {
			m.selected = m.filtered[m.cursor]
		}
		return m, tea.Quit
	case tea.KeyUp, tea.KeyCtrlP:
		if m.cursor > 0 {
			m.cursor--
		}
		return m, nil
	case tea.KeyDown, tea.KeyCtrlN:
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
		}
		return m, nil
	case tea.KeyBackspace:
		if len(m.query) > 0 {
			m.query = m.query[:len(m.query)-1]
			m.refilter()
		}
		return m, nil
	case tea.KeyRunes, tea.KeySpace:
		m.query += string(keyMsg.Runes)
		if keyMsg.Type == tea.KeySpace && len(keyMsg.Runes) == 0 {
			m.query += " "
		}
		m.refilter()
		return m, nil
	}
	return m, nil
}

func (m *model) refilter() {
	m.filtered = Filter(m.all, m.query)
	if m.cursor >= len(m.filtered) {
		m.cursor = len(m.filtered) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

var (
	queryLabelStyle = lipgloss.NewStyle().Bold(true)
	cursorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)
	selectedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	dimStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

func (m model) View() string {
	var b strings.Builder
	b.WriteString(queryLabelStyle.Render("QUERY> "))
	b.WriteString(m.query)
	b.WriteString("\n\n")

	if len(m.filtered) == 0 {
		b.WriteString(dimStyle.Render("  (no matches)"))
		b.WriteString("\n")
	}
	for i, it := range m.filtered {
		if i == m.cursor {
			b.WriteString(cursorStyle.Render("> "))
			b.WriteString(selectedStyle.Render(it))
		} else {
			b.WriteString("  ")
			b.WriteString(it)
		}
		if it == m.defaultItem && m.defaultItem != "" {
			b.WriteString(" ")
			b.WriteString(dimStyle.Render("(cwd)"))
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("Enter: attach   Esc/Ctrl-C: cancel   ↑/↓ or Ctrl-P/N: move"))
	b.WriteString("\n")
	return b.String()
}
