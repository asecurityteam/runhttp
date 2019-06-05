package runhttp

import (
	"context"
	"runtime"
	"sync"
	"time"
)

const (
	statMemstatsAlloc        = "go_expvar.memstats.alloc"
	statMemstatsFrees        = "go_expvar.memstats.frees"
	statMemstatsHeapAlloc    = "go_expvar.memstats.heap_alloc"
	statMemstatsHeapIdle     = "go_expvar.memstats.heap_idle"
	statMemstatsHeapInuse    = "go_expvar.memstats.heap_inuse"
	statMemstatsHeapObjects  = "go_expvar.memstats.heap_objects"
	statMemstatsHeapReleased = "go_expvar.memstats.heap_released"
	statMemstatsHeapSys      = "go_expvar.memstats.heap_sys"
	statMemstatsLookups      = "go_expvar.memstats.lookups"
	statMemstatsMallocs      = "go_expvar.memstats.mallocs"
	statMemstatsNumGC        = "go_expvar.memstats.num_gc"
	statMemstatsPauseNS      = "go_expvar.memstats.pause_ns"
	statMemstatsPauseTotalNS = "go_expvar.memstats.pause_total_ns"
	statMemstatsTotalAlloc   = "go_expvar.memstats.total_alloc"
	statGoroutinesExists     = "go_expvar.goroutines.exists"
	expvarInterval           = 5 * time.Second
)

// Expvar tracks the memory usage of the Go runtime and collects metrics instrumented from Go’s expvar package
type Expvar struct {
	Stat                     Stat
	MemstatsAllocName        string
	MemstatsFreesName        string
	MemstatsHeapAllocName    string
	MemstatsHeapIdleName     string
	MemstatsHeapInuseName    string
	MemstatsHeapObjectsName  string
	MemstatsHeapReleasedName string
	MemstatsHeapSysName      string
	MemstatsLookupsName      string
	MemstatsMallocsName      string
	MemstatsNumGCName        string
	MemstatsPauseNSName      string
	MemstatsPauseTotalNSName string
	MemstatsTotalAllocName   string
	GoroutinesExistsName     string
	Interval                 time.Duration
	statMut                  *sync.Mutex
	stopMut                  *sync.Mutex
	stop                     bool
}

// Report loops on a time interval and pushes a set of gauge metrics.
func (e *Expvar) Report() {
	ticker := time.NewTicker(e.Interval)
	defer ticker.Stop()
	for range ticker.C {
		e.report()
		e.stopMut.Lock()
		if e.stop {
			e.stopMut.Unlock()
			return
		}
		e.stopMut.Unlock()
	}
}

// Close the reporting loop.
func (e *Expvar) Close() {
	e.stopMut.Lock()
	defer e.stopMut.Unlock()
	e.stop = true
}

func (e *Expvar) report() {
	memstats := new(runtime.MemStats)
	runtime.ReadMemStats(memstats)
	numGoroutines := runtime.NumGoroutine()

	e.statMut.Lock()
	defer e.statMut.Unlock()

	e.Stat.Gauge(e.MemstatsAllocName, float64(memstats.Alloc))
	e.Stat.Gauge(e.MemstatsFreesName, float64(memstats.Frees))
	e.Stat.Gauge(e.MemstatsHeapAllocName, float64(memstats.HeapAlloc))
	e.Stat.Gauge(e.MemstatsHeapIdleName, float64(memstats.HeapIdle))
	e.Stat.Gauge(e.MemstatsHeapInuseName, float64(memstats.HeapInuse))
	e.Stat.Gauge(e.MemstatsHeapObjectsName, float64(memstats.HeapObjects))
	e.Stat.Gauge(e.MemstatsHeapReleasedName, float64(memstats.HeapReleased))
	e.Stat.Gauge(e.MemstatsHeapSysName, float64(memstats.HeapSys))
	e.Stat.Gauge(e.MemstatsLookupsName, float64(memstats.Lookups))
	e.Stat.Gauge(e.MemstatsMallocsName, float64(memstats.Mallocs))
	e.Stat.Gauge(e.MemstatsNumGCName, float64(memstats.NumGC))
	e.Stat.Gauge(e.MemstatsPauseTotalNSName, float64(memstats.PauseTotalNs))
	e.Stat.Gauge(e.MemstatsTotalAllocName, float64(memstats.TotalAlloc))
	e.Stat.Gauge(e.GoroutinesExistsName, float64(numGoroutines))

	// TODO PauseNS
}

// ExpvarConfig is a container for internal expvar metrics settings.
type ExpvarConfig struct {
	Alloc            string        `description:"Name of the metric tracking allocated bytes"`
	Frees            string        `description:"Name of the metric tracking number of frees"`
	HeapAlloc        string        `description:"Name of the metric tracking allocated bytes"`
	HeapIdle         string        `description:"Name of the metric tracking bytes in unused spans"`
	HeapInuse        string        `description:"Name of the metric tracking bytes in in-use spans"`
	HeapObjects      string        `description:"Name of the metric tracking total number of object allocated"`
	HeapReleased     string        `description:"Name of the metric tracking bytes realeased to the OS"`
	HeapSys          string        `description:"Name of the metric tracking bytes obtained from the system"`
	Lookups          string        `description:"Name of the metric tracking number of pointer lookups"`
	Mallocs          string        `description:"Name of the metric tracking number of mallocs"`
	NumGC            string        `description:"Name of the metric tracking number of garbage collections"`
	PauseNS          string        `description:"Name of the metric tracking duration of GC pauses"`
	PauseTotalNS     string        `description:"Name of the metric tracking total GC pause duration over lifetime process"`
	TotalAlloc       string        `description:"Name of the metric tracking allocated bytes (even if freed)"`
	GoroutinesExists string        `description:"Name of the metric tracking number of active go routines"`
	ReportInterval   time.Duration `description:"Interval on which metrics are reported."`
}

// Name of the configuration root.
func (*ExpvarConfig) Name() string {
	return "expvar"
}

// Description returns the help information for the configuration root.
func (*ExpvarConfig) Description() string {
	return "Expvar metric names."
}

// ExpvarComponent implements the settings.Component interface for expvar metrics.
type ExpvarComponent struct{}

// Settings returns a configuration with all defaults set.
func (*ExpvarComponent) Settings() *ExpvarConfig {
	return &ExpvarConfig{
		Alloc:            statMemstatsAlloc,
		Frees:            statMemstatsFrees,
		HeapAlloc:        statMemstatsHeapAlloc,
		HeapIdle:         statMemstatsHeapIdle,
		HeapInuse:        statMemstatsHeapInuse,
		HeapObjects:      statMemstatsHeapObjects,
		HeapReleased:     statMemstatsHeapReleased,
		HeapSys:          statMemstatsHeapSys,
		Lookups:          statMemstatsLookups,
		Mallocs:          statMemstatsMallocs,
		NumGC:            statMemstatsNumGC,
		PauseNS:          statMemstatsPauseNS,
		PauseTotalNS:     statMemstatsPauseTotalNS,
		TotalAlloc:       statMemstatsTotalAlloc,
		GoroutinesExists: statGoroutinesExists,
		ReportInterval:   expvarInterval,
	}
}

// New produces a ServerFn bound to the given configuration.
func (*ExpvarComponent) New(_ context.Context, conf *ExpvarConfig) (func() *Expvar, error) {
	return func() *Expvar {
		return &Expvar{
			MemstatsAllocName:        conf.Alloc,
			MemstatsFreesName:        conf.Frees,
			MemstatsHeapAllocName:    conf.HeapAlloc,
			MemstatsHeapIdleName:     conf.HeapIdle,
			MemstatsHeapInuseName:    conf.HeapInuse,
			MemstatsHeapObjectsName:  conf.HeapObjects,
			MemstatsHeapReleasedName: conf.HeapReleased,
			MemstatsHeapSysName:      conf.HeapSys,
			MemstatsLookupsName:      conf.Lookups,
			MemstatsMallocsName:      conf.Mallocs,
			MemstatsNumGCName:        conf.NumGC,
			MemstatsPauseNSName:      conf.PauseNS,
			MemstatsPauseTotalNSName: conf.PauseTotalNS,
			MemstatsTotalAllocName:   conf.TotalAlloc,
			GoroutinesExistsName:     conf.GoroutinesExists,
			Interval:                 conf.ReportInterval,
			statMut:                  &sync.Mutex{},
			stopMut:                  &sync.Mutex{},
		}
	}, nil

}