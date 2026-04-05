package domain

import "github.com/google/uuid"

type Configuration struct {
	BaseModel

	ApplicationID    uuid.UUID  `json:"application_id" db:"application_id"`
	Name             string     `json:"name" db:"name"`
	Env              string     `json:"env" db:"env"`
	SourcePath       string     `json:"source_path" db:"source_path"`
	LatestRevisionNo int        `json:"latest_revision_no" db:"latest_revision_no"`
	LatestRevisionID *uuid.UUID `json:"latest_revision_id,omitempty" db:"latest_revision_id"`
}

func (Configuration) CollectionName() string { return "configurations" }
