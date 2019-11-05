package cmds

import (
	"context"
	"fmt"
	"net/url"
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
	host        string
	owner       string
	repo        string
}

func (repo GithubRepo) BaseURL() *url.URL {
	text := ""
	if repo.host == "github.com" {
		text = fmt.Sprintf("https://api.%v/", repo.host)
	} else {
		text = fmt.Sprintf("https://%s/api/v3/", repo.host)
	}
	u, err := url.Parse(text)
	util.Check(err)
	return u
}

func NewGithubRepo(url string) (repo GithubRepo) {
	repo.accessToken = os.Getenv("ACCESS_TOKEN")

	if strings.HasPrefix(url, "http") {
		parts := strings.Split(url, "/")
		repo.host = parts[2]
		repo.owner = parts[3]
		repo.repo = strings.Split(parts[4], ".")[0]
	} else {
		parts := strings.Split(url, ":")
		repo.host = strings.Split(parts[0], "@")[1]
		parts = strings.Split(parts[1], "/")
		repo.owner = parts[0]
		repo.repo = strings.Split(parts[1], ".")[0]
	}
	return
}

func gitHubRepo() GithubRepo {
	return NewGithubRepo(repoUrl())
}

func repoUrl() string {
	git := exec.Command("git", "config", "--get", "remote.origin.url")
	bytes, err := git.Output()
	util.Check(err)
	return string(bytes)
}

func newClient(repo GithubRepo, cmd *cobra.Command) (context.Context, *github.Client) {
	if repo.host == "" || repo.accessToken == "" || repo.owner == "" || repo.repo == "" {
		_ = cmd.Usage()
		os.Exit(1)
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: repo.accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	client.BaseURL = repo.BaseURL()
	return ctx, client
}
