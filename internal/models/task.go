package models

type Task struct {
	ID          int64  `json:"id" db:"id" validate:"omitempty"`
	Name        string `json:"name" db:"name" validate:"lte=64"`
	Description string `json:"description" db:"description" validate:"omitempty,lte=256"`
	ProjectID   int64  `json:"project_id" db:"project_id" validate:"omitempty"`
	Finished    bool   `json:"finished" db:"finished"`
}

type UserProductivity struct {
	TaskID       int64 `json:"task_id"`
	SpentHours   int   `json:"spent_hours"`
	SpentMinutes int   `json:"spent_minutes"`
}
