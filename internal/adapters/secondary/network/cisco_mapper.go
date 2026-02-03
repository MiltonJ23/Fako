package network

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/MiltonJ23/Fako/internal/core/domain"
)

var ciscoBGPConfigTemplate string // go embed : templates/cisco_bgp

type CiscoMapper struct{}

func NewCiscoMapper() *CiscoMapper {
	return &CiscoMapper{}
}

func (c *CiscoMapper) GenerateApplyCommands(r *domain.Resource) ([]domain.RemoteCommand, error) {
	var commands []domain.RemoteCommand

	switch string(r.Kind) {
	case "cisco-bgp":
		// first of all let's fetch the properties from the intent YAML
		asn, _ := getString(r.Config, "asn")
		routerID, _ := getString(r.Config, "routerID")
		rawNeighbors, _ := getList(r.Config, "neighbors")

		var neighbors []domain.NeighborData
		for _, neighbor := range rawNeighbors {
			neighborIP, _ := getString(neighbor, "ip")
			neighborAS, _ := getString(neighbor, "asn")
			neighbors = append(neighbors, domain.NeighborData{NeighborIP: neighborIP, NeighborAS: neighborAS})
		}

		data := domain.CiscoBGPData{ASN: asn, RouterID: routerID, Neighbors: neighbors} // we got the data
		// now we are going to parse the data into the template

		ciscoTmpl, parsingTmplError := template.New("cisco-bgp").Parse(ciscoBGPConfigTemplate)
		if parsingTmplError != nil {
			return nil, fmt.Errorf("error parsing template: %s", parsingTmplError.Error())
		}
		var buffer bytes.Buffer
		executeTemplateError := ciscoTmpl.Execute(&buffer, data)
		if executeTemplateError != nil {
			return nil, fmt.Errorf("error executing template: %s", executeTemplateError.Error())
		}
		// since i don't have the cisco router to test, i will first of all print something to say the command was executed
		commands = append(commands, domain.RemoteCommand{
			Description: "Apply Cisco BGP Configurations",
			Cmd:         fmt.Sprintf("configure terminal\n%s\nend", buffer.String()),
			//TODO: manage the ping-pong of configurations for Cisco here
		})
	default:
		fmt.Printf("The resource kind %s is not supported", r.Kind)
	}
	return commands, nil
}

func (c *CiscoMapper) GenerateDeleteCommands(r *domain.Resource) ([]domain.RemoteCommand, error) {
	// TODO: Implement the proper GenerateDeleteCommand method
	// TODO: Make sure to implement a graceful Shutdown
	return nil, nil
}
