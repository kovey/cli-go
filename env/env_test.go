package env

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// ---- Load tests ----

func TestLoad(t *testing.T) {
	t.Run("KEY=VALUE format", func(t *testing.T) {
		dir := t.TempDir()
		filePath := filepath.Join(dir, ".env.test")
		content := "TEST_SIMPLE_KEY=test_value_1\n"
		writeTempEnv(t, filePath, content)

		err := Load(filePath)
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		val, ok := os.LookupEnv("TEST_SIMPLE_KEY")
		if !ok {
			t.Fatal("TEST_SIMPLE_KEY not set")
		}
		if val != "test_value_1" {
			t.Errorf("expected 'test_value_1', got '%s'", val)
		}
		os.Unsetenv("TEST_SIMPLE_KEY")
	})

	t.Run("KEY = VALUE (spaces around =)", func(t *testing.T) {
		dir := t.TempDir()
		filePath := filepath.Join(dir, ".env.test")
		content := "TEST_SPACE_KEY = test_space_value\n"
		writeTempEnv(t, filePath, content)

		err := Load(filePath)
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		val, ok := os.LookupEnv("TEST_SPACE_KEY")
		if !ok {
			t.Fatal("TEST_SPACE_KEY not set")
		}
		if val != "test_space_value" {
			t.Errorf("expected 'test_space_value', got '%s'", val)
		}
		os.Unsetenv("TEST_SPACE_KEY")
	})

	t.Run("comments (#)", func(t *testing.T) {
		dir := t.TempDir()
		filePath := filepath.Join(dir, ".env.test")
		content := "# This is a hash comment\nHASH_KEY=hash_value\n"
		writeTempEnv(t, filePath, content)

		err := Load(filePath)
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		val, ok := os.LookupEnv("HASH_KEY")
		if !ok {
			t.Fatal("HASH_KEY not set")
		}
		if val != "hash_value" {
			t.Errorf("expected 'hash_value', got '%s'", val)
		}
		os.Unsetenv("HASH_KEY")
	})

	t.Run("comments (;)", func(t *testing.T) {
		dir := t.TempDir()
		filePath := filepath.Join(dir, ".env.test")
		content := "; This is a semicolon comment\nSEMI_KEY=semi_value\n"
		writeTempEnv(t, filePath, content)

		err := Load(filePath)
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		val, ok := os.LookupEnv("SEMI_KEY")
		if !ok {
			t.Fatal("SEMI_KEY not set")
		}
		if val != "semi_value" {
			t.Errorf("expected 'semi_value', got '%s'", val)
		}
		os.Unsetenv("SEMI_KEY")
	})

	t.Run("comments (--)", func(t *testing.T) {
		dir := t.TempDir()
		filePath := filepath.Join(dir, ".env.test")
		content := "-- This is a double-dash comment\nDASH_KEY=dash_value\n"
		writeTempEnv(t, filePath, content)

		err := Load(filePath)
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		val, ok := os.LookupEnv("DASH_KEY")
		if !ok {
			t.Fatal("DASH_KEY not set")
		}
		if val != "dash_value" {
			t.Errorf("expected 'dash_value', got '%s'", val)
		}
		os.Unsetenv("DASH_KEY")
	})

	t.Run("comments (//)", func(t *testing.T) {
		dir := t.TempDir()
		filePath := filepath.Join(dir, ".env.test")
		content := "// This is a slash-slash comment\nSLASH_KEY=slash_value\n"
		writeTempEnv(t, filePath, content)

		err := Load(filePath)
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		val, ok := os.LookupEnv("SLASH_KEY")
		if !ok {
			t.Fatal("SLASH_KEY not set")
		}
		if val != "slash_value" {
			t.Errorf("expected 'slash_value', got '%s'", val)
		}
		os.Unsetenv("SLASH_KEY")
	})

	t.Run("empty lines", func(t *testing.T) {
		dir := t.TempDir()
		filePath := filepath.Join(dir, ".env.test")
		content := "\n\nEMPTY_TEST_KEY=empty_test_value\n\n\n"
		writeTempEnv(t, filePath, content)

		err := Load(filePath)
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		val, ok := os.LookupEnv("EMPTY_TEST_KEY")
		if !ok {
			t.Fatal("EMPTY_TEST_KEY not set")
		}
		if val != "empty_test_value" {
			t.Errorf("expected 'empty_test_value', got '%s'", val)
		}
		os.Unsetenv("EMPTY_TEST_KEY")
	})

	t.Run("values containing =", func(t *testing.T) {
		dir := t.TempDir()
		filePath := filepath.Join(dir, ".env.test")
		content := "EQ_KEY=val1=val2=val3\n"
		writeTempEnv(t, filePath, content)

		err := Load(filePath)
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		val, ok := os.LookupEnv("EQ_KEY")
		if !ok {
			t.Fatal("EQ_KEY not set")
		}
		if val != "val1=val2=val3" {
			t.Errorf("expected 'val1=val2=val3', got '%s'", val)
		}
		os.Unsetenv("EQ_KEY")
	})

	t.Run("values containing = with spaces", func(t *testing.T) {
		dir := t.TempDir()
		filePath := filepath.Join(dir, ".env.test")
		content := "EQ_SPACE_KEY = val1 = val2 = val3\n"
		writeTempEnv(t, filePath, content)

		err := Load(filePath)
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		val, ok := os.LookupEnv("EQ_SPACE_KEY")
		if !ok {
			t.Fatal("EQ_SPACE_KEY not set")
		}
		if val != "val1=val2=val3" {
			t.Errorf("expected 'val1=val2=val3', got '%s'", val)
		}
		os.Unsetenv("EQ_SPACE_KEY")
	})

	t.Run("invalid format (no =)", func(t *testing.T) {
		dir := t.TempDir()
		filePath := filepath.Join(dir, ".env.test")
		content := "INVALID_LINE_NO_EQUALS\n"
		writeTempEnv(t, filePath, content)

		err := Load(filePath)
		if err == nil {
			t.Fatal("expected error for invalid format, got nil")
		}
	})

	t.Run("invalid format stops parsing", func(t *testing.T) {
		dir := t.TempDir()
		filePath := filepath.Join(dir, ".env.test")
		content := "VALID_KEY=valid_value\nINVALID_LINE\nAFTER_KEY=should_not_be_set\n"
		writeTempEnv(t, filePath, content)

		err := Load(filePath)
		if err == nil {
			t.Fatal("expected error for invalid format, got nil")
		}

		// VALID_KEY should have been set before the error
		val, ok := os.LookupEnv("VALID_KEY")
		if !ok {
			t.Fatal("VALID_KEY should have been set before error line")
		}
		if val != "valid_value" {
			t.Errorf("expected 'valid_value', got '%s'", val)
		}
		os.Unsetenv("VALID_KEY")

		// AFTER_KEY should NOT be set (parsing stopped)
		_, ok = os.LookupEnv("AFTER_KEY")
		if ok {
			t.Fatal("AFTER_KEY should NOT be set (parsing stopped at error)")
		}
	})

	t.Run("last line without newline", func(t *testing.T) {
		dir := t.TempDir()
		filePath := filepath.Join(dir, ".env.test")
		// No trailing newline
		content := "FINAL_LINE_KEY=final_line_value"
		writeTempEnv(t, filePath, content)

		err := Load(filePath)
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		val, ok := os.LookupEnv("FINAL_LINE_KEY")
		if !ok {
			t.Fatal("FINAL_LINE_KEY not set")
		}
		if val != "final_line_value" {
			t.Errorf("expected 'final_line_value', got '%s'", val)
		}
		os.Unsetenv("FINAL_LINE_KEY")
	})

	t.Run("file not found", func(t *testing.T) {
		err := Load("/nonexistent/path/.env")
		if err == nil {
			t.Fatal("expected error for missing file, got nil")
		}
	})

	t.Run("empty value (KEY=)", func(t *testing.T) {
		dir := t.TempDir()
		filePath := filepath.Join(dir, ".env.test")
		content := "EMPTY_VAL_KEY=\n"
		writeTempEnv(t, filePath, content)

		err := Load(filePath)
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		val, ok := os.LookupEnv("EMPTY_VAL_KEY")
		if !ok {
			t.Fatal("EMPTY_VAL_KEY not set")
		}
		if val != "" {
			t.Errorf("expected empty string, got '%s'", val)
		}
		os.Unsetenv("EMPTY_VAL_KEY")
	})

	t.Run("lines with leading and trailing whitespace", func(t *testing.T) {
		dir := t.TempDir()
		filePath := filepath.Join(dir, ".env.test")
		content := "   \t  WS_KEY  =  ws_value  \t  \n"
		writeTempEnv(t, filePath, content)

		err := Load(filePath)
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		val, ok := os.LookupEnv("WS_KEY")
		if !ok {
			t.Fatal("WS_KEY not set")
		}
		if val != "ws_value" {
			t.Errorf("expected 'ws_value', got '%s'", val)
		}
		os.Unsetenv("WS_KEY")
	})

	t.Run("multiple key-value pairs", func(t *testing.T) {
		dir := t.TempDir()
		filePath := filepath.Join(dir, ".env.test")
		content := "MULTI_KEY1=value1\nMULTI_KEY2=value2\nMULTI_KEY3=value3\n"
		writeTempEnv(t, filePath, content)

		err := Load(filePath)
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		expected := map[string]string{
			"MULTI_KEY1": "value1",
			"MULTI_KEY2": "value2",
			"MULTI_KEY3": "value3",
		}
		for key, expectedVal := range expected {
			val, ok := os.LookupEnv(key)
			if !ok {
				t.Errorf("%s not set", key)
				continue
			}
			if val != expectedVal {
				t.Errorf("expected '%s' for %s, got '%s'", expectedVal, key, val)
			}
			os.Unsetenv(key)
		}
	})
}

// ---- Get tests ----

func TestGet(t *testing.T) {
	const key = "TEST_GET_KEY"
	const expected = "get_test_value"

	os.Setenv(key, expected)
	t.Cleanup(func() { os.Unsetenv(key) })

	val, err := Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != expected {
		t.Errorf("expected '%s', got '%s'", expected, val)
	}
}

func TestGetInt(t *testing.T) {
	const key = "TEST_GET_INT_KEY"
	const expected = 42

	os.Setenv(key, "42")
	t.Cleanup(func() { os.Unsetenv(key) })

	val, err := GetInt(key)
	if err != nil {
		t.Fatalf("GetInt failed: %v", err)
	}
	if val != expected {
		t.Errorf("expected %d, got %d", expected, val)
	}
}

func TestGetInt_Negative(t *testing.T) {
	const key = "TEST_GET_INT_NEG_KEY"
	const expected = -123

	os.Setenv(key, "-123")
	t.Cleanup(func() { os.Unsetenv(key) })

	val, err := GetInt(key)
	if err != nil {
		t.Fatalf("GetInt failed: %v", err)
	}
	if val != expected {
		t.Errorf("expected %d, got %d", expected, val)
	}
}

func TestGetInt_InvalidFormat(t *testing.T) {
	const key = "TEST_GET_INT_INVALID_KEY"

	os.Setenv(key, "not_a_number")
	t.Cleanup(func() { os.Unsetenv(key) })

	_, err := GetInt(key)
	if err == nil {
		t.Fatal("expected error for invalid int format, got nil")
	}
}

func TestGetFloat(t *testing.T) {
	const key = "TEST_GET_FLOAT_KEY"
	const expected = 3.14

	os.Setenv(key, "3.14")
	t.Cleanup(func() { os.Unsetenv(key) })

	val, err := GetFloat(key)
	if err != nil {
		t.Fatalf("GetFloat failed: %v", err)
	}
	if val != expected {
		t.Errorf("expected %f, got %f", expected, val)
	}
}

func TestGetFloat_Negative(t *testing.T) {
	const key = "TEST_GET_FLOAT_NEG_KEY"
	const expected = -2.718

	os.Setenv(key, "-2.718")
	t.Cleanup(func() { os.Unsetenv(key) })

	val, err := GetFloat(key)
	if err != nil {
		t.Fatalf("GetFloat failed: %v", err)
	}
	if val != expected {
		t.Errorf("expected %f, got %f", expected, val)
	}
}

func TestGetFloat_InvalidFormat(t *testing.T) {
	const key = "TEST_GET_FLOAT_INVALID_KEY"

	os.Setenv(key, "not_a_float")
	t.Cleanup(func() { os.Unsetenv(key) })

	_, err := GetFloat(key)
	if err == nil {
		t.Fatal("expected error for invalid float format, got nil")
	}
}

func TestGetBool(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected bool
	}{
		{"true", "true", true},
		{"TRUE", "TRUE", true},
		{"1", "1", true},
		{"false", "false", false},
		{"FALSE", "FALSE", false},
		{"0", "0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const key = "TEST_GET_BOOL_KEY"
			os.Setenv(key, tt.envValue)
			t.Cleanup(func() { os.Unsetenv(key) })

			val, err := GetBool(key)
			if err != nil {
				t.Fatalf("GetBool(%q) failed: %v", tt.envValue, err)
			}
			if val != tt.expected {
				t.Errorf("GetBool(%q): expected %v, got %v", tt.envValue, tt.expected, val)
			}
		})
	}
}

func TestGetBool_InvalidFormat(t *testing.T) {
	const key = "TEST_GET_BOOL_INVALID_KEY"

	os.Setenv(key, "not_a_bool")
	t.Cleanup(func() { os.Unsetenv(key) })

	_, err := GetBool(key)
	if err == nil {
		t.Fatal("expected error for invalid bool format, got nil")
	}
}

// ---- Key-not-found tests ----

func TestGet_KeyNotFound(t *testing.T) {
	val, err := Get("NONEXISTENT_KEY_12345")
	if err != Err_Key_Not_Found {
		t.Errorf("expected Err_Key_Not_Found, got %v", err)
	}
	if val != "" {
		t.Errorf("expected empty string for missing key, got '%s'", val)
	}
}

func TestGetInt_KeyNotFound(t *testing.T) {
	val, err := GetInt("NONEXISTENT_INT_KEY_12345")
	if err != Err_Key_Not_Found {
		t.Errorf("expected Err_Key_Not_Found, got %v", err)
	}
	if val != 0 {
		t.Errorf("expected 0 for missing key, got %d", val)
	}
}

func TestGetFloat_KeyNotFound(t *testing.T) {
	val, err := GetFloat("NONEXISTENT_FLOAT_KEY_12345")
	if err != Err_Key_Not_Found {
		t.Errorf("expected Err_Key_Not_Found, got %v", err)
	}
	if val != 0.0 {
		t.Errorf("expected 0.0 for missing key, got %f", val)
	}
}

func TestGetBool_KeyNotFound(t *testing.T) {
	val, err := GetBool("NONEXISTENT_BOOL_KEY_12345")
	if err != Err_Key_Not_Found {
		t.Errorf("expected Err_Key_Not_Found, got %v", err)
	}
	if val != false {
		t.Errorf("expected false for missing key, got %v", val)
	}
}

// ---- HasEnv tests ----

func TestHasEnv(t *testing.T) {
	dir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current dir: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("failed to restore working directory: %v", err)
		}
	})

	// Save and restore loadTime to avoid affecting other tests
	oldLoadTime := loadTime
	t.Cleanup(func() { loadTime = oldLoadTime })

	// Set loadTime to 0 so any file with mod time > epoch passes
	loadTime = 0

	// No .env file yet
	if HasEnv() {
		t.Fatal("HasEnv should return false when .env does not exist")
	}

	// Create .env file
	if err := os.WriteFile(Default, []byte("TEST_HASENV_KEY=value\n"), 0644); err != nil {
		t.Fatalf("failed to create .env: %v", err)
	}

	// File exists, mod time > 0 (loadTime), so HasEnv should return true
	if !HasEnv() {
		t.Fatal("HasEnv should return true when .env exists and mod time > loadTime")
	}

	// Set loadTime to far future, mod time will be less than loadTime
	loadTime = time.Now().Unix() + 10000

	// File exists but mod time is before loadTime
	if HasEnv() {
		t.Fatal("HasEnv should return false when .env mod time <= loadTime")
	}
}

// ---- CheckDefault tests ----

func TestCheckDefault(t *testing.T) {
	dir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current dir: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("failed to restore working directory: %v", err)
		}
	})

	// No .env file yet
	if CheckDefault() {
		t.Fatal("CheckDefault should return false when .env does not exist")
	}

	// Create .env file
	if err := os.WriteFile(Default, []byte("TEST_CHECKDEFAULT_KEY=value\n"), 0644); err != nil {
		t.Fatalf("failed to create .env: %v", err)
	}

	// .env file exists
	if !CheckDefault() {
		t.Fatal("CheckDefault should return true when .env exists")
	}

	// Remove file and create a directory named .env
	os.Remove(Default)
	if err := os.Mkdir(Default, 0755); err != nil {
		t.Fatalf("failed to create .env directory: %v", err)
	}

	// .env is a directory now
	if CheckDefault() {
		t.Fatal("CheckDefault should return false when .env is a directory")
	}
}

// ---- LoadDefault tests ----

func TestLoadDefault(t *testing.T) {
	dir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current dir: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("failed to restore working directory: %v", err)
		}
	})

	// Save and restore loadTime
	oldLoadTime := loadTime
	t.Cleanup(func() { loadTime = oldLoadTime })

	// Create .env file
	if err := os.WriteFile(Default, []byte("LOADDEFAULT_KEY=loaddefault_value\n"), 0644); err != nil {
		t.Fatalf("failed to create .env: %v", err)
	}

	now := time.Now()
	if err := LoadDefault(now); err != nil {
		t.Fatalf("LoadDefault failed: %v", err)
	}

	// Verify env var was set
	val, ok := os.LookupEnv("LOADDEFAULT_KEY")
	if !ok {
		t.Fatal("LOADDEFAULT_KEY not set")
	}
	if val != "loaddefault_value" {
		t.Errorf("expected 'loaddefault_value', got '%s'", val)
	}

	// Verify loadTime was updated
	if loadTime != now.Unix() {
		t.Errorf("expected loadTime=%d, got %d", now.Unix(), loadTime)
	}

	os.Unsetenv("LOADDEFAULT_KEY")
}

func TestLoadDefault_FileNotFound(t *testing.T) {
	dir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current dir: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("failed to restore working directory: %v", err)
		}
	})

	// No .env file exists
	err = LoadDefault(time.Now())
	if err == nil {
		t.Fatal("expected error when .env does not exist, got nil")
	}
}

// ---- Helper ----

func writeTempEnv(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
}
