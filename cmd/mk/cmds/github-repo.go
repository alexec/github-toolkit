package cmds

import (
	"context"
	"os"

	"github.com/google/go-github/v28/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

type GithubRepo struct {
	accessToken string
	owner       string
	repo        string
}

func gitHubRepo(cmd *cobra.Command) GithubRepo {
	repo := GithubRepo{}
	cmd.Flags().StringVar(&repo.accessToken, "access-token", os.Getenv("ACCESS_TOKEN"), "Github personal access token")
	cmd.Flags().StringVar(&repo.owner, "owner", os.Getenv("OWNER"), "Github owner (aka org)")
	cmd.Flags().StringVar(&repo.repo, "repo", os.Getenv("REPO"), "Github repo")
	return repo
}

func newClient(repo GithubRepo, cmd *cobra.Command) (context.Context, *github.Client) {
	if repo.accessToken == "" || repo.owner == "" || repo.repo == "" {
		_ = cmd.Usage()
		os.Exit(1)
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: repo.accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return ctx, client
}
