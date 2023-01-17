package processingEngine

import (
	"io/ioutil"
	"os"
	"testing"
)

// TestProcessingEngineRunValid verifies that Run() returns the expected exit code, stdout, and stderr
// when the binary is valid
func TestProcessingEngineRunValid(t *testing.T) {
	// Create a temporary directory
	tempDir, err := ioutil.TempDir("", "run-valid")
	if err != nil {
		t.Fatalf("error creating temporary directory: %s", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a temporary binary
	binPath := tempDir + "/bin"
	if err := ioutil.WriteFile(binPath, []byte(`#!/bin/bash
echo $KEY1
echo $KEY2 >&2
exit 42
`), 0755); err != nil {
		t.Fatalf("error creating temporary binary: %s", err)
	}

	// Create a temporary environment file
	envFilePath := tempDir + "/env"
	if err := ioutil.WriteFile(envFilePath, []byte(`
#KEY1=THISLINESHOULDBEIGNORED
KEY1=VALUE1
KEY2=VALUE2
`), 0644); err != nil {
		t.Fatalf("error creating temporary environment file: %s", err)
	}

	// Create a ProcessingEngine
	pe := NewProcessingEngine(binPath, envFilePath, []string{"arg1", "arg2"})

	// Run the ProcessingEngine
	exitCode, err := pe.Run()
	if err != nil {
		t.Fatalf("error running ProcessingEngine: %s", err)
	}

	// Verify the exit code
	if exitCode != 42 {
		t.Errorf("exit code: got %d, want %d", exitCode, 42)
	}

	// Verify the stdout with env-var KEY1
	expectedStdout := "VALUE1\n"
	if pe.GetStdout() != expectedStdout {
		t.Errorf("stdout: got %q, want %q", pe.GetStdout(), expectedStdout)
	}

	// Verify the stderr with env-var KEY2
	expectedStdout = "VALUE2\n"
	if pe.GetStderr() != expectedStdout {
		t.Errorf("stderr: got %q, want %q", pe.GetStderr(), expectedStdout)
	}

}
