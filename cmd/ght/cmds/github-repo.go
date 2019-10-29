package cmds

import (
	"context"
	"os"
	"os/exec"
	"strings"

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

func NewGithubRepo(url string) (repo GithubRepo) {
	repo.accessToken = os.Getenv("ACCESS_TOKEN")

	if strings.HasPrefix(url, "http") {
		parts := strings.Split(url, "/")
		repo.owner = parts[3]
		repo.repo = strings.Split(parts[4], ".")[0]
	} else {
		parts := strings.Split(strings.Split(url, ":")[1], "/")
		repo.owner = parts[0]
		repo.repo = strings.Split(parts[1], ".")[0]
	}
	return
}

func gitHubRepo(cmd *cobra.Command) GithubRepo {
	repo := NewGithubRepo(repoUrl())
	cmd.Flags().StringVar(&repo.accessToken, "access-token", repo.accessToken, "Github personal access token, create one at https://github.com/settings/tokens")
	cmd.Flags().StringVar(&repo.owner, "owner", repo.owner, "Github owner (aka org)")
	cmd.Flags().StringVar(&repo.repo, "repo", repo.repo, "Github repo")
	return repo
}

func repoUrl() string {
	git := exec.Command("git", "config", "--get", "remote.origin.url")
	bytes, err := git.Output()
	util.Check(err)
	return string(bytes)
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
