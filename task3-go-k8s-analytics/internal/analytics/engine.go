package analytics

import (
	"math"
	"sync"
	"time"

	"task3-go-k8s-analytics/internal/model"
)

type Snapshot struct {
	RollingAverageCPU float64   `json:"rolling_avg_cpu"`
	RollingAverageRPS float64   `json:"rolling_avg_rps"`
	RollingLatencyMs  float64   `json:"rolling_latency_ms"`
	LastZScoreCPU     float64   `json:"last_zscore_cpu"`
	LastZScoreRPS     float64   `json:"last_zscore_rps"`
	LastZScoreLatency float64   `json:"last_zscore_latency"`
	Anomaly           bool      `json:"anomaly"`
	SampleSize        int       `json:"sample_size"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Engine struct {
	windowSize int
	metrics    []model.Metric
	mu         sync.RWMutex
	threshold  float64
}

func NewEngine(windowSize int, threshold float64) *Engine {
	return &Engine{
		windowSize: windowSize,
		threshold:  threshold,
		metrics:    make([]model.Metric, 0, windowSize),
	}
}

func (e *Engine) Add(m model.Metric) Snapshot {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.metrics = append(e.metrics, m)
	if len(e.metrics) > e.windowSize {
		e.metrics = e.metrics[len(e.metrics)-e.windowSize:]
	}
	return e.snapshotLocked()
}

func (e *Engine) Snapshot() Snapshot {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.snapshotLocked()
}

func (e *Engine) snapshotLocked() Snapshot {
	count := len(e.metrics)
	if count == 0 {
		return Snapshot{SampleSize: 0, UpdatedAt: time.Now()}
	}

	var sumCPU, sumRPS, sumLatency float64
	for _, m := range e.metrics {
		sumCPU += m.CPU
		sumRPS += m.RPS
		sumLatency += m.LatencyMs
	}
	avgCPU := sumCPU / float64(count)
	avgRPS := sumRPS / float64(count)
	avgLatency := sumLatency / float64(count)

	stdCPU, stdRPS, stdLatency := stdDev(e.metrics, avgCPU, avgRPS, avgLatency)
	last := e.metrics[count-1]

	zCPU := zScore(last.CPU, avgCPU, stdCPU)
	zRPS := zScore(last.RPS, avgRPS, stdRPS)
	zLatency := zScore(last.LatencyMs, avgLatency, stdLatency)

	anomaly := math.Abs(zCPU) > e.threshold || math.Abs(zRPS) > e.threshold || math.Abs(zLatency) > e.threshold

	return Snapshot{
		RollingAverageCPU: avgCPU,
		RollingAverageRPS: avgRPS,
		RollingLatencyMs:  avgLatency,
		LastZScoreCPU:     zCPU,
		LastZScoreRPS:     zRPS,
		LastZScoreLatency: zLatency,
		Anomaly:           anomaly,
		SampleSize:        count,
		UpdatedAt:         time.Now(),
	}
}

func stdDev(metrics []model.Metric, meanCPU, meanRPS, meanLatency float64) (cpu, rps, latency float64) {
	n := float64(len(metrics))
	if n == 0 {
		return 0, 0, 0
	}
	for _, m := range metrics {
		cpu += math.Pow(m.CPU-meanCPU, 2)
		rps += math.Pow(m.RPS-meanRPS, 2)
		latency += math.Pow(m.LatencyMs-meanLatency, 2)
	}
	return math.Sqrt(cpu / n), math.Sqrt(rps / n), math.Sqrt(latency / n)
}

func zScore(value, mean, std float64) float64 {
	if std == 0 {
		return 0
	}
	return (value - mean) / std
}
