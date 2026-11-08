package adb

import (
	"testing"
)

func TestFindProgram(t *testing.T) {
	// Should find 'go' as it is installed in environment
	path, err := findProgram("go")
	if err != nil || path == "" {
		t.Fatalf("expected to find 'go', got error: %v", err)
	}

	// Should fail on non-existent binary
	nonExistent := "non_existent_binary_xyz_123"
	_, err = findProgram(nonExistent)
	if err == nil {
		t.Fatalf("expected error for non-existent binary, got nil")
	}
}

func TestNewClientDiscovery(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Skipf("skipping NewClient check if adb/scrcpy not in environment: %v", err)
	}
	if client.ADBPath == "" || client.ScrcpyPath == "" {
		t.Fatalf("expected valid paths, got ADB: %s, Scrcpy: %s", client.ADBPath, client.ScrcpyPath)
	}
}

