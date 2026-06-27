package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExecFilePath(t *testing.T) {
	// Defer recover to catch any panic and convert it to a test failure.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("ExecFilePath() panicked: %v", r)
		}
	}()

	path := ExecFilePath()
	if path == "" {
		t.Error("ExecFilePath() returned an empty string")
	}
}

func TestRunDir(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("RunDir() panicked: %v", r)
		}
	}()

	dir := RunDir()
	if dir == "" {
		t.Error("RunDir() returned an empty string")
	}

	if !filepath.IsAbs(dir) {
		t.Errorf("RunDir() returned a non-absolute path: %s", dir)
	}
}

func TestIsFile(t *testing.T) {
	// Test 1: an existing file — use this test file itself.
	// The test binary runs with the package directory as its working directory,
	// so "util_test.go" resolves correctly.
	existingFile := "util_test.go"
	if !IsFile(existingFile) {
		t.Errorf("IsFile(%q) = false, want true for an existing file", existingFile)
	}

	// Test 2: a file path that definitely does not exist.
	nonExistent := filepath.Join(os.TempDir(), "cli-go-util-test-nonexistent-12345")
	if IsFile(nonExistent) {
		t.Errorf("IsFile(%q) = true, want false for a non-existent path", nonExistent)
	}

	// Test 3: a directory path must return false.
	dirPath := "."
	if IsFile(dirPath) {
		t.Errorf("IsFile(%q) = true, want false for a directory path", dirPath)
	}

	// Extra safety: use an absolute directory path as well.
	absDir := os.TempDir()
	if IsFile(absDir) {
		t.Errorf("IsFile(%q) = true, want false for an absolute directory path", absDir)
	}
}

func TestCurrentDir(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("CurrentDir() panicked: %v", r)
		}
	}()

	path := CurrentDir()
	if path == "" {
		t.Error("CurrentDir() returned an empty string")
	}

	if !filepath.IsAbs(path) {
		t.Errorf("CurrentDir() returned a non-absolute path: %s", path)
	}
}

func TestIsRunWithGoRunCmd(t *testing.T) {
	// IsRunWithGoRunCmd only panics if RunDir panics, and RunDir panics
	// only if os.Args[0] can't be resolved. That won't happen under test,
	// but we guard anyway.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("IsRunWithGoRunCmd() panicked: %v", r)
		}
	}()

	result := IsRunWithGoRunCmd()
	// We can't force go-run mode in a compiled test, but we can assert
	// the function returns a boolean without panicking.
	if result != true && result != false {
		t.Errorf("IsRunWithGoRunCmd() returned a non-boolean value: %v", result)
	}
}
