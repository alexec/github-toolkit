package cmds;

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v28/github"
	"github.com/hako/durafmt"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"gopkg.in/go-playground/colors.v1"

	"github.com/alexec/github-issue-cards/cmd/mk/util"
)

type GitRepo struct {
	accessToken string
	owner       string
	repo        string
}

func NewCardsCmd() *cobra.Command {

	var repo GitRepo
	var state string
	var labels []string
	var excludeLabels []string
	var since time.Duration
	var milestone string

	var cmd = &cobra.Command{
		Use: "cards",
		Example: `
	export ACCESS_TOKEN=db015666.. ;# Create an access token at:  https://github.com/settings/tokens

	# enhancements backlog 
	mk cards --owner argoproj --repo argo-cd --label enhancement --exclude-label wontfix --milestone none 
	
	# bugs backlog
	mk cards --owner argoproj --repo argo-cd --label bug --exclude-label wontfix --milestone none 

	# help wanted backlog
	mk cards --owner argoproj --repo argo-cd --label 'help wanted' --exclude-label wontfix' --milestone none 

	# open issues in milestone v1.3
	mk cards --owner argoproj --repo argo-cd  --milestone v1.3

	# issues opened in the last day
	mk cards --owner argoproj --repo argo-cd  --state all --since 24h
`,
		Run: func(cmd *cobra.Command, args []string) {

			if repo.accessToken == "" {
				_ = cmd.Usage()
				os.Exit(1)
			}

			ctx := context.Background()
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: repo.accessToken},
			)
			tc := oauth2.NewClient(ctx, ts)
			client := github.NewClient(tc)

			switch milestone {
			case "none", "*":
			default:
				milestones, _, err := client.Issues.ListMilestones(ctx, repo.owner, repo.repo, &github.MilestoneListOptions{})
				util.Check(err)
				for _, m := range milestones {
					if m.GetTitle() == milestone {
						milestone = fmt.Sprintf("%v", m.GetNumber())
						break
					}
				}
			}

			// https://developer.github.com/v3/issues/#list-issues-for-a-repository
			issues, _, err := client.Issues.ListByRepo(ctx, repo.owner, repo.repo, &github.IssueListByRepoOptions{
				State:       state,
				Labels:      labels,
				Sort:        "update",
				Milestone:   milestone,
				Since:       time.Now().Add(-since),
				ListOptions: github.ListOptions{PerPage: 100},
			})
			util.Check(err)
			fmt.Printf(`<!doctype html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">

    <title>Max %v Cards</title>
<style type='text/css'>
a {
	color: black;
}
.open path {
	fill: green;
}
.closed path {
	fill: red;
}
.merged path {
	fill: purple;
}
</style>
  </head>
  <body>
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
					util.Check(err)
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
					signpost := `<svg aria-label="Milestone" class="octicon octicon-milestone" viewBox="0 0 14 16" version="1.1" width="14" height="16" role="img"><path fill-rule="evenodd" d="M8 2H6V0h2v2zm4 5H2c-.55 0-1-.45-1-1V4c0-.55.45-1 1-1h10l2 2-2 2zM8 4H6v2h2V4zM6 16h2V8H6v8z"></path></svg>`
					milestone = fmt.Sprintf(`%s %s`, signpost, i.GetMilestone().GetTitle())
				}
				icon := `<svg class="octicon octicon-issue-closed closed" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M7 10h2v2H7v-2zm2-6H7v5h2V4zm1.5 1.5l-1 1L12 9l4-4.5-1-1L12 7l-1.5-1.5zM8 13.7A5.71 5.71 0 0 1 2.3 8c0-3.14 2.56-5.7 5.7-5.7 1.83 0 3.45.88 4.5 2.2l.92-.92A6.947 6.947 0 0 0 8 1C4.14 1 1 4.14 1 8s3.14 7 7 7 7-3.14 7-7l-1.52 1.52c-.66 2.41-2.86 4.19-5.48 4.19v-.01z"></path></svg>`
				if i.GetState() == "open" {
					icon = `<svg class="octicon octicon-issue-opened open" viewBox="0 0 14 16" version="1.1" width="14" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M7 2.3c3.14 0 5.7 2.56 5.7 5.7s-2.56 5.7-5.7 5.7A5.71 5.71 0 0 1 1.3 8c0-3.14 2.56-5.7 5.7-5.7zM7 1C3.14 1 0 4.14 0 8s3.14 7 7 7 7-3.14 7-7-3.14-7-7-7zm1 3H6v5h2V4zm0 6H6v2h2v-2z"></path></svg>`
				}
				if i.IsPullRequest() {
					// icon = `<svg class="octicon octicon-git-pull-request" viewBox="0 0 12 16" version="1.1" width="12" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M11 11.28V5c-.03-.78-.34-1.47-.94-2.06C9.46 2.35 8.78 2.03 8 2H7V0L4 3l3 3V4h1c.27.02.48.11.69.31.21.2.3.42.31.69v6.28A1.993 1.993 0 0 0 10 15a1.993 1.993 0 0 0 1-3.72zm-1 2.92c-.66 0-1.2-.55-1.2-1.2 0-.65.55-1.2 1.2-1.2.65 0 1.2.55 1.2 1.2 0 .65-.55 1.2-1.2 1.2zM4 3c0-1.11-.89-2-2-2a1.993 1.993 0 0 0-1 3.72v6.56A1.993 1.993 0 0 0 2 15a1.993 1.993 0 0 0 1-3.72V4.72c.59-.34 1-.98 1-1.72zm-.8 10c0 .66-.55 1.2-1.2 1.2-.65 0-1.2-.55-1.2-1.2 0-.65.55-1.2 1.2-1.2.65 0 1.2.55 1.2 1.2zM2 4.2C1.34 4.2.8 3.65.8 3c0-.65.55-1.2 1.2-1.2.65 0 1.2.55 1.2 1.2 0 .65-.55 1.2-1.2 1.2z"></path></svg>`
					icon = `<svg class="octicon octicon-git-merge merged" viewBox="0 0 12 16" version="1.1" width="12" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M10 7c-.73 0-1.38.41-1.73 1.02V8C7.22 7.98 6 7.64 5.14 6.98c-.75-.58-1.5-1.61-1.89-2.44A1.993 1.993 0 0 0 2 .99C.89.99 0 1.89 0 3a2 2 0 0 0 1 1.72v6.56c-.59.35-1 .99-1 1.72 0 1.11.89 2 2 2a1.993 1.993 0 0 0 1-3.72V7.67c.67.7 1.44 1.27 2.3 1.69.86.42 2.03.63 2.97.64v-.02c.36.61 1 1.02 1.73 1.02 1.11 0 2-.89 2-2 0-1.11-.89-2-2-2zm-6.8 6c0 .66-.55 1.2-1.2 1.2-.65 0-1.2-.55-1.2-1.2 0-.65.55-1.2 1.2-1.2.65 0 1.2.55 1.2 1.2zM2 4.2C1.34 4.2.8 3.65.8 3c0-.65.55-1.2 1.2-1.2.65 0 1.2.55 1.2 1.2 0 .65-.55 1.2-1.2 1.2zm8 6c-.66 0-1.2-.55-1.2-1.2 0-.65.55-1.2 1.2-1.2.65 0 1.2.55 1.2 1.2 0 .65-.55 1.2-1.2 1.2z"></path></svg>`
					if i.GetState() == "open" {
						icon = `<svg class="octicon octicon-git-pull-request open" viewBox="0 0 12 16" version="1.1" width="12" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M11 11.28V5c-.03-.78-.34-1.47-.94-2.06C9.46 2.35 8.78 2.03 8 2H7V0L4 3l3 3V4h1c.27.02.48.11.69.31.21.2.3.42.31.69v6.28A1.993 1.993 0 0 0 10 15a1.993 1.993 0 0 0 1-3.72zm-1 2.92c-.66 0-1.2-.55-1.2-1.2 0-.65.55-1.2 1.2-1.2.65 0 1.2.55 1.2 1.2 0 .65-.55 1.2-1.2 1.2zM4 3c0-1.11-.89-2-2-2a1.993 1.993 0 0 0-1 3.72v6.56A1.993 1.993 0 0 0 2 15a1.993 1.993 0 0 0 1-3.72V4.72c.59-.34 1-.98 1-1.72zm-.8 10c0 .66-.55 1.2-1.2 1.2-.65 0-1.2-.55-1.2-1.2 0-.65.55-1.2 1.2-1.2.65 0 1.2.55 1.2 1.2zM2 4.2C1.34 4.2.8 3.65.8 3c0-.65.55-1.2 1.2-1.2.65 0 1.2.55 1.2 1.2 0 .65-.55 1.2-1.2 1.2z"></path></svg>`

					}
				}

				issueType := "issue"
				if i.IsPullRequest() {
					issueType = "pull"
				}
				fmt.Printf(`<div class="card">
  <div class="card-body">
    <h5 class="card-title"><a href="https://github.com/%s/%s/%s/%v">%s %s</a></h5>
    <h6 class="card-subtitle mb-2">%s</h6>
    <h6 class="card-subtitle mb-2 text-muted">#%v opened %v ago by %v %s</h6>
    <p class="card-text">%s %s</p>
  </div>
</div>`,
					repo.owner,
					repo.repo,
					issueType,
					i.GetNumber(),
					icon,
					i.GetTitle(),
					labels,
					i.GetNumber(),
					durafmt.ParseShort(time.Since(i.GetCreatedAt())),
					i.GetUser().GetLogin(),
					milestone,
					reactions,
					comments,
				)
			}
			fmt.Println(`  </div>
</body>
</html>`)

			if len(issues) >= 100 {
				panic("100 or more issues, we do not support pagination, so we do not support this number of issues")
			}
		},
	}

	cmd.Flags().StringVar(&repo.accessToken, "access-token", os.Getenv("ACCESS_TOKEN"), "Github personal access token")
	cmd.Flags().StringVar(&repo.owner, "owner", "", "Github owner (aka org)")
	cmd.Flags().StringVar(&repo.repo, "repo", "", "Github repo")
	cmd.Flags().StringVar(&state, "state", "open", "Github issue state, 'all', 'open' or 'closed'")
	cmd.Flags().StringArrayVar(&labels, "label", []string{}, "Github labels")
	cmd.Flags().StringArrayVar(&excludeLabels, "exclude-label", []string{}, "Github labels no exclude")
	cmd.Flags().DurationVar(&since, "since", 20*24*365*time.Hour, "Github issue since, e.g. 24h")
	cmd.Flags().StringVar(&milestone, "milestone", "*", "Github milestone, can be 'none', '*', or the title")
	_ = cmd.MarkFlagRequired("owner")
	_ = cmd.MarkFlagRequired("repo")

	return cmd
}
