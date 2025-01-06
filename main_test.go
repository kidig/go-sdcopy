package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCopyFile(t *testing.T) {
	// Create temporary source and destination directories
	srcDir, err := os.MkdirTemp("", "src")
	if err != nil {
		t.Fatalf("failed to create temp source directory: %v", err)
	}
	defer os.RemoveAll(srcDir)

	dstDir, err := os.MkdirTemp("", "dst")
	if err != nil {
		t.Fatalf("failed to create temp destination directory: %v", err)
	}
	defer os.RemoveAll(dstDir)

	// Create a temporary source file
	srcFile, err := os.CreateTemp(srcDir, "testfile")
	if err != nil {
		t.Fatalf("failed to create temp source file: %v", err)
	}
	defer srcFile.Close()

	// Write some content to the source file
	content := []byte("Hello, World!")
	if _, err := srcFile.Write(content); err != nil {
		t.Fatalf("failed to write to temp source file: %v", err)
	}

	// Set modification and access times
	mtime := time.Date(2023, time.March, 14, 0, 0, 0, 0, time.UTC)
	atime := time.Date(2023, time.March, 14, 0, 0, 0, 0, time.UTC)
	if err := os.Chtimes(srcFile.Name(), atime, mtime); err != nil {
		t.Fatalf("failed to set file times: %v", err)
	}

	// Get file info
	info, err := srcFile.Stat()
	if err != nil {
		t.Fatalf("failed to get file info: %v", err)
	}

	// Call copyFile function
	err = copyFile(srcFile.Name(), dstDir, info)
	if err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	// Check if the file was copied
	dstFile := filepath.Join(dstDir, filepath.Base(srcFile.Name()))
	if _, err := os.Stat(dstFile); os.IsNotExist(err) {
		t.Fatalf("destination file does not exist: %v", dstFile)
	}

	// Check if the content is the same
	copiedContent, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("failed to read destination file: %v", err)
	}
	if string(copiedContent) != string(content) {
		t.Errorf("expected content %s, got %s", string(content), string(copiedContent))
	}

	// Check if the timestamps are preserved
	dstInfo, err := os.Stat(dstFile)
	if err != nil {
		t.Fatalf("failed to get destination file info: %v", err)
	}
	if !dstInfo.ModTime().Equal(mtime) {
		t.Errorf("expected modification time %v, got %v", mtime, dstInfo.ModTime())
	}
}

func TestFileExists(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Test that the file exists
	if !fileExists(tempFile.Name()) {
		t.Errorf("Expected file to exist: %s", tempFile.Name())
	}

	// Test that a non-existent file does not exist
	if fileExists("nonexistentfile.txt") {
		t.Errorf("Expected file to not exist: nonexistentfile.txt")
	}
}

func TestResolveDestinationPath(t *testing.T) {
	d := time.Date(2023, time.March, 14, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		input    string
		expected string
	}{
		{"Today is {year}-{month}-{day}", "Today is 2023-03-14"},
		{"Year: {year}, Month: {month}, Day: {day}", "Year: 2023, Month: 03, Day: 14"},
		{"No placeholders here", "No placeholders here"},
		{"Unknown {placeholder}", "Unknown "},
	}

	for _, test := range tests {
		result := resolveDestinationPath(test.input, d)
		if result != test.expected {
			t.Errorf("expected %s, got %s", test.expected, result)
		}
	}
}
