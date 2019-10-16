package cmds;

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/v28/github"
	"github.com/peterbourgon/diskv"
	"github.com/spf13/cobra"

	"github.com/alexec/github-issue-cards/cmd/mk/util"
)

func NewReleaseNoteCmd() *cobra.Command {

	var repo GithubRepo
	var commits string

	var cmd = &cobra.Command{
		Use: "relnote",
		Example: `	
	export ACCESS_TOKEN=db015666.. ;# Create an access token at:  https://github.com/settings/tokens
	export OWNER=argoproj
	export REPO=argo-cd
	
	# Get a list of commits:
	git log --format=%H v1.2.0..HEAD . > commits
	
	# Use that list for your note:
	mk relnote --commits "$(cat commits | tr "\n" ,)"
`,
		Run: func(cmd *cobra.Command, args []string) {

			ctx, client := newClient(repo, cmd)
			contributors := map[string]int{}
			var enhancements []string
			var bugFixes []string
			var other []string

			_ = os.MkdirAll("/tmp/relnote/commit", 777)
			_ = os.MkdirAll("/tmp/relnote/issue", 777)
			cache := diskv.New(diskv.Options{
				BasePath:     "/tmp/relnote",
				Transform:    func(s string) []string { return []string{} },
				CacheSizeMax: 1024 * 1024,
			})

			fmt.Println("<!--")
			for _, sha := range strings.Split(commits, ",") {
				if sha == "" {
					continue
				}
				key := "commit/" + sha
				data, err := cache.Read(key)
				commit := &github.Commit{}
				if err == nil {
					err = json.Unmarshal(data, commit)
					util.Check(err)
				} else {
					commit, _, err = client.Git.GetCommit(ctx, repo.owner, repo.repo, sha)
					util.Check(err)
					marshal, err := json.Marshal(commit)
					util.Check(err)
					err = cache.Write(key, marshal)
					util.Check(err)
				}
				// extract the issue and add to the note
				message := strings.SplitN(commit.GetMessage(), "\n", 2)[0]
				fmt.Println(message)
				issues := map[int]bool{}
				for _, text := range regexp.MustCompile("#[0-9]+").FindAllString(message, -1) {
					id, err := strconv.Atoi(strings.TrimPrefix(text, "#"))
					util.Check(err)
					_, ok := issues[id]
					issues[id] = true
					if !ok {
						key := fmt.Sprintf("issue/%v", id)
						data, err = cache.Read(key)
						issue := &github.Issue{}
						if err == nil {
							err := json.Unmarshal(data, issue)
							util.Check(err)
						} else {
							issue, _, err = client.Issues.Get(ctx, repo.owner, repo.repo, id)
							util.Check(err)
							data, err := json.Marshal(issue)
							util.Check(err)
							err = cache.Write(key, data)
							util.Check(err)
						}

						labels := map[string]bool{}
						for _, l := range issue.Labels {
							labels[*l.Name] = true
						}
						if issue.IsPullRequest() {
							continue
						}
						if labels["enhancement"] {
							enhancements = append(enhancements, message)
						} else if labels["bug"] {
							bugFixes = append(bugFixes, message)
						} else {
							other = append(other, message)
						}
					}
				}
				// add the author as a contributor
				name := *commit.Author.Name
				num, ok := contributors[name]
				if ok {
					contributors[name] = num + 1
				} else {
					contributors[name] = 1
				}

			}
			fmt.Println("-->")
			fmt.Println("#### New Features")
			fmt.Println()
			fmt.Println("TODO")
			fmt.Println()
			fmt.Println("#### Enhancements")
			fmt.Println()
			for _, i := range enhancements {
				fmt.Printf("* %s\n", i)
			}
			fmt.Println()
			fmt.Println("#### Bug Fixes")
			fmt.Println()
			for _, i := range bugFixes {
				fmt.Printf("- %s\n", i)
			}
			fmt.Println()
			fmt.Println("#### Other")
			fmt.Println()
			for _, i := range other {
				fmt.Printf("- %s\n", i)
			}
			fmt.Println()
			fmt.Println("#### Contributors")
			fmt.Println()
			fmt.Println()
			for name, num := range contributors {
				fmt.Printf("* %s <!-- num=%v -->\n", name, num)
			}
		},
	}

	repo = gitHubRepo(cmd)
	cmd.Flags().StringVar(&commits, "commits", "string", "List of commits")
	_ = cmd.MarkFlagRequired("commit")

	return cmd
}
