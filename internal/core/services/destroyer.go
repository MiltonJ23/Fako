package services

import (
	"context"
	"fmt"

	"github.com/MiltonJ23/Fako/internal/core/domain"
	"github.com/MiltonJ23/Fako/internal/core/ports"
)

func Destroy(ctx context.Context, plan []*domain.Resource, driver ports.NetworkDriver, repo ports.StateRepository) error {
	// we are first of all loading the state
	currentState, LoadingStateError := repo.Load()
	if LoadingStateError != nil {
		return fmt.Errorf("-> an error occured when loading the state: %s", LoadingStateError)
	}
	fmt.Println("-> Starting Network Destruction")
	for i, resource := range plan {
		fmt.Printf("->Step %d/%d: Destroying %s... \n", i+1, len(plan), resource.ID)
		_, exists := currentState.Resources[resource.ID]
		if !exists {
			fmt.Println("-> [SKIPPED] Resource not found in state(was never created in the first place)")
			continue
		}
		// now let's delete it properly
		ResourceDeletionError := driver.DeleteResource(ctx, resource)
		if ResourceDeletionError != nil {
			return fmt.Errorf("failed to destroy resource %s : %s", resource.ID, ResourceDeletionError)
		}
		// we are then going to remove the List of current resources
		delete(currentState.Resources, resource.ID)

		SavingStateError := repo.Save(currentState)
		if SavingStateError != nil {
			return fmt.Errorf("-> an error occured when saving the state: %s", SavingStateError)
		}
	}
	fmt.Println("-> Finishing Network Destruction")
	return nil
}
