package cmd

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"hellogang/internal/install"
)

var (
	uninstallShell string
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove HelloGang from shell startup",
	Long: `Removes HelloGang from your shell's startup configuration.
This undoes what the 'install' command set up.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var shell install.ShellType

		switch uninstallShell {
		case "powershell", "ps":
			shell = install.ShellPowerShell
		case "cmd":
			shell = install.ShellCMD
		case "bash", "git-bash":
			shell = install.ShellBash
		case "auto", "":
			shell = install.DetectShell()
			fmt.Printf("🔍 Detected shell: %s\n", shell)
		case "prompt":
			shell = install.PromptForShell()
		default:
			return fmt.Errorf("unknown shell type: %s (use: powershell, cmd, bash, or auto)", uninstallShell)
		}

		if err := install.Uninstall(install.InstallOptions{
			Shell: shell,
		}); err != nil {
			return err
		}

		if runtime.GOOS == "windows" {
			reader := bufio.NewReader(os.Stdin)
			// For CMD, always prompt about startup removal since we added it automatically
			if shell == install.ShellCMD {
				fmt.Print("\nWould you also like to remove HelloGang from Windows startup? [Y/n]: ")
			} else {
				fmt.Print("\nWould you also like to remove HelloGang from Windows startup? [y/N]: ")
			}
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(strings.ToLower(input))
			if input != "n" && input != "no" {
				if err := install.UninstallStartupApp(install.InstallOptions{}); err != nil {
					fmt.Printf("⚠️  Could not remove from startup: %v\n", err)
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	uninstallCmd.Flags().StringVarP(&uninstallShell, "shell", "s", "auto",
		"Shell to uninstall from (powershell, cmd, bash, auto, prompt)")
}
