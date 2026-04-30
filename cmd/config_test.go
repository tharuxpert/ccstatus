package cmd

import (
	"regexp"
	"strings"
	"testing"
	"unicode/utf8"

	"ccstatus/internal/config"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(input string) string {
	return ansiPattern.ReplaceAllString(input, "")
}

func runeIndex(input, substr string) int {
	byteIndex := strings.Index(input, substr)
	if byteIndex == -1 {
		return -1
	}
	return utf8.RuneCountInString(input[:byteIndex])
}

func TestConfigViewKeepsDescriptionOutsideOptionRows(t *testing.T) {
	model := configModel{
		options: []configOption{
			{
				key:         "session",
				label:       "Session Usage",
				description: "Show current session usage percentage",
				enabled:     true,
			},
			{
				key:         "weekly",
				label:       "Weekly Usage",
				description: "Show weekly usage percentage",
				enabled:     true,
			},
		},
		cursor:      0,
		originalCfg: config.DefaultCCStatusConfig(),
	}

	view := model.View()
	sessionIndex := strings.Index(view, "Session Usage")
	weeklyIndex := strings.Index(view, "Weekly Usage")
	descriptionIndex := strings.Index(view, "Show current session usage percentage")

	if sessionIndex == -1 || weeklyIndex == -1 || descriptionIndex == -1 {
		t.Fatalf("expected option labels and description in view:\n%s", view)
	}

	if descriptionIndex < weeklyIndex {
		t.Fatal("expected description to render after option rows, not inside the item list")
	}
}

func TestConfigViewKeepsToggleColumnAligned(t *testing.T) {
	model := configModel{
		options: []configOption{
			{
				key:         "session",
				label:       "Session Usage",
				description: "Show current session usage percentage",
				enabled:     true,
			},
			{
				key:         "weekly",
				label:       "Weekly Usage",
				description: "Show weekly usage percentage",
				enabled:     true,
			},
		},
		cursor:      0,
		originalCfg: config.DefaultCCStatusConfig(),
	}

	lines := strings.Split(stripANSI(model.View()), "\n")
	var toggleColumns []int
	for _, line := range lines {
		if strings.Contains(line, "Usage") && strings.Contains(line, "ON") {
			toggleColumns = append(toggleColumns, runeIndex(line, "● ON"))
		}
	}

	if len(toggleColumns) != 2 {
		t.Fatalf("expected two option rows, got columns %v", toggleColumns)
	}
	if toggleColumns[0] != toggleColumns[1] {
		t.Fatalf("expected toggle column to stay aligned, got %v", toggleColumns)
	}
}

func TestConfigViewShowsHasChangesOnSaveLine(t *testing.T) {
	model := configModel{
		options: []configOption{
			{
				key:         "session",
				label:       "Session Usage",
				description: "Show current session usage percentage",
				enabled:     true,
			},
		},
		cursor:      1,
		originalCfg: config.DefaultCCStatusConfig(),
		hasChanges:  true,
	}

	lines := strings.Split(stripANSI(model.View()), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Save changes") {
			if !strings.Contains(line, "(has changes)") {
				t.Fatalf("expected save line to include change marker, got %q", line)
			}
			return
		}
	}

	t.Fatalf("expected save line in view:\n%s", model.View())
}
