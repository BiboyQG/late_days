package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	LateDays interface {
		Create(ctx context.Context, lateDay *LateDay) error
		Update(ctx context.Context, lateDay *LateDay) error
		GetByStudentID(ctx context.Context, studentID string) (*LateDay, error)
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		LateDays: &LateDayStore{db: db},
	}
}
