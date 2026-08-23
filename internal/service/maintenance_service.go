package service

import (
	"context"
	"fmt"
	"predictive-maintenance/internal/model"
	"predictive-maintenance/internal/repository"
	"predictive-maintenance/internal/util"
	"time"
)

type MaintenanceService struct {
	Repo *repository.MaintenanceRepository
}

func NewMaintenanceService(r *repository.MaintenanceRepository) *MaintenanceService {
	return &MaintenanceService{r}
}
func (s *MaintenanceService) CreatePlan(ctx context.Context, p model.MaintenancePlan) (uint64, error) {
	if p.DeviceID == 0 || p.Type == "" {
		return 0, fmt.Errorf("deviceId and type are required")
	}
	if p.Type == "periodic" && p.CycleValue <= 0 {
		return 0, fmt.Errorf("cycleValue must be positive")
	}
	return s.Repo.CreatePlan(ctx, p)
}
func (s *MaintenanceService) GenerateTasks(ctx context.Context) (int, error) {
	plans, e := s.Repo.DuePlans(ctx)
	if e != nil {
		return 0, e
	}
	created := 0
	for _, p := range plans {
		task := model.MaintenanceTask{PlanID: p.ID, DeviceID: p.DeviceID, Title: "scheduled maintenance", Description: p.Description, Priority: model.PriorityMedium, Status: model.TaskPending, ScheduledDate: util.NowUTC()}
		if _, e = s.Repo.CreateTask(ctx, task); e != nil {
			return created, e
		}
		created++
	}
	return created, nil
}
func (s *MaintenanceService) NextDate(last time.Time, cycle string, value int) time.Time {
	return util.AddCycle(last, cycle, value)
}
func (s *MaintenanceService) Transition(ctx context.Context, id, from, to int) error {
	if from == to {
		return fmt.Errorf("status unchanged")
	}
	if from == model.TaskPending && to != model.TaskRunning {
		return fmt.Errorf("pending task can only start")
	}
	if from == model.TaskRunning && (to != model.TaskCompleted && to != 4) {
		return fmt.Errorf("running task can only complete or cancel")
	}
	return s.Repo.UpdateTaskStatus(ctx, uint64(id), from, to)
}
