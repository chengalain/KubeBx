package cluster

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	kindVersion    = "v0.20.0"
	kubectlVersion = "v1.29.0"
)

// GetKubeBxDir returns the KubeBx home directory
func GetKubeBxDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	kubebxDir := filepath.Join(home, ".kubebx")
	binDir := filepath.Join(kubebxDir, "bin")

	// Create directories if they don't exist
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create kubebx directory: %w", err)
	}

	return kubebxDir, nil
}

// CheckAndInstallDependencies checks and installs Kind and kubectl if needed
func CheckAndInstallDependencies() error {
	// Check Docker (mandatory)
	fmt.Println("Checking Docker...")
	if err := checkCommand("docker", "--version"); err != nil {
		return fmt.Errorf("docker is not installed or not running.\n\nDocker is required to run KubeBx. Please install it from:\n  https://docs.docker.com/get-docker/")
	}

	if err := checkCommand("docker", "info"); err != nil {
		return fmt.Errorf("docker daemon is not running. Please start Docker")
	}
	fmt.Println("✓ Docker is running")

	// Check/install Kind
	fmt.Println("\nChecking Kind...")
	if err := ensureKind(); err != nil {
		return err
	}
	fmt.Println("✓ Kind is ready")

	// Check/install kubectl
	fmt.Println("\nChecking kubectl...")
	if err := ensureKubectl(); err != nil {
		return err
	}
	fmt.Println("✓ kubectl is ready")

	return nil
}

// ensureKind ensures Kind is available (system or downloaded)
func ensureKind() error {
	// Check if Kind is in PATH
	if err := checkCommand("kind", "version"); err == nil {
		return nil
	}

	// Check if we have it in ~/.kubebx/bin
	kubebxDir, err := GetKubeBxDir()
	if err != nil {
		return err
	}

	kindPath := filepath.Join(kubebxDir, "bin", "kind")
	if runtime.GOOS == "windows" {
		kindPath += ".exe"
	}

	if err := checkCommandPath(kindPath, "version"); err == nil {
		return nil
	}

	// Download Kind
	fmt.Println("  Downloading Kind...")
	if err := downloadKind(kindPath); err != nil {
		return fmt.Errorf("failed to download Kind: %w", err)
	}

	return nil
}

// ensureKubectl ensures kubectl is available (system or downloaded)
func ensureKubectl() error {
	// Check if kubectl is in PATH
	if err := checkCommand("kubectl", "version", "--client"); err == nil {
		return nil
	}

	// Check if we have it in ~/.kubebx/bin
	kubebxDir, err := GetKubeBxDir()
	if err != nil {
		return err
	}

	kubectlPath := filepath.Join(kubebxDir, "bin", "kubectl")
	if runtime.GOOS == "windows" {
		kubectlPath += ".exe"
	}

	if err := checkCommandPath(kubectlPath, "version", "--client"); err == nil {
		return nil
	}

	// Download kubectl
	fmt.Println("  Downloading kubectl...")
	if err := downloadKubectl(kubectlPath); err != nil {
		return fmt.Errorf("failed to download kubectl: %w", err)
	}

	return nil
}

// downloadKind downloads Kind binary for the current platform
func downloadKind(destPath string) error {
	url := getKindDownloadURL()
	return downloadBinary(url, destPath)
}

// downloadKubectl downloads kubectl binary for the current platform
func downloadKubectl(destPath string) error {
	url := getKubectlDownloadURL()
	return downloadBinary(url, destPath)
}

// getKindDownloadURL returns the download URL for Kind based on OS and arch
func getKindDownloadURL() string {
	baseURL := "https://kind.sigs.k8s.io/dl/" + kindVersion + "/kind-"

	switch runtime.GOOS {
	case "linux":
		return baseURL + "linux-" + runtime.GOARCH
	case "darwin":
		return baseURL + "darwin-" + runtime.GOARCH
	case "windows":
		return baseURL + "windows-" + runtime.GOARCH
	default:
		return baseURL + "linux-amd64"
	}
}

// getKubectlDownloadURL returns the download URL for kubectl based on OS and arch
func getKubectlDownloadURL() string {
	baseURL := "https://dl.k8s.io/release/" + kubectlVersion + "/bin/"

	os := runtime.GOOS
	arch := runtime.GOARCH

	url := baseURL + os + "/" + arch + "/kubectl"
	if runtime.GOOS == "windows" {
		url += ".exe"
	}

	return url
}

// downloadBinary downloads a binary from a URL and makes it executable
func downloadBinary(url, destPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Make executable
	if err := os.Chmod(destPath, 0755); err != nil {
		return fmt.Errorf("failed to make executable: %w", err)
	}

	return nil
}

// GetKindCommand returns the path to the kind binary (system or local)
func GetKindCommand() string {
	// Try system Kind first
	if err := checkCommand("kind", "version"); err == nil {
		return "kind"
	}

	// Use local Kind
	kubebxDir, err := GetKubeBxDir()
	if err != nil {
		return "kind" // Fallback
	}

	kindPath := filepath.Join(kubebxDir, "bin", "kind")
	if runtime.GOOS == "windows" {
		kindPath += ".exe"
	}

	return kindPath
}

// GetKubectlCommand returns the path to the kubectl binary (system or local)
func GetKubectlCommand() string {
	// Try system kubectl first
	if err := checkCommand("kubectl", "version", "--client"); err == nil {
		return "kubectl"
	}

	// Use local kubectl
	kubebxDir, err := GetKubeBxDir()
	if err != nil {
		return "kubectl" // Fallback
	}

	kubectlPath := filepath.Join(kubebxDir, "bin", "kubectl")
	if runtime.GOOS == "windows" {
		kubectlPath += ".exe"
	}

	return kubectlPath
}

func checkCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func checkCommandPath(path string, args ...string) error {
	cmd := exec.Command(path, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
