package web

import (
	"fmt"
	"runtime"
	"time"

	"github.com/labstack/echo/v4"

	"chain-demo/cloudserver/hq-echo/util/conv"
)

func DashboardHandler(c echo.Context) error {
	updateSystemStatus()

	c.Set("tmpl", "web/dashboard")
	c.Set("data", map[string]interface{}{
		"title":     "Dashboard",
		"SysStatus": sysStatus,
	})

	return nil
}

var (
	startTime = time.Now()
)

var sysStatus struct {
	Uptime       string
	NumGoroutine int

	// General statistics.
	MemAllocated string // bytes allocated and still in use
	MemTotal     string // bytes allocated (even if freed)
	MemSys       string // bytes obtained from system (sum of XxxSys below)
	Lookups      uint64 // number of pointer lookups
	MemMallocs   uint64 // number of mallocs
	MemFrees     uint64 // number of frees

	// Main allocation heap statistics.
	HeapAlloc    string // bytes allocated and still in use
	HeapSys      string // bytes obtained from system
	HeapIdle     string // bytes in idle spans
	HeapInuse    string // bytes in non-idle span
	HeapReleased string // bytes released to the OS
	HeapObjects  uint64 // total number of allocated objects

	// Low-level fixed-size structure allocator statistics.
	//	Inuse is bytes used now.
	//	Sys is bytes obtained from system.
	StackInuse  string // bootstrap stacks
	StackSys    string
	MSpanInuse  string // mspan structures
	MSpanSys    string
	MCacheInuse string // mcache structures
	MCacheSys   string
	BuckHashSys string // profiling bucket hash table
	GCSys       string // GC metadata
	OtherSys    string // other system allocations

	// Garbage collector statistics.
	NextGC       string // next run in HeapAlloc time (bytes)
	LastGC       string // last run in absolute time (ns)
	PauseTotalNs string
	PauseNs      string // circular buffer of recent GC pause times, most recent at [(NumGC+255)%256]
	NumGC        uint32
}

func updateSystemStatus() {
	sysStatus.Uptime = conv.TimeSincePro(startTime)

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	sysStatus.NumGoroutine = runtime.NumGoroutine()

	sysStatus.MemAllocated = conv.FileSize(int64(m.Alloc))
	sysStatus.MemTotal = conv.FileSize(int64(m.TotalAlloc))
	sysStatus.MemSys = conv.FileSize(int64(m.Sys))
	sysStatus.Lookups = m.Lookups
	sysStatus.MemMallocs = m.Mallocs
	sysStatus.MemFrees = m.Frees

	sysStatus.HeapAlloc = conv.FileSize(int64(m.HeapAlloc))
	sysStatus.HeapSys = conv.FileSize(int64(m.HeapSys))
	sysStatus.HeapIdle = conv.FileSize(int64(m.HeapIdle))
	sysStatus.HeapInuse = conv.FileSize(int64(m.HeapInuse))
	sysStatus.HeapReleased = conv.FileSize(int64(m.HeapReleased))
	sysStatus.HeapObjects = m.HeapObjects

	sysStatus.StackInuse = conv.FileSize(int64(m.StackInuse))
	sysStatus.StackSys = conv.FileSize(int64(m.StackSys))
	sysStatus.MSpanInuse = conv.FileSize(int64(m.MSpanInuse))
	sysStatus.MSpanSys = conv.FileSize(int64(m.MSpanSys))
	sysStatus.MCacheInuse = conv.FileSize(int64(m.MCacheInuse))
	sysStatus.MCacheSys = conv.FileSize(int64(m.MCacheSys))
	sysStatus.BuckHashSys = conv.FileSize(int64(m.BuckHashSys))
	sysStatus.GCSys = conv.FileSize(int64(m.GCSys))
	sysStatus.OtherSys = conv.FileSize(int64(m.OtherSys))

	sysStatus.NextGC = conv.FileSize(int64(m.NextGC))
	sysStatus.LastGC = fmt.Sprintf("%.3fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000)
	sysStatus.PauseTotalNs = fmt.Sprintf("%.3fs", float64(m.PauseTotalNs)/1000/1000/1000)
	sysStatus.PauseNs = fmt.Sprintf("%.6fs", float64(m.PauseNs[(m.NumGC+255)%256])/1000/1000/1000)
	sysStatus.NumGC = m.NumGC
}
