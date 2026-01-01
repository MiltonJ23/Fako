package ports

import "github.com/MiltonJ23/Fako/internal/core/domain"

// StateRepository defines how we persist the state of the application
type StateRepository interface {
	// Load reads the state from storage
	Load() (*domain.State, error)
	// Save store the state from memory to storage
	Save(state *domain.State) error
}
