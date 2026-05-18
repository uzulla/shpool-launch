package tui

import (
	"reflect"
	"testing"
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
