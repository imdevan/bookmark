## Feature 1: fix quote escaping

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

## Feature 2: refactor bookmarks file


move alias description from BM comment.
description should be saved as a comment after alias value

# BM: vb|/home/devy/.bookmarks/|vim bookmarks|2026-11-11T18:02:23-08:00|2026-11-11T18:03:26-08:00|vb|||/home/devy/.bookmarks/bookmarks.sh
alias vb='cd '\''/home/devy/.bookmarks/'\'' && tmux rename-window '\''vb'\'' && nvim '\''/home/devy/.bookmarks/bookmarks.sh'\'''

