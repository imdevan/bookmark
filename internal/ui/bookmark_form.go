package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// BookmarkFormModel is a form for creating a new bookmark.
type BookmarkFormModel struct {
	form  *huh.Form
	theme Theme
	alias string
	path  string
	desc  string
}

// NewBookmarkFormModel creates a new bookmark form.
func NewBookmarkFormModel(theme Theme) BookmarkFormModel {
	m := BookmarkFormModel{
		theme: theme,
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Alias").
				Description("Short name for the bookmark").
				Value(&m.alias).
				Validate(huh.ValidateNotEmpty()),

			huh.NewInput().
				Title("Path").
				Description("Directory path to bookmark").
				Value(&m.path).
				Validate(huh.ValidateNotEmpty()),

			huh.NewInput().
				Title("Description (optional)").
				Value(&m.desc),
		),
	).WithTheme(bookmarkFormHuhTheme(theme))

	return m
}

// Init initializes the form.
func (m BookmarkFormModel) Init() tea.Cmd {
	return m.form.Init()
}

// Update handles messages for the form.
func (m BookmarkFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	form, formCmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmd = formCmd
	}

	if m.form.State == huh.StateCompleted {
		return m, tea.Quit
	}

	return m, cmd
}

// View renders the form.
func (m BookmarkFormModel) View() string {
	if m.form.State == huh.StateCompleted {
		return ""
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.theme.Headings).
		Render("Create New Bookmark")

	help := lipgloss.NewStyle().
		Foreground(m.theme.Muted).
		Render("Press Esc to cancel")

	content := strings.Join([]string{
		title,
		"",
		m.form.View(),
		"",
		help,
	}, "\n")

	return lipgloss.NewStyle().
		Margin(1, 1).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Border).
		Render(content)
}

// Values returns the form values.
func (m BookmarkFormModel) Values() (alias, path, desc string) {
	return strings.TrimSpace(m.alias), strings.TrimSpace(m.path), strings.TrimSpace(m.desc)
}

// IsCompleted returns true if the form was completed successfully.
func (m BookmarkFormModel) IsCompleted() bool {
	return m.form.State == huh.StateCompleted
}

func bookmarkFormHuhTheme(theme Theme) *huh.Theme {
	huhTheme := huh.ThemeBase()
	huhTheme.Focused.Base = lipgloss.NewStyle()
	huhTheme.Blurred.Base = lipgloss.NewStyle()
	huhTheme.Focused.Title = lipgloss.NewStyle().
		Foreground(theme.Secondary).
		Bold(true)
	huhTheme.Focused.Description = lipgloss.NewStyle().
		Foreground(theme.Muted)
	huhTheme.Focused.TextInput.Cursor = lipgloss.NewStyle().
		Foreground(theme.Secondary)
	huhTheme.Focused.TextInput.Placeholder = lipgloss.NewStyle().
		Foreground(theme.Muted)
	huhTheme.Focused.TextInput.Prompt = lipgloss.NewStyle().
		Foreground(theme.Secondary)
	return huhTheme
}
