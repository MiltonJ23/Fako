package services

import (
	"context"
	"fmt"

	"github.com/MiltonJ23/Fako/internal/core/domain"
	"github.com/MiltonJ23/Fako/internal/core/ports"
)

// Enforce takes a plan and executes it using the provided driver
func Enforce(ctx context.Context, plan []*domain.Resource, driver ports.NetworkDriver) error { // we accept the interface not the struct , this is dependencyInjection i tried

	fmt.Println("-> Starting Network Enforcement .....")

	for i, res := range plan {
		fmt.Printf("Step %d/%d: Processing %s...\n", i+1, len(plan), res.ID)
		// now we are going to delegate the details to the driver
		err := driver.ApplyResource(ctx, res)
		if err != nil {
			return fmt.Errorf("failed to apply resource %s: %w", res.ID, err)
		}
	}
	fmt.Println("-> Enforcement Complete. Network is synced.")
	return nil
}
