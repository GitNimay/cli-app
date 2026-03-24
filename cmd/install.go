package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"hellogang/internal/config"
	"hellogang/internal/install"
)

var (
	installShell string
	installForce bool
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install HelloGang to run on shell startup",
	Long: `Installs HelloGang so it runs automatically every time you open
a new terminal session. Supports PowerShell, CMD, and Git Bash.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// --- Auto-detect Name ---
		var defaultName string
		if currentUser, err := user.Current(); err == nil {
			defaultName = currentUser.Name
			if defaultName == "" {
				defaultName = currentUser.Username
			}
			// Get first name if there's a space
			if idx := strings.Index(defaultName, " "); idx != -1 {
				defaultName = defaultName[:idx]
			}
			// Strip domain from username if any (e.g., DOMAIN\\user)
			if idx := strings.Index(defaultName, "\\"); idx != -1 {
				defaultName = defaultName[idx+1:]
			}
		}
		if defaultName == "" {
			defaultName = os.Getenv("USERNAME")
			if defaultName == "" {
				defaultName = "User"
			}
		}

		name := defaultName

		fmt.Printf("✨ What is your name? [%s]: ", defaultName)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			name = input
		}

		if name != "" {
			if err := config.SetName(name); err != nil {
				fmt.Printf("⚠️  Could not save name to config: %v\n", err)
			}
		}

		fmt.Printf("👋 Hello %s! Automatically configuring HelloGang for your terminal...\n", name)

		var shell install.ShellType

		switch installShell {
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
			return fmt.Errorf("unknown shell type: %s (use: powershell, cmd, bash, or auto)", installShell)
		}

		if err := install.Install(install.InstallOptions{
			Shell: shell,
			Force: installForce,
		}); err != nil {
			return err
		}

		if runtime.GOOS == "windows" {
			// For CMD, also add to Windows Startup because CMD AutoRun is unreliable
			if shell == install.ShellCMD {
				fmt.Print("\n⚡ Adding to Windows Startup for reliable CMD support? [Y/n]: ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(strings.ToLower(input))
				if input != "n" && input != "no" {
					if err := install.InstallStartupApp(install.InstallOptions{
						Force: installForce,
					}); err != nil {
						fmt.Printf("⚠️  Could not add to startup: %v\n", err)
					}
				}
			} else {
				fmt.Print("\nWould you like to run HelloGang as a startup application? [y/N]: ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(strings.ToLower(input))
				if input == "y" || input == "yes" {
					if err := install.InstallStartupApp(install.InstallOptions{
						Force: installForce,
					}); err != nil {
						fmt.Printf("⚠️  Could not add to startup: %v\n", err)
					}
				}
			}
		}

		fmt.Println("\n🎉 Installation complete! Open a new terminal to see HelloGang.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVarP(&installShell, "shell", "s", "auto",
		"Shell to install to (powershell, cmd, bash, auto, prompt)")
	installCmd.Flags().BoolVarP(&installForce, "force", "f", false,
		"Force reinstall even if already installed")
}
