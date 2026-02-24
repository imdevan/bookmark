package domain

import "time"

// Bookmark represents a saved directory location with metadata.
type Bookmark struct {
	Alias           string    `toml:"alias"`
	Path            string    `toml:"path"`
	Description     string    `toml:"description,omitempty"`
	CreatedAt       time.Time `toml:"created_at"`
	UpdatedAt       time.Time `toml:"updated_at"`
	TmuxWindowName  string    `toml:"tmux_window_name,omitempty"`
	Execute         string    `toml:"execute,omitempty"`
	PostJumpScript  string    `toml:"post_jump_script,omitempty"`
	File            string    `toml:"file,omitempty"`
}

// BookmarkStore represents the TOML file structure.
type BookmarkStore struct {
	Bookmarks []Bookmark `toml:"bookmarks"`
}
