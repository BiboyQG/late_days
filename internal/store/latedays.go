package store

import (
	"context"
	"database/sql"
)

type LateDay struct {
	ID           int    `json:"id"`
	StudentEmail string `json:"student_email"`
	StudentID    string `json:"student_id"`
	Days         int    `json:"days"`
}

type LateDays interface {
	GetByStudentID(ctx context.Context, studentID string) (*LateDay, error)
	Create(ctx context.Context, lateDay *LateDay) error
	Update(ctx context.Context, lateDay *LateDay) error
}

type LateDayStore struct {
	db *sql.DB
}

func (s *LateDayStore) GetByStudentID(ctx context.Context, studentID string) (*LateDay, error) {
	query := `
		SELECT id, student_email, student_id, days
		FROM late_days
		WHERE student_id = $1
	`

	var lateDay LateDay

	err := s.db.QueryRowContext(ctx, query, studentID).Scan(&lateDay.ID, &lateDay.StudentEmail, &lateDay.StudentID, &lateDay.Days)
	return &lateDay, err
}

func (s *LateDayStore) Create(ctx context.Context, lateDay *LateDay) error {
	query := `
		INSERT INTO late_days (student_email, student_id, days)
		VALUES ($1, $2, $3) RETURNING id
	`

	err := s.db.QueryRowContext(ctx, query, lateDay.StudentEmail, lateDay.StudentID, lateDay.Days).Scan(&lateDay.ID)
	return err
}

func (s *LateDayStore) Update(ctx context.Context, lateDay *LateDay) error {
	query := `
		UPDATE late_days
		SET student_email = $1, student_id = $2, days = $3
		WHERE id = $4
	`

	_, err := s.db.ExecContext(ctx, query, lateDay.StudentEmail, lateDay.StudentID, lateDay.Days, lateDay.ID)
	return err
}