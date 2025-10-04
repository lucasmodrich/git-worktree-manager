package version

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/lucasmodrich/git-worktree-manager/internal/config"
)

// UpgradeToLatest downloads and installs the latest version
func UpgradeToLatest(currentVersion, latestVersion string) error {
	// Parse and compare versions
	current, err := ParseVersion(currentVersion)
	if err != nil {
		return fmt.Errorf("invalid current version: %w", err)
	}

	latest, err := ParseVersion(latestVersion)
	if err != nil {
		return fmt.Errorf("invalid latest version: %w", err)
	}

	if !latest.GreaterThan(current) {
		return fmt.Errorf("already on latest version %s", currentVersion)
	}

	// Determine binary name for current platform
	binaryName := getBinaryName()

	// Download binary
	binaryURL := fmt.Sprintf("https://github.com/lucasmodrich/git-worktree-manager/releases/download/v%s/%s", latestVersion, binaryName)

	tempBinary := filepath.Join(os.TempDir(), "git-worktree-manager-new")
	if err := downloadFile(binaryURL, tempBinary); err != nil {
		return fmt.Errorf("failed to download binary: %w", err)
	}
	defer os.Remove(tempBinary)

	// Download checksum file
	checksumURL := fmt.Sprintf("https://github.com/lucasmodrich/git-worktree-manager/releases/download/v%s/checksums.txt", latestVersion)
	checksumFile := filepath.Join(os.TempDir(), "checksums.txt")
	if err := downloadFile(checksumURL, checksumFile); err != nil {
		return fmt.Errorf("failed to download checksums: %w", err)
	}
	defer os.Remove(checksumFile)

	// Verify checksum
	if err := verifyChecksum(tempBinary, checksumFile, binaryName); err != nil {
		return fmt.Errorf("checksum verification failed: %w", err)
	}

	// Download additional files
	installDir := config.GetInstallDir()
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	files := []string{"README.md", "VERSION", "LICENSE"}
	for _, file := range files {
		url := fmt.Sprintf("https://raw.githubusercontent.com/lucasmodrich/git-worktree-manager/refs/heads/main/%s", file)
		dest := filepath.Join(installDir, file)
		if err := downloadFile(url, dest); err != nil {
			// Non-fatal - continue
			fmt.Printf("Warning: failed to download %s: %v\n", file, err)
		}
	}

	// Get current binary path
	currentBinary := config.GetBinaryPath()

	// Set executable permissions on new binary
	if err := os.Chmod(tempBinary, 0755); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// Replace binary atomically
	if err := os.Rename(tempBinary, currentBinary); err != nil {
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	return nil
}

// getBinaryName returns the binary name for the current platform
func getBinaryName() string {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	// Map Go arch names to release naming
	archName := arch
	if arch == "amd64" {
		archName = "x86_64"
	}

	// Capitalize OS name
	osNameCap := strings.Title(osName)

	binaryName := fmt.Sprintf("git-worktree-manager_%s_%s", osNameCap, archName)

	if osName == "windows" {
		binaryName += ".exe"
	}

	return binaryName
}

// downloadFile downloads a file from a URL to a local path
func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, url)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// verifyChecksum verifies the SHA256 checksum of a file
func verifyChecksum(filePath, checksumFile, binaryName string) error {
	// Calculate actual checksum
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}
	actualChecksum := fmt.Sprintf("%x", hash.Sum(nil))

	// Read expected checksum from file
	checksumData, err := os.ReadFile(checksumFile)
	if err != nil {
		return err
	}

	// Parse checksums.txt (format: "checksum  filename")
	lines := strings.Split(string(checksumData), "\n")
	var expectedChecksum string
	for _, line := range lines {
		if strings.Contains(line, binaryName) {
			parts := strings.Fields(line)
			if len(parts) >= 1 {
				expectedChecksum = parts[0]
				break
			}
		}
	}

	if expectedChecksum == "" {
		return fmt.Errorf("checksum not found for %s", binaryName)
	}

	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
	}

	return nil
}
