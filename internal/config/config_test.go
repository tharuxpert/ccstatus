package config

import "testing"

func TestSetStatuslineCommandCreatesCommandType(t *testing.T) {
	settings := Settings{}

	SetStatuslineCommand(settings, "ccstatus")

	statusline := GetStatuslineObject(settings)
	if statusline == nil {
		t.Fatal("expected statusLine object")
	}

	if got := statusline["type"]; got != "command" {
		t.Fatalf("expected statusLine type command, got %v", got)
	}
	if got := statusline["command"]; got != "ccstatus" {
		t.Fatalf("expected statusLine command ccstatus, got %v", got)
	}
}

func TestSetStatuslineCommandPreservesOtherFields(t *testing.T) {
	settings := Settings{
		StatusLineKey: map[string]any{
			"padding": float64(0),
		},
	}

	SetStatuslineCommand(settings, "ccstatus")

	statusline := GetStatuslineObject(settings)
	if got := statusline["padding"]; got != float64(0) {
		t.Fatalf("expected padding to be preserved, got %v", got)
	}
	if got := statusline["type"]; got != "command" {
		t.Fatalf("expected statusLine type command, got %v", got)
	}
	if got := statusline["command"]; got != "ccstatus" {
		t.Fatalf("expected statusLine command ccstatus, got %v", got)
	}
}

func TestRemoveStatuslineRemovesCommandType(t *testing.T) {
	settings := Settings{
		StatusLineKey: map[string]any{
			"type":    "command",
			"command": "ccstatus",
		},
	}

	RemoveStatusline(settings)

	if _, ok := settings[StatusLineKey]; ok {
		t.Fatal("expected empty statusLine object to be removed")
	}
}

func TestRemoveStatuslinePreservesOtherFields(t *testing.T) {
	settings := Settings{
		StatusLineKey: map[string]any{
			"type":    "command",
			"command": "ccstatus",
			"padding": float64(0),
		},
	}

	RemoveStatusline(settings)

	statusline := GetStatuslineObject(settings)
	if statusline == nil {
		t.Fatal("expected statusLine object to remain")
	}
	if _, ok := statusline["type"]; ok {
		t.Fatal("expected statusLine type to be removed")
	}
	if _, ok := statusline["command"]; ok {
		t.Fatal("expected statusLine command to be removed")
	}
	if got := statusline["padding"]; got != float64(0) {
		t.Fatalf("expected padding to be preserved, got %v", got)
	}
}
