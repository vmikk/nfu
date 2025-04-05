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
	// Common Go duration suffixes: "ns", "us" (or "µs"), "ms", "s", "m", "h"
	// Try Go's standard time.ParseDuration first
	duration, err := time.ParseDuration(durationStr)
	if err == nil {
		return duration, nil
	}

	// For more complex formats, use regex to extract value and unit
	re := regexp.MustCompile(`^([\d\.]+)\s*([a-zA-Z]+)$`)
	matches := re.FindStringSubmatch(durationStr)

	if len(matches) != 3 {
		return 0, fmt.Errorf("unsupported duration format: %s", durationStr)
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
		return time.Duration(value * float64(time.Nanosecond)), nil
	case "us", "µs", "microsecond", "microseconds":
		return time.Duration(value * float64(time.Microsecond)), nil
	case "ms", "millisecond", "milliseconds":
		return time.Duration(value * float64(time.Millisecond)), nil
	case "s", "sec", "second", "seconds":
		return time.Duration(value * float64(time.Second)), nil
	case "m", "min", "minute", "minutes":
		return time.Duration(value * float64(time.Minute)), nil
	case "h", "hr", "hour", "hours":
		return time.Duration(value * float64(time.Hour)), nil
	case "d", "day", "days":
		return time.Duration(value * 24 * float64(time.Hour)), nil
	default:
		return 0, fmt.Errorf("unknown time unit: %s", unit)
	}
}

