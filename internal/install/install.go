package install

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ShellType represents the type of shell
type ShellType string

const (
	ShellPowerShell ShellType = "powershell"
	ShellCMD        ShellType = "cmd"
	ShellBash       ShellType = "bash"
	ShellUnknown    ShellType = "unknown"
)

// InstallOptions holds installation options
type InstallOptions struct {
	Shell    ShellType
	Force    bool
	ExecPath string
}

// DetectShell detects the current shell type
func DetectShell() ShellType {
	// Check environment variables
	shell := os.Getenv("SHELL")
	term := os.Getenv("TERM_PROGRAM")
	psModulePath := os.Getenv("PSModulePath")

	// Check for Git Bash
	if strings.Contains(strings.ToLower(shell), "bash") ||
		strings.Contains(strings.ToLower(os.Getenv("MSYSTEM")), "mingw") {
		return ShellBash
	}

	// Check for PowerShell
	if term == "pwsh" || term == "powershell" ||
		strings.Contains(psModulePath, "PowerShell") ||
		strings.Contains(psModulePath, "WindowsPowerShell") {
		return ShellPowerShell
	}

	// Check for Zsh / Fish on Unix
	if strings.Contains(strings.ToLower(shell), "zsh") ||
		strings.Contains(strings.ToLower(shell), "fish") {
		return ShellBash // treat as bash-like
	}

	// Default to CMD on Windows, Bash on others
	if runtime.GOOS == "windows" {
		return ShellCMD
	}

	return ShellBash
}

// Install installs the greeting to run on shell startup
func Install(opts InstallOptions) error {
	// Get executable path if not provided
	if opts.ExecPath == "" {
		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %w", err)
		}
		opts.ExecPath = execPath
	}

	// Auto-detect shell if not specified
	if opts.Shell == ShellUnknown {
		opts.Shell = DetectShell()
	}

	switch opts.Shell {
	case ShellPowerShell:
		return installPowerShell(opts)
	case ShellCMD:
		return installCMD(opts)
	case ShellBash:
		return installBash(opts)
	default:
		return fmt.Errorf("unsupported shell type: %s", opts.Shell)
	}
}

// installPowerShell installs to PowerShell profile
func installPowerShell(opts InstallOptions) error {
	// Get PowerShell profile path
	profilePath, err := getPowerShellProfilePath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	profileDir := filepath.Dir(profilePath)
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		return fmt.Errorf("failed to create profile directory: %w", err)
	}

	// Check if already installed
	content, err := os.ReadFile(profilePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read profile: %w", err)
	}

	commandLine := fmt.Sprintf("& '%s'", opts.ExecPath)
	if strings.Contains(string(content), commandLine) {
		if !opts.Force {
			fmt.Println("✅ HelloGang is already installed in PowerShell profile.")
			return nil
		}
	}

	// Append to profile
	f, err := os.OpenFile(profilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open profile: %w", err)
	}
	defer f.Close()

	installLine := fmt.Sprintf("\n# HelloGang - Auto-start greeting\n%s\n", commandLine)
	if _, err := f.WriteString(installLine); err != nil {
		return fmt.Errorf("failed to write to profile: %w", err)
	}

	fmt.Println("✅ Successfully installed HelloGang to PowerShell profile!")
	fmt.Printf("   Profile: %s\n", profilePath)
	return nil
}

// getPowerShellProfilePath returns the PowerShell profile path
func getPowerShellProfilePath() (string, error) {
	// Try PowerShell Core first (pwsh)
	cmd := exec.Command("pwsh", "-NoProfile", "-Command", "echo $PROFILE")
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		return strings.TrimSpace(string(output)), nil
	}

	// Fallback to Windows PowerShell
	cmd = exec.Command("powershell", "-NoProfile", "-Command", "echo $PROFILE")
	output, err = cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get PowerShell profile path: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// installBash installs to .bashrc
func installBash(opts InstallOptions) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	bashrcPath := filepath.Join(home, ".bashrc")

	// Check if already installed
	content, err := os.ReadFile(bashrcPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read .bashrc: %w", err)
	}

	commandLine := fmt.Sprintf(`'%s'`, opts.ExecPath)
	if strings.Contains(string(content), commandLine) && !opts.Force {
		fmt.Println("✅ HelloGang is already installed in .bashrc.")
		return nil
	}

	// Append to .bashrc
	f, err := os.OpenFile(bashrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open .bashrc: %w", err)
	}
	defer f.Close()

	installLine := fmt.Sprintf("\n# HelloGang - Auto-start greeting\n%s\n", commandLine)
	if _, err := f.WriteString(installLine); err != nil {
		return fmt.Errorf("failed to write to .bashrc: %w", err)
	}

	fmt.Println("✅ Successfully installed HelloGang to .bashrc!")
	fmt.Printf("   File: %s\n", bashrcPath)
	return nil
}

// Uninstall removes the greeting from shell startup
func Uninstall(opts InstallOptions) error {
	if opts.Shell == ShellUnknown {
		opts.Shell = DetectShell()
	}

	switch opts.Shell {
	case ShellPowerShell:
		return uninstallPowerShell(opts)
	case ShellCMD:
		return uninstallCMD(opts)
	case ShellBash:
		return uninstallBash(opts)
	default:
		return fmt.Errorf("unsupported shell type: %s", opts.Shell)
	}
}

// uninstallPowerShell removes from PowerShell profile
func uninstallPowerShell(opts InstallOptions) error {
	profilePath, err := getPowerShellProfilePath()
	if err != nil {
		return err
	}

	content, err := os.ReadFile(profilePath)
	if err != nil {
		return fmt.Errorf("failed to read profile: %w", err)
	}

	newContent := removeLines(string(content), "hellogang")

	if string(content) == newContent {
		fmt.Println("ℹ️  HelloGang was not found in PowerShell profile.")
		return nil
	}

	if err := os.WriteFile(profilePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write profile: %w", err)
	}

	fmt.Println("✅ Successfully uninstalled HelloGang from PowerShell profile.")
	return nil
}

// uninstallBash removes from .bashrc
func uninstallBash(opts InstallOptions) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	bashrcPath := filepath.Join(home, ".bashrc")
	content, err := os.ReadFile(bashrcPath)
	if err != nil {
		return fmt.Errorf("failed to read .bashrc: %w", err)
	}

	newContent := removeLines(string(content), "hellogang")

	if string(content) == newContent {
		fmt.Println("ℹ️  HelloGang was not found in .bashrc.")
		return nil
	}

	if err := os.WriteFile(bashrcPath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write .bashrc: %w", err)
	}

	fmt.Println("✅ Successfully uninstalled HelloGang from .bashrc.")
	return nil
}

// removeLines removes lines containing a string
func removeLines(content, target string) string {
	lines := strings.Split(content, "\n")
	var newLines []string

	for _, line := range lines {
		if strings.Contains(line, target) {
			continue
		}
		// Also remove comment lines we added
		if strings.Contains(line, "# HelloGang") {
			continue
		}
		newLines = append(newLines, line)
	}

	return strings.Join(newLines, "\n")
}

// removeCommand removes a command from a compound command string
func removeCommand(content, target string) string {
	parts := strings.Split(content, " & ")
	var newParts []string

	for _, part := range parts {
		if strings.Contains(part, target) {
			continue
		}
		newParts = append(newParts, strings.TrimSpace(part))
	}

	return strings.Join(newParts, " & ")
}

// PromptForShell asks the user which shell to install to
func PromptForShell() ShellType {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\nWhich shell would you like to install HelloGang to?")
	fmt.Println("  1. PowerShell")
	fmt.Println("  2. Command Prompt (CMD)")
	fmt.Println("  3. Git Bash")
	fmt.Print("\nEnter choice [1-3]: ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	switch input {
	case "1":
		return ShellPowerShell
	case "2":
		return ShellCMD
	case "3":
		return ShellBash
	default:
		return DetectShell()
	}
}