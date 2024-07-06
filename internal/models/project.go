package models

import "database/sql/driver"

type Project struct {
	ID          int64   `json:"id" db:"id" validate:"omitempty"`
	Name        string  `json:"name" db:"name" validate:"lte=64"`
	Description *string `json:"description" db:"description" validate:"omitempty,lte=1024"`
	CreatorID   int64   `json:"creator_id" db:"creator_id" validate:"omitempty"`
}

func (project *Project) Columns() []string {
	return []string{"id", "name", "description", "creator_id"}
}

func (project *Project) Fields() []driver.Value {
	return []driver.Value{project.ID, project.Name, project.Description, project.CreatorID}
}
