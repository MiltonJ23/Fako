package network

import (
	"fmt"
	"os"

	"github.com/MiltonJ23/Fako/internal/core/ports"
)

func GetDriver(deviceType string) (ports.NetworkDriver, error) {
	switch deviceType {
	case "mock":
		return NewMockDriver(), nil
	case "linux-local":
		return NewLinuxDriver(), nil
	case "ssh-target":
		host := os.Getenv("FAKO_TARGET_HOST")
		user := os.Getenv("FAKO_TARGET_USER")
		keyPath := os.Getenv("FAKO_TARGET_KEY")
		passPhrase := os.Getenv("FAKO_TARGET_PASSPHRASE")

		if host == "" || user == "" || keyPath == "" {
			return nil, fmt.Errorf("missing ssh configuration, checked the environment variables (FAKO_TARGET_*)")
		}
		return NewSSHDriver(host, user, keyPath, passPhrase)
	default:
		return nil, fmt.Errorf("device type %s not supported", deviceType)
	}
}
