package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const _100K = 100000
const _500K = 500000
const _1M = 1000000

func TestRoot(t *testing.T) {
	tcs := []struct {
		name    string
		args    []string
		wantErr error
		cleanup func(t *testing.T)
	}{
		{
			// will use events.json file located at project root level.
			name:    "when no flags parsed should use defaults",
			args:    []string{},
			wantErr: nil,
			cleanup: func(t *testing.T) {
				err := os.Remove("../result.txt")
				require.NoError(t, err)
			},
		},
		{
			name:    "when unable to open input file should error",
			args:    []string{"--input_file=../noExistingFile.json"},
			wantErr: ErrParseInputFile,
		},

		{
			// will use events.json file located at project root level.
			name:    "when invalid window should error",
			args:    []string{"--window=-1"},
			wantErr: ErrInvalidWindow,
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rootCmd.SetArgs(tc.args)
			err := rootCmd.Execute()
			require.ErrorIs(t, err, tc.wantErr)
			if tc.cleanup != nil {
				tc.cleanup(t)
			}
		})
	}
}

// Test regression. If this start failing, re-check your changes!
func TestRootRegression(t *testing.T) {
	tcs := []struct {
		name          string
		args          []string
		wantErr       error
		setup         func()
		cleanup       func(t *testing.T)
		checkResponse func(t *testing.T)
	}{
		{
			// will use testInput.json file located at project root level.
			name:    "when no errors should create result output",
			args:    []string{"--input_file=./testInput.json", "--window=10"},
			wantErr: nil,
			cleanup: func(t *testing.T) {
				// clean on the result file and keep the test input/result files.
				// it's ok to remove result.txt everytime cause tests default run is sequencial.
				err := os.Remove("../result.txt")
				require.NoError(t, err)
			},
			checkResponse: func(t *testing.T) {
				// parse got got.
				file, err := os.Open("../result.txt")
				require.NoError(t, err)
				defer file.Close()
				got := parseResult(file, t)

				// parse want result.
				file, err = os.Open("./testResult.txt")
				require.NoError(t, err)
				want := parseResult(file, t)

				require.Equal(t, want, got)
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rootCmd.SetArgs(tc.args)
			err := rootCmd.Execute()
			require.NoError(t, err)

			tc.checkResponse(t)
			tc.cleanup(t)
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
		// re-test with "most possible" iddle machine:
		// executed in:54.73s
		// from now we need to decide where to benchmark!
		{
			// will use created heavy-load.json file located at project root level.
			name:    "when input file has 100 entries should successfuly process",
			args:    []string{"--input_file=./heavy-load.json", "--window=5"},
			wantErr: ErrParseInputFile,
			setup: func() {
				generateSampelFile("./heavy-load.json", _1M)
			},
			cleanup: func(t *testing.T) {
				err := os.Remove("./heavy-load.json")
				require.NoError(t, err)

				err = os.Remove("../result.txt")
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

func parseResult(file *os.File, t *testing.T) []output {
	result := []output{}
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		var line output
		bs, err := json.Marshal([]byte(fileScanner.Text()))
		require.NoError(t, err)

		json.Unmarshal(bs, &line)
		result = append(result, line)
	}

	return result
}

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
