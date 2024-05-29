package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// event represents a translation event.
type event struct {
	Timestamp customTime `json:"timestamp"`
	Duration  int        `json:"duration"`
}

// customTime represents an alias for time.
type customTime struct {
	time.Time
}

// UnmarshalJSON is a custom marshaller for the type customTime.
// This gets event timespamp as time.Time to avoid coversions in later steps.
func (t *customTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" || string(data) == `""` {
		return nil
	}

	// Fractional seconds are handled implicitly by Parse.
	// Parse the time string
	// Remove quotes from JSON string.
	timeStr := strings.Trim(string(data), `"`)
	tt, err := parseTime(timeStr)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return err
	}

	*t = customTime{tt}
	return err
}

func parseTime(timeStr string) (time.Time, error) {
	// Layout string corresponding to the input format.
	layout := "2006-01-02 15:04:05.999999"

	// Parse the time string.
	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return time.Time{}, err
	}

	return parsedTime, nil

}

// parseInputFile opens the given input file and marshall into the event struct type.
func parseInputFile(filename string) ([]event, error) {
	var data []event
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return data, err
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println(err)
		return data, err
	}

	return data, nil
}

// writeOutput write the final outputfile ordered by event timestamp.
// result file will always live at root level.
func writeOutput(data map[time.Time]output) error {
	f, err := os.Create("../result.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	// This is needed due to when we range over a map we will get random order
	// affcting the final output result file.
	sorted := sortResultData(data)

	for _, row := range sorted {
		bs, err := json.Marshal(row)
		if err != nil {
			return err
		}

		// we will append a new line at the end of each row for output readability.
		bs = append(bs, []byte("\n")...)
		_, err = f.Write(bs)
		if err != nil {
			return err
		}
	}

	return err
}

// sortResultData sorts ascendetly the input map by key
// where key it the event timestamp, and returns an array of output.
func sortResultData(data map[time.Time]output) []output {
	var result []output

	// first create slice from map.
	for _, v := range data {
		result = append(result, v)
	}

	// use builtin sort func.
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.Before(result[j].Date)
	})

	return result
}

// MarshalJSON is a custom marshaller for the type output.
// This is need to remove the 'Z' from time format.
func (t output) MarshalJSON() ([]byte, error) {
	customStruct := struct {
		// TODO: fix, not really json tags here.
		Date            string  `json:"date"`
		AvgDeliveryTime float32 `json:"average_delivery_time"`
	}{
		// "2006-01-02 15:04:05" is the layout format.
		Date:            t.Date.Format("2006-01-02 15:04:05"),
		AvgDeliveryTime: t.AvgDeliveryTime,
	}
	return json.Marshal(customStruct)
}
