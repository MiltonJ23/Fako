package domain

import "time"

const (
	StatusCreated = "CREATED"
	StatusFailed  = "FAILED"
	StatusPending = "PENDING"
)

// ResourceState represents the actual state of a resource
type ResourceState struct {
	ID          string       `json:"id"`
	Kind        ResourceType `json:"kind"`
	Status      string       `json:"status"`
	LastApplied time.Time    `json:"last_applied"`
}

type State struct {
	Resources map[string]ResourceState `json:"resources"`
}

func NewState() *State {
	return &State{
		Resources: make(map[string]ResourceState),
	}
}
