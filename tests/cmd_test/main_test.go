package cmd_test

//
//import (
//	"bytes"
//	"os/exec"
//	"strings"
//	"testing"
//)
//
//// TestRunApp tests the main application loop
//func TestRunApp(t *testing.T) {
//	cmd := exec.Command("go", "run", "main.go")
//
//	var out bytes.Buffer
//	var errOut bytes.Buffer
//	cmd.Stdout = &out
//	cmd.Stderr = &errOut
//
//	// Start the command
//	if err := cmd.Start(); err != nil {
//		t.Fatalf("Failed to start command: %v", err)
//	}
//
//	// Simulate user input
//	// Using go routines to simulate user input (note: this is a simplification)
//	go func() {
//		// Wait for the application to prompt
//		cmd.Stdin.Write([]byte("1\n")) // For SignUp
//		cmd.Stdin.Write([]byte("2\n")) // For Login
//		cmd.Stdin.Write([]byte("3\n")) // For Exit
//	}()
//
//	if err := cmd.Wait(); err != nil {
//		t.Fatalf("Command finished with error: %v\n%s", err, errOut.String())
//	}
//
//	// Check if the output contains expected strings
//	output := out.String()
//	if !strings.Contains(output, "Welcome") {
//		t.Errorf("Expected output to contain 'Welcome', but got %s", output)
//	}
//
//	if !strings.Contains(output, "SignUp") || !strings.Contains(output, "Login") || !strings.Contains(output, "Exit") {
//		t.Errorf("Expected options to be present in the output, but got %s", output)
//	}
//}
