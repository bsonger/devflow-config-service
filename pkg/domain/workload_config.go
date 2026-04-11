package domain

import "github.com/google/uuid"

type WorkloadConfig struct {
	BaseModel

	ApplicationID uuid.UUID      `json:"application_id" db:"application_id"`
	EnvironmentID string         `json:"environment_id,omitempty" db:"environment_id"`
	Name          string         `json:"name" db:"name"`
	Replicas      int            `json:"replicas" db:"replicas"`
	Resources     map[string]any `json:"resources,omitempty" db:"resources"`
	Probes        map[string]any `json:"probes,omitempty" db:"probes"`
	Env           []EnvVar       `json:"env,omitempty" db:"env"`
	WorkloadType  string         `json:"workload_type,omitempty" db:"workload_type"`
	Strategy      string         `json:"strategy,omitempty" db:"strategy"`
}

type WorkloadConfigInput struct {
	ApplicationID uuid.UUID      `json:"application_id"`
	EnvironmentID string         `json:"environment_id,omitempty"`
	Name          string         `json:"name"`
	Replicas      int            `json:"replicas"`
	Resources     map[string]any `json:"resources,omitempty"`
	Probes        map[string]any `json:"probes,omitempty"`
	Env           []EnvVar       `json:"env,omitempty"`
	WorkloadType  string         `json:"workload_type,omitempty"`
	Strategy      string         `json:"strategy,omitempty"`
}
