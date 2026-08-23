package service

import (
	"context"
	"fmt"
	"predictive-maintenance/internal/model"
	"sync"
	"time"
)

type Notification struct {
	ID        uint64
	Kind      string
	DeviceID  uint64
	Title     string
	Content   string
	CreatedAt time.Time
	Read      bool
}
type NotificationService struct {
	mu    sync.RWMutex
	items []Notification
	next  uint64
}

func NewNotificationService() *NotificationService {
	return &NotificationService{items: []Notification{}}
}
func (s *NotificationService) Publish(ctx context.Context, kind string, device uint64, title, content string) Notification {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next++
	n := Notification{ID: s.next, Kind: kind, DeviceID: device, Title: title, Content: content, CreatedAt: time.Now().UTC()}
	s.items = append(s.items, n)
	return n
}
func (s *NotificationService) PublishAbnormal(ctx context.Context, a model.AbnormalRecord) Notification {
	return s.Publish(ctx, "abnormal", a.DeviceID, fmt.Sprintf("设备 %d 异常", a.DeviceID), a.SensorType)
}
func (s *NotificationService) List(unreadOnly bool) []Notification {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []Notification{}
	for _, n := range s.items {
		if unreadOnly && n.Read {
			continue
		}
		out = append(out, n)
	}
	return out
}
func (s *NotificationService) MarkRead(id uint64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.items {
		if s.items[i].ID == id {
			s.items[i].Read = true
			return true
		}
	}
	return false
}
func (s *NotificationService) UnreadCount() int {
	n := 0
	for _, x := range s.List(true) {
		_ = x
		n++
	}
	return n
}
