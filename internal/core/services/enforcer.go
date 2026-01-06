package services

import (
	"context"
	"fmt"
	"time"

	"github.com/MiltonJ23/Fako/internal/core/domain"
	"github.com/MiltonJ23/Fako/internal/core/ports"
)

const (
	MaxRetries = 3
	BaseDelay  = 1 * time.Second
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
		// First of all, let's see if the user pressed ctrl C
		if ctx.Err() != nil {
			return fmt.Errorf("-> operation canceled by user")
		}
		fmt.Printf("Step %d/%d: Processing %s... ", i+1, len(plan), resource.ID)

		// we are then going to check if the resource we are iterating on  already exist
		existing, found := currentState.Resources[resource.ID]
		if found {
			if existing.Status == domain.StatusCreated {
				fmt.Println("[Skipping] Resource already exist...")
				continue
			}
		}

		var ApplyError error
		for attempt := 1; attempt <= MaxRetries; attempt++ {
			// We try to apply first
			ApplyError = driver.ApplyResource(ctx, resource)
			if ApplyError == nil {
				break // This would actually mean that the application was succesfful and that we can exit now
			}

			if attempt == MaxRetries {
				break // Man, enough trials
			}
			// reaching here means, the application wasn't successful
			delay := BaseDelay * time.Duration(attempt)
			fmt.Printf("\n -> Attempt %d failed: %v. Retrying is %s...", attempt, ApplyError, delay)

			// we will then wait meanwhile also be on guard to catch a terminaison signal
			select {
			case <-time.After(delay):
			// it is supposed to continue with the execution
			case <-ctx.Done():
				return fmt.Errorf("-> cancellation occured during retry operations")
			}
		}

		if ApplyError != nil {
			currentState.Resources[resource.ID] = domain.ResourceState{
				ID: resource.ID, Kind: resource.Kind, Status: domain.StatusFailed, LastApplied: time.Now(),
			}
			SavingError := repo.Save(currentState)
			if SavingError != nil {
				return fmt.Errorf("->an error occurred while saving the state %v", SavingError)
			}
			return fmt.Errorf("->unable to apply the resource %s:%v", resource.ID, ApplyError)
		}

		currentState.Resources[resource.ID] = domain.ResourceState{
			ID: resource.ID, Kind: resource.Kind, Status: domain.StatusCreated, LastApplied: time.Now(),
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
