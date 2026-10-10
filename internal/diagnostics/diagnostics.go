package diagnostics

import (
	"runtime"
	"time"
)

type Snapshot struct {
	GoVersion     string `json:"go_version"`
	GOOS          string `json:"os"`
	GOARCH        string `json:"arch"`
	Goroutines    int    `json:"goroutines"`
	HeapAlloc     uint64 `json:"heap_alloc_bytes"`
	SystemMemory  uint64 `json:"go_sys_bytes"`
	UptimeSeconds int64  `json:"uptime_seconds"`
}

func SnapshotNow(started time.Time) Snapshot {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return Snapshot{GoVersion: runtime.Version(), GOOS: runtime.GOOS, GOARCH: runtime.GOARCH, Goroutines: runtime.NumGoroutine(), HeapAlloc: m.HeapAlloc, SystemMemory: m.Sys, UptimeSeconds: int64(time.Since(started).Seconds())}
}
