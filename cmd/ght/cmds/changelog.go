package cmds

import (
	"fmt"
	"os"

	"github.com/google/go-github/v28/github"
	"github.com/spf13/cobra"

	"github.com/alexec/github-toolkit/cmd/ght/util"
)

func NewChangeLogCmd() *cobra.Command {

	var repo githubRepo

	var cmd = &cobra.Command{
		Use:   "changelog",
		Short: "Print a CHANGELOG.md",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 0 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}
			ctx, client := newClient(repo, cmd)

			var releases, _, err = client.Repositories.ListReleases(ctx, repo.owner, repo.repo, &github.ListOptions{
				PerPage: 100,
			})
			util.Check(err)
			fmt.Println("# Changelog")
			fmt.Println()
			for _, release := range releases {
				if *release.Prerelease {
					continue
				}
				fmt.Printf("## %s (%v)\n", *release.Name, release.PublishedAt.Format("2006-01-02"))
				fmt.Println()
				fmt.Println(*release.Body)
				fmt.Println()
			}
		},
	}

	repo = gitHubRepo()
	return cmd
}
