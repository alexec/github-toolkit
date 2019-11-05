package cmds

import (
	"testing"
)

func TestNewReleaseNoteCmd(t *testing.T) {
	cmd := NewReleaseNoteCmd()
	cmd.Run(cmd, []string{"v0.1..v0.2"})
}