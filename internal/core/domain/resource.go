package domain

import "errors"

// ResourceType will enable us to make the difference between a VLAN and an interface
type ResourceType string

const (
	TypeVLAN      ResourceType = "VLAN"
	TypeInterface ResourceType = "INTERFACE"
)

// Resource represents a single node in our node graph
type Resource struct {
	ID        string                 `yaml:"id"`         // this is the unique identifier of a resource, let's say: firewall-10
	Kind      ResourceType           `yaml:"kind"`       // as it can suggest, the type of the resource
	DependsOn []string               `yaml:"depends_on"` // this is the list of resource the current one depends on
	Config    map[string]interface{} `yaml:"config"`     // the actual specs of the resource
}

// Intent represents the full desired state of the network defined by the user
type Intent struct {
	Resources []Resource `yaml:"resources"`
}

// Validate ensures the resource has the least required field
func (r *Resource) Validate() error {
	if r.ID == "" {
		return errors.New("resource missing required field: id")
	}
	if r.Kind == "" {
		return errors.New("resource missing required field: kind")
	}
	return nil
}
