package models

import "database/sql/driver"

type Task struct {
	ID          int64  `json:"id" db:"id" validate:"omitempty"`
	Name        string `json:"name" db:"name" validate:"lte=64"`
	Description string `json:"description" db:"description" validate:"omitempty,lte=256"`
	ProjectID   int64  `json:"project_id" db:"project_id" validate:"omitempty"`
	Finished    bool   `json:"finished" db:"finished"`
}

func (task *Task) Columns() []string {
	return []string{"id", "name", "description", "project_id", "finished"}
}

func (task *Task) Fields() []driver.Value {
	return []driver.Value{task.ID, task.Name, task.Description, task.ProjectID, task.Finished}
}

type UserProductivity struct {
	TaskID       int64 `json:"task_id"`
	SpentHours   int   `json:"spent_hours"`
	SpentMinutes int   `json:"spent_minutes"`
}
