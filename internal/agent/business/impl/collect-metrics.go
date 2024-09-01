package impl

import (
	"math/big"
	"math/rand"
	"runtime"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/agent/business"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/agent/constants"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shopspring/decimal"
)

func (s *svc) collectMetrics() {
	metrics := make([]Metric, 0, constants.MetricsCount)

	metrics = append(metrics, s.getMetrics()...)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.pollCount = s.pollCount.Add(decimal.NewFromInt(1))

	metrics = append(metrics, Metric{
		ID:   "PollCount",
		Type: string(business.CounterMT),
		Val:  s.pollCount,
	})

	s.metrics = metrics
}

func (s *svc) getMetrics() []Metric {
	metrics := make([]Metric, 0, constants.MetricsCount)

	var rms runtime.MemStats
	runtime.ReadMemStats(&rms)

	metrics = append(
		metrics,
		[]Metric{
			{
				ID:   "Alloc",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.Alloc), 0),
			},
			{
				ID:   "BuckHashSys",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.BuckHashSys), 0),
			},
			{
				ID:   "Frees",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.Frees), 0),
			},
			{
				ID:   "GCCPUFraction",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromFloat(rms.GCCPUFraction),
			},
			{
				ID:   "GCSys",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.GCSys), 0),
			},
			{
				ID:   "HeapAlloc",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.HeapAlloc), 0),
			},
			{
				ID:   "HeapIdle",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.HeapIdle), 0),
			},
			{
				ID:   "HeapInuse",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.HeapInuse), 0),
			},
			{
				ID:   "HeapObjects",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.HeapObjects), 0),
			},
			{
				ID:   "HeapReleased",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.HeapReleased), 0),
			},
			{
				ID:   "HeapSys",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.HeapSys), 0),
			},
			{
				ID:   "LastGC",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.LastGC), 0),
			},
			{
				ID:   "Lookups",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.Lookups), 0),
			},
			{
				ID:   "MCacheInuse",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.MCacheInuse), 0),
			},
			{
				ID:   "MCacheSys",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.MCacheSys), 0),
			},
			{
				ID:   "MSpanInuse",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.MSpanInuse), 0),
			},
			{
				ID:   "MSpanSys",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.MSpanSys), 0),
			},
			{
				ID:   "Mallocs",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.Mallocs), 0),
			},
			{
				ID:   "NextGC",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.NextGC), 0),
			},
			{
				ID:   "NumForcedGC",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(uint64(rms.NumForcedGC)), 0),
			},
			{
				ID:   "NumGC",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(uint64(rms.NumGC)), 0),
			},
			{
				ID:   "OtherSys",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.OtherSys), 0),
			},
			{
				ID:   "PauseTotalNs",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.PauseTotalNs), 0),
			},
			{
				ID:   "StackInuse",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.StackInuse), 0),
			},
			{
				ID:   "StackSys",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.StackSys), 0),
			},
			{
				ID:   "Sys",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.Sys), 0),
			},
			{
				ID:   "TotalAlloc",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(rms.TotalAlloc), 0),
			},
			{
				ID:   "RandomValue",
				Type: string(business.GaugeMT),
				Val:  decimal.NewFromFloat(rand.Float64()),
			},
		}...,
	)

	return metrics
}

func (s *svc) collectExtraMetrics() {
	metrics := make([]Metric, 0, constants.ExtraMetricsCount)

	v, err := mem.VirtualMemory()
	if err != nil {
		s.logger.Errorw(
			"can't get virtual memory",
			"reason", err,
		)
		return
	}

	cpuPercentages, err := cpu.Percent(0, false)
	if err != nil {
		s.logger.Errorw(
			"can't get CPU info",
			"reason", err,
		)
		return
	}

	metrics = append(metrics, []Metric{
		{
			ID:   "TotalMemory",
			Type: string(business.GaugeMT),
			Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(v.Total/1024/1024), 0),
		},
		{
			ID:   "FreeMemory",
			Type: string(business.GaugeMT),
			Val:  decimal.NewFromBigInt(new(big.Int).SetUint64(v.Free/1024/1024), 0),
		},
		{
			ID:   "CPUutilization1",
			Type: string(business.GaugeMT),
			Val:  decimal.NewFromFloat(cpuPercentages[0]),
		},
	}...)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.extraMetrics = metrics
}
