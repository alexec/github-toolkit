package main

import (
	"github.com/spf13/cobra"

	"github.com/alexec/github-issue-cards/cmd/mk/cmds"
	"github.com/alexec/github-issue-cards/cmd/mk/util"
)

func main() {
	cmd := &cobra.Command{
		Use:   "gt",
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
	}
	cmd.AddCommand(cmds.NewCardsCmd())
	cmd.AddCommand(cmds.NewReleaseNoteCmd())
	util.Check(cmd.Execute())
}
