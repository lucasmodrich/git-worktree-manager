package config

import (
	"path/filepath"
	"runtime"
)

// GetBinaryPath returns the full path to the gwtm binary in the install directory
func GetBinaryPath() string {
	installDir := GetInstallDir()
	name := "gwtm"
	if runtime.GOOS == "windows" {
		name = "gwtm.exe"
	}
	return filepath.Join(installDir, name)
}
