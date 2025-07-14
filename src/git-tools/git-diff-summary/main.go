package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

// FileInfo holds information about a file diff
type FileInfo struct {
	File       string `json:"file"`
	Author     string `json:"author"`
	Date       string `json:"date"`
	Deleted    bool   `json:"deleted,omitempty"`
	Status     string `json:"status,omitempty"`
	Insertions int    `json:"insertions,omitempty"`
	Deletions  int    `json:"deletions,omitempty"`
}

// Config holds command configuration
type Config struct {
	SrcBranch  string
	DstBranch  string
	FilePath   string
	UseFormat  bool
	ShowStat   bool
	OutputJSON bool
	MaxWorkers int
}

// GitError represents an error that occurred during git command execution
// This provides better type safety and error context
type GitError struct {
	Command string
	Err     error
}

func (e GitError) Error() string {
	return fmt.Sprintf("git command failed: %s: %v", e.Command, e.Err)
}

func (e GitError) Unwrap() error {
	return e.Err
}

const bashCompletionScript = `# Bash completion for git-diff-summary
# Git subcommand completion function
_git_diff_summary() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local prev="${COMP_WORDS[COMP_CWORD-1]}"
    
    # Determine which argument we're completing
    local git_cmd_idx=0
    for (( i=0; i < ${#COMP_WORDS[@]}; i++ )); do
        if [[ "${COMP_WORDS[i]}" == "git" ]]; then
            git_cmd_idx=$i
            break
        fi
    done
    
    local arg_idx=$((COMP_CWORD - git_cmd_idx - 2))
    
    # Handle options
    if [[ "$cur" == -* ]]; then
        local opts="--format --stat --json --help"
        COMPREPLY=($(compgen -W "$opts" -- "$cur"))
        return
    fi
    
    if [[ $arg_idx -eq 0 ]]; then
        # First argument: source branch completion
        local branches=$(git branch -r --format="%(refname:short)" 2>/dev/null | grep -v HEAD)
        COMPREPLY=($(compgen -W "$branches" -- "$cur"))
    elif [[ $arg_idx -eq 1 ]]; then
        # Second argument: destination branch completion
        local branches=$(git branch -r --format="%(refname:short)" 2>/dev/null | grep -v HEAD)
        COMPREPLY=($(compgen -W "$branches" -- "$cur"))
    else
        # Third argument and beyond: file path completion
        COMPREPLY=($(compgen -d -- "$cur"))
    fi
}

# Register completion for direct script usage
complete -F _git_diff_summary git-diff-summary

# Git subcommand completion (automatically detected by git-completion.bash)
# The function name _git_diff_summary is automatically recognized by git
`

func main() {
	var config Config
	config.MaxWorkers = 10 // Default number of goroutines

	var rootCmd = &cobra.Command{
		Use:   "gh-diff-summary <SRC_BRANCH> <DST_BRANCH> [file_path]",
		Short: "Show summary of files that differ between branches with their last commit info",
		Long: `git diff-summary - Show summary of files that differ between branches with their last commit info

USAGE:
    gh-diff-summary [OPTIONS] <SRC_BRANCH> <DST_BRANCH> [file_path]

ARGUMENTS:
    <SRC_BRANCH>    Source branch (e.g., origin/main)
    <DST_BRANCH>    Destination branch (e.g., origin/develop)
    [file_path]     Optional path to filter files (e.g., db/migrate, app/models)

OUTPUT FORMATS:
    Default:  filename: author (YYYY/MM/DD)
    Format:   filename                           author      (YYYY/MM/DD)
    Stat:     A  filename | 25 +++++++++++++++++++++++++
    JSON:     [{"file": "...", "author": "...", "date": "...", "status": "added", "insertions": 6, "deletions": 2}]`,
		Args: cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			config.SrcBranch = args[0]
			config.DstBranch = args[1]
			if len(args) > 2 {
				config.FilePath = args[2]
			}
			return runDiffSummary(config)
		},
	}

	rootCmd.Flags().BoolVarP(&config.OutputJSON, "json", "", false, "Output results in JSON format")
	rootCmd.Flags().BoolVarP(&config.ShowStat, "stat", "", false, "Show file change statistics with status (A/M/D/R/C)")
	rootCmd.Flags().BoolVarP(&config.UseFormat, "format", "f", false, "Format output with aligned columns")
	rootCmd.Flags().IntVarP(&config.MaxWorkers, "workers", "w", 10, "Number of concurrent workers for git operations")

	// Add completion subcommand
	completionCmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate bash completion script",
		Long:  "Generate bash completion script for git-diff-summary",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(bashCompletionScript)
		},
	}
	rootCmd.AddCommand(completionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runDiffSummary(config Config) error {
	branchRange := fmt.Sprintf("%s..%s", config.SrcBranch, config.DstBranch)

	if !config.OutputJSON {
		fmt.Printf("Getting files that differ between %s and %s in path: %s\n",
			config.SrcBranch, config.DstBranch, getPathDisplay(config.FilePath))
		fmt.Println()
	}

	// Get file differences
	files, err := getDiffFiles(branchRange, config.FilePath, config.ShowStat)
	if err != nil {
		return fmt.Errorf("failed to get diff files: %w", err)
	}

	if len(files) == 0 {
		if config.OutputJSON {
			fmt.Println("[]")
		} else {
			fmt.Printf("No files found between %s and %s in path: %s\n",
				config.SrcBranch, config.DstBranch, getPathDisplay(config.FilePath))
		}
		return nil
	}

	// Get commit info for all files concurrently
	fileInfos, err := getFileInfos(files, config.MaxWorkers)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Output results
	return outputResults(fileInfos, config)
}

func getPathDisplay(path string) string {
	if path == "" {
		return "(all files)"
	}
	return path
}

func getDiffFiles(branchRange, filePath string, showStat bool) ([]FileInfo, error) {
	var cmd *exec.Cmd
	var statFiles []FileInfo

	if showStat {
		// Get file status and stats
		if filePath != "" {
			cmd = exec.Command("git", "diff", "--name-status", branchRange, "--", filePath)
		} else {
			cmd = exec.Command("git", "diff", "--name-status", branchRange)
		}

		statusOutput, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("git diff --name-status failed: %w", err)
		}

		// Get numeric stats
		if filePath != "" {
			cmd = exec.Command("git", "diff", "--numstat", branchRange, "--", filePath)
		} else {
			cmd = exec.Command("git", "diff", "--numstat", branchRange)
		}

		numstatOutput, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("git diff --numstat failed: %w", err)
		}

		statFiles = parseStatOutput(string(statusOutput), string(numstatOutput))
	} else {
		// Get simple file list
		if filePath != "" {
			cmd = exec.Command("git", "diff", "--name-only", branchRange, "--", filePath)
		} else {
			cmd = exec.Command("git", "diff", "--name-only", branchRange)
		}

		output, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("git diff --name-only failed: %w", err)
		}

		statFiles = parseNameOnlyOutput(string(output))
	}

	return statFiles, nil
}

func parseStatOutput(statusOutput, numstatOutput string) []FileInfo {
	var files []FileInfo
	statMap := make(map[string]FileInfo)

	// Parse numstat output for insertions/deletions
	scanner := bufio.NewScanner(strings.NewReader(numstatOutput))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) >= 3 {
			filename := parts[2]
			insertions := parseIntWithDefault(parts[0], 0)
			deletions := parseIntWithDefault(parts[1], 0)

			statMap[filename] = FileInfo{
				File:       filename,
				Insertions: insertions,
				Deletions:  deletions,
			}
		}
	}

	// Parse status output
	scanner = bufio.NewScanner(strings.NewReader(statusOutput))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) >= 2 {
			status := parts[0]
			filename := parts[1]

			fileInfo := FileInfo{
				File:   filename,
				Status: string(status[0]), // First character of status
			}

			// Handle renames/copies
			if strings.HasPrefix(status, "R") || strings.HasPrefix(status, "C") {
				if len(parts) >= 3 {
					fileInfo.File = fmt.Sprintf("%s => %s", parts[1], parts[2])
					filename = parts[2] // Use new filename for stats lookup
				}
			}

			// Set deleted flag
			if strings.HasPrefix(status, "D") {
				fileInfo.Deleted = true
			}

			// Copy stats if available
			if stat, exists := statMap[filename]; exists {
				fileInfo.Insertions = stat.Insertions
				fileInfo.Deletions = stat.Deletions
			}

			files = append(files, fileInfo)
		}
	}

	return files
}

func parseNameOnlyOutput(output string) []FileInfo {
	var files []FileInfo

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Check if file exists to determine if it was deleted
		deleted := false
		cmd := exec.Command("git", "ls-files", "--error-unmatch", line)
		if cmd.Run() != nil {
			deleted = true
		}

		status := "M"
		if deleted {
			status = "D"
		}

		files = append(files, FileInfo{
			File:    line,
			Status:  status,
			Deleted: deleted,
		})
	}

	return files
}

func getFileInfos(files []FileInfo, maxWorkers int) ([]FileInfo, error) {
	if len(files) == 0 {
		return files, nil
	}

	// Create work channels
	jobs := make(chan int, len(files))
	results := make(chan FileInfo, len(files))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				fileInfo := files[j]
				author, date := getCommitInfoWithBranch(fileInfo.File, fileInfo.Deleted)
				fileInfo.Author = author
				fileInfo.Date = date
				results <- fileInfo
			}
		}()
	}

	// Send jobs
	for i := range files {
		jobs <- i
	}
	close(jobs)

	// Wait for completion
	wg.Wait()
	close(results)

	// Collect results
	var fileInfos []FileInfo
	for result := range results {
		fileInfos = append(fileInfos, result)
	}

	// Sort by filename to maintain consistent order
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].File < fileInfos[j].File
	})

	return fileInfos, nil
}

func getCommitInfoWithBranch(filename string, deleted bool) (string, string) {
	// Try multiple approaches to get commit info
	approaches := [][]string{}

	if deleted {
		// For deleted files, try various approaches
		approaches = append(approaches, []string{
			"git", "log", "-1", "--format=%an|%cd", "--date=format:%Y/%m/%d", "--diff-filter=D", "--follow", "--", filename,
		})
		approaches = append(approaches, []string{
			"git", "log", "-1", "--format=%an|%cd", "--date=format:%Y/%m/%d", "--all", "--", filename,
		})
		approaches = append(approaches, []string{
			"git", "log", "-1", "--format=%an|%cd", "--date=format:%Y/%m/%d", "--", filename,
		})
	}

	// Always try these approaches (for existing files or as fallback)
	approaches = append(approaches, []string{
		"git", "log", "-1", "--format=%an|%cd", "--date=format:%Y/%m/%d", "--", filename,
	})
	approaches = append(approaches, []string{
		"git", "log", "-1", "--format=%an|%cd", "--date=format:%Y/%m/%d", filename,
	})
	approaches = append(approaches, []string{
		"git", "log", "-1", "--format=%an|%cd", "--date=format:%Y/%m/%d", "--follow", "--", filename,
	})
	approaches = append(approaches, []string{
		"git", "log", "-1", "--format=%an|%cd", "--date=format:%Y/%m/%d", "--all", "--", filename,
	})

	// Try each approach
	for _, args := range approaches {
		cmd := exec.Command(args[0], args[1:]...)
		if output, err := cmd.Output(); err == nil && len(output) > 0 {
			if line := strings.TrimSpace(string(output)); line != "" {
				if parts := strings.Split(line, "|"); len(parts) >= 2 {
					return parts[0], parts[1]
				}
			}
		}
	}

	return "(unknown)", "(unknown)"
}

func outputResults(fileInfos []FileInfo, config Config) error {
	if config.OutputJSON {
		return outputJSON(fileInfos)
	}

	if config.ShowStat {
		return outputStat(fileInfos, config.UseFormat)
	}

	return outputRegular(fileInfos, config.UseFormat)
}

func outputJSON(fileInfos []FileInfo) error {
	data, err := json.MarshalIndent(fileInfos, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(data))
	return nil
}

func outputStat(fileInfos []FileInfo, useFormat bool) error {
	var maxFileLen, maxAuthorLen int
	var totalInsertions, totalDeletions, totalFiles int

	// Calculate max lengths and totals
	for _, info := range fileInfos {
		if len(info.File) > maxFileLen {
			maxFileLen = len(info.File)
		}
		if len(info.Author) > maxAuthorLen {
			maxAuthorLen = len(info.Author)
		}
		totalInsertions += info.Insertions
		totalDeletions += info.Deletions
		totalFiles++
	}

	// Output each file
	for _, info := range fileInfos {
		statBar := generateStatBar(info.Insertions, info.Deletions)
		changes := info.Insertions + info.Deletions

		if useFormat {
			fmt.Printf("%s  %-*s | %3d %s  %-*s  (%s)\n",
				info.Status, maxFileLen, info.File, changes, statBar,
				maxAuthorLen, info.Author, info.Date)
		} else {
			fmt.Printf("%s  %-*s | %3d %s\n",
				info.Status, maxFileLen, info.File, changes, statBar)
		}
	}

	// Show summary
	fmt.Printf(" %d files changed, %d insertions(+), %d deletions(-)\n",
		totalFiles, totalInsertions, totalDeletions)

	return nil
}

func outputRegular(fileInfos []FileInfo, useFormat bool) error {
	var maxFileLen, maxAuthorLen int

	if useFormat {
		// Calculate max lengths for formatting
		for _, info := range fileInfos {
			if len(info.File) > maxFileLen {
				maxFileLen = len(info.File)
			}
			if len(info.Author) > maxAuthorLen {
				maxAuthorLen = len(info.Author)
			}
		}
	}

	// Output each file
	for _, info := range fileInfos {
		if useFormat {
			if info.Deleted {
				fmt.Printf("%-*s  %-*s  (%s) [DELETED]\n",
					maxFileLen, info.File, maxAuthorLen, info.Author, info.Date)
			} else {
				fmt.Printf("%-*s  %-*s  (%s)\n",
					maxFileLen, info.File, maxAuthorLen, info.Author, info.Date)
			}
		} else {
			if info.Deleted {
				fmt.Printf("%s: [DELETED] %s (%s)\n", info.File, info.Author, info.Date)
			} else {
				fmt.Printf("%s: %s (%s)\n", info.File, info.Author, info.Date)
			}
		}
	}

	return nil
}

func generateStatBar(insertions, deletions int) string {
	const maxWidth = 50
	total := insertions + deletions

	if total == 0 {
		return ""
	}

	plusChars := (insertions * maxWidth) / total
	minusChars := (deletions * maxWidth) / total

	// Ensure at least one character if there are changes
	if insertions > 0 && plusChars == 0 {
		plusChars = 1
	}
	if deletions > 0 && minusChars == 0 {
		minusChars = 1
	}

	plus := strings.Repeat("+", plusChars)
	minus := strings.Repeat("-", minusChars)

	return plus + minus
}

// checkDeletedFilesBatch checks which files are deleted in a batch operation
// This is a performance optimization over checking files individually
func checkDeletedFilesBatch(files []string) map[string]bool {
	if len(files) == 0 {
		return make(map[string]bool)
	}
	
	args := append([]string{"ls-files", "--error-unmatch"}, files...)
	cmd := exec.Command("git", args...)
	output, _ := cmd.Output()
	
	existing := make(map[string]bool)
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		existing[scanner.Text()] = true
	}
	
	deleted := make(map[string]bool)
	for _, file := range files {
		deleted[file] = !existing[file]
	}
	return deleted
}

// parseIntWithDefault parses a string to int with a default value for errors
// This provides better error handling than ignoring strconv.Atoi errors
func parseIntWithDefault(s string, defaultValue int) int {
	if s == "-" {
		return defaultValue
	}
	if val, err := strconv.Atoi(s); err == nil {
		return val
	}
	return defaultValue
}

// optimizedGetFileInfos is a memory-optimized version of getFileInfos
// It pre-allocates slice capacity to reduce memory allocations
func optimizedGetFileInfos(files []FileInfo, maxWorkers int) ([]FileInfo, error) {
	if len(files) == 0 {
		return files, nil
	}

	// Create work channels
	jobs := make(chan int, len(files))
	results := make(chan FileInfo, len(files))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				fileInfo := files[j]
				author, date := getCommitInfoWithBranch(fileInfo.File, fileInfo.Deleted)
				fileInfo.Author = author
				fileInfo.Date = date
				results <- fileInfo
			}
		}()
	}

	// Send jobs
	for i := range files {
		jobs <- i
	}
	close(jobs)

	// Wait for completion
	wg.Wait()
	close(results)

	// Collect results with pre-allocated capacity
	fileInfos := make([]FileInfo, 0, len(files))
	for result := range results {
		fileInfos = append(fileInfos, result)
	}

	// Sort by filename to maintain consistent order
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].File < fileInfos[j].File
	})

	return fileInfos, nil
}
