package cmd

import (
	"strings"
	"testing"

	"ccstatus/internal/config"
)

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
