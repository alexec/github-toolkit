package cmds;

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/google/go-github/v28/github"
	"github.com/peterbourgon/diskv"
	"github.com/spf13/cobra"

	"github.com/alexec/github-toolkit/cmd/ght/util"
)

func NewReleaseNoteCmd() *cobra.Command {

	var repo githubRepo
	cache := true

	var cmd = &cobra.Command{
		Use:   "relnote REVISION_RANGE",
		Short: "Create release note based on Github issue.",
		Example: `  # Create the note:
  ght relnote v1.3.0-rc3..v1.3.0-rc4`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}
			revisionRange := args[0]

			ctx, client := newClient(repo, cmd)

			base := filepath.Join("/tmp", "relnote", repo.owner, repo.repo)

			_ = os.MkdirAll(base, 0777)
			diskCache := diskv.New(diskv.Options{
				BasePath:     base,
				Transform:    func(s string) []string { return []string{} },
				CacheSizeMax: 1024 * 1024,
			})

			output, err := exec.Command("git", "log", "--format=%H", revisionRange, "--", ".").Output()
			util.Check(err)
			var issues []int
			contributors := map[string]int{}
			var other []string
			for _, sha := range strings.Split(string(output), "\n") {
				if sha == "" {
					continue
				}
				key := "commit." + sha
				data, err := diskCache.Read(key)
				commit := &github.Commit{}
				if cache && err == nil {
					util.Check(json.Unmarshal(data, commit))
				} else {
					commit, _, err = client.Git.GetCommit(ctx, repo.owner, repo.repo, sha)
					util.Check(err)
					marshal, err := json.Marshal(commit)
					util.Check(err)
					util.Check(diskCache.Write(key, marshal))
				}
				// extract the issue and add to the note
				message := strings.SplitN(commit.GetMessage(), "\n", 2)[0]

				foundIssues := findIssues(message)
				if len(foundIssues) == 0 {
					other = append(other, message)
				} else {
					issues = append(issues, foundIssues...)
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

			done := make(map[int]bool)
			var enhancements []string
			var bugFixes []string
			var pullRequests []string
			for ; len(issues) > 0; {
				var id int
				id, issues = issues[len(issues)-1], issues[:len(issues)-1]
				_, ok := done[id]
				done[id] = true
				if !ok {
					key := fmt.Sprintf("issue.%v", id)
					data, err := diskCache.Read(key)
					issue := &github.Issue{}
					if err == nil {
						err := json.Unmarshal(data, issue)
						util.Check(err)
					} else {
						issue, _, err = client.Issues.Get(ctx, repo.owner, repo.repo, id)
						util.Check(err)
						data, err := json.Marshal(issue)
						util.Check(err)
						err = diskCache.Write(key, data)
						util.Check(err)
					}

					labels := map[string]bool{}
					for _, l := range issue.Labels {
						labels[*l.Name] = true
					}
					message := fmt.Sprintf("#%v %s", id, *issue.Title)
					if labels["enhancement"] {
						enhancements = append(enhancements, message)
					} else if labels["bug"] {
						bugFixes = append(bugFixes, message)
					} else if issue.IsPullRequest() {
						relatedIssues := findIssues(*issue.Body)
						if len(relatedIssues) > 0 {
							issues = append(issues, relatedIssues...)
						} else {
							pullRequests = append(pullRequests, message)
						}
					} else {
						other = append(other, message)
					}
				}
			}

			if len(enhancements) > 0 {
				fmt.Println("#### Enhancements")
				fmt.Println()
				sort.Strings(enhancements)
				for _, i := range enhancements {
					fmt.Printf("* %s\n", i)
				}
				fmt.Println()
			}
			if len(bugFixes) > 0 {
				fmt.Println("#### Bug Fixes")
				fmt.Println()
				sort.Strings(bugFixes)
				for _, i := range bugFixes {
					fmt.Printf("- %s\n", i)
				}
				fmt.Println()
			}
			if len(other) > 0 {
				fmt.Println("#### Other")
				fmt.Println()
				sort.Strings(other)
				for _, i := range other {
					fmt.Printf("- %s\n", i)
				}
				fmt.Println()
			}
			if len(pullRequests) > 0 {
				fmt.Println("#### Pull Requests")
				fmt.Println()
				sort.Strings(pullRequests)
				for _, i := range pullRequests {
					fmt.Printf("- %s\n", i)
				}
				fmt.Println()
			}
			fmt.Println("#### Contributors")
			fmt.Println()
			var names []string
			for name := range contributors {
				names = append(names, name)
			}
			sort.Strings(names)
			for _, name := range names {
				fmt.Printf("* %s <!-- num=%v -->\n", name, contributors[name])
			}
		},
	}

	repo = gitHubRepo()
	cmd.Flags().Bool("cache", true, "Use a cache")

	return cmd
}

func findIssues(message string) []int {
	var issues []int
	for _, text := range regexp.MustCompile("#[0-9]+").FindAllString(message, 1) {
		id, err := strconv.Atoi(strings.TrimPrefix(text, "#"))
		util.Check(err)
		issues = append(issues, id)
	}
	return issues
}
