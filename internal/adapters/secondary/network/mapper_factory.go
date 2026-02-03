package network

import (
	"fmt"

	"github.com/MiltonJ23/Fako/internal/core/domain"
)

// GetMapper is going to select the correct template for the right equipment (Wheter Cisco / Juniper / Microtik..>)
func GetMapper(osType string) (domain.CommandMapper, error) {
	switch osType {
	case "linux":
		return NewLinuxMapper(), nil
	case "cisco":
		return NewCiscoMapper(), nil
	default:
		return nil, fmt.Errorf("unsupported OS: %s", osType)
	}
}
