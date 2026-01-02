package services

import (
	"context"
	"fmt"
	"time"

	"github.com/MiltonJ23/Fako/internal/core/domain"
	"github.com/MiltonJ23/Fako/internal/core/ports"
)

// Enforce takes a plan and executes it using the provided driver
func Enforce(ctx context.Context, plan []*domain.Resource, driver ports.NetworkDriver, repo ports.StateRepository) error { // we accept the interface not the struct , this is dependencyInjection i tried
	// we are going to load the current memory , i mean the memory of the state of the application
	currentState, LoadingStateError := repo.Load()
	if LoadingStateError != nil {
		return fmt.Errorf("unable to load the state %v", LoadingStateError)
	}
	fmt.Println("-> Starting Network Enforcement ......")
	for i, resource := range plan {
		fmt.Printf("Step %d/%d: Processing %s... ", i+1, len(plan), resource.ID)

		// we are then going to check if the resource we are iterating on is already existing
		existing, found := currentState.Resources[resource.ID]
		if found {
			if existing.Status == "CREATED " {
				fmt.Println("[Skipping] Resource already exist...")
				continue
			}
		}

		// reaching here means , the resource doesn't exist already
		// We are then going to create the resource
		CreatingResourceError := driver.ApplyResource(ctx, resource)
		if CreatingResourceError != nil {
			currentState.Resources[resource.ID] = domain.ResourceState{
				ID: resource.ID, Kind: resource.Kind, Status: "FAILED", LastApplied: time.Now(),
			}
			SavingError := repo.Save(currentState)
			if SavingError != nil {
				return fmt.Errorf("->an error occurred while saving the state %v", SavingError)
			}
			return fmt.Errorf("->unable to apply the resource %s:%v", resource.ID, CreatingResourceError)
		}

		currentState.Resources[resource.ID] = domain.ResourceState{
			ID: resource.ID, Kind: resource.Kind, Status: "CREATED", LastApplied: time.Now(),
		}
		// Now let's save the thing to disk
		SavingError := repo.Save(currentState)
		if SavingError != nil {
			return fmt.Errorf("->an error occured while saving the state %v", SavingError)
		}
		fmt.Println("")
	}
	fmt.Println("-> [Enforcement] Completed")
	return nil
}
