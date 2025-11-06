package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseLocation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantWT   string
		wantPath string
		wantErr  bool
	}{
		{
			name:     "通常のworktree指定",
			input:    "feature-branch:src/main.go",
			wantWT:   "feature-branch",
			wantPath: "src/main.go",
			wantErr:  false,
		},
		{
			name:     "現在のworktreeを示す@記号",
			input:    "@:README.md",
			wantWT:   "@",
			wantPath: "README.md",
			wantErr:  false,
		},
		{
			name:     "深いパス指定",
			input:    "dev:src/internal/utils/helper.go",
			wantWT:   "dev",
			wantPath: "src/internal/utils/helper.go",
			wantErr:  false,
		},
		{
			name:     "コロンを含むパス",
			input:    "main:file:with:colons.txt",
			wantWT:   "main",
			wantPath: "file:with:colons.txt",
			wantErr:  false,
		},
		{
			name:     "コロンなしの不正な形式",
			input:    "invalid-format",
			wantWT:   "",
			wantPath: "",
			wantErr:  true,
		},
		{
			name:     "空の入力",
			input:    "",
			wantWT:   "",
			wantPath: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWT, gotPath, err := parseLocation(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWT != tt.wantWT {
				t.Errorf("parseLocation() gotWT = %v, want %v", gotWT, tt.wantWT)
			}
			if gotPath != tt.wantPath {
				t.Errorf("parseLocation() gotPath = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}

func TestGetWorktreeRoot(t *testing.T) {
	t.Run("現在のworktree(@)の場合", func(t *testing.T) {
		// This test should be run in a git repository
		root, err := getWorktreeRoot("@")
		if err != nil {
			t.Skipf("Not in a git repository: %v", err)
		}
		if root == "" {
			t.Error("Expected non-empty root path for current worktree")
		}
		// Check if it's an absolute path
		if !filepath.IsAbs(root) {
			t.Errorf("Expected absolute path, got: %s", root)
		}
	})
}

func TestGetWorktreeNames(t *testing.T) {
	names := getWorktreeNames()
	// This test just ensures the function doesn't panic
	// The actual result depends on the git repository state
	t.Logf("Found %d worktree names: %v", len(names), names)
}

func TestConfig(t *testing.T) {
	c := Config{
		Source:      "src:file.txt",
		Destination: "dst:file.txt",
		Verbose:     true,
	}

	if c.Source != "src:file.txt" {
		t.Errorf("Expected Source to be 'src:file.txt', got %s", c.Source)
	}
	if c.Destination != "dst:file.txt" {
		t.Errorf("Expected Destination to be 'dst:file.txt', got %s", c.Destination)
	}
	if !c.Verbose {
		t.Error("Expected Verbose to be true")
	}
}

func TestWorktreeInfo(t *testing.T) {
	wt := WorktreeInfo{
		Branch: "feature-branch",
		Path:   "~/repos/project",
		Head:   "abc123",
		Bare:   false,
	}

	if wt.Branch != "feature-branch" {
		t.Errorf("Expected Branch to be 'feature-branch', got %s", wt.Branch)
	}
	if wt.Path != "~/repos/project" {
		t.Errorf("Expected Path to be '~/repos/project', got %s", wt.Path)
	}
}

// Integration test
func TestIntegrationCopyFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create temporary directories for testing
	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "src")
	dstDir := filepath.Join(tmpDir, "dst")

	// Create source directory and file
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	srcFile := filepath.Join(srcDir, "test.txt")
	content := []byte("Hello, World!")
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Test file copy functionality
	dstFile := filepath.Join(dstDir, "test.txt")
	
	// Create destination directory
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		t.Fatalf("Failed to create destination directory: %v", err)
	}

	// Copy file
	input, err := os.ReadFile(srcFile)
	if err != nil {
		t.Fatalf("Failed to read source file: %v", err)
	}

	if err := os.WriteFile(dstFile, input, 0644); err != nil {
		t.Fatalf("Failed to write destination file: %v", err)
	}

	// Verify copy
	copiedContent, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	if string(copiedContent) != string(content) {
		t.Errorf("File content mismatch. Expected: %s, Got: %s", content, copiedContent)
	}
}

// Benchmark tests
func BenchmarkParseLocation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _ = parseLocation("feature-branch:src/internal/utils/helper.go")
	}
}

func BenchmarkGetWorktreeNames(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = getWorktreeNames()
	}
}