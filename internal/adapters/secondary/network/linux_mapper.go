package network

import (
	"fmt"

	"github.com/MiltonJ23/Fako/internal/core/domain"
)

type LinuxMapper struct {
}

func NewLinuxMapper() *LinuxMapper {
	return &LinuxMapper{}
}

// getString will help extract string properties from our configuration slice
func getString(properties map[string]interface{}, key string) (string, error) {
	// first thing first
	value, exists := properties[key]
	if !exists {
		return "", fmt.Errorf("key %s does not exist", key)
	}
	// after that we check if the value is a string, so we cast it to a string / or we extract the string value out of it
	strValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("key %s is not a string", key)
	}
	return strValue, nil
}

func (l *LinuxMapper) GenerateApplyCommands(r *domain.Resource) ([]domain.RemoteCommand, error) {
	// First
	var cmds []domain.RemoteCommand

	switch string(r.Kind) {
	case "linux-interface":
		interfaceName := r.ID
		interfaceIpAddress, IpAddressFetchingError := getString(r.Config, "ip")
		if IpAddressFetchingError != nil {
			return nil, fmt.Errorf("an error occured trying to fetch ip address from configuration file for resource %s : %v", interfaceName, IpAddressFetchingError)
		}

		cmds = append(cmds, []domain.RemoteCommand{
			{
				Cmd:         fmt.Sprintf("ip link show %s > /dev/null 2>&1 || sudo ip link add %s type dummy", interfaceName, interfaceName),
				Description: fmt.Sprintf("Ensure interface %s exists", interfaceName),
			},
			{
				Cmd:         fmt.Sprintf("sudo ip addr add %s dev %s || true", interfaceIpAddress, interfaceName),
				Description: fmt.Sprintf("Ensure Ip address is set for interface  %s", interfaceName),
			},
			{
				Cmd:         fmt.Sprintf("sudo ip link set %s up ", interfaceName),
				Description: fmt.Sprintf("Ensure interface is UP"),
			},
		}...)
	case "linux-route":
		routeDestination := r.Config["destination"].(string) // we make a cast
		gateway := r.Config["gateway"].(string)

		device := r.Config["device"].(string)

		// now let's build the command string
		commandString := fmt.Sprintf("sudo ip route add %s via %s dev %s || sudo ip route replace %s via %s dev %s", routeDestination, gateway, device, routeDestination, gateway, device)

		cmds = append(cmds, domain.RemoteCommand{
			Cmd:         commandString,
			Description: fmt.Sprintf("Add route to %s ", routeDestination),
		})
	default:
		return nil, fmt.Errorf("unsupported resource type for Linux %s ", r.Kind)
	}

	return cmds, nil
}

func (l *LinuxMapper) GenerateDeleteCommands(r *domain.Resource) ([]domain.RemoteCommand, error) {
	return nil, nil
}
