package cmds

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/google/go-github/v28/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/alexec/github-toolkit/cmd/ght/util"
)

type GithubRepo struct {
	accessToken string
	owner       string
	repo        string
}

func gitHubRepo(cmd *cobra.Command) GithubRepo {

	git := exec.Command("git", "config", "--get", "remote.origin.url")
	bytes, err := git.Output()
	util.Check(err)
	fmt.Printf(string(bytes))

	repo := GithubRepo{}
	cmd.Flags().StringVar(&repo.accessToken, "access-token", os.Getenv("ACCESS_TOKEN"), "Github personal access token, create one at https://github.com/settings/tokens")
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
