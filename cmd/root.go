package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

const (
	INPUT_FILE_FLAG = "input_file"
	WINDOW_FLAG     = "window"
)

var (
	inputFile string
	window    int32
)

var ErrInvalidWindow = errors.New("window must be a positive integer")
var ErrParseInputFile = errors.New("unable to parse input file")

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "calculator-cli",
	Short: "Calculates the simple moving average(sma) from input data in a given period of time",
	Long: `Calculator-cli will calculate the simple moving average(sma) from a input file in
	in the .json format, the file should be indentified with --input_file flag.
	The time window to be considered in the sma calculation, e.g. 10 min, should be identified by
	flag --window.
	The output will be printed in the stdout.
	calculator_cli --input_file events.json --window_size 10`,
	// SilenceUsage will stop displayinh usage(--help) when error from Execute.
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if window <= 0 {
			return ErrInvalidWindow
		}

		now := time.Now().UTC()
		data, err := parseInputFile(inputFile)
		if err != nil {
			return ErrParseInputFile
		}

		// simpleMovingAverage(data, window)
		smaFIFO(data, window)
		then := time.Now()
		diff := then.Sub(now)
		// this is to measure successful execution time.
		// depending on the outcome we could try to benchmark and improve if possible.
		fmt.Printf("executed in:%.2fs\n", diff.Seconds())
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// define your flags and configuration settings.
	// TODO: we'r defaulting/expecting input json to be at root level
	rootCmd.Flags().StringVar(&inputFile, "input_file", "../events.json", "The input file with recored events")
	rootCmd.Flags().Int32Var(&window, "window", 10, "The time window considered in the sma calculation")
	// TODO: define if we want them to be required of if we can default.
	// default is a good option!
}
