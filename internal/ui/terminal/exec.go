package terminal

import (
	"os"
	"os/exec"
	"sync/atomic"

	tea "github.com/charmbracelet/bubbletea"
)

// ExecProcess hands the terminal to cmd for the length of its run.
//
// Bubble Tea's program output is wrapped for IME cursor placement, but
// os/exec only gives a child a TTY when Stdout is the actual *os.File.
func ExecProcess(cmd *exec.Cmd, fn tea.ExecCallback) tea.Cmd {
	if cmd.Stdout == nil {
		cmd.Stdout = os.Stdout
	}
	return tea.ExecProcess(cmd, AfterHandoff(fn))
}

var handedBack atomic.Bool

// AfterHandoff notes that a program had the terminal and gave it back.
// Bubble Tea switches mouse reporting off for the handoff and its
// RestoreTerminal does not switch it back on, and it restores before it runs
// fn, so the root can arm the mouse again on the next update (TakeHandedBack).
func AfterHandoff(fn tea.ExecCallback) tea.ExecCallback {
	return func(err error) tea.Msg {
		handedBack.Store(true)
		return fn(err)
	}
}

// TakeHandedBack reports a handoff that ended since the last call.
func TakeHandedBack() bool { return handedBack.Swap(false) }
