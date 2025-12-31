package ports

import (
	"context"

	"github.com/MiltonJ23/Fako/internal/core/domain"
)

// NetworkDriver describe the behavior of any device driver . Whether its SSH, Simulation or API, will have to implement these methods
type NetworkDriver interface {
	// we will have to pass  the context, so that if the user press ctrl c we can stop whatever we were doing

	//ApplyResource makes the actual change on  the network
	ApplyResource(ctx context.Context, r *domain.Resource) error

	//DeleteResource revert  the change made by applyResource , Crazy right
	DeleteResource(ctx context.Context, r *domain.Resource) error
}
