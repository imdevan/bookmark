# Context

This project is a Go CLI application for managing directory bookmarks with shell integration. The tool allows users to quickly save and navigate to frequently-used directories, with support for custom aliases, tmux integration, and post-navigation scripts.

# Definitions

- **Bookmark**: A saved directory path with an associated alias for quick navigation
- **Alias**: A short name used to reference a bookmark (e.g., `bm-proj` for `/home/user/projects`)
- **Post-jump script**: Custom commands executed after navigating to a bookmarked directory
- **CRUD**: Create, Read, Update, Delete operations for managing bookmarks
- **Navigation tool**: The underlying command used for directory changes (cd, z, zoxide, etc.)

# Data Structures

## Bookmark File Format (TOML)

```toml
# ~/.bookmarks/bookmarks.toml

[[bookmarks]]
alias = "proj"
path = "/home/user/projects/myapp"
description = "Main application project"
created_at = "2026-02-24T10:30:00Z"
updated_at = "2026-02-24T10:30:00Z"
tmux_window_name = "myapp"
post_jump_script = "source .env && echo 'Welcome to myapp'"

[[bookmarks]]
alias = "docs"
path = "/home/user/projects/myapp/docs"
description = "Documentation folder"
created_at = "2026-02-24T11:00:00Z"
updated_at = "2026-02-24T11:00:00Z"
```

# Command Examples

## 1. Add Bookmark (Current Directory)

```bash
# Auto-generated alias
$ cd /home/user/projects/my-cool-project
$ bookmark
✓ Bookmark created: mcp → /home/user/my/cool/project

# Auto-generated alias with same tmux name
$ cd /home/user/projects/my-cool-project
$ bookmark -t
✓ Bookmark created: mcp → /home/user/my/cool/project


# Custom alias
$ cd /home/user/projects/webapp
$ bookmark web
✓ Bookmark created: web → /home/user/projects/webapp

# Custom alias with same tmux name
$ cd /home/user/projects/webapp
$ bookmark web
✓ Bookmark created: web → /home/user/projects/webapp

# Custom alias with custom tmux name
$ cd /home/user/projects/webapp
$ bookmark web -t wweb
✓ Bookmark created: web → /home/user/projects/webapp

# With description
$ bookmark web --description "Main web application"
✓ Bookmark created: web → /home/user/projects/webapp

# Overwrite confirmation
$ bookmark web
? Bookmark 'web -> ~/projects/webapp' already exists. Overwrite? (y/N) y
use ui/confirmation
✓ Bookmark updated: web → /home/user/projects/webapp
```

## 2. Interactive Browser

```bash
$ bookmark -i
# or
$ bookmark --interactive

# Output (Bubble Tea UI):
use ui/list
```

## 3. List Bookmarks

```bash
$ bookmark list
proj     /home/user/projects/myapp          Main application project
docs     /home/user/projects/myapp/docs     Documentation folder
web      /home/user/projects/webapp
scripts  /home/user/scripts

## 4. Navigate to Bookmark

bookmarks are saved as aliases that are managed by the users shell

```bash
$ proj
# Outputs shell command to stdout:
cd /home/user/projects/myapp

# With tmux integration
$ bookmark go proj -t myapp
tmux rename-window 'myapp' && cd /home/user/projects/myapp

# With post-jump script
$ bookmark go proj
cd /home/user/projects/myapp && source .env && echo 'Welcome to myapp'
```

## 5. Delete Bookmark

```bash
$ bookmark delete web
? Delete bookmark 'web'? (y/N) y
use ui/confirmation
✓ Bookmark deleted: web

$ bookmark delete web --force
✓ Bookmark deleted: web
```

## 6. Edit Bookmark

```bash
$ bookmark edit proj
# Opens bookmark_file in editor at bookmark location

$ bookmark edit proj --description "Updated description"
✓ Bookmark updated: proj

$ bookmark edit proj --path /new/path
✓ Bookmark updated: proj → /new/path
```

## 7. Shell Integration

```bash
# Generate shell function
$ bookmark shell-init bash

echo "source $bookmark_file" >> ."$bookarm_shell"rc
```

# v0.1.0

- [ ] 1. Core Bookmark Management
---
  - [x] 1.1 Root command to bookmark current folder
    - notes: `bookmark [alias]` command saves current directory with auto-generated or custom alias
    - example: `bookmark` → creates alias "mcp" for `/my/cool/project`
    - example: `bookmark web` → creates alias "web" for current directory
    - output: `✓ Bookmark created: web → /home/user/projects/webapp`

  - [x] 1.2 Auto-generate alias from directory name
    - notes: Default naming convention uses first letters of each "word" in current dir
    - example: `/my/cool/project` → `mcp`
    - example: `/home/user/dev` → `hud`
    - example: `/projects/web-app` → `pwa`

  - [x] 1.3 Optional custom alias via argument
    - notes: Allow user to pass custom bookmark string: `bookmark my-alias`
    - example: `bookmark web` saves current dir as "web"
    - validation: alias must be alphanumeric + hyphens/underscores only

  - [x] 1.4 Confirmation before overwriting existing bookmark
    - notes: Prompt user if alias already exists before replacing
    - example: `? Bookmark 'web' already exists. Overwrite? (y/N)`
    - behavior: default to No, require explicit confirmation

  - [ ] 1.5 -y to accept overwrite from flag
    - bypass confirmations

  - [ ] 1.6 -f file flag
    - [ ] - f and file to add open file in editor to bookmark, executed after navigation
      - file bookmarks will navigate to folder and open file in configured editor

  - [ ] 1.7 -e flag
    - no args: open bookmarks_file in configured editor
    - existing bookmark: open in editor at bookmark line
    - non-existing bookmark: create and open in editor

- [ ] 2. Interactive Bookmark Browser
---
  - [x] 2.1 `-i` flag to view filterable list of bookmarks
    - notes: Display all bookmarks with search/filter capability using Bubble Tea inline UI
    - example input: `bookmark -i` or `bookmark --interactive`
    - example output:
      ```
      ┌─ Bookmarks ────────────────────────────────────┐
      │ Search: proj_                                  │
      ├────────────────────────────────────────────────┤
      │ > proj     ~/projects/myapp                    │
      │   docs     ~/projects/myapp/docs               │
      └────────────────────────────────────────────────┘
      ```
    - implementation: Use `internal/ui/list.go` pattern with Bubble Tea

    - [ ] 2.1.1 home_symbol = "~" (default) use as standin for $HOME directory
    - [ ] 2.1.2 default_sort_by = "newist" (default) | "latest" | "A to Z" | "Z to A" use as standin for $HOME directory


  - [ ] 2.2 CRUD operations in interactive mode
    - notes: Support Create, Read, Update, Delete actions from the list view
    - keybindings:
      - `Enter`: Navigate to selected bookmark
      - `e`: Edit selected bookmark
      - `d`: Delete selected bookmark (with confirmation)
      - `n`: Create new bookmark
      - `q`: or `Esc`: Quit
    - example: Press `d` → `? Delete bookmark 'proj'? (y/N)`

    - [ ] 2.2.1 enter: Enter`: Navigate to selected bookmark
    - [ ] 2.2.2 e: open selected book mark in editor # DO NOT IMPLEMENT
            - resource bookmarks on save (if in nvim)
            - alternative: open alias in input box to edit
    - [ ] 2.2.3 d/D: delete selected bookmark (d: confirmation | D: no conf)
    - [ ] 2.2.4 n: Create new bookmark
    - [x] 2.2.5 q`: or `Esc`: Quit



  - [ ] 2.3 Navigate to selected bookmark
    - notes: Select bookmark from list to navigate to that directory
    - behavior: Output shell command to stdout for evaluation
    - example output: `cd /home/user/projects/myapp`

- [ ] 3. Advanced Features
---
  - [ ] 3.1 `-t` flag for tmux window naming
    - notes: Optional flag to define tmux window name when navigating to bookmark
    - example: `bookmark go proj -t myapp`
    - output: `tmux rename-window 'myapp' && cd /home/user/projects/myapp`
    - storage: Save `tmux_window_name` field in bookmark TOML

  - [ ] 3.2 Post-jump script execution
    - notes: Define and execute custom scripts after navigation
    - example config in bookmark:
      ```toml
      post_jump_script = "source .env && echo 'Welcome!'"
      ```
    - output: `cd /path && source .env && echo 'Welcome!'`
    - validation: Escape shell special characters for safety

  - [ ] 3.3 Bookmark descriptions via comments
    - notes: Support adding descriptions/comments to bookmarks for documentation
    - example: `bookmark web --description "Main web application"`
    - storage: Save `description` field in bookmark TOML
    - display: Show in list view and interactive browser

- [ ] 4. Configuration System
---
  - [ ] 4.1 Navigation tool selection
    - notes: Config option to choose navigation method: none, cd, z, zoxide, etc.
    - config field: `navigation_tool = "cd"`
    - valid values: `"cd"`, `"z"`, `"zoxide"`, `"none"`
    - behavior: Changes output command format (e.g., `z /path` vs `cd /path`)

  - [ ] 4.2 Shell type configuration
    - notes: Define which shell the user uses (bash, zsh, fish, etc.)
    - config field: `shell = "bash"`
    - valid values: `"bash"`, `"zsh"`, `"fish"`
    - usage: Affects shell-init command output format

  - [ ] 4.3 Bookmark storage location
    - notes: Configurable bookmark file location
    - config field: `bookmarks_file = "~/.bookmarks/bookmarks.sh"`
    - default: `~/.bookmarks/bookmarks.sh`
    - validation: Expand tilde, create parent directories if needed

- [ ] 5. Shell Integration
---
  - [ ] 5.1 Generate shell-specific aliases
    - notes: Output shell commands that can be sourced for navigation
    - command: `bookmark shell-init <shell>`
    - example: `bookmark shell-init bash`
    - output for bash:
      ```bash
      bm() {
        local output=$(bookmark go "$@")
        if [ $? -eq 0 ]; then
          eval "$output"
        fi
      }
      ```
    - output for zsh:
      ```zsh
      bm() {
        local output=$(bookmark go "$@")
        if [[ $? -eq 0 ]]; then
          eval "$output"
        fi
      }
      ```
    - output for fish:
      ```fish
      function bm
        set output (bookmark go $argv)
        if test $status -eq 0
          eval $output
        end
      end
      ```

  - [ ] 5.2 Shell function generation
    - notes: Create wrapper functions for seamless shell integration
    - usage instructions: `eval "$(bookmark shell-init bash)"`
    - add to shell rc: `echo 'eval "$(bookmark shell-init bash)"' >> ~/.bashrc`
    - behavior: Function wraps `bookmark go` and evaluates output in current shell

#  v0.2.0

  - [ ] 1 Per-shell bookmark locations
    - notes: Optional: support different bookmark files for different shells
    - config example:
      ```toml
      [shell_bookmarks]
      bash = "~/.bookmarks/.toml"
      zsh = "~/.bookmarks/zsh.toml"
      ```
    - behavior: Falls back to `bookmarks_file` if shell-specific not set

- [ ] 2 pin via comment? 
    ```
    alias marker = "cd ~/marker" # Go to marker - pin
    ```
