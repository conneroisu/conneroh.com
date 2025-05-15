// Package cwg contains a WaitGroup with debugging capabilities.
package cwg

import (
	"fmt"
	"runtime"
	"sync"
)

// DebugEntry represents a single Add operation on the WaitGroup.
type DebugEntry struct {
	location string
	delta    int
	active   bool // Whether this entry is still active (not Done)
	id       int  // Unique ID for tracking
}

// DebugWaitGroup wraps sync.WaitGroup to add debugging capabilities.
type DebugWaitGroup struct {
	wg        sync.WaitGroup
	mutex     sync.Mutex
	entries   []*DebugEntry
	nextID    int         // For generating unique IDs
	activeMap map[int]int // Maps ID to index in entries
}

// Add adds delta to the WaitGroup counter and records the caller.
func (dwg *DebugWaitGroup) Add(delta int) {
	dwg.mutex.Lock()
	defer dwg.mutex.Unlock()

	// Initialize the active map if needed
	if dwg.activeMap == nil {
		dwg.activeMap = make(map[int]int)
	}

	// Record the caller's information
	_, file, line, _ := runtime.Caller(1)
	location := fmt.Sprintf("%s:%d", file, line)

	// Create a new entry
	id := dwg.nextID
	dwg.nextID++

	entry := &DebugEntry{
		location: location,
		delta:    delta,
		active:   true,
		id:       id,
	}

	// Add to our tracking
	dwg.entries = append(dwg.entries, entry)
	entryIndex := len(dwg.entries) - 1

	// Track active entries
	for range delta {
		subID := dwg.nextID
		dwg.nextID++
		dwg.activeMap[subID] = entryIndex
	}

	// Call the actual WaitGroup Add
	dwg.wg.Add(delta)
}

// Done decrements the WaitGroup counter by one and marks one entry as done.
func (dwg *DebugWaitGroup) Done() {
	dwg.mutex.Lock()
	defer dwg.mutex.Unlock()

	// Find an active entry to mark as done
	var keyToRemove int
	for k := range dwg.activeMap {
		keyToRemove = k

		break
	}

	// Remove it from the active map
	if len(dwg.activeMap) > 0 {
		delete(dwg.activeMap, keyToRemove)
	}

	// Call the actual WaitGroup Done
	dwg.wg.Done()
}

// Wait blocks until the WaitGroup counter is zero.
func (dwg *DebugWaitGroup) Wait() {
	dwg.wg.Wait()
}

// PrintActiveDebugInfo prints only the active (not Done) entries.
func (dwg *DebugWaitGroup) PrintActiveDebugInfo() {
	dwg.mutex.Lock()
	defer dwg.mutex.Unlock()

	fmt.Printf("WaitGroup has %d active goroutines that were started at:\n", len(dwg.activeMap))

	// Create a frequency map of active entries
	activeEntries := make(map[int]int) // entryIndex -> count
	for _, entryIndex := range dwg.activeMap {
		activeEntries[entryIndex]++
	}

	// Print active entries
	count := 1
	for entryIndex, frequency := range activeEntries {
		entry := dwg.entries[entryIndex]
		fmt.Printf("  %d. %s (remaining: %d)\n", count, entry.location, frequency)
		count++
	}
}

// Count returns the number of active entries.
func (dwg *DebugWaitGroup) Count() int {
	return len(dwg.entries)
}
