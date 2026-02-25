package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"bookmark/internal/bookmark"
	"bookmark/internal/config"
	"bookmark/internal/domain"
)

func newListCmd() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all bookmarks",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := cmd.Flags().GetString("cwd")
			if err != nil {
				cwd = "."
			}

			manager := config.NewManager(cwd)
			var cfg domain.Config
			if configPath != "" {
				cfg, err = manager.LoadWithOverride(configPath)
			} else {
				cfg, err = manager.Load()
			}
			if err != nil {
				cfg = domain.DefaultConfig()
			}

			bmManager := bookmark.NewManager(cfg.BookmarkFile(), cfg.Shell, cfg.NavigationTool, cfg.Editor, cfg.FunctionAlias, cfg.InteractiveAlias)
			bookmarks, err := bmManager.Load()
			if err != nil {
				return err
			}

			if len(bookmarks) == 0 {
				cmd.Println("No bookmarks found")
				return nil
			}

			// Find max alias length for alignment
			maxAlias := 0
			maxPath := 0
			for _, bm := range bookmarks {
				if len(bm.Alias) > maxAlias {
					maxAlias = len(bm.Alias)
				}
				if len(bm.Path) > maxPath {
					maxPath = len(bm.Path)
				}
			}

			for _, bm := range bookmarks {
				line := fmt.Sprintf("%-*s  %-*s", maxAlias, bm.Alias, maxPath, bm.Path)
				if bm.Description != "" {
					line += "  " + bm.Description
				}
				cmd.Println(line)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "", "config file path")

	return cmd
}
