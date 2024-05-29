package cmd

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

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

func TestRootWithHeavyLoad(t *testing.T) {
	// we will skip this from the "make test" call
	// since it's used only insights on improvements.
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	tcs := []struct {
		name    string
		args    []string
		wantErr error
		setup   func()
		cleanup func(t *testing.T)
	}{
		// 100000 entries
		// executed in: ~84.409s
		// from now we need to decide where to benchmark!
		{
			name:    "when input file has 100 entries should successfuly process",
			args:    []string{"--input_file=../heavy-load.json", "--window=5"},
			wantErr: ErrParseInputFile,
			setup: func() {
				// create sample file with 100 entries
				generateSampelFile("../heavy-load.json", 100000)
			},
			cleanup: func(t *testing.T) {
				err := os.Remove("../heavy-load.json")
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// create sample data
			tc.setup()
			// Set the flags before executing the command.
			rootCmd.SetArgs(tc.args)

			// Execute the command.
			err := rootCmd.Execute()
			require.NoError(t, err)
			tc.cleanup(t)
		})
	}
}

// test utillity code:
// We will create sample data to test the solution under 'heavy load'.

type testEvent struct {
	Timestamp      string `json:"timestamp"`
	TranslationID  string `json:"translation_id"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
	ClientName     string `json:"client_name"`
	EventName      string `json:"event_name"`
	NrWords        int    `json:"nr_words"`
	Duration       int    `json:"duration"`
}

func generateRandomEvent(baseTime time.Time, id int) testEvent {
	clients := []string{"airliberty", "taxi-eats", "flyhigh", "quicktrans"}
	languages := []string{"en", "fr", "de", "es"}
	events := []string{"translation_delivered", "translation_requested"}

	// Increment the timestamp by a random number of minutes
	timestamp := baseTime.Add(time.Duration(rand.Intn(60)+1) * time.Minute)
	return testEvent{
		Timestamp:      timestamp.Format("2006-01-02 15:04:05.999999"),
		TranslationID:  fmt.Sprintf("%d", id),
		SourceLanguage: languages[rand.Intn(len(languages))],
		TargetLanguage: languages[rand.Intn(len(languages))],
		ClientName:     clients[rand.Intn(len(clients))],
		EventName:      events[rand.Intn(len(events))],
		NrWords:        rand.Intn(100) + 1,
		Duration:       rand.Intn(120) + 1,
	}
}

func generateSampelFile(filename string, numEntries int) {
	starTime := time.Now().UTC()
	events := make([]testEvent, numEntries)

	for i := 0; i < numEntries; i++ {
		events[i] = generateRandomEvent(starTime, i)
		// increase event timestamp over time.
		starTime = starTime.Add(time.Minute)
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		// if we can't create test files we should stop.
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(events); err != nil {
		fmt.Println(err)
		// if we can't encode events we should stop.
		panic(err)
	}
}
