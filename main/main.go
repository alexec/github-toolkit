package main;

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
	"gopkg.in/go-playground/colors.v1"
)

func main() {
	// https://github.com/settings/tokens
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("ACCESS_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	issues, _, err := client.Issues.ListByRepo(ctx, "argoproj", "argo-cd", &github.IssueListByRepoOptions{
		State:  "open",
		Labels: []string{"enhancement"},
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

    <title>%v cards</title>
  </head>
  <body>
<div class="container">
<div class="card-columns">
`, len(issues))
	for _, i := range issues {

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

		fmt.Printf(`<div class="card">
  <div class="card-body">
    <h5 class="card-title">%s</h5></h5>
    <h6 class="card-subtitle mb-2">%s</h6>
    <h6 class="card-subtitle mb-2 text-muted">#%v opened on %v by %v</h6>
    <p class="card-text">%s %s</p>
  </div>
</div>`,
			i.GetTitle(),
			labels,
			i.GetNumber(),
			fmt.Sprintf("%s %v", i.GetCreatedAt().Month(), i.GetCreatedAt().Day()),
			i.GetUser().GetLogin(),
			reactions,
			comments,
		)
	}
	fmt.Println(`  </div></div>
</body>
</html>`)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
