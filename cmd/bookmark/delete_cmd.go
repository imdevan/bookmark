package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"bookmark/internal/bookmark"
	"bookmark/internal/config"
	"bookmark/internal/domain"
	"bookmark/internal/ui"
)

func newDeleteCmd() *cobra.Command {
	var configPath string
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <alias>",
		Short: "Delete a bookmark",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			alias := args[0]

			cwd, _ := os.Getwd()
			manager := config.NewManager(cwd)
			var cfg domain.Config
			var err error
			if configPath != "" {
				cfg, err = manager.LoadWithOverride(configPath)
			} else {
				cfg, err = manager.Load()
			}
			if err != nil {
				cfg = domain.DefaultConfig()
			}

			bmManager := bookmark.NewManager(cfg.BookmarkFile, cfg.Shell, cfg.NavigationTool)

			// Check if bookmark exists
			bm, err := bmManager.Get(alias)
			if err == bookmark.ErrBookmarkNotFound {
				return fmt.Errorf("bookmark '%s' not found", alias)
			}
			if err != nil {
				return err
			}

			// Confirm deletion unless --force
			if !force {
				theme := ui.ThemeFromConfig(cfg)
				confirmModel := ui.NewConfirmationModel(
					"Delete Bookmark",
					fmt.Sprintf("Delete bookmark '%s → %s'?", bm.Alias, bm.Path),
					theme,
				)

				p := tea.NewProgram(confirmModel, tea.WithoutSignalHandler())
				result, err := p.Run()
				if err != nil {
					return err
				}

				if confirmResult, ok := result.(ui.ConfirmationModel); ok {
					if !confirmResult.ChoiceValue() {
						cmd.Println("Cancelled")
						return nil
					}
				}
			}

			if err := bmManager.Delete(alias); err != nil {
				return err
			}

			cmd.Printf("✓ Bookmark deleted: %s\n", alias)
			return nil
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "", "config file path")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "skip confirmation")

	return cmd
}
