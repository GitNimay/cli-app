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
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GHPR struct {
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	UpdatedAt time.Time `json:"updatedAt"`
	State     string    `json:"state"`
}

type GHRun struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	URL        string `json:"url"`
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

	repoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#58A6FF"))
	prStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#3FB950"))
	metaStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8B949E"))
)

// cleanStr safely strips any carriage returns or weird ANSI if needed; for JSON it's usually clean
func cleanStr(s string) string { return strings.TrimSpace(s) }

var gsCmd = &cobra.Command{
	Use:   "gs",
	Short: "Show latest GitHub stats (repos, PRs, actions, releases) nicely formatted",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := exec.LookPath("gh"); err != nil {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#F85149")).Render("gh cli is not installed. Please install it with: winget install --id GitHub.cli"))
			return nil
		}

		fmt.Println(titleStyle.Render("📊 GitHub Stats Overview"))

		// 1. Repos
		out, err := exec.Command("gh", "repo", "list", "--limit", "5", "--json", "name,url,updatedAt").Output()
		if err == nil {
			var repos []GHRepo
			if json.Unmarshal(out, &repos) == nil && len(repos) > 0 {
				fmt.Println(lipgloss.NewStyle().Bold(true).Render("📁 Latest Repositories:"))
				for _, r := range repos {
					timeAgo := time.Since(r.UpdatedAt).Round(time.Hour)
					fmt.Printf("  • %s %s [%s]\n", repoStyle.Render(cleanStr(r.Name)), metaStyle.Render(timeAgo.String()+" ago"), hyperlink(r.URL, "Open Link"))
				}
				fmt.Println()
			}
		}

		// 2. PRs
		outPR, errPR := exec.Command("gh", "pr", "list", "--author", "@me", "--limit", "4", "--json", "title,url,updatedAt,state").Output()
		if errPR == nil {
			var prs []GHPR
			if json.Unmarshal(outPR, &prs) == nil {
				fmt.Println(lipgloss.NewStyle().Bold(true).Render("🔄 Latest Pull Requests:"))
				if len(prs) == 0 {
					fmt.Println(metaStyle.Render("  No recent pull requests."))
				} else {
					for _, pr := range prs {
						timeAgo := time.Since(pr.UpdatedAt).Round(time.Hour)
						state := lipgloss.NewStyle().Foreground(lipgloss.Color("#A371F7")).Render(pr.State)
						if strings.ToLower(pr.State) == "open" {
							state = prStyle.Render(pr.State)
						}
						fmt.Printf("  • %s %s [%s] %s\n", prStyle.Render(cleanStr(pr.Title)), metaStyle.Render(timeAgo.String()+" ago"), state, hyperlink(pr.URL, "Open Link"))
					}
				}
				fmt.Println()
			}
		}

		// 3. Actions (current repo)
		outRun, errRun := exec.Command("gh", "run", "list", "--limit", "4", "--json", "name,status,conclusion,url").Output()
		if errRun == nil {
			var runs []GHRun
			if json.Unmarshal(outRun, &runs) == nil && len(runs) > 0 {
				fmt.Println(lipgloss.NewStyle().Bold(true).Render("⚡ Recent Actions (Current Repo):"))
				for _, r := range runs {
					status := r.Status
					if r.Conclusion != "" {
						status = r.Conclusion
					}
					fmt.Printf("  • %s [%s] %s\n", repoStyle.Render(cleanStr(r.Name)), status, hyperlink(r.URL, "Open Link"))
				}
				fmt.Println()
			}
		}

		// 4. Releases (current repo)
		outRel, errRel := exec.Command("gh", "release", "list", "--limit", "3", "--json", "name,tagName,url,publishedAt").Output()
		if errRel == nil {
			var rels []GHRelease
			if json.Unmarshal(outRel, &rels) == nil && len(rels) > 0 {
				fmt.Println(lipgloss.NewStyle().Bold(true).Render("🚀 Recent Releases (Current Repo):"))
				for _, r := range rels {
					timeAgo := time.Since(r.PublishedAt).Round(time.Hour)
					name := r.Name
					if name == "" {
						name = r.TagName
					}
					fmt.Printf("  • %s (%s) %s %s\n", prStyle.Render(cleanStr(name)), r.TagName, metaStyle.Render(timeAgo.String()+" ago"), hyperlink(r.URL, "Open Link"))
				}
				fmt.Println()
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(gsCmd)
}
