package bookmark

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"bookmark/internal/domain"
)

func TestGenerateAlias(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		separator string
		lowercase bool
		want      string
	}{
		{
			name:      "simple path",
			path:      "/home/user/my-cool-project",
			separator: "",
			lowercase: true,
			want:      "mcp",
		},
		{
			name:      "with separator",
			path:      "/home/user/my-cool-project",
			separator: "-",
			lowercase: true,
			want:      "m-c-p",
		},
		{
			name:      "uppercase",
			path:      "/home/user/my-cool-project",
			separator: "",
			lowercase: false,
			want:      "MCP",
		},
		{
			name:      "underscore separated",
			path:      "/home/user/web_app_project",
			separator: "",
			lowercase: true,
			want:      "wap",
		},
		{
			name:      "single word",
			path:      "/home/user/projects",
			separator: "",
			lowercase: true,
			want:      "p",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateAlias(tt.path, tt.separator, tt.lowercase)
			if got != tt.want {
				t.Errorf("GenerateAlias() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_AddAndGet(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "bookmarks.sh")
	
	m := NewManager(filePath, "bash", "cd", "nvim", "true")
	
	bookmark := domain.Bookmark{
		Alias:       "test",
		Path:        "/home/user/test",
		Description: "Test bookmark",
	}
	
	if err := m.Add(bookmark); err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	
	got, err := m.Get("test")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	
	if got.Alias != bookmark.Alias || got.Path != bookmark.Path {
		t.Errorf("Get() = %v, want %v", got, bookmark)
	}
	
	if got.CreatedAt.IsZero() || got.UpdatedAt.IsZero() {
		t.Error("timestamps should be set")
	}
}

func TestManager_Update(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "bookmarks.sh")
	
	m := NewManager(filePath, "bash", "cd", "nvim", "true")
	
	bookmark := domain.Bookmark{
		Alias: "test",
		Path:  "/home/user/test",
	}
	
	if err := m.Add(bookmark); err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	
	first, _ := m.Get("test")
	time.Sleep(1 * time.Second)
	
	// Update
	bookmark.Path = "/home/user/updated"
	if err := m.Add(bookmark); err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	
	updated, err := m.Get("test")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	
	if updated.Path != "/home/user/updated" {
		t.Errorf("path not updated")
	}
	
	if !updated.CreatedAt.Equal(first.CreatedAt) {
		t.Error("CreatedAt should not change on update")
	}
	
	if !updated.UpdatedAt.After(first.UpdatedAt) {
		t.Errorf("UpdatedAt should be newer: first=%v, updated=%v", first.UpdatedAt, updated.UpdatedAt)
	}
}

func TestManager_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "bookmarks.sh")
	
	m := NewManager(filePath, "bash", "cd", "nvim", "true")
	
	bookmark := domain.Bookmark{
		Alias: "test",
		Path:  "/home/user/test",
	}
	
	if err := m.Add(bookmark); err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	
	if err := m.Delete("test"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	
	_, err := m.Get("test")
	if err != ErrBookmarkNotFound {
		t.Errorf("expected ErrBookmarkNotFound, got %v", err)
	}
}

func TestManager_Exists(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "bookmarks.sh")
	
	m := NewManager(filePath, "bash", "cd", "nvim", "true")
	
	exists, err := m.Exists("test")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if exists {
		t.Error("bookmark should not exist")
	}
	
	bookmark := domain.Bookmark{
		Alias: "test",
		Path:  "/home/user/test",
	}
	
	if err := m.Add(bookmark); err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	
	exists, err = m.Exists("test")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Error("bookmark should exist")
	}
}

func TestIsValidAlias(t *testing.T) {
	tests := []struct {
		alias string
		valid bool
	}{
		{"test", true},
		{"test123", true},
		{"test_abc", true},
		{"Test", true},
		{"test-123", true},       // hyphen is allowed (kebab-case)
		{"my-bookmark", true},    // kebab-case is valid
		{"-", true},              // just hyphen is technically valid
		{"", false},
		{"test#tag", false},      // hash is blacklisted
		{"test space", false},
		{"test@home", false},
		{"test.dot", false},
		{"#", false},             // just hash
		{"tag#1", false},         // contains hash
	}
	
	for _, tt := range tests {
		t.Run(tt.alias, func(t *testing.T) {
			got := isValidAlias(tt.alias)
			if got != tt.valid {
				t.Errorf("isValidAlias(%q) = %v, want %v", tt.alias, got, tt.valid)
			}
		})
	}
}

func TestExpandPath(t *testing.T) {
	home, _ := os.UserHomeDir()
	
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "tilde expansion",
			path: "~/test",
			want: filepath.Join(home, "test"),
		},
		{
			name: "no expansion needed",
			path: "/absolute/path",
			want: "/absolute/path",
		},
		{
			name: "empty path",
			path: "",
			want: "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandPath(tt.path)
			if got != tt.want {
				t.Errorf("expandPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
