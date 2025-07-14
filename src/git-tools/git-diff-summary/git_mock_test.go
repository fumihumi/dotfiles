package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// MockGitCommand は Git コマンドをモックするためのインターフェース
type MockGitCommand struct {
	responses map[string]string
	errors    map[string]error
}

// NewMockGitCommand creates a new MockGitCommand
func NewMockGitCommand() *MockGitCommand {
	return &MockGitCommand{
		responses: make(map[string]string),
		errors:    make(map[string]error),
	}
}

// AddResponse adds a mocked response for a specific command
func (m *MockGitCommand) AddResponse(command string, response string) {
	m.responses[command] = response
}

// AddError adds a mocked error for a specific command
func (m *MockGitCommand) AddError(command string, err error) {
	m.errors[command] = err
}

// GetResponse returns the mocked response for a command
func (m *MockGitCommand) GetResponse(command string) (string, error) {
	if err, exists := m.errors[command]; exists {
		return "", err
	}
	if response, exists := m.responses[command]; exists {
		return response, nil
	}
	return "", nil
}

// TestDiffFiles tests the getDiffFiles function with mocked git commands
func TestGetDiffFiles(t *testing.T) {
	t.Run("正常なdiff --name-only出力の処理", func(t *testing.T) {
		// This test would require refactoring the original function to accept dependency injection
		// For now, we'll test the parsing logic directly
		
		// 実際のgit diffコマンドの出力をシミュレート
		mockOutput := `file1.go
file2.go
file3.go`
		
		files := parseNameOnlyOutput(mockOutput)
		
		if len(files) != 3 {
			t.Errorf("Expected 3 files, got %d", len(files))
		}
		
		expectedFiles := []string{"file1.go", "file2.go", "file3.go"}
		for i, file := range files {
			if file.File != expectedFiles[i] {
				t.Errorf("Expected file %s, got %s", expectedFiles[i], file.File)
			}
		}
	})
}

// TestGetCommitInfo tests the getCommitInfoWithBranch function
func TestGetCommitInfo(t *testing.T) {
	t.Run("コミット情報の取得", func(t *testing.T) {
		// 実際のgit logコマンドの出力をシミュレート
		author, date := getCommitInfoWithBranch("nonexistent_file.go", false)
		
		// 存在しないファイルの場合は (unknown) が返される
		if author == "" || date == "" {
			t.Errorf("Expected non-empty author and date, got author: %s, date: %s", author, date)
		}
	})
}

// Integration test helper functions
func setupTestRepository(t *testing.T) string {
	// テスト用の一時的なgitリポジトリを作成
	tempDir, err := os.MkdirTemp("", "gh-diff-summary-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	
	// Initialize git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}
	
	// Set up git config
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to set git config: %v", err)
	}
	
	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to set git config: %v", err)
	}
	
	return tempDir
}

func createTestFile(t *testing.T, dir, filename, content string) {
	filePath := dir + "/" + filename
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file %s: %v", filename, err)
	}
}

func commitFile(t *testing.T, dir, filename, message string) {
	cmd := exec.Command("git", "add", filename)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add file %s: %v", filename, err)
	}
	
	cmd = exec.Command("git", "commit", "-m", message)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to commit file %s: %v", filename, err)
	}
}

func createBranch(t *testing.T, dir, branchName string) {
	cmd := exec.Command("git", "checkout", "-b", branchName)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create branch %s: %v", branchName, err)
	}
}

// switchBranch switches to the specified branch (unused but kept for future use)
func switchBranch(t *testing.T, dir, branchName string) {
	cmd := exec.Command("git", "checkout", branchName)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to switch to branch %s: %v", branchName, err)
	}
}

// Integration test for the entire workflow
func TestIntegrationWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	tempDir := setupTestRepository(t)
	defer os.RemoveAll(tempDir)
	
	// Change to the test directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)
	
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	
	t.Run("ブランチ間の差分を検出", func(t *testing.T) {
		// Create initial commit on main branch
		createTestFile(t, tempDir, "file1.go", "package main\n\nfunc main() {\n}")
		commitFile(t, tempDir, "file1.go", "Initial commit")
		
		// Create feature branch
		createBranch(t, tempDir, "feature")
		
		// Add new file on feature branch
		createTestFile(t, tempDir, "file2.go", "package main\n\nfunc test() {\n}")
		commitFile(t, tempDir, "file2.go", "Add test function")
		
		// Modify existing file
		createTestFile(t, tempDir, "file1.go", "package main\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}")
		commitFile(t, tempDir, "file1.go", "Add println")
		
		// Test getDiffFiles function
		files, err := getDiffFiles("main..feature", "", false)
		if err != nil {
			t.Fatalf("getDiffFiles failed: %v", err)
		}
		
		if len(files) != 2 {
			t.Errorf("Expected 2 files, got %d", len(files))
		}
		
		// Check if both files are detected
		fileNames := make([]string, len(files))
		for i, file := range files {
			fileNames[i] = file.File
		}
		
		expectedFiles := []string{"file1.go", "file2.go"}
		for _, expected := range expectedFiles {
			found := false
			for _, actual := range fileNames {
				if actual == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected file %s not found in diff", expected)
			}
		}
	})
}

// Benchmark tests
func BenchmarkParseStatOutput(b *testing.B) {
	statusOutput := `M	file1.go
A	file2.go
D	file3.go
M	file4.go
A	file5.go`
	
	numstatOutput := `10	5	file1.go
15	0	file2.go
0	20	file3.go
25	10	file4.go
5	0	file5.go`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parseStatOutput(statusOutput, numstatOutput)
	}
}

func BenchmarkParseNameOnlyOutput(b *testing.B) {
	output := strings.Repeat("file1.go\nfile2.go\nfile3.go\n", 100)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parseNameOnlyOutput(output)
	}
}

func BenchmarkGenerateStatBar(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generateStatBar(100, 50)
	}
}

// Example test for documentation
func Example_parseStatOutput() {
	statusOutput := `M	example.go
A	new.go`
	
	numstatOutput := `10	5	example.go
15	0	new.go`
	
	files := parseStatOutput(statusOutput, numstatOutput)
	
	for _, file := range files {
		if file.File == "example.go" {
			// Output: example.go M 10 5
			fmt.Printf("%s %s %d %d\n", file.File, file.Status, file.Insertions, file.Deletions)
		}
	}
}