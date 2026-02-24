package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"bookmark/internal/bookmark"
	"bookmark/internal/config"
	"bookmark/internal/domain"
	pkg "bookmark/internal/package"
	"bookmark/internal/ui"
)

// Metadata loaded from package.toml at build time
var (
	version = pkg.Version()
	name    = pkg.Name()
	short   = pkg.Short()
)

type rootOptions struct {
	configPath  string
	showVersion bool
	interactive bool
	tmux        bool
	tmuxName    string
	description string
	yes         bool
}

var rootCmd = newRootCmd()

// Execute is the CLI entrypoint.
func Execute() error {
	return rootCmd.Execute()
}

func newRootCmd() *cobra.Command {
	opts := &rootOptions{}
	cmd := &cobra.Command{
		Use:   name + " [alias]",
		Short: short,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.showVersion {
				ver := resolvedVersion()
				cmd.Printf("%s\n", ver)
				return nil
			}

			// Load config
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			manager := config.NewManager(cwd)
			var cfg domain.Config
			if opts.configPath != "" {
				cfg, err = manager.LoadWithOverride(opts.configPath)
			} else {
				cfg, err = manager.Load()
			}
			if err != nil {
				cfg = domain.DefaultConfig()
			}

			// Interactive mode
			if opts.interactive || (len(args) == 0 && cfg.InteractiveDefault && !opts.tmux && opts.description == "") {
				return runInteractive(cmd, opts, cfg)
			}

			// Add bookmark mode
			return runAddBookmark(cmd, args, opts, cfg, cwd)
		},
	}

	cmd.Flags().StringVarP(&opts.configPath, "config", "c", "", "config file path")
	cmd.Flags().BoolVarP(&opts.showVersion, "version", "v", false, "print version information")
	cmd.Flags().BoolVarP(&opts.interactive, "interactive", "i", false, "interactive bookmark browser")
	cmd.Flags().BoolVarP(&opts.tmux, "tmux", "t", false, "set tmux window name (same as alias)")
	cmd.Flags().StringVar(&opts.tmuxName, "tmux-name", "", "custom tmux window name")
	cmd.Flags().StringVarP(&opts.description, "description", "d", "", "bookmark description")
	cmd.Flags().BoolVarP(&opts.yes, "yes", "y", false, "skip confirmation prompts")

	cmd.AddCommand(newConfigCmd())
	cmd.AddCommand(newCompletionCmd())
	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newDeleteCmd())

	return cmd
}

func resolvedVersion() string {
	ver := version
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ver
	}
	if ver == "dev" && strings.TrimSpace(info.Main.Version) != "" && info.Main.Version != "(devel)" {
		ver = info.Main.Version
	}
	return ver
}

func runAddBookmark(cmd *cobra.Command, args []string, opts *rootOptions, cfg domain.Config, cwd string) error {
	bmManager := bookmark.NewManager(cfg.BookmarkFile, cfg.Shell, cfg.NavigationTool)

	// Generate or use provided alias
	alias := generateAlias(args, cwd, cfg)

	// Check if bookmark exists and handle confirmation
	exists, err := bmManager.Exists(alias)
	if err != nil {
		return err
	}

	if exists && !opts.yes && !confirmOverwrite(cmd, bmManager, alias, cfg) {
		cmd.Println("Cancelled")
		return nil
	}

	// Create and save bookmark
	bm := buildBookmark(alias, cwd, opts)
	if err := bmManager.Add(bm); err != nil {
		return err
	}

	printSuccess(cmd, alias, cwd, exists)
	return nil
}

func generateAlias(args []string, cwd string, cfg domain.Config) string {
	if len(args) > 0 {
		return args[0]
	}
	return bookmark.GenerateAlias(cwd, cfg.AutoAliasSeparator, cfg.AutoAliasLowercase)
}

func confirmOverwrite(cmd *cobra.Command, bmManager *bookmark.Manager, alias string, cfg domain.Config) bool {
	existing, _ := bmManager.Get(alias)
	theme := ui.ThemeFromConfig(cfg)
	confirmModel := ui.NewConfirmationModel(
		"Overwrite Bookmark",
		fmt.Sprintf("Bookmark '%s → %s' already exists. Overwrite?", alias, existing.Path),
		theme,
	)

	p := tea.NewProgram(confirmModel, tea.WithoutSignalHandler())
	result, err := p.Run()
	if err != nil {
		return false
	}

	if confirmResult, ok := result.(ui.ConfirmationModel); ok {
		return confirmResult.ChoiceValue()
	}
	return false
}

func buildBookmark(alias, cwd string, opts *rootOptions) domain.Bookmark {
	bm := domain.Bookmark{
		Alias:       alias,
		Path:        cwd,
		Description: opts.description,
	}

	// Handle tmux settings
	if opts.tmux {
		bm.TmuxWindowName = alias
	}
	if opts.tmuxName != "" {
		bm.TmuxWindowName = opts.tmuxName
	}

	return bm
}

func printSuccess(cmd *cobra.Command, alias, cwd string, isUpdate bool) {
	action := "created"
	if isUpdate {
		action = "updated"
	}
	cmd.Printf("✓ Bookmark %s: %s → %s\n", action, alias, cwd)
}

func runInteractive(cmd *cobra.Command, opts *rootOptions, cfg domain.Config) error {
	bmManager := bookmark.NewManager(cfg.BookmarkFile, cfg.Shell, cfg.NavigationTool)
	bookmarks, err := bmManager.Load()
	if err != nil {
		return err
	}

	if len(bookmarks) == 0 {
		cmd.Println("No bookmarks found. Add one with: bookmark [alias]")
		return nil
	}

	return runBookmarkListing(bookmarks, cfg, bmManager)
}

func runBookmarkListing(bookmarks []domain.Bookmark, cfg domain.Config, bmManager *bookmark.Manager) error {
	items := make([]list.Item, 0, len(bookmarks))
	for _, bm := range bookmarks {
		items = append(items, bookmarkItem{Bookmark: bm})
	}

	theme := ui.ThemeFromConfig(cfg)
	delegate := ui.NewListDelegate(theme, ui.ListDelegateOptions{
		Spacing: cfg.ListSpacing,
	})

	listModel := ui.NewListModel(items, delegate, 80, 20, theme)
	listModel.Title = "Bookmarks"
	listModel.SetShowStatusBar(true)
	listModel.SetFilteringEnabled(true)

	model := bookmarkListModel{
		list:       listModel,
		theme:      theme,
		responsive: ui.NewResponsiveManager(80),
		manager:    bmManager,
	}

	model.list.AdditionalShortHelpKeys = model.getShortHelpKeys
	model.list.AdditionalFullHelpKeys = model.allHelpKeys

	p := tea.NewProgram(model, tea.WithoutSignalHandler())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run interactive list: %w", err)
	}

	return nil
}

type bookmarkItem struct {
	Bookmark domain.Bookmark
}

func (b bookmarkItem) Title() string {
	return b.Bookmark.Alias
}

func (b bookmarkItem) Description() string {
	desc := b.Bookmark.Path
	if b.Bookmark.Description != "" {
		desc = b.Bookmark.Description + " • " + desc
	}
	return desc
}

func (b bookmarkItem) FilterValue() string {
	return b.Bookmark.Alias + " " + b.Bookmark.Path + " " + b.Bookmark.Description
}

type bookmarkListModel struct {
	list          list.Model
	theme         ui.Theme
	responsive    *ui.ResponsiveManager
	manager       *bookmark.Manager
	message       string
	confirmMode   bool
	confirmModel  *ui.ConfirmationModel
	pendingAction string
	pendingItem   bookmarkItem
}

func (m bookmarkListModel) allHelpKeys() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "copy cd command"),
		),
		key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete bookmark"),
		),
	}
}

func (m bookmarkListModel) getShortHelpKeys() []key.Binding {
	allKeys := m.allHelpKeys()

	var splitAt int
	switch m.responsive.Breakpoint() {
	case ui.BreakpointXL:
		splitAt = 2
	case ui.BreakpointLG:
		splitAt = 1
	default:
		splitAt = 1
	}

	return allKeys[:splitAt]
}

func (m bookmarkListModel) getFullHelpKeys() []key.Binding {
	allKeys := m.allHelpKeys()

	var splitAt int
	switch m.responsive.Breakpoint() {
	case ui.BreakpointXS:
		splitAt = 1
	case ui.BreakpointSM:
		splitAt = 1
	case ui.BreakpointMD:
		splitAt = 2
	default:
		return []key.Binding{}
	}

	return allKeys[splitAt:]
}

func (m bookmarkListModel) Init() tea.Cmd {
	return nil
}

func (m bookmarkListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.confirmMode && m.confirmModel != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			updated, cmd := m.confirmModel.Update(msg)
			if updatedConfirm, ok := updated.(ui.ConfirmationModel); ok {
				m.confirmModel = &updatedConfirm
				if cmd != nil {
					if _, isQuit := cmd().(tea.QuitMsg); isQuit {
						confirmed := m.confirmModel.ChoiceValue()
						m.confirmMode = false
						if confirmed {
							return m.executeAction()
						} else {
							m.message = fmt.Sprintf("%s cancelled", m.pendingAction)
							m.pendingAction = ""
							return m, nil
						}
					}
				}
			}
			return m, cmd
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.responsive.SetWidth(msg.Width)
		width, height := m.responsive.GetListDimensions(msg.Width, msg.Height)
		m.list.SetSize(width, height)
		m.list.AdditionalShortHelpKeys = m.getShortHelpKeys
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if item, ok := m.list.SelectedItem().(bookmarkItem); ok {
				cdCmd := fmt.Sprintf("cd %s", item.Bookmark.Path)
				fmt.Println(cdCmd)
				return m, tea.Quit
			}
		case "d":
			if item, ok := m.list.SelectedItem().(bookmarkItem); ok {
				m.pendingAction = "Delete"
				m.pendingItem = item
				confirmModel := ui.NewConfirmationModel(
					"Delete Bookmark",
					fmt.Sprintf("Delete bookmark '%s'?", item.Bookmark.Alias),
					m.theme,
				)
				m.confirmModel = &confirmModel
				m.confirmMode = true
				return m, confirmModel.Init()
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m bookmarkListModel) View() string {
	if m.confirmMode && m.confirmModel != nil {
		return m.confirmModel.View()
	}

	listView := m.list.View()

	if m.message != "" {
		listView = listView + "\n\n" + m.message
	}

	return m.responsive.AdaptiveFrameStyle(m.theme).Render(listView)
}

func (m bookmarkListModel) executeAction() (tea.Model, tea.Cmd) {
	switch m.pendingAction {
	case "Delete":
		if err := m.manager.Delete(m.pendingItem.Bookmark.Alias); err != nil {
			m.message = fmt.Sprintf("✗ Failed to delete: %s", err)
		} else {
			m.message = fmt.Sprintf("✓ Deleted: %s", m.pendingItem.Bookmark.Alias)

			// Remove from list
			items := m.list.Items()
			filtered := make([]list.Item, 0, len(items))
			for _, item := range items {
				if bm, ok := item.(bookmarkItem); ok {
					if bm.Bookmark.Alias != m.pendingItem.Bookmark.Alias {
						filtered = append(filtered, item)
					}
				}
			}
			m.list.SetItems(filtered)
		}
	}
	m.pendingAction = ""
	return m, nil
}
