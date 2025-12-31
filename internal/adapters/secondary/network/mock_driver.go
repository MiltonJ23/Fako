package network

import (
	"context"
	"fmt"
	"time"

	"github.com/MiltonJ23/Fako/internal/core/domain"
)

type MockDriver struct {
}

func NewMockDriver() *MockDriver {
	return &MockDriver{}
}

// now let's implement the function of the contract

func (m *MockDriver) ApplyResource(ctx context.Context, r *domain.Resource) error {
	// let's simulate the access and all
	select {
	case <-time.After(500 * time.Millisecond):
	// This right here means everything was ok, or would be precise to say it simulates  everything is okay
	case <-ctx.Done():
		return ctx.Err()
	}

	fmt.Printf("-> [Mock Driver] Connecting to device... Configured %s (%s)\n", r.Kind, r.ID)
	return nil
}

func (m *MockDriver) DeleteResource(ctx context.Context, r *domain.Resource) error {
	fmt.Printf("->  [Mock Driver] Removing %s (%s)\n", r.Kind, r.ID)
	return nil
}
