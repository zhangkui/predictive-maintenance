package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type AuditEvent struct {
	ID         uint64         `json:"id"`
	ActorID    uint64         `json:"actorId"`
	Action     string         `json:"action"`
	Resource   string         `json:"resource"`
	ResourceID uint64         `json:"resourceId"`
	Detail     map[string]any `json:"detail"`
	CreatedAt  time.Time      `json:"createdAt"`
}
type AuditService struct{ DB *sql.DB }

func NewAuditService(db *sql.DB) *AuditService { return &AuditService{db} }
func (s *AuditService) EnsureTable() error {
	_, e := s.DB.Exec("CREATE TABLE IF NOT EXISTS audit_events(id INTEGER PRIMARY KEY AUTOINCREMENT,actor_id INTEGER,action TEXT,resource TEXT,resource_id INTEGER,detail TEXT,created_at TEXT)")
	return e
}
func (s *AuditService) Record(ctx context.Context, e AuditEvent) error {
	if err := s.EnsureTable(); err != nil {
		return err
	}
	b, _ := json.Marshal(e.Detail)
	_, err := s.DB.ExecContext(ctx, "INSERT INTO audit_events(actor_id,action,resource,resource_id,detail,created_at)VALUES(?,?,?,?,?,?)", e.ActorID, e.Action, e.Resource, e.ResourceID, string(b), time.Now().UTC().Format(time.RFC3339))
	return err
}
func (s *AuditService) List(ctx context.Context, limit int) ([]AuditEvent, error) {
	if err := s.EnsureTable(); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 100
	}
	rows, err := s.DB.QueryContext(ctx, "SELECT id,actor_id,action,resource,resource_id,detail,created_at FROM audit_events ORDER BY id DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []AuditEvent{}
	for rows.Next() {
		var e AuditEvent
		var detail, created string
		if err = rows.Scan(&e.ID, &e.ActorID, &e.Action, &e.Resource, &e.ResourceID, &detail, &created); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(detail), &e.Detail)
		e.CreatedAt, _ = time.Parse(time.RFC3339, created)
		out = append(out, e)
	}
	return out, rows.Err()
}
