Convert this template project into a fully functioning go cli app. 

Features: 
- root command bookmark current folder
- optional: pass bookmark string to generate an alias
  - default naming convention based on the first letters of each "word" in current dir
- -i flag: view filter able list of bookmarks
  - with crud interaction
- optional: -t flag to define tmux name to rename window when navigating to location
- optional: define post jump script that is run in bookmark 
  - via: `alias bm="navigation && ...`
- define description via comment after alias
- handles confirming before overwriting existing bookmark

Config:
- which tool to use to navigate: none, cd, z, etc
- define which shell the user uses
- define where the users bookmarks are located (default to ~/.bookmarks/...)
- possible: define per shell locations for bookmarks

