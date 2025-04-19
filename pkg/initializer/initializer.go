package initializer

import (
	"fmt"
	"os/exec"
	"runtime"
)

// InitEnvironment initializes the development environment
func InitEnvironment() error {
	// Check and install Node.js
	if err := checkAndInstallNode(); err != nil {
		return fmt.Errorf("failed to setup Node.js: %v", err)
	}

	// Check and install Wails
	if err := checkAndInstallWails(); err != nil {
		return fmt.Errorf("failed to setup Wails: %v", err)
	}

	return nil
}

func checkAndInstallNode() error {
	// Check if Node.js is installed
	_, err := exec.Command("node", "--version").Output()
	if err == nil {
		fmt.Println("✓ Node.js is already installed")
		return nil
	}

	fmt.Println("Installing Node.js...")

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		// For macOS, use Homebrew
		_, err := exec.Command("brew", "--version").Output()
		if err != nil {
			return fmt.Errorf("Homebrew is not installed. Please install Homebrew first: https://brew.sh")
		}
		cmd = exec.Command("brew", "install", "node")
	case "linux":
		// For Linux, use apt (Ubuntu/Debian) or other package managers as needed
		cmd = exec.Command("sudo", "apt", "install", "-y", "nodejs", "npm")
	case "windows":
		return fmt.Errorf("For Windows, please install Node.js manually from https://nodejs.org")
	default:
		return fmt.Errorf("unsupported operating system")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install Node.js: %v\n%s", err, output)
	}

	fmt.Println("✓ Node.js installed successfully")
	return nil
}

func checkAndInstallWails() error {
	// Check if Wails is installed
	_, err := exec.Command("wails", "version").Output()
	if err == nil {
		fmt.Println("✓ Wails is already installed")
		return nil
	}

	fmt.Println("Installing Wails...")
	cmd := exec.Command("go", "install", "github.com/wailsapp/wails/v2/cmd/wails@latest")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install Wails: %v\n%s", err, output)
	}

	fmt.Println("✓ Wails installed successfully")
	return nil
}
