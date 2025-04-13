package main

import (
	"bufio"
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

func main() {
	// Test durations to validate ParseDuration function
	testDurations := []string{
		"1h 21m 27s",
		"1m 53s",
		"2m",
		"42.9s",
		"21h 40m 51s",
		"3.5d",
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
