package network

import (
	"fmt"

	"github.com/MiltonJ23/Fako/internal/core/ports"
)

func GetDriver(deviceType string) (ports.NetworkDriver, error) {
	switch deviceType {
	case "mock":
		return NewMockDriver(), nil
	case "linux-local":
		return NewLinuxDriver(), nil
	default:
		return nil, fmt.Errorf("device type %s not supported", deviceType)
	}
}
