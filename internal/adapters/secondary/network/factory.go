package network

import (
	"fmt"
	"os"

	"github.com/MiltonJ23/Fako/internal/core/domain"
	"github.com/MiltonJ23/Fako/internal/core/ports"
)

func GetDriver(deviceType string, osType string, dryRun bool) (ports.NetworkDriver, error) {
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

		var secretPass *domain.Secret
		if passPhrase != "" {
			secretPass = domain.NewSecret(passPhrase)
		}
		// After we got the passphrase from the environment, we unset the passphrase environment variable
		os.Unsetenv("FAKO_TARGET_PASSPHRASE")
		if host == "" || user == "" || keyPath == "" {
			return nil, fmt.Errorf("missing ssh configuration, checked the environment variables (FAKO_TARGET_*)")
		}

		selectedMapper, SelectionMapperError := GetMapper(osType)
		if SelectionMapperError != nil {
			return nil, fmt.Errorf("an error happened while selecting the Mapper:%v", SelectionMapperError.Error())
		}

		return NewSSHDriver(host, user, keyPath, secretPass, selectedMapper, dryRun)
		
	default:
		return nil, fmt.Errorf("device type %s not supported", deviceType)
	}
}
