package cmds;

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/v28/github"
	"github.com/peterbourgon/diskv"
	"github.com/spf13/cobra"

	"github.com/alexec/github-toolkit/cmd/ght/util"
)

func NewReleaseNoteCmd() *cobra.Command {

	var repo GithubRepo

	var cmd = &cobra.Command{
		Use:   "relnote REVISION_RANGE",
		Short: "Create release note based on Github issue.",
		Long:  `TODO`,
		Example: `	
	# Create the note:
	ght relnote release-1.3..HEAD
`,
		Run: func(cmd *cobra.Command, args []string) {

			revisionRange := args[0]

			ctx, client := newClient(repo, cmd)
			contributors := map[string]int{}
			var enhancements []string
			var bugFixes []string
			var pullRequests []string
			var other []string

			_ = os.MkdirAll("/tmp/relnote/commit", 777)
			_ = os.MkdirAll("/tmp/relnote/issue", 777)
			cache := diskv.New(diskv.Options{
				BasePath:     "/tmp/relnote",
				Transform:    func(s string) []string { return []string{} },
				CacheSizeMax: 1024 * 1024,
			})

			output, err := exec.Command("git", "log", "--format=%H", revisionRange, "--", ".").Output()
			util.Check(err)
			for _, sha := range strings.Split(string(output), "\n") {
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
				issues := map[int]bool{}
				for _, id := range findIssues(message) {
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
						if labels["enhancement"] {
							enhancements = append(enhancements, message)
						} else if labels["bug"] {
							bugFixes = append(bugFixes, message)
						} else if issue.IsPullRequest() {
							// TODO - we should be better at attributing PRs to non-PR issues
							pullRequests = append(pullRequests, message)
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
			if len(enhancements) > 0 {
				fmt.Println("#### Enhancements")
				fmt.Println()
				for _, i := range enhancements {
					fmt.Printf("* %s\n", i)
				}
				fmt.Println()
			}
			if len(bugFixes) > 0 {
				fmt.Println("#### Bug Fixes")
				fmt.Println()
				for _, i := range bugFixes {
					fmt.Printf("- %s\n", i)
				}
				fmt.Println()
			}
			if len(other) > 0 {
				fmt.Println("#### Other")
				fmt.Println()
				for _, i := range other {
					fmt.Printf("- %s\n", i)
				}
				fmt.Println()
			}
			if len(pullRequests) > 0 {
				fmt.Println("#### Pull Requests")
				fmt.Println()
				for _, i := range pullRequests {
					fmt.Printf("- %s\n", i)
				}
				fmt.Println()
			}
			fmt.Println("#### Contributors")
			fmt.Println()
			for name, num := range contributors {
				fmt.Printf("* %s <!-- num=%v -->\n", name, num)
			}
		},
	}

	repo = gitHubRepo(cmd)
	_ = cmd.MarkFlagRequired("commit")

	return cmd
}

func findIssues(message string) []int {
	var issues []int
	for _, text:= range regexp.MustCompile("#[0-9]+").FindAllString(message, 1) {
		id, err := strconv.Atoi(strings.TrimPrefix(text, "#"))
		util.Check(err)
		issues = append(issues, id)
	}
	return issues
}
