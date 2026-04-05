package model

import (
	"time"

	"github.com/google/uuid"
)

type ConfigurationRevision struct {
	ID              uuid.UUID `json:"id" db:"id"`
	ConfigurationID uuid.UUID `json:"configuration_id" db:"configuration_id"`
	RevisionNo      int       `json:"revision_no" db:"revision_no"`
	Files           []File    `json:"files" db:"files"`
	ContentHash     string    `json:"content_hash" db:"content_hash"`
	Message         string    `json:"message,omitempty" db:"message"`
	CreatedBy       string    `json:"created_by,omitempty" db:"created_by"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

func (ConfigurationRevision) CollectionName() string { return "configuration_revisions" }
