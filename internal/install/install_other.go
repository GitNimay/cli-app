//go:build !windows

package install

import "fmt"

// installCMD is not supported on non-Windows platforms
func installCMD(opts InstallOptions) error {
	return fmt.Errorf("CMD AutoRun is only supported on Windows. Use --shell bash instead")
}

// uninstallCMD is not supported on non-Windows platforms
func uninstallCMD(opts InstallOptions) error {
	return fmt.Errorf("CMD AutoRun is only supported on Windows")
}
