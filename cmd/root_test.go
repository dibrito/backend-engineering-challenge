package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoot(t *testing.T) {

	tcs := []struct {
		name    string
		args    []string
		wantErr error
	}{
		{
			name:    "when no flags parsed should use defaults",
			args:    []string{},
			wantErr: nil,
		},
		{
			name:    "when unable to open input file should error",
			args:    []string{"--input_file=../noExistingFile.json"},
			wantErr: ErrParseInputFile,
		},

		{
			name:    "when unable to open input file should error",
			args:    []string{"--window=-1"},
			wantErr: ErrInvalidWindow,
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Set the flags before executing the command.
			rootCmd.SetArgs(tc.args)

			// Execute the command.
			err := rootCmd.Execute()
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}
