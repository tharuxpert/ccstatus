package cmd

import (
	"fmt"
	"strings"

	"ccstatus/internal/config"
	"ccstatus/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure statusline display options",
	Long: `Configure what information is displayed in the statusline.

Toggle options on/off by pressing Enter on the selected item.
Use arrow keys to navigate, then select Save or Cancel.`,
	RunE: runConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

// Styles for the config UI
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("6")). // Cyan
			MarginBottom(1)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedStyle = lipgloss.NewStyle().
			PaddingLeft(0).
			Foreground(lipgloss.Color("6")). // Cyan
			Bold(true)

	toggleOnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")). // Green
			Bold(true)

	toggleOffStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")) // Dim gray

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")). // Light gray - more visible
			Italic(true)

	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")). // Medium gray
			MarginTop(1).
			MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("246")). // Visible gray
			MarginTop(1)

	actionLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250")) // Light gray for unselected actions

	saveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")). // Green
			Bold(true)

	cancelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")). // Red
			Bold(true)
)

// configOption represents a single toggle option
type configOption struct {
	key         string
	label       string
	description string
	enabled     bool
}

// configModel is the bubbletea model for the config screen
type configModel struct {
	options       []configOption
	cursor        int
	originalCfg   *config.CCStatusConfig
	saved         bool
	cancelled     bool
	hasChanges    bool
}

func initialModel() (configModel, error) {
	cfg, err := config.LoadCCStatusConfig()
	if err != nil {
		// Use defaults on error
		cfg = config.DefaultCCStatusConfig()
	}

	options := []configOption{
		{
			key:         "session",
			label:       "Session Usage",
			description: "Show current session usage percentage",
			enabled:     cfg.ShowSessionUsage,
		},
		{
			key:         "weekly",
			label:       "Weekly Usage",
			description: "Show weekly usage percentage",
			enabled:     cfg.ShowWeeklyUsage,
		},
		{
			key:         "reset",
			label:       "Reset Times",
			description: "Show when usage limits reset",
			enabled:     cfg.ShowResetTimes,
		},
		{
			key:         "git",
			label:       "Git Branch",
			description: "Show current git branch name",
			enabled:     cfg.ShowGitBranch,
		},
	}

	return configModel{
		options:     options,
		cursor:      0,
		originalCfg: cfg,
		saved:       false,
		cancelled:   false,
		hasChanges:  false,
	}, nil
}

func (m configModel) Init() tea.Cmd {
	return nil
}

func (m configModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.cancelled = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			// Total items = options + 2 (Save, Cancel)
			maxCursor := len(m.options) + 1
			if m.cursor < maxCursor {
				m.cursor++
			}

		case "enter", " ":
			if m.cursor < len(m.options) {
				// Toggle the option
				m.options[m.cursor].enabled = !m.options[m.cursor].enabled
				m.hasChanges = m.checkForChanges()
			} else if m.cursor == len(m.options) {
				// Save
				m.saved = true
				return m, tea.Quit
			} else {
				// Cancel
				m.cancelled = true
				return m, tea.Quit
			}

		case "s":
			// Quick save
			m.saved = true
			return m, tea.Quit

		case "esc":
			m.cancelled = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m configModel) checkForChanges() bool {
	return m.options[0].enabled != m.originalCfg.ShowSessionUsage ||
		m.options[1].enabled != m.originalCfg.ShowWeeklyUsage ||
		m.options[2].enabled != m.originalCfg.ShowResetTimes ||
		m.options[3].enabled != m.originalCfg.ShowGitBranch
}

func (m configModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString("\n")
	b.WriteString(titleStyle.Render("  ◆ Statusline Configuration"))
	b.WriteString("\n")
	b.WriteString(dividerStyle.Render("  " + strings.Repeat("─", 44)))
	b.WriteString("\n\n")

	// Options
	for i, opt := range m.options {
		cursor := "  "
		if m.cursor == i {
			cursor = selectedStyle.Render("→ ")
		}

		// Toggle indicator
		var toggle string
		if opt.enabled {
			toggle = toggleOnStyle.Render("● ON ")
		} else {
			toggle = toggleOffStyle.Render("○ OFF")
		}

		// Option label
		label := opt.label
		if m.cursor == i {
			label = selectedStyle.Render(label)
		}

		// Build the line
		line := fmt.Sprintf("%s%-20s %s", cursor, label, toggle)
		b.WriteString(line)
		b.WriteString("\n")

		// Description (shown for selected item)
		if m.cursor == i {
			b.WriteString(descStyle.Render(fmt.Sprintf("      %s", opt.description)))
			b.WriteString("\n")
		}
	}

	// Divider before actions
	b.WriteString("\n")
	b.WriteString(dividerStyle.Render("  " + strings.Repeat("─", 44)))
	b.WriteString("\n\n")

	// Save option
	saveCursor := "  "
	saveLabel := "Save changes"
	if m.cursor == len(m.options) {
		saveCursor = selectedStyle.Render("→ ")
		saveLabel = saveStyle.Render("Save changes")
	} else {
		saveLabel = actionLabelStyle.Render("Save changes")
	}
	if m.hasChanges {
		b.WriteString(fmt.Sprintf("%s%s %s\n", saveCursor, saveLabel, helpStyle.Render("(has changes)")))
	} else {
		b.WriteString(fmt.Sprintf("%s%s\n", saveCursor, saveLabel))
	}

	// Cancel option
	cancelCursor := "  "
	cancelLabel := "Cancel"
	if m.cursor == len(m.options)+1 {
		cancelCursor = selectedStyle.Render("→ ")
		cancelLabel = cancelStyle.Render("Cancel")
	} else {
		cancelLabel = actionLabelStyle.Render("Cancel")
	}
	b.WriteString(fmt.Sprintf("%s%s\n", cancelCursor, cancelLabel))

	// Help text
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("  ↑/↓ Navigate • Enter Toggle/Select • s Save • Esc Cancel"))
	b.WriteString("\n")

	return b.String()
}

func (m configModel) getConfig() *config.CCStatusConfig {
	return &config.CCStatusConfig{
		ShowSessionUsage: m.options[0].enabled,
		ShowWeeklyUsage:  m.options[1].enabled,
		ShowResetTimes:   m.options[2].enabled,
		ShowGitBranch:    m.options[3].enabled,
	}
}

func runConfig(cmd *cobra.Command, args []string) error {
	model, err := initialModel()
	if err != nil {
		ui.ErrorMessage("Failed to load config", err.Error())
		return nil
	}

	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		ui.ErrorMessage("Error running config UI", err.Error())
		return nil
	}

	m := finalModel.(configModel)

	fmt.Println()

	if m.saved {
		cfg := m.getConfig()
		if err := config.SaveCCStatusConfig(cfg); err != nil {
			ui.ErrorMessage("Failed to save config", err.Error())
			return nil
		}
		ui.SuccessMessage("Configuration saved!", "Your statusline preferences have been updated.")
	} else if m.cancelled {
		ui.WarningMessage("Cancelled", "No changes were made.")
	}

	fmt.Println()
	return nil
}
