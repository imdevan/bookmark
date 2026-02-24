package editor

import "errors"

// Adapter launches the configured editor.
type Adapter struct {
	Command string
}

// New returns an editor adapter using the given command.
func New(command string) *Adapter {
	return &Adapter{Command: command}
}

// Open launches the editor with the provided file path.
func (a Adapter) Open(path string) error {
	command := ResolveCommand(a.Command)
	if command == "" {
		return errors.New("editor command is required")
	}
	return runEditorCommand(command, []string{path})
}

// OpenAtLine opens a file at a specific line number.
func (a Adapter) OpenAtLine(path string, line int) error {
	command := ResolveCommand(a.Command)
	if command == "" {
		return errors.New("editor command is required")
	}
	
	if IsVim(command) {
		return OpenVimAtLine(command, path, line)
	}
	if IsNano(command) {
		return OpenNanoAtLine(command, path, line)
	}
	if IsVSCode(command) {
		return OpenVSCodeAtLine(command, path, line)
	}
	if IsEmacs(command) {
		return OpenEmacsAtLine(command, path, line)
	}
	
	// For unknown editors, just open normally
	return a.Open(path)
}

// OpenAtEnd opens a file and positions the cursor at the end when supported.
func (a Adapter) OpenAtEnd(path string) error {
	command := ResolveCommand(a.Command)
	if command == "" {
		return errors.New("editor command is required")
	}
	if IsVim(command) {
		return OpenVimAtEnd(command, path)
	}
	return a.Open(path)
}
