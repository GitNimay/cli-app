package cmd

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	vercelTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#FFFFFF")).
		Padding(0, 1).
		MarginBottom(1).
		MarginTop(1)

	vercelItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#0070F3"))
	vercelDimStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
)

// urlPattern is used to extract URLs easily
var urlPattern = regexp.MustCompile(`https://[^\s]+`)

// cleanVercelLine removes unnecessary control characters or spinner leftovers
func cleanVercelLine(line string) string {
	line = strings.TrimSpace(line)
	// Some Vercel CLI lines start with "⠋" or other spinners
	for _, pfx := range []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏", "Vercel CLI", ">"} {
		if strings.HasPrefix(line, pfx) {
			line = strings.TrimSpace(strings.TrimPrefix(line, pfx))
		}
	}
	return line
}

var vsCmd = &cobra.Command{
	Use:   "vs",
	Short: "Show latest Vercel stats (projects, deployments) nicely formatted",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if vercel is installed
		if _, err := exec.LookPath("vercel"); err != nil {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#F85149")).Render("vercel cli is not installed. Please install it with: npm i -g vercel"))
			return nil
		}

		fmt.Println(vercelTitleStyle.Render("▲ Vercel Stats Overview"))

		// Check whoami
		outWho, errWho := exec.Command("vercel", "whoami").Output()
		if errWho == nil {
			user := strings.TrimSpace(string(outWho))
			fmt.Printf("👤 Logged in as: %s\n\n", vercelItemStyle.Render(user))
		}

		// Fetch Projects - using CombinedOutput because Vercel often prints to stderr
		fmt.Println(lipgloss.NewStyle().Bold(true).Render("🚀 Latest Projects:"))
		outProj, _ := exec.Command("vercel", "project", "ls").CombinedOutput()
		lines := strings.Split(string(outProj), "\n")
		foundProjects := false
		for _, rawLine := range lines {
			line := cleanVercelLine(rawLine)
			// Print lines that look like table rows
			if line != "" && !strings.Contains(line, "Fetching") && !strings.Contains(line, "Projects found under") && !strings.Contains(line, "Project Name") && !strings.Contains(line, "50.") {
				// Highlight URLs if present
				if match := urlPattern.FindString(line); match != "" {
					line = strings.Replace(line, match, vercelItemStyle.Render(match), 1)
					// Can also add ANSI hyperlink: \x1b]8;;%[1]s\x1b\\%[1]s\x1b]8;;\x1b\\
					link := fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", match, "Open Link")
					fmt.Printf("  • %s [%s]\n", line, link)
				} else {
					fmt.Printf("  %s\n", vercelDimStyle.Render(line))
				}
				foundProjects = true
			}
		}
		if !foundProjects {
			fmt.Println(vercelDimStyle.Render("  No projects found or unable to fetch."))
		}
		fmt.Println()

		// Attempt to fetch deployments if inside a Vercel project (--yes prevents prompt)
		outDeploy, _ := exec.Command("vercel", "ls", "--yes").CombinedOutput()
		dLines := strings.Split(string(outDeploy), "\n")
		hasDeploys := false
		for _, rawLine := range dLines {
			line := cleanVercelLine(rawLine)
			if match := urlPattern.FindString(line); match != "" {
				if !hasDeploys {
					fmt.Println(lipgloss.NewStyle().Bold(true).Render("🌐 Recent Deployments (Current Project):"))
					hasDeploys = true
				}
				link := fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", match, "Open Link")
				fmt.Printf("  • %s [%s]\n", vercelDimStyle.Render(line), link)
			}
		}
		if hasDeploys {
			fmt.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(vsCmd)
}
