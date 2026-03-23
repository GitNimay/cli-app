//go:build windows

package install

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// installCMD installs to CMD AutoRun registry (Windows only)
func installCMD(opts InstallOptions) error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER,
		`Software\Microsoft\Command Processor`, registry.SET_VALUE|registry.QUERY_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	commandLine := fmt.Sprintf(`"%s"`, opts.ExecPath)

	// Check existing value
	existing, _, err := key.GetStringValue("AutoRun")
	if err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("failed to read registry: %w", err)
	}

	if strings.Contains(existing, commandLine) && !opts.Force {
		fmt.Println("✅ HelloGang is already installed in CMD AutoRun.")
		return nil
	}

	// Append to existing or set new value
	var newValue string
	if existing != "" {
		newValue = existing + " & " + commandLine
	} else {
		newValue = commandLine
	}

	if err := key.SetStringValue("AutoRun", newValue); err != nil {
		return fmt.Errorf("failed to set registry value: %w", err)
	}

	fmt.Println("✅ Successfully installed HelloGang to CMD AutoRun!")
	fmt.Println("   Registry: HKCU\\Software\\Microsoft\\Command Processor\\AutoRun")
	return nil
}

// uninstallCMD removes from CMD AutoRun (Windows only)
func uninstallCMD(opts InstallOptions) error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER,
		`Software\Microsoft\Command Processor`, registry.SET_VALUE|registry.QUERY_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	existing, _, err := key.GetStringValue("AutoRun")
	if err != nil {
		fmt.Println("ℹ️  No AutoRun registry key found.")
		return nil
	}

	newValue := removeCommand(existing, "hellogang")

	if newValue == "" {
		if err := key.DeleteValue("AutoRun"); err != nil {
			return fmt.Errorf("failed to delete registry value: %w", err)
		}
	} else {
		if err := key.SetStringValue("AutoRun", newValue); err != nil {
			return fmt.Errorf("failed to set registry value: %w", err)
		}
	}

	fmt.Println("✅ Successfully uninstalled HelloGang from CMD AutoRun.")
	return nil
}
