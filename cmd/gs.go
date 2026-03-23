package cmd

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// ANSI hyperlink helper
func hyperlink(url, text string) string {
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, text)
}

type GHRepo struct {
	Name          string    `json:"name"`
	NameWithOwner string    `json:"nameWithOwner"`
	URL           string    `json:"url"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type GHPR struct {
	Title      string    `json:"title"`
	URL        string    `json:"url"`
	UpdatedAt  time.Time `json:"updatedAt"`
	State      string    `json:"state"`
	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`
}

type GHCommit struct {
	Commit struct {
		Message string `json:"message"`
		Author  struct {
			Date time.Time `json:"date"`
		} `json:"author"`
	} `json:"commit"`
	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`
	URL string `json:"url"`
}

type GHRun struct {
	DatabaseId uint64 `json:"databaseId"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
}

type GHRunJob struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
}

type GHRelease struct {
	Name        string    `json:"name"`
	TagName     string    `json:"tagName"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"publishedAt"`
}

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#2E8B57")).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		MarginBottom(1).
		MarginTop(1)

	repoStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#58A6FF"))
	prStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#3FB950"))
	metaStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#8B949E"))
	commitStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#E3B341"))
	jobStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#8FBC8F"))
)

func formatTimeAgo(t time.Time) string {
	diff := time.Since(t)
	if diff.Hours() > 24 {
		return fmt.Sprintf("%dd ago", int(diff.Hours()/24))
	}
	if diff.Hours() >= 1 {
		return fmt.Sprintf("%dh ago", int(diff.Hours()))
	}
	if diff.Minutes() >= 1 {
		return fmt.Sprintf("%dm ago", int(diff.Minutes()))
	}
	return "just now"
}

func shortMsg(s string, max int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}

var gsCmd = &cobra.Command{
	Use:   "gs",
	Short: "Show latest GitHub stats (repos, PRs, actions, releases) nicely formatted",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := exec.LookPath("gh"); err != nil {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#F85149")).Render("gh cli is not installed. Please install it with: winget install --id GitHub.cli"))
			return nil
		}

		fmt.Println(titleStyle.Render("📊 GitHub Stats Overview"))

		var latestRepo string

		// 1. Repos
		out, err := exec.Command("gh", "repo", "list", "--limit", "4", "--json", "name,nameWithOwner,url,updatedAt").Output()
		if err == nil {
			var repos []GHRepo
			if json.Unmarshal(out, &repos) == nil && len(repos) > 0 {
				latestRepo = repos[0].NameWithOwner
				fmt.Println(lipgloss.NewStyle().Bold(true).Render("📁 Latest Repositories:"))
				for _, r := range repos {
					fmt.Printf("  • %s %s [%s]\n", repoStyle.Render(r.Name), metaStyle.Render(formatTimeAgo(r.UpdatedAt)), hyperlink(r.URL, "Open Link"))
				}
				fmt.Println()
			}
		}

		// 2. Commits (Global)
		outCom, errCom := exec.Command("gh", "search", "commits", "--author", "@me", "--sort", "committer-date", "--limit", "4", "--json", "commit,url,repository").Output()
		if errCom == nil {
			var commits []GHCommit
			if json.Unmarshal(outCom, &commits) == nil && len(commits) > 0 {
				fmt.Println(lipgloss.NewStyle().Bold(true).Render("📝 Latest Commits:"))
				for _, c := range commits {
					fmt.Printf("  • %s in %s %s [%s]\n", commitStyle.Render(shortMsg(c.Commit.Message, 40)), repoStyle.Render(c.Repository.Name), metaStyle.Render(formatTimeAgo(c.Commit.Author.Date)), hyperlink(c.URL, "Open Link"))
				}
				fmt.Println()
			}
		}

		// 3. PRs (Global)
		outPR, errPR := exec.Command("gh", "search", "prs", "--author", "@me", "--limit", "4", "--json", "title,url,updatedAt,state,repository").Output()
		if errPR == nil {
			var prs []GHPR
			if json.Unmarshal(outPR, &prs) == nil && len(prs) > 0 {
				fmt.Println(lipgloss.NewStyle().Bold(true).Render("🔄 Latest Pull Requests:"))
				for _, pr := range prs {
					state := lipgloss.NewStyle().Foreground(lipgloss.Color("#A371F7")).Render(pr.State)
					if strings.ToLower(pr.State) == "open" {
						state = prStyle.Render(pr.State)
					}
					fmt.Printf("  • %s (%s) %s %s [%s]\n", prStyle.Render(shortMsg(pr.Title, 35)), pr.Repository.Name, metaStyle.Render(formatTimeAgo(pr.UpdatedAt)), state, hyperlink(pr.URL, "Open Link"))
				}
				fmt.Println()
			}
		}

		// 4. Actions & Releases (Show for the tracked latest repo)
		if latestRepo != "" {
			// Get Latest Run ID
			outRun, errRun := exec.Command("gh", "run", "list", "--repo", latestRepo, "--limit", "1", "--json", "databaseId,name,status,conclusion").Output()
			if errRun == nil {
				var runs []GHRun
				if json.Unmarshal(outRun, &runs) == nil && len(runs) > 0 {
					run := runs[0]
					fmt.Printf("%s %s\n", lipgloss.NewStyle().Bold(true).Render("⚡ Last Workflow ("+latestRepo+"):"), repoStyle.Render(run.Name))

					// Fetch Jobs for this run
					outJobs, errJobs := exec.Command("gh", "run", "view", fmt.Sprintf("%d", run.DatabaseId), "--repo", latestRepo, "--json", "jobs").Output()
					if errJobs == nil {
						var data struct {
							Jobs []GHRunJob `json:"jobs"`
						}
						if json.Unmarshal(outJobs, &data) == nil && len(data.Jobs) > 0 {
							for _, j := range data.Jobs {
								icon := "⏳"
								concStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#A371F7"))
								if j.Conclusion == "success" {
									icon = "✅"
									concStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#3FB950"))
								} else if j.Conclusion == "failure" {
									icon = "❌"
									concStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F85149"))
								}
								fmt.Printf("    %s %s [%s]\n", icon, jobStyle.Render(j.Name), concStyle.Render(j.Conclusion))
							}
						}
					}
					fmt.Println()
				}
			}

			outRel, errRel := exec.Command("gh", "release", "list", "--repo", latestRepo, "--limit", "2", "--json", "name,tagName,url,publishedAt").Output()
			if errRel == nil {
				var rels []GHRelease
				if json.Unmarshal(outRel, &rels) == nil && len(rels) > 0 {
					fmt.Printf("%s\n", lipgloss.NewStyle().Bold(true).Render("🚀 Recent Releases ("+latestRepo+"):"))
					for _, r := range rels {
						name := r.Name
						if name == "" {
							name = r.TagName
						}
						fmt.Printf("  • %s (%s) %s [%s]\n", prStyle.Render(name), r.TagName, metaStyle.Render(formatTimeAgo(r.PublishedAt)), hyperlink(r.URL, "Open Link"))
					}
					fmt.Println()
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(gsCmd)
}
