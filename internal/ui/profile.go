package ui

import (
	"fmt"
	"time"
)

// Simple profiling helper for debugging
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("[PROFILE] %s took %v\n", name, elapsed)
}
