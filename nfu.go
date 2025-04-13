package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TraceRecord represents a single row from the execution trace file
type TraceRecord struct {
	TaskID     string
	Status     string
	Realtime   time.Duration
	CPUPercent float64
	PeakRSS    string
	PeakVmem   string
}

// ParseDuration parses time strings with various suffixes to time.Duration
// Handles formats like "3.6s", "218ms", "1h", "10m", etc.
func ParseDuration(durationStr string) (time.Duration, error) {
	// First try to handle standard durations with time.ParseDuration
	duration, err := time.ParseDuration(durationStr)
	if err == nil {
		return duration, nil
	}

	// Handle complex formats with multiple units like "1h 21m 27s"
	parts := strings.Fields(durationStr)
	var totalDuration time.Duration

	for _, part := range parts {
		// Try to parse each part separately
		partDuration, err := time.ParseDuration(part)
		if err != nil {
			// If parsing fails, it might be due to a non-standard format
			re := regexp.MustCompile(`^([\d\.]+)\s*([a-zA-Z]+)$`)
			matches := re.FindStringSubmatch(part)

			if len(matches) != 3 {
				return 0, fmt.Errorf("unsupported duration format: %s", part)
			}

			valueStr := matches[1]
			unit := strings.ToLower(matches[2])

			value, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				return 0, fmt.Errorf("error parsing duration value %s: %w", valueStr, err)
			}

			// Convert to time.Duration based on unit
			switch unit {
			case "ns", "nanosecond", "nanoseconds":
				partDuration = time.Duration(value * float64(time.Nanosecond))
			case "us", "µs", "microsecond", "microseconds":
				partDuration = time.Duration(value * float64(time.Microsecond))
			case "ms", "millisecond", "milliseconds":
				partDuration = time.Duration(value * float64(time.Millisecond))
			case "s", "sec", "second", "seconds":
				partDuration = time.Duration(value * float64(time.Second))
			case "m", "min", "minute", "minutes":
				partDuration = time.Duration(value * float64(time.Minute))
			case "h", "hr", "hour", "hours":
				partDuration = time.Duration(value * float64(time.Hour))
			case "d", "day", "days":
				partDuration = time.Duration(value * 24 * float64(time.Hour))
			default:
				return 0, fmt.Errorf("unknown time unit: %s", unit)
			}
		}

		totalDuration += partDuration
	}

	return totalDuration, nil
}

// testDurationParsing tests the ParseDuration function with various formats
func testDurationParsing() {
	testDurations := []string{
		"3.5d",
		"21h 40m 51s",
		"1h 21m 27s",
		"2m",
		"1m 53s",
		"42.9s",
		"500ms",
	}

	fmt.Println("Testing ParseDuration function:")
	fmt.Println("-------------------------------")
	for _, durStr := range testDurations {
		dur, err := ParseDuration(durStr)
		if err != nil {
			fmt.Printf("Error parsing '%s': %v\n", durStr, err)
			continue
		}

		// Calculate minutes for easier comparison
		minutes := dur.Minutes()

		fmt.Printf("Original: %-15s | Parsed: %-15v | Minutes: %.2f\n",
			durStr, dur, minutes)
	}
	fmt.Println("-------------------------------")
}

// calculateTotalDuration calculates the total duration from a file
func calculateTotalDuration(filePath string) (time.Duration, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Skip header line
	if !scanner.Scan() {
		return 0, fmt.Errorf("error reading header line: %w", scanner.Err())
	}
	header := scanner.Text()

	// Parse header to find the duration column index
	columns := strings.Split(header, "\t")
	durationIdx := -1
	for i, col := range columns {
		if col == "duration" {
			durationIdx = i
			break
		}
	}

	if durationIdx == -1 {
		return 0, fmt.Errorf("duration column not found in input file")
	}

	var totalDuration time.Duration

	// Process each data line
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")

		// Skip lines with insufficient columns
		if len(fields) <= durationIdx {
			continue
		}

		// Parse duration
		durationStr := fields[durationIdx]
		duration, err := ParseDuration(durationStr)
		if err != nil {
			fmt.Printf("Warning: error parsing duration '%s': %v\n", durationStr, err)
			continue
		}

		totalDuration += duration
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error scanning file: %w", err)
	}

	return totalDuration, nil
}

func main() {
	// Define and parse command line flags
	testFlag := flag.Bool("t", false, "Run tests for duration parsing")
	flag.BoolVar(testFlag, "test", false, "Run tests for duration parsing")

	inputFlag := flag.String("i", "", "Path to the input file")
	flag.StringVar(inputFlag, "input", "", "Path to the input file")

	flag.Parse()

	// If test flag is provided, run test function
	if *testFlag {
		testDurationParsing()
		return
	}

	// Check if input flag is provided
	if *inputFlag == "" {
		fmt.Println("Please provide an input file path using -i or --input flag")
		flag.Usage()
		os.Exit(1)
	}

	// Calculate total duration from the input file
	totalDuration, err := calculateTotalDuration(*inputFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Print the total duration in various formats
	fmt.Printf("Total duration: %v\n", totalDuration)

	// Convert to human-readable format
	hours := int(totalDuration.Hours())
	minutes := int(totalDuration.Minutes()) % 60
	seconds := int(totalDuration.Seconds()) % 60

	fmt.Printf("Total duration: %dh %dm %ds\n", hours, minutes, seconds)
	fmt.Printf("Total minutes: %.2f\n", totalDuration.Minutes())
}
