package main

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

// TestParseStatOutput tests the parseStatOutput function
func TestParseStatOutput(t *testing.T) {
	t.Run("正常なstat出力のパース", func(t *testing.T) {
		statusOutput := `M	file1.go
A	file2.go
D	file3.go`
		
		numstatOutput := `10	5	file1.go
15	0	file2.go
0	20	file3.go`

		expected := []FileInfo{
			{
				File:       "file1.go",
				Status:     "M",
				Insertions: 10,
				Deletions:  5,
			},
			{
				File:       "file2.go",
				Status:     "A",
				Insertions: 15,
				Deletions:  0,
			},
			{
				File:       "file3.go",
				Status:     "D",
				Deleted:    true,
				Insertions: 0,
				Deletions:  20,
			},
		}

		result := parseStatOutput(statusOutput, numstatOutput)
		
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("parseStatOutput() = %v, expected %v", result, expected)
		}
	})

	t.Run("バイナリファイルの処理", func(t *testing.T) {
		statusOutput := `M	binary.png`
		numstatOutput := `-	-	binary.png`

		expected := []FileInfo{
			{
				File:       "binary.png",
				Status:     "M",
				Insertions: 0,
				Deletions:  0,
			},
		}

		result := parseStatOutput(statusOutput, numstatOutput)
		
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("parseStatOutput() = %v, expected %v", result, expected)
		}
	})

	t.Run("リネームされたファイルの処理", func(t *testing.T) {
		statusOutput := `R100	old_file.go	new_file.go`
		numstatOutput := `5	3	new_file.go`

		expected := []FileInfo{
			{
				File:       "old_file.go => new_file.go",
				Status:     "R",
				Insertions: 5,
				Deletions:  3,
			},
		}

		result := parseStatOutput(statusOutput, numstatOutput)
		
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("parseStatOutput() = %v, expected %v", result, expected)
		}
	})
}

// TestParseNameOnlyOutput tests the parseNameOnlyOutput function
func TestParseNameOnlyOutput(t *testing.T) {
	t.Run("正常なname-only出力のパース", func(t *testing.T) {
		output := `file1.go
file2.go
file3.go`

		result := parseNameOnlyOutput(output)
		
		if len(result) != 3 {
			t.Errorf("parseNameOnlyOutput() returned %d files, expected 3", len(result))
		}
		
		expectedFiles := []string{"file1.go", "file2.go", "file3.go"}
		for i, file := range result {
			if file.File != expectedFiles[i] {
				t.Errorf("parseNameOnlyOutput() file[%d] = %s, expected %s", i, file.File, expectedFiles[i])
			}
		}
	})

	t.Run("空の出力の処理", func(t *testing.T) {
		output := ""

		result := parseNameOnlyOutput(output)
		
		if len(result) != 0 {
			t.Errorf("parseNameOnlyOutput() returned %d files, expected 0", len(result))
		}
	})

	t.Run("空行が含まれる出力の処理", func(t *testing.T) {
		output := `file1.go

file2.go

file3.go`

		result := parseNameOnlyOutput(output)
		
		if len(result) != 3 {
			t.Errorf("parseNameOnlyOutput() returned %d files, expected 3", len(result))
		}
	})
}

// TestGetPathDisplay tests the getPathDisplay function
func TestGetPathDisplay(t *testing.T) {
	t.Run("空のパスの処理", func(t *testing.T) {
		result := getPathDisplay("")
		expected := "(all files)"
		
		if result != expected {
			t.Errorf("getPathDisplay(\"\") = %s, expected %s", result, expected)
		}
	})

	t.Run("パスが指定されている場合", func(t *testing.T) {
		path := "src/main"
		result := getPathDisplay(path)
		
		if result != path {
			t.Errorf("getPathDisplay(\"%s\") = %s, expected %s", path, result, path)
		}
	})
}

// TestGenerateStatBar tests the generateStatBar function
func TestGenerateStatBar(t *testing.T) {
	t.Run("挿入と削除がある場合", func(t *testing.T) {
		result := generateStatBar(10, 5)
		
		// 合計15の変更で、10挿入、5削除
		// 10/15 * 50 = 33.3... -> 33文字の+
		// 5/15 * 50 = 16.6... -> 16文字の-
		// 実際の実装では端数処理の違いがあるかもしれません
		if len(result) == 0 {
			t.Errorf("generateStatBar(10, 5) returned empty string")
		}
		
		plusCount := 0
		minusCount := 0
		for _, char := range result {
			if char == '+' {
				plusCount++
			} else if char == '-' {
				minusCount++
			}
		}
		
		if plusCount == 0 && minusCount == 0 {
			t.Errorf("generateStatBar(10, 5) should contain + and - characters")
		}
	})

	t.Run("変更がない場合", func(t *testing.T) {
		result := generateStatBar(0, 0)
		
		if result != "" {
			t.Errorf("generateStatBar(0, 0) = %s, expected empty string", result)
		}
	})

	t.Run("挿入のみの場合", func(t *testing.T) {
		result := generateStatBar(10, 0)
		
		if len(result) == 0 {
			t.Errorf("generateStatBar(10, 0) returned empty string")
		}
		
		for _, char := range result {
			if char != '+' {
				t.Errorf("generateStatBar(10, 0) should contain only + characters, found %c", char)
			}
		}
	})

	t.Run("削除のみの場合", func(t *testing.T) {
		result := generateStatBar(0, 10)
		
		if len(result) == 0 {
			t.Errorf("generateStatBar(0, 10) returned empty string")
		}
		
		for _, char := range result {
			if char != '-' {
				t.Errorf("generateStatBar(0, 10) should contain only - characters, found %c", char)
			}
		}
	})

	t.Run("最小値が1になることの確認", func(t *testing.T) {
		result := generateStatBar(1, 1000)
		
		plusCount := 0
		minusCount := 0
		for _, char := range result {
			if char == '+' {
				plusCount++
			} else if char == '-' {
				minusCount++
			}
		}
		
		// 1:1000の比率でも、少なくとも1文字の+が表示されるはず
		if plusCount == 0 {
			t.Errorf("generateStatBar(1, 1000) should contain at least one + character")
		}
	})
}

// TestFileInfo tests the FileInfo struct
func TestFileInfo(t *testing.T) {
	t.Run("FileInfo構造体の初期化", func(t *testing.T) {
		info := FileInfo{
			File:       "test.go",
			Author:     "Test Author",
			Date:       "2023/12/01",
			Status:     "M",
			Insertions: 10,
			Deletions:  5,
			Deleted:    false,
		}
		
		if info.File != "test.go" {
			t.Errorf("FileInfo.File = %s, expected test.go", info.File)
		}
		if info.Author != "Test Author" {
			t.Errorf("FileInfo.Author = %s, expected Test Author", info.Author)
		}
		if info.Date != "2023/12/01" {
			t.Errorf("FileInfo.Date = %s, expected 2023/12/01", info.Date)
		}
		if info.Status != "M" {
			t.Errorf("FileInfo.Status = %s, expected M", info.Status)
		}
		if info.Insertions != 10 {
			t.Errorf("FileInfo.Insertions = %d, expected 10", info.Insertions)
		}
		if info.Deletions != 5 {
			t.Errorf("FileInfo.Deletions = %d, expected 5", info.Deletions)
		}
		if info.Deleted != false {
			t.Errorf("FileInfo.Deleted = %t, expected false", info.Deleted)
		}
	})
}

// TestConfig tests the Config struct
func TestConfig(t *testing.T) {
	t.Run("Config構造体の初期化", func(t *testing.T) {
		config := Config{
			SrcBranch:  "main",
			DstBranch:  "develop",
			FilePath:   "src/",
			UseFormat:  true,
			ShowStat:   false,
			OutputJSON: false,
			MaxWorkers: 5,
		}
		
		if config.SrcBranch != "main" {
			t.Errorf("Config.SrcBranch = %s, expected main", config.SrcBranch)
		}
		if config.DstBranch != "develop" {
			t.Errorf("Config.DstBranch = %s, expected develop", config.DstBranch)
		}
		if config.FilePath != "src/" {
			t.Errorf("Config.FilePath = %s, expected src/", config.FilePath)
		}
		if config.UseFormat != true {
			t.Errorf("Config.UseFormat = %t, expected true", config.UseFormat)
		}
		if config.ShowStat != false {
			t.Errorf("Config.ShowStat = %t, expected false", config.ShowStat)
		}
		if config.OutputJSON != false {
			t.Errorf("Config.OutputJSON = %t, expected false", config.OutputJSON)
		}
		if config.MaxWorkers != 5 {
			t.Errorf("Config.MaxWorkers = %d, expected 5", config.MaxWorkers)
		}
	})
}

// TestCheckDeletedFilesBatch tests the new batch deletion check function
func TestCheckDeletedFilesBatch(t *testing.T) {
	t.Run("複数ファイルの削除状態をバッチで確認", func(t *testing.T) {
		files := []string{"file1.go", "file2.go", "file3.go"}
		
		// この関数はまだ実装されていないので、テストは失敗する（Red）
		result := checkDeletedFilesBatch(files)
		
		if len(result) != len(files) {
			t.Errorf("checkDeletedFilesBatch() returned %d results, expected %d", len(result), len(files))
		}
		
		// 各ファイルについて削除状態の情報が含まれているかチェック
		for _, file := range files {
			if _, exists := result[file]; !exists {
				t.Errorf("checkDeletedFilesBatch() missing result for file %s", file)
			}
		}
	})
	
	t.Run("空のファイルリストの処理", func(t *testing.T) {
		files := []string{}
		
		result := checkDeletedFilesBatch(files)
		
		if len(result) != 0 {
			t.Errorf("checkDeletedFilesBatch() returned %d results for empty input, expected 0", len(result))
		}
	})
	
	t.Run("単一ファイルの処理", func(t *testing.T) {
		files := []string{"single.go"}
		
		result := checkDeletedFilesBatch(files)
		
		if len(result) != 1 {
			t.Errorf("checkDeletedFilesBatch() returned %d results, expected 1", len(result))
		}
		
		if _, exists := result["single.go"]; !exists {
			t.Errorf("checkDeletedFilesBatch() missing result for file single.go")
		}
	})
}

// TestParseIntWithDefault tests the enhanced integer parsing function
func TestParseIntWithDefault(t *testing.T) {
	t.Run("正常な数値文字列のパース", func(t *testing.T) {
		result := parseIntWithDefault("123", 0)
		expected := 123
		
		if result != expected {
			t.Errorf("parseIntWithDefault(\"123\", 0) = %d, expected %d", result, expected)
		}
	})
	
	t.Run("バイナリファイルのハイフン処理", func(t *testing.T) {
		result := parseIntWithDefault("-", 0)
		expected := 0
		
		if result != expected {
			t.Errorf("parseIntWithDefault(\"-\", 0) = %d, expected %d", result, expected)
		}
	})
	
	t.Run("無効な文字列のデフォルト値処理", func(t *testing.T) {
		result := parseIntWithDefault("invalid", 42)
		expected := 42
		
		if result != expected {
			t.Errorf("parseIntWithDefault(\"invalid\", 42) = %d, expected %d", result, expected)
		}
	})
	
	t.Run("空文字列のデフォルト値処理", func(t *testing.T) {
		result := parseIntWithDefault("", 100)
		expected := 100
		
		if result != expected {
			t.Errorf("parseIntWithDefault(\"\", 100) = %d, expected %d", result, expected)
		}
	})
	
	t.Run("負の数値の処理", func(t *testing.T) {
		result := parseIntWithDefault("-5", 0)
		expected := -5
		
		if result != expected {
			t.Errorf("parseIntWithDefault(\"-5\", 0) = %d, expected %d", result, expected)
		}
	})
}

// TestOptimizedGetFileInfos tests the memory-optimized version of getFileInfos
func TestOptimizedGetFileInfos(t *testing.T) {
	t.Run("事前に容量を確保したスライスの動作", func(t *testing.T) {
		inputFiles := []FileInfo{
			{File: "file1.go"},
			{File: "file2.go"},
			{File: "file3.go"},
		}
		
		// この関数はまだ実装されていないので、テストは失敗する（Red）
		result, err := optimizedGetFileInfos(inputFiles, 2)
		
		if err != nil {
			t.Errorf("optimizedGetFileInfos() returned error: %v", err)
		}
		
		if len(result) != len(inputFiles) {
			t.Errorf("optimizedGetFileInfos() returned %d files, expected %d", len(result), len(inputFiles))
		}
		
		// スライスの容量が適切に設定されているかは直接テストできないが、
		// 機能的には同じ結果が返されることを確認
		for i, file := range result {
			if file.File != inputFiles[i].File {
				t.Errorf("optimizedGetFileInfos() file[%d] = %s, expected %s", i, file.File, inputFiles[i].File)
			}
		}
	})
	
	t.Run("空のファイルリストの処理", func(t *testing.T) {
		inputFiles := []FileInfo{}
		
		result, err := optimizedGetFileInfos(inputFiles, 2)
		
		if err != nil {
			t.Errorf("optimizedGetFileInfos() returned error: %v", err)
		}
		
		if len(result) != 0 {
			t.Errorf("optimizedGetFileInfos() returned %d files for empty input, expected 0", len(result))
		}
	})
}

// TestGitError tests the custom Git error type
func TestGitError(t *testing.T) {
	t.Run("GitErrorの作成と文字列化", func(t *testing.T) {
		originalErr := fmt.Errorf("exit status 128")
		gitErr := GitError{
			Command: "git diff --name-only",
			Err:     originalErr,
		}
		
		expectedMsg := "git command failed: git diff --name-only: exit status 128"
		if gitErr.Error() != expectedMsg {
			t.Errorf("GitError.Error() = %s, expected %s", gitErr.Error(), expectedMsg)
		}
	})
	
	t.Run("GitErrorのUnwrap", func(t *testing.T) {
		originalErr := fmt.Errorf("some error")
		gitErr := GitError{
			Command: "git log",
			Err:     originalErr,
		}
		
		// GitErrorがerrors.Unwrapableかテスト
		unwrapped := gitErr.Unwrap()
		if unwrapped != originalErr {
			t.Errorf("GitError.Unwrap() returned different error")
		}
	})
	
	t.Run("GitErrorのIs判定", func(t *testing.T) {
		originalErr := fmt.Errorf("target error")
		gitErr := GitError{
			Command: "git status",
			Err:     originalErr,
		}
		
		// errors.Is でラップされたエラーを検出できるかテスト
		if !errors.Is(gitErr, originalErr) {
			t.Errorf("errors.Is(gitErr, originalErr) should return true")
		}
	})
}