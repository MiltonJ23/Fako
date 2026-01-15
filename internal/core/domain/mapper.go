package domain

// CommandMapper is the contract every driver type (linux, cisco, juniper, ...) must implement
type CommandMapper interface {

	// GenerateApplyCommand translate a resource into a list of command Shell/CLI
	GenerateApplyCommand(r *Resource) ([]RemoteCommand, error)

	// GenerateDeleteCommand will generate the command list to take down a resource configuration
	GenerateDeleteCommand(r *Resource) ([]RemoteCommand, error)
}
