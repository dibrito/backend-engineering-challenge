package cmd

import "testing"

func TestRoot(t *testing.T) {
	// Set the flags before executing the command
	rootCmd.SetArgs([]string{"--input_file=../events.json", "--window=10"})

	// Execute the command
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("Command execution failed: %v", err)
	}
}
