package system

import (
	"fmt"
	"runtime"
	"time"

	"github.com/qtoad/mygo-admin/modules/config"
	"github.com/qtoad/mygo-admin/modules/language"
	"github.com/qtoad/mygo-admin/modules/util"
)

var (
	startTime = time.Now()
)

type AppStatus struct {
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

func GetAppStatus() AppStatus {
	var app AppStatus
	app.Uptime = util.TimeSincePro(startTime, language.Lang[config.GetLanguage()])

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	app.NumGoroutine = runtime.NumGoroutine()

	app.MemAllocated = util.FileSize(m.Alloc)
	app.MemTotal = util.FileSize(m.TotalAlloc)
	app.MemSys = util.FileSize(m.Sys)
	app.Lookups = m.Lookups
	app.MemMallocs = m.Mallocs
	app.MemFrees = m.Frees

	app.HeapAlloc = util.FileSize(m.HeapAlloc)
	app.HeapSys = util.FileSize(m.HeapSys)
	app.HeapIdle = util.FileSize(m.HeapIdle)
	app.HeapInuse = util.FileSize(m.HeapInuse)
	app.HeapReleased = util.FileSize(m.HeapReleased)
	app.HeapObjects = m.HeapObjects

	app.StackInuse = util.FileSize(m.StackInuse)
	app.StackSys = util.FileSize(m.StackSys)
	app.MSpanInuse = util.FileSize(m.MSpanInuse)
	app.MSpanSys = util.FileSize(m.MSpanSys)
	app.MCacheInuse = util.FileSize(m.MCacheInuse)
	app.MCacheSys = util.FileSize(m.MCacheSys)
	app.BuckHashSys = util.FileSize(m.BuckHashSys)
	app.GCSys = util.FileSize(m.GCSys)
	app.OtherSys = util.FileSize(m.OtherSys)

	app.NextGC = util.FileSize(m.NextGC)
	app.LastGC = fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000)
	app.PauseTotalNs = fmt.Sprintf("%.1fs", float64(m.PauseTotalNs)/1000/1000/1000)
	app.PauseNs = fmt.Sprintf("%.3fs", float64(m.PauseNs[(m.NumGC+255)%256])/1000/1000/1000)
	app.NumGC = m.NumGC

	return app
}

type SysStatus struct {
	CpuLogicalCore int
	CpuCore        int
	OSPlatform     string
	OSFamily       string
	OSVersion      string
	Load1          float64
	Load5          float64
	Load15         float64
	MemTotal       string
	MemAvailable   string
	MemUsed        string
}
