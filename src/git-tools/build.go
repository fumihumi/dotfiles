package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	colorGreen  = "\033[0;32m"
	colorRed    = "\033[0;31m"
	colorYellow = "\033[1;33m"
	colorReset  = "\033[0m"
)

func main() {
	// Default to build if no arguments
	if len(os.Args) < 2 {
		buildAll()
		return
	}

	command := os.Args[1]
	switch command {
	case "build":
		buildAll()
	case "check":
		checkStatus()
	case "completion":
		if len(os.Args) < 3 {
			logError("Usage: ./build.sh completion <install|show>")
			os.Exit(1)
		}
		handleCompletion(os.Args[2])
	case "clean":
		cleanAll()
	case "help", "-h", "--help":
		printUsage()
	default:
		logError("Unknown command: %s", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Git Tools Build System")
	fmt.Println()
	fmt.Println("Usage: ./build.sh [command]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  build              Build all git tools (default)")
	fmt.Println("  check              Check build status and show what needs rebuilding")
	fmt.Println("  completion show    Show completion setup instructions")
	fmt.Println("  completion install Install completion to ~/.bashrc")
	fmt.Println("  clean              Remove built binaries")
	fmt.Println("  help               Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  ./build.sh                    # Build all tools")
	fmt.Println("  ./build.sh check              # Show what needs to be built")
	fmt.Println("  ./build.sh completion install # Install bash completion")
	fmt.Println("  ./build.sh clean              # Remove all built binaries")
}

func buildAll() {
	// Get the actual script directory, not the temp build directory
	scriptDir := getScriptDir()
	projectRoot := filepath.Join(scriptDir, "../..")
	binDir := filepath.Join(projectRoot, ".bash/bin")
	completionDir := filepath.Join(projectRoot, ".bash/completion.d")

	// Create directories
	os.MkdirAll(binDir, 0755)
	os.MkdirAll(completionDir, 0755)

	// Scan all subdirectories
	entries, err := os.ReadDir(scriptDir)
	if err != nil {
		logError("Failed to read directory: %v", err)
		os.Exit(1)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		toolName := entry.Name()
		toolDir := filepath.Join(scriptDir, toolName)

		// Check if it's a Go project
		if _, err := os.Stat(filepath.Join(toolDir, "go.mod")); err == nil {
			buildGoProject(toolName, toolDir, binDir)
			generateCompletionLoader(toolName, completionDir, "go")
		} else if _, err := os.Stat(filepath.Join(toolDir, toolName+".sh")); err == nil {
			// Check if it's a Bash script
			installBashScript(toolName, toolDir, binDir)
			if _, err := os.Stat(filepath.Join(toolDir, "completion.bash")); err == nil {
				generateCompletionLoader(toolName, completionDir, "bash")
			}
		}
	}

	logInfo("Build completed successfully!")
	logInfo("Binaries installed to: %s", binDir)
	logInfo("Completion scripts in: %s", completionDir)
}

func buildGoProject(name, srcDir, binDir string) {
	logInfo("Building Go project: %s", name)
	
	outputPath := filepath.Join(binDir, name)
	cmd := exec.Command("go", "build", "-o", outputPath, ".")
	cmd.Dir = srcDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		logError("Failed to build %s: %v", name, err)
		return
	}
	
	logInfo("Successfully built %s", name)
}

func installBashScript(name, srcDir, binDir string) {
	logInfo("Installing bash script: %s", name)
	
	srcPath := filepath.Join(srcDir, name+".sh")
	dstPath := filepath.Join(binDir, name)
	
	if err := copyFile(srcPath, dstPath); err != nil {
		logError("Failed to copy %s: %v", name, err)
		return
	}
	
	// Make executable
	if err := os.Chmod(dstPath, 0755); err != nil {
		logError("Failed to chmod %s: %v", name, err)
		return
	}
	
	logInfo("Successfully installed %s", name)
}

func generateCompletionLoader(name, completionDir, toolType string) {
	completionFile := filepath.Join(completionDir, name+".rc")
	
	var content string
	switch toolType {
	case "go":
		// First, check if the tool uses Cobra's standard completion (with shell argument)
		// or legacy completion (without argument)
		content = fmt.Sprintf(`# Bash completion loader for %s
if command -v %s &> /dev/null; then
    # Try with bash argument first (Cobra standard), fallback to no argument (legacy)
    if %s completion bash &>/dev/null 2>&1; then
        eval "$(%s completion bash 2>/dev/null || true)"
    else
        eval "$(%s completion 2>/dev/null || true)"
    fi
fi
`, name, name, name, name, name)
	case "bash":
		content = fmt.Sprintf(`# Bash completion loader for %s
GIT_TOOLS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../src/git-tools" && pwd)"
if [[ -f "${GIT_TOOLS_DIR}/%s/completion.bash" ]]; then
    source "${GIT_TOOLS_DIR}/%s/completion.bash"
fi
`, name, name, name)
	}
	
	if err := os.WriteFile(completionFile, []byte(content), 0644); err != nil {
		logError("Failed to write completion loader for %s: %v", name, err)
	}
}

func handleCompletion(subcommand string) {
	scriptDir := getScriptDir()
	projectRoot := filepath.Join(scriptDir, "../..")
	completionDir := filepath.Join(projectRoot, ".bash/completion.d")
	
	completionSnippet := fmt.Sprintf(`
# Git tools completion
if [[ -d "%s" ]]; then
    for completion in "%s"/*.rc; do
        [[ -f "$completion" ]] && source "$completion"
    done
fi`, completionDir, completionDir)
	
	switch subcommand {
	case "show":
		fmt.Println("Add the following to your ~/.bashrc or ~/.bash_profile:")
		fmt.Println(completionSnippet)
	case "install":
		bashrcPath := filepath.Join(os.Getenv("HOME"), ".bashrc")
		
		// Check if already installed
		content, err := os.ReadFile(bashrcPath)
		if err == nil && strings.Contains(string(content), "Git tools completion") {
			logInfo("Completion already installed in %s", bashrcPath)
			return
		}
		
		// Append to .bashrc
		f, err := os.OpenFile(bashrcPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			logError("Failed to open %s: %v", bashrcPath, err)
			return
		}
		defer f.Close()
		
		if _, err := f.WriteString("\n" + completionSnippet + "\n"); err != nil {
			logError("Failed to write to %s: %v", bashrcPath, err)
			return
		}
		
		logInfo("Completion installed to %s", bashrcPath)
		logInfo("Run 'source ~/.bashrc' to activate")
	default:
		logError("Unknown completion subcommand: %s", subcommand)
	}
}

func cleanAll() {
	scriptDir := getScriptDir()
	projectRoot := filepath.Join(scriptDir, "../..")
	binDir := filepath.Join(projectRoot, ".bash/bin")
	
	// Get list of tools
	entries, err := os.ReadDir(scriptDir)
	if err != nil {
		logError("Failed to read directory: %v", err)
		return
	}
	
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		
		toolName := entry.Name()
		binPath := filepath.Join(binDir, toolName)
		
		if _, err := os.Stat(binPath); err == nil {
			if err := os.Remove(binPath); err != nil {
				logError("Failed to remove %s: %v", binPath, err)
			} else {
				logInfo("Removed %s", toolName)
			}
		}
	}
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	
	_, err = io.Copy(dstFile, srcFile)
	return err
}

func logInfo(format string, args ...interface{}) {
	fmt.Printf("%s[INFO]%s %s\n", colorGreen, colorReset, fmt.Sprintf(format, args...))
}

func logError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s[ERROR]%s %s\n", colorRed, colorReset, fmt.Sprintf(format, args...))
}

func logWarning(format string, args ...interface{}) {
	fmt.Printf("%s[WARN]%s %s\n", colorYellow, colorReset, fmt.Sprintf(format, args...))
}

func getScriptDir() string {
	// The build.sh script changes to the correct directory before running
	// So we can trust the current working directory
	wd, _ := os.Getwd()
	return wd
}

type ToolStatus struct {
	Name       string
	Type       string // "go" or "bash"
	Installed  bool
	NeedsRebuild bool
	Reason     string
}

func checkStatus() {
	scriptDir := getScriptDir()
	projectRoot := filepath.Join(scriptDir, "../..")
	binDir := filepath.Join(projectRoot, ".bash/bin")
	
	var allTools []ToolStatus
	var needsRebuild []ToolStatus
	var notInstalled []ToolStatus
	
	// Scan all subdirectories
	entries, err := os.ReadDir(scriptDir)
	if err != nil {
		logError("Failed to read directory: %v", err)
		return
	}
	
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		
		toolName := entry.Name()
		toolDir := filepath.Join(scriptDir, toolName)
		binPath := filepath.Join(binDir, toolName)
		
		status := ToolStatus{Name: toolName}
		
		// Check tool type and status
		if _, err := os.Stat(filepath.Join(toolDir, "go.mod")); err == nil {
			status.Type = "go"
			status.Installed = fileExists(binPath)
			
			if status.Installed {
				// Check if source is newer than binary
				srcFiles := []string{
					filepath.Join(toolDir, "main.go"),
					filepath.Join(toolDir, "go.mod"),
				}
				
				binInfo, _ := os.Stat(binPath)
				for _, srcFile := range srcFiles {
					if srcInfo, err := os.Stat(srcFile); err == nil {
						if srcInfo.ModTime().After(binInfo.ModTime()) {
							status.NeedsRebuild = true
							status.Reason = fmt.Sprintf("%s modified", filepath.Base(srcFile))
							break
						}
					}
				}
			} else {
				status.Reason = "not installed"
			}
		} else if _, err := os.Stat(filepath.Join(toolDir, toolName+".sh")); err == nil {
			status.Type = "bash"
			status.Installed = fileExists(binPath)
			
			if status.Installed {
				srcPath := filepath.Join(toolDir, toolName+".sh")
				if srcInfo, _ := os.Stat(srcPath); srcInfo != nil {
					if binInfo, _ := os.Stat(binPath); binInfo != nil {
						if srcInfo.ModTime().After(binInfo.ModTime()) {
							status.NeedsRebuild = true
							status.Reason = "script modified"
						}
					}
				}
			} else {
				status.Reason = "not installed"
			}
		} else {
			continue // Skip directories without recognized tool files
		}
		
		allTools = append(allTools, status)
		if !status.Installed {
			notInstalled = append(notInstalled, status)
		} else if status.NeedsRebuild {
			needsRebuild = append(needsRebuild, status)
		}
	}
	
	// Display results
	fmt.Println("=== Git Tools Status ===")
	fmt.Println()
	
	// Show all tools with their status
	fmt.Println("All tools:")
	for _, tool := range allTools {
		statusStr := "✓ up to date"
		if !tool.Installed {
			statusStr = "✗ " + tool.Reason
		} else if tool.NeedsRebuild {
			statusStr = "⚠ " + tool.Reason
		}
		fmt.Printf("  %-20s [%s] %s\n", tool.Name, tool.Type, statusStr)
	}
	
	fmt.Println()
	
	// Summary
	if len(notInstalled) == 0 && len(needsRebuild) == 0 {
		logInfo("All tools are up to date!")
	} else {
		if len(notInstalled) > 0 {
			fmt.Printf("Tools not installed: %d\n", len(notInstalled))
			for _, tool := range notInstalled {
				fmt.Printf("  - %s\n", tool.Name)
			}
		}
		
		if len(needsRebuild) > 0 {
			fmt.Printf("\nTools needing rebuild: %d\n", len(needsRebuild))
			for _, tool := range needsRebuild {
				fmt.Printf("  - %s (%s)\n", tool.Name, tool.Reason)
			}
		}
		
		fmt.Println()
		logWarning("Run './build.sh' to update all tools")
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}