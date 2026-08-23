package service

import (
	"fmt"
	"predictive-maintenance/internal/detector"
	"predictive-maintenance/internal/model"
	"sync"
)

type Rule struct {
	ID         uint64          `json:"id"`
	Name       string          `json:"name"`
	SensorType string          `json:"sensorType"`
	Config     detector.Config `json:"config"`
	Enabled    bool            `json:"enabled"`
}
type RuleService struct {
	mu    sync.RWMutex
	rules map[uint64]Rule
	next  uint64
}

func NewRuleService() *RuleService { return &RuleService{rules: map[uint64]Rule{}} }
func (s *RuleService) Create(r Rule) (Rule, error) {
	if !model.ValidSensorType(r.SensorType) {
		return r, fmt.Errorf("invalid sensor type")
	}
	if r.Config.MaxValue <= r.Config.MinValue {
		return r, fmt.Errorf("invalid range")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next++
	r.ID = s.next
	r.Enabled = true
	s.rules[r.ID] = r
	return r, nil
}
func (s *RuleService) Get(id uint64) (Rule, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.rules[id]
	return r, ok
}
func (s *RuleService) List() []Rule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Rule, 0, len(s.rules))
	for _, r := range s.rules {
		out = append(out, r)
	}
	return out
}
func (s *RuleService) Enable(id uint64, enabled bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.rules[id]
	if !ok {
		return false
	}
	r.Enabled = enabled
	s.rules[id] = r
	return true
}
func (s *RuleService) Evaluate(r Rule, sample detector.Sample, history []detector.Sample) detector.Result {
	if !r.Enabled {
		return detector.Result{Score: 100, Rule: "disabled"}
	}
	return detector.CompositeDetector{Detectors: []detector.Detector{detector.ThresholdDetector{}, detector.RateDetector{}}}.Detect(sample, history, r.Config)
}
