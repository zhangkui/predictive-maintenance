package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

type Metrics struct {
	requests  uint64
	errors    uint64
	ingested  uint64
	anomalies uint64
	started   time.Time
	mu        sync.RWMutex
	latencies []time.Duration
}

func New() *Metrics             { return &Metrics{started: time.Now()} }
func (m *Metrics) Request()     { atomic.AddUint64(&m.requests, 1) }
func (m *Metrics) Error()       { atomic.AddUint64(&m.errors, 1) }
func (m *Metrics) Ingest(n int) { atomic.AddUint64(&m.ingested, uint64(n)) }
func (m *Metrics) Anomaly()     { atomic.AddUint64(&m.anomalies, 1) }
func (m *Metrics) Observe(d time.Duration) {
	m.mu.Lock()
	m.latencies = append(m.latencies, d)
	if len(m.latencies) > 1000 {
		m.latencies = m.latencies[len(m.latencies)-1000:]
	}
	m.mu.Unlock()
}
func (m *Metrics) Snapshot() map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var total time.Duration
	for _, d := range m.latencies {
		total += d
	}
	avg := int64(0)
	if len(m.latencies) > 0 {
		avg = total.Milliseconds() / int64(len(m.latencies))
	}
	return map[string]any{"requests": atomic.LoadUint64(&m.requests), "errors": atomic.LoadUint64(&m.errors), "ingested": atomic.LoadUint64(&m.ingested), "anomalies": atomic.LoadUint64(&m.anomalies), "avgLatencyMs": avg, "uptimeSeconds": int(time.Since(m.started).Seconds())}
}
