package ports

import "github.com/MiltonJ23/Fako/internal/core/domain"

// CommandMapper is the contract every driver type (linux, cisco, juniper, ...) must implement
type CommandMapper interface {

	// GenerateApplyCommand translate a resource into a list of command Shell/CLI
	GenerateApplyCommand(r *domain.Resource) ([]domain.RemoteCommand, error)

	// GenerateDeleteCommand will generate the command list to take down a resource configuration
	GenerateDeleteCommand(r *domain.Resource) ([]domain.RemoteCommand, error)
}
