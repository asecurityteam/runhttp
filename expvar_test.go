package runhttp

import (
	"context"
	"math/rand"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigName(t *testing.T) {
	assert.Equal(t, "expvar", (&ExpvarConfig{}).Name())
}

func TestConfigDescription(t *testing.T) {
	assert.Equal(t, "Expvar metric names", (&ExpvarConfig{}).Description())
}

func TestReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := &runtime.MemStats{
		Alloc:        randGen.Uint64(),
		Frees:        randGen.Uint64(),
		HeapAlloc:    randGen.Uint64(),
		HeapIdle:     randGen.Uint64(),
		HeapInuse:    randGen.Uint64(),
		HeapObjects:  randGen.Uint64(),
		HeapReleased: randGen.Uint64(),
		HeapSys:      randGen.Uint64(),
		Lookups:      randGen.Uint64(),
		Mallocs:      randGen.Uint64(),
		NumGC:        randGen.Uint32(),
		PauseTotalNs: randGen.Uint64(),
		TotalAlloc:   randGen.Uint64(),
	}

	routines := randGen.Int()

	mockStats := NewMockStat(ctrl)

	gomock.InOrder(
		mockStats.EXPECT().Gauge(statMemstatsAlloc, float64(ms.Alloc)).Times(1),
		mockStats.EXPECT().Gauge(statMemstatsFrees, float64(ms.Frees)).Times(1),
		mockStats.EXPECT().Gauge(statMemstatsHeapAlloc, float64(ms.HeapAlloc)).Times(1),
		mockStats.EXPECT().Gauge(statMemstatsHeapIdle, float64(ms.HeapIdle)).Times(1),
		mockStats.EXPECT().Gauge(statMemstatsHeapInuse, float64(ms.HeapInuse)).Times(1),
		mockStats.EXPECT().Gauge(statMemstatsHeapObjects, float64(ms.HeapObjects)).Times(1),
		mockStats.EXPECT().Gauge(statMemstatsHeapReleased, float64(ms.HeapReleased)).Times(1),
		mockStats.EXPECT().Gauge(statMemstatsHeapSys, float64(ms.HeapSys)).Times(1),
		mockStats.EXPECT().Gauge(statMemstatsLookups, float64(ms.Lookups)).Times(1),
		mockStats.EXPECT().Gauge(statMemstatsMallocs, float64(ms.Mallocs)).Times(1),
		mockStats.EXPECT().Gauge(statMemstatsNumGC, float64(ms.NumGC)).Times(1),
		mockStats.EXPECT().Gauge(statMemstatsPauseTotalNS, float64(ms.PauseTotalNs)).Times(1),
		mockStats.EXPECT().Gauge(statMemstatsTotalAlloc, float64(ms.TotalAlloc)).Times(1),
		mockStats.EXPECT().Gauge(statGoroutinesExists, float64(routines)).Times(1),
	)

	expvar := newExpvar(t, fakeReadMemstats(ms), fakeNumGoroutine(routines))
	expvar.Stat = mockStats
	expvar.lastNumGC = ms.NumGC // do not collect pauseGC
	expvar.report()
}

func TestPauseGC(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pauseNS := [256]uint64{}
	for i := 0; i < 256; i++ {
		pauseNS[i] = uint64(i)
	}

	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := &runtime.MemStats{
		NumGC:   uint32(randGen.Int()),
		PauseNs: pauseNS,
	}

	routines := randGen.Int()

	mockStats := NewMockStat(ctrl)

	mockStats.EXPECT().Gauge(gomock.Any(), gomock.Any()).AnyTimes()
	mockStats.EXPECT().Histogram(statMemstatsPauseNS, gomock.Any()).Times(int((ms.NumGC+255)%256) + 1)

	expvar := newExpvar(t, fakeReadMemstats(ms), fakeNumGoroutine(routines))
	expvar.Stat = mockStats
	expvar.report()
}

func TestPauseGCWithWrap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pauseNS := [256]uint64{}
	for i := 0; i < 256; i++ {
		pauseNS[i] = uint64(i)
	}

	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := &runtime.MemStats{
		NumGC:   uint32(randGen.Int()),
		PauseNs: pauseNS,
	}

	routines := randGen.Int()

	mockStats := NewMockStat(ctrl)

	mockStats.EXPECT().Gauge(gomock.Any(), gomock.Any()).AnyTimes()
	mockStats.EXPECT().Histogram(statMemstatsPauseNS, gomock.Any()).Times(int((ms.NumGC+255)%256) + 2)

	expvar := newExpvar(t, fakeReadMemstats(ms), fakeNumGoroutine(routines))
	expvar.Stat = mockStats
	expvar.lastNumGC = 255
	expvar.report()
}

func newExpvar(t *testing.T, fakeReadMemStats func(*runtime.MemStats), fakeNumGoroutine func() int) *Expvar {
	c := &ExpvarComponent{}
	expvarFn, err := c.New(context.Background(), c.Settings())
	require.NoError(t, err)
	expvar := expvarFn()
	expvar.numGoroutine = fakeNumGoroutine
	expvar.readMemStats = fakeReadMemStats
	return expvar
}

func fakeNumGoroutine(n int) func() int {
	return func() int {
		return n
	}
}

func fakeReadMemstats(ms *runtime.MemStats) func(*runtime.MemStats) {
	return func(ms2 *runtime.MemStats) {
		src := reflect.Indirect(reflect.ValueOf(ms))
		dst := reflect.ValueOf(ms2).Elem()
		for i := 0; i < src.NumField(); i++ {
			dstField := dst.Field(i)
			if dstField.CanSet() {
				dstField.Set(reflect.ValueOf(src.Field(i).Interface()))
			}
		}
	}
}
