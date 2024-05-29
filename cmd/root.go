package cmd

import (
	"fmt"
	"os"

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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("calc called")
		fmt.Printf("reading file:%s\n", inputFile)
		fmt.Printf("time window:%d\n", window)

		if window <= 0 {
			fmt.Printf("window: %v must be a positive integer\n", window)
			return
		}

		data, err := parseInputFile(inputFile)
		if err != nil {
			fmt.Printf("unable to parse input file: %s:%v\n", inputFile, err)
			return
		}

		simpleMovingAverage(data, window)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// define your flags and configuration settings.
	rootCmd.Flags().StringVar(&inputFile, "input_file", "events.json", "The input file with recored events")
	rootCmd.Flags().Int32Var(&window, "window", 10, "The time window considered in the sma calculation")
	// TODO: define if we want them to be required of if we can default.
	// default is a good option!
}
