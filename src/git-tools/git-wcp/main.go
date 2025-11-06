package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

type Config struct {
	Source      string
	Destination string
	Verbose     bool
}

type WorktreeInfo struct {
	Branch string `json:"branch"`
	Path   string `json:"path"`
	Head   string `json:"head"`
	Bare   bool   `json:"bare"`
}

var config Config

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "git-wcp <source> <destination> [path]",
	Short: "Copy files between git worktrees",
	Long: `Copy files between git worktrees.

Format:
  <worktree>:<path>  - File in specified worktree
  @:<path>           - File in current worktree
  <worktree>         - Worktree name only (when used with 3 args)

Usage patterns:
  1. git wcp <worktreeA>:<pathA> <worktreeB>           - Copy to same path
  2. git wcp <worktreeA> <worktreeB> <path>            - Copy path between worktrees
  3. git wcp <worktreeA>:<pathA> <worktreeB>:<pathB>   - Copy to different path

Examples:
  git wcp master:README.md feature-branch               - Copy README.md to same location
  git wcp master feature-branch README.md               - Copy README.md between worktrees
  git wcp master:README.md feature-branch:docs/README.md - Copy to different location`,
	Args: cobra.RangeArgs(2, 3),
	RunE: runCopy,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		switch len(args) {
		case 0:
			// First argument: can be worktree or worktree:path
			return getLocationCompletions(toComplete), cobra.ShellCompDirectiveNoSpace
		case 1:
			// Second argument: can be worktree or worktree:path
			// If first arg has no colon, only complete worktree names
			if !strings.Contains(args[0], ":") {
				return getWorktreeCompletions(toComplete), cobra.ShellCompDirectiveNoSpace
			}
			return getLocationCompletions(toComplete), cobra.ShellCompDirectiveNoSpace
		case 2:
			// Third argument: only if first two are worktree names without paths
			if !strings.Contains(args[0], ":") && !strings.Contains(args[1], ":") {
				// Complete file paths
				return nil, cobra.ShellCompDirectiveDefault
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		default:
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
	},
}

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:
  $ source <(git-wcp completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ git-wcp completion bash > /etc/bash_completion.d/git-wcp

  # macOS:
  $ git-wcp completion bash > $(brew --prefix)/etc/bash_completion.d/git-wcp

Zsh:
  $ source <(git-wcp completion zsh)

  # To load completions for each session, execute once:
  $ git-wcp completion zsh > "${fpath[1]}/_git-wcp"
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			// Generate both standard cobra completion and git subcommand completion
			cmd.Root().GenBashCompletion(os.Stdout)
			// Add git subcommand completion
			fmt.Println("\n# Git subcommand completion")
			fmt.Println("# This allows 'git wcp' to work with tab completion")
			fmt.Println("_git_wcp() {")
			fmt.Println("    local cur=\"${COMP_WORDS[COMP_CWORD]}\"")
			fmt.Println("    local prev=\"${COMP_WORDS[COMP_CWORD-1]}\"")
			fmt.Println("    ")
			fmt.Println("    # Find the position of 'git' in the command")
			fmt.Println("    local git_idx=0")
			fmt.Println("    for (( i=0; i < ${#COMP_WORDS[@]}; i++ )); do")
			fmt.Println("        if [[ \"${COMP_WORDS[i]}\" == \"git\" ]]; then")
			fmt.Println("            git_idx=$i")
			fmt.Println("            break")
			fmt.Println("        fi")
			fmt.Println("    done")
			fmt.Println("    ")
			fmt.Println("    # Adjust COMP_WORDS to work with __start_git-wcp")
			fmt.Println("    local new_comp_words=(\"git-wcp\")")
			fmt.Println("    for (( i=$((git_idx + 2)); i < ${#COMP_WORDS[@]}; i++ )); do")
			fmt.Println("        new_comp_words+=(\"${COMP_WORDS[i]}\")")
			fmt.Println("    done")
			fmt.Println("    ")
			fmt.Println("    # Set up environment for __start_git-wcp")
			fmt.Println("    COMP_WORDS=(\"${new_comp_words[@]}\")")
			fmt.Println("    COMP_CWORD=$((${#new_comp_words[@]} - 1))")
			fmt.Println("    COMP_LINE=\"${new_comp_words[*]}\"")
			fmt.Println("    COMP_POINT=${#COMP_LINE}")
			fmt.Println("    ")
			fmt.Println("    __start_git-wcp")
			fmt.Println("}")
			fmt.Println("")
			fmt.Println("# Note: File path completion after colon is not supported in bash completion")
			fmt.Println("# This is a limitation shared with docker cp and similar tools")
			fmt.Println("# Users need to type the full path after the colon")
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		}
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&config.Verbose, "verbose", "v", false, "Show verbose output")
	rootCmd.AddCommand(completionCmd)
}

func runCopy(cmd *cobra.Command, args []string) error {
	var srcWorktree, srcPath, dstWorktree, dstPath string
	var err error

	switch len(args) {
	case 2:
		// Pattern 1 or 3: <worktree>:<path> <worktree> or <worktree>:<path> <worktree>:<path>
		if strings.Contains(args[0], ":") {
			// Source has path
			srcWorktree, srcPath, err = parseLocation(args[0])
			if err != nil {
				return fmt.Errorf("invalid source format: %w", err)
			}

			if strings.Contains(args[1], ":") {
				// Pattern 3: Both have paths
				dstWorktree, dstPath, err = parseLocation(args[1])
				if err != nil {
					return fmt.Errorf("invalid destination format: %w", err)
				}
			} else {
				// Pattern 1: Destination is just worktree, use same path
				dstWorktree = args[1]
				dstPath = srcPath
			}
		} else {
			return fmt.Errorf("invalid format: when using 2 arguments, source must contain path (worktree:path)")
		}

	case 3:
		// Pattern 2: <worktree> <worktree> <path>
		if strings.Contains(args[0], ":") || strings.Contains(args[1], ":") {
			return fmt.Errorf("invalid format: when using 3 arguments, first two must be worktree names only")
		}
		srcWorktree = args[0]
		dstWorktree = args[1]
		srcPath = args[2]
		dstPath = args[2]

	default:
		return fmt.Errorf("unexpected number of arguments")
	}

	// Get worktree roots
	srcRoot, err := getWorktreeRoot(srcWorktree)
	if err != nil {
		return fmt.Errorf("source worktree '%s' not found: %w", srcWorktree, err)
	}

	dstRoot, err := getWorktreeRoot(dstWorktree)
	if err != nil {
		return fmt.Errorf("destination worktree '%s' not found: %w", dstWorktree, err)
	}

	// Build full paths
	srcFull := filepath.Join(srcRoot, srcPath)
	dstFull := filepath.Join(dstRoot, dstPath)

	// Check if source exists
	if _, err := os.Stat(srcFull); os.IsNotExist(err) {
		return fmt.Errorf("source file not found: %s", srcFull)
	}

	// Create destination directory if needed
	dstDir := filepath.Dir(dstFull)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Copy file
	input, err := os.ReadFile(srcFull)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	if err := os.WriteFile(dstFull, input, 0644); err != nil {
		return fmt.Errorf("failed to write destination file: %w", err)
	}

	// Print result
	var copyDescription string
	switch len(args) {
	case 2:
		if strings.Contains(args[1], ":") {
			copyDescription = fmt.Sprintf("%s -> %s", args[0], args[1])
		} else {
			copyDescription = fmt.Sprintf("%s -> %s:%s", args[0], args[1], srcPath)
		}
	case 3:
		copyDescription = fmt.Sprintf("%s:%s -> %s:%s", args[0], args[2], args[1], args[2])
	}
	
	fmt.Printf("Copied: %s\n", copyDescription)
	if config.Verbose {
		fmt.Printf("  From: %s\n", srcFull)
		fmt.Printf("  To:   %s\n", dstFull)
	}

	return nil
}

func parseLocation(location string) (worktree, path string, err error) {
	parts := strings.SplitN(location, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("expected format: <worktree>:<path>")
	}
	return parts[0], parts[1], nil
}

func getWorktreeCompletions(toComplete string) []string {
	var suggestions []string
	
	// Special entries first (@ for current worktree)
	if toComplete == "" || strings.HasPrefix("@", toComplete) {
		suggestions = append(suggestions, "@")
	}

	// Get worktree names (already sorted)
	worktrees := getWorktreeNames()
	for _, wt := range worktrees {
		if toComplete == "" || strings.HasPrefix(wt, toComplete) {
			suggestions = append(suggestions, wt)
		}
	}

	return suggestions
}

func getLocationCompletions(toComplete string) []string {
	if strings.Contains(toComplete, ":") {
		// Already has worktree part, complete file paths
		parts := strings.SplitN(toComplete, ":", 2)
		worktree := parts[0]
		pathPrefix := ""
		if len(parts) > 1 {
			pathPrefix = parts[1]
		}

		root, err := getWorktreeRoot(worktree)
		if err != nil {
			return nil
		}

		return getFileCompletions(root, pathPrefix, worktree)
	}

	// Complete worktree names without automatic colon
	var suggestions []string

	// Special entries first (@ for current worktree)
	if toComplete == "" || strings.HasPrefix("@", toComplete) {
		suggestions = append(suggestions, "@")
	}

	// Get worktree names (already sorted)
	worktrees := getWorktreeNames()
	for _, wt := range worktrees {
		if toComplete == "" || strings.HasPrefix(wt, toComplete) {
			suggestions = append(suggestions, wt)
		}
	}

	return suggestions
}

func getFileCompletions(root, prefix, worktree string) []string {
	var completions []string

	searchPath := filepath.Join(root, prefix)
	dir := searchPath
	basename := ""

	// If searchPath is not a directory, use its parent
	if info, err := os.Stat(searchPath); err != nil || !info.IsDir() {
		dir = filepath.Dir(searchPath)
		basename = filepath.Base(prefix)
		// Update prefix to be the directory part
		if prefix != "" && !strings.HasSuffix(prefix, "/") {
			newPrefix := filepath.Dir(prefix)
			if newPrefix == "." {
				prefix = ""
			} else {
				prefix = newPrefix
				if prefix != "" && !strings.HasSuffix(prefix, "/") {
					prefix += "/"
				}
			}
		}
	} else if prefix != "" && !strings.HasSuffix(prefix, "/") {
		// If it's a directory and doesn't end with /, add it
		prefix += "/"
	}

	// Read directory contents
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	// Build completions
	for _, entry := range entries {
		name := entry.Name()

		// Filter by basename if we have one
		if basename != "" && !strings.HasPrefix(name, basename) {
			continue
		}

		// Skip hidden files unless prefix ends with . or basename starts with .
		if strings.HasPrefix(name, ".") && !strings.HasSuffix(prefix, ".") && !strings.HasPrefix(basename, ".") {
			continue
		}

		completionPath := prefix + name
		if entry.IsDir() {
			completionPath += "/"
		}
		completions = append(completions, worktree+":"+completionPath)
	}

	return completions
}

func getWorktreeRoot(worktreeName string) (string, error) {
	if worktreeName == "@" {
		// Get current worktree root
		cmd := exec.Command("git", "rev-parse", "--show-toplevel")
		output, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf("not in a git repository")
		}
		return strings.TrimSpace(string(output)), nil
	}

	// Try to get worktree info from gwq
	worktrees, err := getWorktreesFromGwq()
	if err == nil {
		for _, wt := range worktrees {
			if wt.Branch == worktreeName {
				// Expand tilde in path
				path := wt.Path
				if strings.HasPrefix(path, "~/") {
					home, err := os.UserHomeDir()
					if err == nil {
						path = filepath.Join(home, path[2:])
					}
				}
				return path, nil
			}
		}
	}

	// Fallback to git worktree list
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to list worktrees: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	var currentPath string
	for _, line := range lines {
		if strings.HasPrefix(line, "worktree ") {
			currentPath = strings.TrimPrefix(line, "worktree ")
		} else if strings.HasPrefix(line, "branch refs/heads/") {
			branch := strings.TrimPrefix(line, "branch refs/heads/")
			if branch == worktreeName {
				return currentPath, nil
			}
		}
	}

	return "", fmt.Errorf("worktree not found")
}

func getWorktreeNames() []string {
	var names []string

	// Try gwq first
	worktrees, err := getWorktreesFromGwq()
	if err == nil {
		for _, wt := range worktrees {
			names = append(names, wt.Branch)
		}
		sort.Strings(names)
		return names
	}

	// Fallback to git worktree list
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return names
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "branch refs/heads/") {
			branch := strings.TrimPrefix(line, "branch refs/heads/")
			names = append(names, branch)
		}
	}

	sort.Strings(names)
	return names
}

func getWorktreesFromGwq() ([]WorktreeInfo, error) {
	// Check if gwq is available
	if _, err := exec.LookPath("gwq"); err != nil {
		return nil, fmt.Errorf("gwq not found")
	}

	cmd := exec.Command("gwq", "list", "--json")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var worktrees []WorktreeInfo
	if err := json.Unmarshal(output, &worktrees); err != nil {
		return nil, err
	}

	return worktrees, nil
}
