package cmds;

import (
	"github.com/spf13/cobra"
)

func NewReleaseNoteCmd() *cobra.Command {

	var repo GithubRepo
	var commits []string

	var cmd = &cobra.Command{
		Use: "relnote",
		Example: `
  # Get a list of commits:
  git log --oneline v1.2.0..HEAD . > commits
`,
		Run: func(cmd *cobra.Command, args []string) {

			_, _ = newClient(repo, cmd)

		},
	}

	repo = gitHubRepo(cmd)
	cmd.Flags().StringArrayVar(&commits, "commit", []string{}, "TODO")

	return cmd
}
