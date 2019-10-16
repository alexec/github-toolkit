package main

import (
	"github.com/spf13/cobra"

	"github.com/alexec/github-issue-cards/cmd/mk/cmds"
	"github.com/alexec/github-issue-cards/cmd/mk/util"
)

func main() {
	cmd := &cobra.Command{
		Use:   "mk",
		Short: "argocd controls a Argo CD server",
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
	}
	cmd.AddCommand(cmds.NewCardsCmd())
	util.Check(cmd.Execute())
}
