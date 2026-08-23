package service

import (
	"context"
	"fmt"
	"predictive-maintenance/internal/model"
	"predictive-maintenance/internal/repository"
)

type AnomalyService struct {
	Repo *repository.AbnormalRepository
}

func NewAnomalyService(r *repository.AbnormalRepository) *AnomalyService { return &AnomalyService{r} }
func (s *AnomalyService) ListPending(ctx context.Context, device uint64) ([]model.AbnormalRecord, error) {
	return s.Repo.Pending(ctx, device)
}
func (s *AnomalyService) Handle(ctx context.Context, id, handler uint64) error {
	if handler == 0 {
		return fmt.Errorf("handler is required")
	}
	return s.Repo.Handle(ctx, id, handler, model.AbnormalHandled)
}
func (s *AnomalyService) Ignore(ctx context.Context, id, handler uint64) error {
	if handler == 0 {
		return fmt.Errorf("handler is required")
	}
	return s.Repo.Handle(ctx, id, handler, 3)
}
func (s *AnomalyService) Priority(severity int) int {
	switch severity {
	case model.SeverityCritical:
		return model.PriorityEmergency
	case model.SeveritySerious:
		return model.PriorityHigh
	case model.SeverityWarning:
		return model.PriorityHigh
	default:
		return model.PriorityLow
	}
}
