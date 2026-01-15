package network

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/MiltonJ23/Fako/internal/core/domain"
)

type LinuxMapper struct {
}

func NewLinuxMapper() *LinuxMapper {
	return &LinuxMapper{}
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
	case "linux-bgp":
		var frrBGPConfigurationTemplate string
		// fetch the properties
		autonomouSystemN, asError := getString(r.Config, "asn")
		if asError != nil {
			return nil, fmt.Errorf("unable to find properties %s for resource %s :%v", "asn", r.ID, asError)
		}
		routerId, routerFetchingError := getString(r.Config, "router_id")
		if routerFetchingError != nil {
			return nil, fmt.Errorf("unable to find properties %s for resource %s :%v", "router_id", r.ID, routerFetchingError)
		}
		hostname, hostnameFetchError := getString(r.Config, "hostname")
		if hostnameFetchError != nil {
			return nil, fmt.Errorf("unable to find properties %s for resource %s :%v", "hostname", r.ID, hostnameFetchError)
		}

		neighborsList, neighborsListFetchingError := getList(r.Config, "neighbors")
		if neighborsListFetchingError != nil {
			return nil, fmt.Errorf("unable to locate the list of bgp neighbors for %s : %v", r.ID, neighborsListFetchingError)
		}
		var neighs []domain.NeighborData
		for _ = range neighborsList {
			neighborIP, _ := getString(r.Config, "ip")
			neighborAS, _ := getString(r.Config, "asn")
			neighs = append(neighs, domain.NeighborData{NeighborIP: neighborIP, NeighborAS: neighborAS})
		}

		data := domain.BGPConfigTemplateDate{
			Hostname:  hostname,
			ASN:       autonomouSystemN,
			RouterID:  routerId,
			Neighbors: neighs,
		}
		// rendering the template
		parsedTemplate, parsingTemplateError := template.New("frr-bgp").Parse(frrBGPConfigurationTemplate)
		if parsingTemplateError != nil {
			return nil, fmt.Errorf("unable to parse the BGP template %s", parsingTemplateError)
		}

		var ConfigBytes bytes.Buffer
		TemplateRenderingError := parsedTemplate.Execute(&ConfigBytes, data)
		if TemplateRenderingError != nil {
			return nil, fmt.Errorf("failed to render BGP Template: %v", TemplateRenderingError)
		}
		bgpConfig := ConfigBytes.String()
		// build the remoteCommands
		cmds = append(cmds, []domain.RemoteCommand{
			{
				Description: fmt.Sprintf("Install FRR if not already present"),
				Cmd:         fmt.Sprintf("dpkg -l frr > /dev/null 2>&1 || sudo apt-get update && sudo apt-get install -y frr"),
			},
			{
				Description: fmt.Sprintf("Enable the BGP daemon"),
				Cmd:         fmt.Sprintf("sudo sed -i 's/bgpd=no/bgpd=yes/g' /etc/frr/daemons"),
			},
			{
				Description: fmt.Sprintf("Enable IPV4 forwarding"),
				Cmd:         fmt.Sprintf("sudo sysctl -w net.ipv4.ip_forward=1"),
			},
			{
				Description: fmt.Sprintf("Write the Frr configuration file from Template"),
				Cmd:         fmt.Sprintf("cat <<EOF | sudo tee /etc/frr/frr.conf\\n%s\\nEOF", bgpConfig),
			},
			{
				Description: fmt.Sprintf("Restart the FRR service"),
				Cmd:         fmt.Sprintf("sudo systemctl restart frr"),
			},
		}...)

	default:
		return nil, fmt.Errorf("unsupported resource type for Linux %s ", r.Kind)
	}

	return cmds, nil
}

func (l *LinuxMapper) GenerateDeleteCommands(r *domain.Resource) ([]domain.RemoteCommand, error) {
	return nil, nil
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

// getList for reading a list of neighbors for networking protocols
func getList(properties map[string]interface{}, key string) ([]map[string]interface{}, error) {
	// first thing first
	value, exists := properties[key]
	if !exists {
		return nil, fmt.Errorf("key %s does not exist", key)
	}
	// convert a generic interface to a slice
	SliValue, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("key %s must be a list", key)
	}

	// now let's fetch the neighbors
	var neighborsList []map[string]interface{}

	for _, element := range SliValue {
		mapItem, ok := element.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("item %s must be a map", mapItem)
		}
		neighborsList = append(neighborsList, mapItem)
	}
	return neighborsList, nil

}
