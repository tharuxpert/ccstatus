// Package ui provides consistent styling and visual elements for the CLI.
package ui

import (
	"fmt"
	"io"
	"strings"
	"text/template"
	"time"

	"github.com/briandowns/spinner"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

// bellSkipper wraps an io.WriteCloser to skip bell characters
type bellSkipper struct {
	io.WriteCloser
}

func (b *bellSkipper) Write(data []byte) (int, error) {
	// Filter out bell character (ASCII 7)
	filtered := make([]byte, 0, len(data))
	for _, c := range data {
		if c != 7 { // Skip bell character
			filtered = append(filtered, c)
		}
	}
	if len(filtered) == 0 {
		return len(data), nil
	}
	_, err := b.WriteCloser.Write(filtered)
	return len(data), err
}

// newBellSkipper creates a readline config with bell disabled
func newBellSkipper() io.WriteCloser {
	return &bellSkipper{readline.Stdout}
}

// Colors
var (
	// Primary colors
	Primary   = color.New(color.FgCyan, color.Bold)
	Secondary = color.New(color.FgWhite)

	// Status colors
	Success = color.New(color.FgGreen)
	Warning = color.New(color.FgYellow)
	Error   = color.New(color.FgRed)
	Info    = color.New(color.FgCyan)

	// Text styles
	Bold   = color.New(color.Bold)
	Dim    = color.New(color.Faint)
	Italic = color.New(color.Italic)

	// Highlighted
	SuccessBold = color.New(color.FgGreen, color.Bold)
	ErrorBold   = color.New(color.FgRed, color.Bold)
	WarningBold = color.New(color.FgYellow, color.Bold)
	InfoBold    = color.New(color.FgCyan, color.Bold)
)

// Icons
const (
	IconCheck     = "\u2714" // ✔
	IconCross     = "\u2718" // ✘
	IconWarning   = "\u26A0" // ⚠
	IconInfo      = "\u2139" // ℹ
	IconArrow     = "\u2192" // →
	IconBullet    = "\u2022" // •
	IconStar      = "\u2605" // ★
	IconBox       = "\u25A0" // ■
	IconCircle    = "\u25CF" // ●
	IconDiamond   = "\u25C6" // ◆
	IconGitBranch = "\u2387" // ⎇
)

// Header prints a styled header
func Header(text string) {
	fmt.Println()
	Primary.Println(text)
	fmt.Println(strings.Repeat("─", len(text)+2))
}

// SubHeader prints a styled sub-header
func SubHeader(text string) {
	fmt.Println()
	Bold.Println(text)
}

// StatusOK prints a success status line
func StatusOK(label, message string) {
	Success.Printf("  %s ", IconCheck)
	Bold.Print(label)
	if message != "" {
		Dim.Printf(" %s %s", IconArrow, message)
	}
	fmt.Println()
}

// StatusError prints an error status line
func StatusError(label, message string) {
	Error.Printf("  %s ", IconCross)
	Bold.Print(label)
	if message != "" {
		fmt.Println()
		Error.Printf("      %s\n", message)
	} else {
		fmt.Println()
	}
}

// StatusWarning prints a warning status line
func StatusWarning(label, message string) {
	Warning.Printf("  %s ", IconWarning)
	Bold.Print(label)
	if message != "" {
		fmt.Println()
		Warning.Printf("      %s\n", message)
	} else {
		fmt.Println()
	}
}

// StatusInfo prints an info status line
func StatusInfo(label, message string) {
	Info.Printf("  %s ", IconInfo)
	fmt.Print(label)
	if message != "" {
		Dim.Printf(" %s %s", IconArrow, message)
	}
	fmt.Println()
}

// Bullet prints a bulleted item
func Bullet(text string) {
	fmt.Printf("  %s %s\n", IconBullet, text)
}

// Step prints a numbered step
func Step(num int, text string) {
	InfoBold.Printf("  %d. ", num)
	fmt.Println(text)
}

// SuccessMessage prints a success message block
func SuccessMessage(title, message string) {
	fmt.Println()
	SuccessBold.Printf("%s %s\n", IconCheck, title)
	if message != "" {
		Dim.Printf("   %s\n", message)
	}
}

// ErrorMessage prints an error message block
func ErrorMessage(title, message string) {
	fmt.Println()
	ErrorBold.Printf("%s %s\n", IconCross, title)
	if message != "" {
		Error.Printf("   %s\n", message)
	}
}

// WarningMessage prints a warning message block
func WarningMessage(title, message string) {
	fmt.Println()
	WarningBold.Printf("%s %s\n", IconWarning, title)
	if message != "" {
		Warning.Printf("   %s\n", message)
	}
}

// InfoBox prints an info box with a message
func InfoBox(lines ...string) {
	const boxWidth = 50
	fmt.Println()
	Info.Println("  ┌" + strings.Repeat("─", boxWidth+2) + "┐")
	for _, line := range lines {
		if len(line) > boxWidth {
			line = line[:boxWidth-3] + "..."
		}
		padding := boxWidth - len(line)
		Info.Printf("  │ %s%s │\n", line, strings.Repeat(" ", padding))
	}
	Info.Println("  └" + strings.Repeat("─", boxWidth+2) + "┘")
}

// CodeBlock prints a styled code/config block
func CodeBlock(content string) {
	const boxWidth = 50
	lines := strings.Split(content, "\n")
	Dim.Println("  ┌" + strings.Repeat("─", boxWidth+2) + "┐")
	for _, line := range lines {
		if len(line) > boxWidth {
			line = line[:boxWidth-3] + "..."
		}
		padding := boxWidth - len(line)
		Dim.Print("  │ ")
		Info.Print(line)
		fmt.Print(strings.Repeat(" ", padding))
		Dim.Println(" │")
	}
	Dim.Println("  └" + strings.Repeat("─", boxWidth+2) + "┘")
}

// Divider prints a horizontal divider
func Divider() {
	Dim.Println("  " + strings.Repeat("─", 50))
}

// NewSpinner creates a styled spinner with the given message
func NewSpinner(message string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 80*time.Millisecond)
	s.Prefix = "  "
	s.Suffix = " " + message
	s.Color("cyan")
	return s
}

// NewProgressSpinner creates a spinner that looks like progress
func NewProgressSpinner(message string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Prefix = "  "
	s.Suffix = " " + message
	s.Color("cyan")
	return s
}

// Confirm prompts the user for yes/no confirmation with interactive selector
func Confirm(prompt string) bool {
	fmt.Println()

	items := []string{"Yes", "No"}

	// Custom templates for styling
	templates := &promptui.SelectTemplates{
		Label:    fmt.Sprintf("  %s {{ . | bold }}", IconWarning),
		Active:   fmt.Sprintf("  %s {{ . | cyan | bold }}", IconArrow),
		Inactive: "    {{ . | faint }}",
		Selected: fmt.Sprintf("  %s {{ . | green }}", IconCheck),
		Help:     Dim.Sprint("  Use ↑/↓ arrows to move, Enter to select"),
	}

	sel := promptui.Select{
		Label:        prompt,
		Items:        items,
		Templates:    templates,
		HideSelected: false,
		HideHelp:     false,
		Stdout:       newBellSkipper(),
	}

	idx, _, err := sel.Run()
	if err != nil {
		return false
	}

	return idx == 0 // "Yes" is at index 0
}

// PromptChoice prompts the user to select from options with interactive selector
func PromptChoice(prompt string, options []string) int {
	fmt.Println()

	// Custom templates for styling
	funcMap := template.FuncMap{
		"cyan":    func(s string) string { return Info.Sprint(s) },
		"green":   func(s string) string { return Success.Sprint(s) },
		"yellow":  func(s string) string { return Warning.Sprint(s) },
		"red":     func(s string) string { return Error.Sprint(s) },
		"bold":    func(s string) string { return Bold.Sprint(s) },
		"faint":   func(s string) string { return Dim.Sprint(s) },
		"primary": func(s string) string { return Primary.Sprint(s) },
	}

	templates := &promptui.SelectTemplates{
		Label:    fmt.Sprintf("  %s {{ . | bold }}", IconDiamond),
		Active:   fmt.Sprintf("  %s {{ . | cyan | bold }}", IconArrow),
		Inactive: "    {{ . | faint }}",
		Selected: fmt.Sprintf("  %s {{ . | green }}", IconCheck),
		Help:     Dim.Sprint("  Use ↑/↓ arrows to move, Enter to select"),
		FuncMap:  funcMap,
	}

	sel := promptui.Select{
		Label:        prompt,
		Items:        options,
		Templates:    templates,
		HideSelected: false,
		HideHelp:     false,
		Size:         len(options),
		Stdout:       newBellSkipper(),
	}

	idx, _, err := sel.Run()
	if err != nil {
		return 1 // Default to first option on error
	}

	return idx + 1 // Return 1-indexed for compatibility
}

// PrintKeyValue prints a key-value pair with styling
func PrintKeyValue(key, value string) {
	Dim.Printf("  %s: ", key)
	fmt.Println(value)
}

// PrintPath prints a file path with styling
func PrintPath(label, path string) {
	Dim.Printf("  %s: ", label)
	Info.Println(path)
}

// Title prints a large title banner
func Title(text string) {
	fmt.Println()
	Primary.Println("  " + strings.Repeat("━", len(text)+4))
	Primary.Printf("  ┃ %s ┃\n", text)
	Primary.Println("  " + strings.Repeat("━", len(text)+4))
}

// CompactTitle prints a simpler title
func CompactTitle(text string) {
	fmt.Println()
	Primary.Printf("  %s %s\n", IconDiamond, text)
	Dim.Println("  " + strings.Repeat("─", len(text)+4))
}
