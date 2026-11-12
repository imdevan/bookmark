## Feature 1: fix quote escaping [x]

current saved aliases are double escaping single quote strings. refactor the string parser such that the external quote is single and all internal quotes are double-quotes; escaped as necessary (if a user were to pass an alias with `alias foo='echo "bar"'` or similar. 

current:
```
# BM: vb|/home/devy/.bookmarks/|vim bookmarks|2026-11-11T18:02:23-08:00|2026-11-11T18:03:26-08:00|vb|||/home/devy/.bookmarks/bookmarks.sh
alias vb='cd '\''/home/devy/.bookmarks/'\'' && tmux rename-window '\''vb'\'' && nvim '\''/home/devy/.bookmarks/bookmarks.sh'\'''
```

should look like this:
```
# BM: vb|/home/devy/.bookmarks/|vim bookmarks|2026-11-11T18:02:23-08:00|2026-11-11T18:03:26-08:00|vb|||/home/devy/.bookmarks/bookmarks.sh
alias vb='cd "/home/devy/.bookmarks/" && tmux rename-window "vb" && nvim "/home/devy/.bookmarks/bookmarks.sh"'
```

- [x] POSIX (bash/zsh/sh) inner args (path, tmux window name, file) are wrapped in double quotes in `escapeShellArg` (`internal/bookmark/manager.go`)
- [x] Escape `\`, `"`, `$` and backtick inside the double-quoted args so they stay literal
- [x] Outer alias stays single-quoted; any user-supplied `'` (e.g. in execute scripts) is still escaped by the alias formatter
- [x] Fish and Nushell escaping unchanged
- [x] Update tests (`manager_escaping_test.go`, `manager_spaces_test.go`, `cmd/bookmark/integration_test.go`) and add tests for `"`/`$`/`\` and the exact alias output

## Feature 2: refactor bookmarks file


move alias description from BM comment.
description should be saved as a comment after alias value

# BM: vb|/home/devy/.bookmarks/|vim bookmarks|2026-11-11T18:02:23-08:00|2026-11-11T18:03:26-08:00|vb|||/home/devy/.bookmarks/bookmarks.sh
alias vb='cd '\''/home/devy/.bookmarks/'\'' && tmux rename-window '\''vb'\'' && nvim '\''/home/devy/.bookmarks/bookmarks.sh'\'''

