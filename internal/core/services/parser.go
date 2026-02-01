package services

import (
	"fmt"

	"github.com/MiltonJ23/Fako/internal/core/domain"
	"gopkg.in/yaml.v3"
)

func ParseAndValidate(data []byte) (*domain.DependencyGraph, error) {
	var intent domain.Intent
	// let's parse the yaml file
	ParsingError := yaml.Unmarshal(data, &intent)
	if ParsingError != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", ParsingError)
	}
	// Parsing passed successfully, now let's initialize the graph
	graph := domain.NewDependencyGraph()
	// The graph was created , now let's add nodes
	for _, res := range intent.Resources {
		ValidationError := res.Validate()
		if ValidationError != nil {
			return nil, ValidationError
		}
		AddingResourceError := graph.AddResource(&res)
		if AddingResourceError != nil {
			return nil, AddingResourceError
		}
	}
	// now let's build the edges--- which means we are now building the dependencies
	EdgeBuildingError := graph.BuildEdges()
	if EdgeBuildingError != nil {
		return nil, EdgeBuildingError
	}
	// Safety check | Cycles
	FindCyclesError := graph.DetectCycles()
	if FindCyclesError != nil {
		return nil, FindCyclesError
	}
	return graph, nil
}
