package domain

type RemoteCommand struct {
	Cmd         string
	Description string
	IgnoreError bool
}
