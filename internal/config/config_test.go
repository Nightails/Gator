package config

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestGetConfigFilePath(t *testing.T) {
	// Test successful case
	t.Run("returns valid config path", func(t *testing.T) {
		path, err := getConfigFilePath()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if path == "" {
			t.Fatal("expected non-empty path")
		}

		// Verify the path contains the config file name
		if !strings.HasSuffix(path, configFileName) {
			t.Errorf("expected path to end with %q, got %q", configFileName, path)
		}

		// Verify the path contains a directory separator before the filename
		if !strings.Contains(path, "/"+configFileName) {
			t.Errorf("expected path to contain '/%s', got %q", configFileName, path)
		}

		// Verify the home directory portion matches what we expect
		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Fatalf("failed to get home directory: %v", err)
		}

		expectedPath := homeDir + "/" + configFileName
		if path != expectedPath {
			t.Errorf("expected path %q, got %q", expectedPath, path)
		}
	})
}

func TestWrite(t *testing.T) {
	t.Run("writes config to file successfully", func(t *testing.T) {
		// Create a test config
		testConfig := Config{
			URL:      "postgres://localhost:5432/testdb",
			UserName: "testuser",
		}

		// Write the config
		err := write(testConfig)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Get the config file path to read it back
		cfgPath, err := getConfigFilePath()
		if err != nil {
			t.Fatalf("failed to get config file path: %v", err)
		}

		// Defer cleanup
		defer os.Remove(cfgPath)

		// Read the file back
		data, err := os.ReadFile(cfgPath)
		if err != nil {
			t.Fatalf("failed to read config file: %v", err)
		}

		// Unmarshal and verify the content
		var readConfig Config
		err = json.Unmarshal(data, &readConfig)
		if err != nil {
			t.Fatalf("failed to unmarshal config: %v", err)
		}

		if readConfig.URL != testConfig.URL {
			t.Errorf("expected URL %q, got %q", testConfig.URL, readConfig.URL)
		}

		if readConfig.UserName != testConfig.UserName {
			t.Errorf("expected UserName %q, got %q", testConfig.UserName, readConfig.UserName)
		}

		// Verify the JSON structure
		if !strings.Contains(string(data), `"db_url"`) {
			t.Error("expected JSON to contain 'db_url' field")
		}

		if !strings.Contains(string(data), `"current_user_name"`) {
			t.Error("expected JSON to contain 'current_user_name' field")
		}
	})

	t.Run("overwrites existing config file", func(t *testing.T) {
		// Write first config
		firstConfig := Config{
			URL:      "postgres://first",
			UserName: "first_user",
		}
		err := write(firstConfig)
		if err != nil {
			t.Fatalf("failed to write first config: %v", err)
		}

		cfgPath, _ := getConfigFilePath()
		defer os.Remove(cfgPath)

		// Write second config (should overwrite)
		secondConfig := Config{
			URL:      "postgres://second",
			UserName: "second_user",
		}
		err = write(secondConfig)
		if err != nil {
			t.Fatalf("failed to write second config: %v", err)
		}

		// Read and verify it's the second config
		data, err := os.ReadFile(cfgPath)
		if err != nil {
			t.Fatalf("failed to read config file: %v", err)
		}

		var readConfig Config
		json.Unmarshal(data, &readConfig)

		if readConfig.URL != secondConfig.URL {
			t.Errorf("expected URL %q, got %q", secondConfig.URL, readConfig.URL)
		}

		if readConfig.UserName != secondConfig.UserName {
			t.Errorf("expected UserName %q, got %q", secondConfig.UserName, readConfig.UserName)
		}
	})
}

func TestRead(t *testing.T) {
	t.Run("reads config from file successfully", func(t *testing.T) {
		// Create a test config and write it first
		testConfig := Config{
			URL:      "postgres://localhost:5432/mydb",
			UserName: "john_doe",
		}

		err := write(testConfig)
		if err != nil {
			t.Fatalf("failed to write test config: %v", err)
		}

		cfgPath, _ := getConfigFilePath()
		defer os.Remove(cfgPath)

		// Now read it back
		readConfig := Read()

		if readConfig.URL != testConfig.URL {
			t.Errorf("expected URL %q, got %q", testConfig.URL, readConfig.URL)
		}

		if readConfig.UserName != testConfig.UserName {
			t.Errorf("expected UserName %q, got %q", testConfig.UserName, readConfig.UserName)
		}
	})

	t.Run("returns empty config when file does not exist", func(t *testing.T) {
		// Ensure the config file doesn't exist
		cfgPath, err := getConfigFilePath()
		if err != nil {
			t.Fatalf("failed to get config file path: %v", err)
		}
		os.Remove(cfgPath) // Remove if it exists

		// Read should return an empty config
		readConfig := Read()

		if readConfig.URL != "" {
			t.Errorf("expected empty URL, got %q", readConfig.URL)
		}

		if readConfig.UserName != "" {
			t.Errorf("expected empty UserName, got %q", readConfig.UserName)
		}
	})

	t.Run("returns empty config when file contains invalid JSON", func(t *testing.T) {
		// Write invalid JSON to the config file
		cfgPath, err := getConfigFilePath()
		if err != nil {
			t.Fatalf("failed to get config file path: %v", err)
		}

		invalidJSON := []byte(`{"db_url": "incomplete json"`)
		err = os.WriteFile(cfgPath, invalidJSON, 0600)
		if err != nil {
			t.Fatalf("failed to write invalid JSON: %v", err)
		}

		defer os.Remove(cfgPath)

		// Read should return an empty config
		readConfig := Read()

		if readConfig.URL != "" {
			t.Errorf("expected empty URL, got %q", readConfig.URL)
		}

		if readConfig.UserName != "" {
			t.Errorf("expected empty UserName, got %q", readConfig.UserName)
		}
	})

	t.Run("reads config with empty fields", func(t *testing.T) {
		// Create a config with empty fields
		testConfig := Config{
			URL:      "",
			UserName: "",
		}

		err := write(testConfig)
		if err != nil {
			t.Fatalf("failed to write test config: %v", err)
		}

		cfgPath, _ := getConfigFilePath()
		defer os.Remove(cfgPath)

		// Read it back
		readConfig := Read()

		if readConfig.URL != "" {
			t.Errorf("expected empty URL, got %q", readConfig.URL)
		}

		if readConfig.UserName != "" {
			t.Errorf("expected empty UserName, got %q", readConfig.UserName)
		}
	})
}

func TestSetUser(t *testing.T) {
	t.Run("sets user name and writes to file", func(t *testing.T) {
		// Create an initial config with a URL
		initialConfig := Config{
			URL:      "postgres://localhost:5432/testdb",
			UserName: "old_user",
		}

		err := write(initialConfig)
		if err != nil {
			t.Fatalf("failed to write initial config: %v", err)
		}

		cfgPath, _ := getConfigFilePath()
		defer os.Remove(cfgPath)

		// Call SetUser to update the username
		initialConfig.SetUser("new_user")

		// Read the config back from file to verify it was persisted
		readConfig := Read()

		if readConfig.UserName != "new_user" {
			t.Errorf("expected UserName %q, got %q", "new_user", readConfig.UserName)
		}

		// Verify the URL is preserved
		if readConfig.URL != initialConfig.URL {
			t.Errorf("expected URL %q to be preserved, got %q", initialConfig.URL, readConfig.URL)
		}
	})

	t.Run("updates user name from empty to populated", func(t *testing.T) {
		// Create config with empty username
		initialConfig := Config{
			URL:      "postgres://localhost:5432/testdb",
			UserName: "",
		}

		err := write(initialConfig)
		if err != nil {
			t.Fatalf("failed to write initial config: %v", err)
		}

		cfgPath, _ := getConfigFilePath()
		defer os.Remove(cfgPath)

		// Set a username
		initialConfig.SetUser("first_user")

		// Read back and verify
		readConfig := Read()

		if readConfig.UserName != "first_user" {
			t.Errorf("expected UserName %q, got %q", "first_user", readConfig.UserName)
		}

		if readConfig.URL != initialConfig.URL {
			t.Errorf("expected URL to be preserved")
		}
	})

	t.Run("overwrites existing user name", func(t *testing.T) {
		// Create initial config
		initialConfig := Config{
			URL:      "postgres://localhost:5432/testdb",
			UserName: "first_user",
		}

		err := write(initialConfig)
		if err != nil {
			t.Fatalf("failed to write initial config: %v", err)
		}

		cfgPath, _ := getConfigFilePath()
		defer os.Remove(cfgPath)

		// Update username multiple times
		initialConfig.SetUser("second_user")

		readConfig := Read()
		if readConfig.UserName != "second_user" {
			t.Errorf("expected UserName %q after first update, got %q", "second_user", readConfig.UserName)
		}

		// Update again
		readConfig.SetUser("third_user")

		finalConfig := Read()
		if finalConfig.UserName != "third_user" {
			t.Errorf("expected UserName %q after second update, got %q", "third_user", finalConfig.UserName)
		}

		// URL should still be preserved
		if finalConfig.URL != initialConfig.URL {
			t.Errorf("expected URL to be preserved through multiple updates")
		}
	})

	t.Run("sets empty username", func(t *testing.T) {
		// Create config with a username
		initialConfig := Config{
			URL:      "postgres://localhost:5432/testdb",
			UserName: "existing_user",
		}

		err := write(initialConfig)
		if err != nil {
			t.Fatalf("failed to write initial config: %v", err)
		}

		cfgPath, _ := getConfigFilePath()
		defer os.Remove(cfgPath)

		// Set username to empty string
		initialConfig.SetUser("")

		// Read back and verify
		readConfig := Read()

		if readConfig.UserName != "" {
			t.Errorf("expected empty UserName, got %q", readConfig.UserName)
		}

		if readConfig.URL != initialConfig.URL {
			t.Errorf("expected URL to be preserved")
		}
	})
}
