package domain

// CommandMapper is the contract every driver type (linux, cisco, juniper, ...) must implement
type CommandMapper interface {
	// GenerateApplyCommands translate a resource into a list of command Shell/CLI
	GenerateApplyCommands(r *Resource) ([]RemoteCommand, error)

	// GenerateDeleteCommands will generate the command list to take down a resource configuration
	GenerateDeleteCommands(r *Resource) ([]RemoteCommand, error)
}
