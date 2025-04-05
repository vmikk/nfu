package main

import (
	"bufio"
	"fmt"
	"os"
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

