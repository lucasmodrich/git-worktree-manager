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
	"time"

	"github.com/lucasmodrich/git-worktree-manager/internal/config"
)

// UpgradeToLatest downloads and installs the latest version, printing progress as it goes.
func UpgradeToLatest(currentVersion, latestVersion string) error {
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

	binaryName := getBinaryName()

	// Download binary
	binaryURL := fmt.Sprintf("https://github.com/lucasmodrich/git-worktree-manager/releases/download/v%s/%s", latestVersion, binaryName)
	tempBinary := filepath.Join(os.TempDir(), "gwtm-new")
	if err := downloadFile(binaryURL, tempBinary); err != nil {
		return fmt.Errorf("failed to download binary: %w", err)
	}
	defer os.Remove(tempBinary)
	fmt.Println("✓ Binary downloaded")

	// Download and verify checksum
	checksumURL := fmt.Sprintf("https://github.com/lucasmodrich/git-worktree-manager/releases/download/v%s/checksums.txt", latestVersion)
	checksumFile := filepath.Join(os.TempDir(), "gwtm-checksums.txt")
	if err := downloadFile(checksumURL, checksumFile); err != nil {
		return fmt.Errorf("failed to download checksums: %w", err)
	}
	defer os.Remove(checksumFile)

	if err := verifyChecksum(tempBinary, checksumFile, binaryName); err != nil {
		return fmt.Errorf("checksum verification failed: %w", err)
	}
	fmt.Println("✓ Checksum verified")

	// Ensure install directory exists
	installDir := config.GetInstallDir()
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	// Download auxiliary files from the release tag (not main branch)
	for _, file := range []string{"README.md", "VERSION", "LICENSE"} {
		url := fmt.Sprintf("https://raw.githubusercontent.com/lucasmodrich/git-worktree-manager/refs/tags/v%s/%s", latestVersion, file)
		dest := filepath.Join(installDir, file)
		if err := downloadFile(url, dest); err != nil {
			fmt.Printf("Warning: failed to download %s: %v\n", file, err)
		} else {
			fmt.Printf("✓ %s downloaded\n", file)
		}
	}

	// Set executable permissions
	if err := os.Chmod(tempBinary, 0755); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// Install binary (handles cross-filesystem moves with copy fallback)
	if err := moveBinary(tempBinary, config.GetBinaryPath()); err != nil {
		return fmt.Errorf("failed to install binary: %w", err)
	}

	return nil
}

// getBinaryName returns the release asset filename for the current platform,
// matching the naming convention produced by GoReleaser.
func getBinaryName() string {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	archName := arch
	if arch == "amd64" {
		archName = "x86_64"
	}

	// Capitalise first letter to match GoReleaser title-case (e.g. "linux" -> "Linux")
	osNameCap := strings.ToUpper(osName[:1]) + osName[1:]

	binaryName := fmt.Sprintf("gwtm_%s_%s", osNameCap, archName)
	if osName == "windows" {
		binaryName += ".exe"
	}
	return binaryName
}

// moveBinary moves src to dst. Falls back to copy+delete when rename fails
// across filesystems (e.g. /tmp on a different device than the install dir).
func moveBinary(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source binary: %w", err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination binary: %w", err)
	}

	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		os.Remove(dst)
		return fmt.Errorf("failed to copy binary: %w", err)
	}

	if err := out.Close(); err != nil {
		os.Remove(dst)
		return fmt.Errorf("failed to finalise binary: %w", err)
	}

	// Restore executable bit lost by os.Create (no-op on Windows where
	// executability is determined by file extension, not permission bits).
	if err := os.Chmod(dst, 0755); err != nil {
		os.Remove(dst)
		return fmt.Errorf("failed to set binary permissions: %w", err)
	}

	os.Remove(src) // best-effort cleanup of temp file
	return nil
}

// downloadFile downloads url to the local path with a 60-second timeout.
func downloadFile(url, path string) error {
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, url)
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// verifyChecksum verifies the SHA256 checksum of filePath against checksums.txt.
func verifyChecksum(filePath, checksumFile, binaryName string) error {
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

	checksumData, err := os.ReadFile(checksumFile)
	if err != nil {
		return err
	}

	// Format: "<checksum>  <filename>"
	var expectedChecksum string
	for _, line := range strings.Split(string(checksumData), "\n") {
		if strings.Contains(line, binaryName) {
			if parts := strings.Fields(line); len(parts) >= 1 {
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
