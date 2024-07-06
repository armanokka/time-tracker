package models

import "time"

type TimeEntry struct {
	ID        int64     `json:"id" db:"id"`
	TaskID    int64     `json:"task_id" db:"task_id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	StartedAt time.Time `json:"started_at" db:"started_at"`
	EndedAt   time.Time `json:"ended_at" db:"ended_at"`
}
