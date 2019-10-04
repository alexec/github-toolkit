package main;

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v28/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"gopkg.in/go-playground/colors.v1"
)

var accessToken string
var owner string
var repo string
var labels []string
var excludeLabels []string
var since string

var rootCmd = &cobra.Command{
	Use: "mkcards",
	Example: `
	export ACCESS_TOKEN=db015666.. ;# Create an access token at:  https://github.com/settings/tokens

	mkcards --owner argoproj --repo argo-cd --exclude-label 'wontfix' --exclude-label 'workaround' --exclude-label 'help wanted' > enhancements.html
	mkcards --owner argoproj --repo argo-cd --label 'bug' --exclude-label 'wontfix' --exclude-label 'workaround' --exclude-label 'help wanted' > bugs.html
`,
	Run: func(cmd *cobra.Command, args []string) {

		if accessToken == "" {
			_ = cmd.Usage()
			os.Exit(1)
		}

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)
		issues, _, err := client.Issues.ListByRepo(ctx, owner, repo, &github.IssueListByRepoOptions{
			State:  "open",
			Labels: labels,
			Sort:   "comments",
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		})
		check(err)
		fmt.Printf(`<!doctype html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">

    <title>Max %v Cards</title>
  </head>
  <body>
<div class="container">
<div class="card-columns">
`, len(issues))
		for _, i := range issues {

			skip := false
			for _, l := range i.Labels {
				for _, e := range excludeLabels {
					if l.GetName() == e {
						skip = true
					}
				}
			}
			if skip {
				continue
			}

			labels := ""
			for _, l := range i.Labels {
				color, err := colors.Parse("#" + l.GetColor())
				check(err)
				fg := "white"
				if color.IsLight() {
					fg = "black"
				}
				labels = labels + fmt.Sprintf(`<span class="badge badge-secondary" style="background-color:%s;color:%s"'>%s</span>`, "#"+l.GetColor(), fg, l.GetName())
			}
			reactions := ""
			if i.GetReactions().GetPlusOne() > 0 {
				reactions = fmt.Sprintf("👍 %v", i.GetReactions().GetPlusOne())
			}
			comments := ""
			if i.GetComments() > 0 {
				comments = fmt.Sprintf("💬 %v", i.GetComments())
			}
			milestone := ""
			if i.GetMilestone() != nil {
				milestone = fmt.Sprintf("%s", i.GetMilestone().GetTitle())
			}

			fmt.Printf(`<div class="card">
  <div class="card-body">
    <h5 class="card-title">%s</h5>
    <h6 class="card-subtitle mb-2">%s</h6>
    <h6 class="card-subtitle mb-2 text-muted">#%v opened on %v by %v</h6>
    <p class="card-text">%s %s %s</p>
  </div>
</div>`,
				i.GetTitle(),
				labels,
				i.GetNumber(),
				i.GetCreatedAt().Format("2 Jan"),
				i.GetUser().GetLogin(),
				reactions,
				comments,
				milestone,
			)
		}
		fmt.Println(`  </div></div>
</body>
</html>`)

		if len(issues) >= 100 {
			panic("100 or more issues, we do not support pagination, so we do not support this number of issues")
		}
	},
}

func init() {
	rootCmd.Flags().StringVar(&accessToken, "access-token", os.Getenv("ACCESS_TOKEN"), "Github personal access token")
	rootCmd.Flags().StringVar(&owner, "owner", "", "Github owner (aka org)")
	rootCmd.Flags().StringVar(&repo, "repo", "", "Github repo")
	rootCmd.Flags().StringArrayVar(&labels, "label", []string{"enhancement"}, "Github labels")
	rootCmd.Flags().StringArrayVar(&excludeLabels, "exclude-label", []string{}, "Github labels no exclude")
	rootCmd.Flags().StringVar(&since, "since", "", "Github issue since, e.g. ")
	_ = rootCmd.MarkFlagRequired("owner")
	_ = rootCmd.MarkFlagRequired("repo")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	check(rootCmd.Execute())
}
