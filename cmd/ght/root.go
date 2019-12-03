package main

import (
	"github.com/spf13/cobra"

	"github.com/alexec/github-toolkit/cmd/ght/cmds"
	"github.com/alexec/github-toolkit/cmd/ght/util"
)

func main() {
	cmd := &cobra.Command{
		Use:   "ght",
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
	}
	cmd.AddCommand(cmds.NewCardsCmd())
	cmd.AddCommand(cmds.NewChangeLogCmd())
	cmd.AddCommand(cmds.NewReleaseNoteCmd())
	util.Check(cmd.Execute())
}
