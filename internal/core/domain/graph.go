package domain

import (
	"errors"
	"fmt"
)

// DependencyGraph manages the relationships between resources
type DependencyGraph struct {
	Nodes map[string]*Resource
	Edges map[string][]string // Key: parent | Value: Children(dependency)
}

// NewDependencyGraph creates a new empty graph
func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		Nodes: make(map[string]*Resource),
		Edges: make(map[string][]string),
	}
}

func (graph *DependencyGraph) AddResource(r Resource) error {
	// let's check if the resource exists already in the resource graph
	_, exist := graph.Nodes[r.ID]
	if exist {
		return errors.New("duplicate resource ID found : " + r.ID)
	}
	// reaching here means the resource is not duplicated
	graph.Nodes[r.ID] = &r

	return nil
}

// BuildEdges connects the nodes based on the `depends_on` fields
func (graph *DependencyGraph) BuildEdges() error {
	for id, resource := range graph.Nodes {
		for _, depID := range resource.DependsOn {
			// we first of all check if the dependency actually exists
			_, exist := graph.Edges[depID]
			if !exist {
				return fmt.Errorf("resource %s depends on missing resource %s", id, depID)
			}
			// Here it means , ti actually exists
			graph.Edges[depID] = append(graph.Edges[depID], id)
		}
	}
	return nil
}

// DetectCycles implements the Depth First Search (DFS) to detect loops
func (graph *DependencyGraph) DetectCycles() error {
	visited := make(map[string]bool)
	recursionStack := make(map[string]bool)

	for NodeID := range graph.Nodes {
		if graph.hasCycle(NodeID, visited, recursionStack) {
			return fmt.Errorf("cycle detected invloving resource %s", NodeID)
		}
	}
	return nil
}

// hasCycle is the recursive helper for DFS
func (graph *DependencyGraph) hasCycle(node string, visited map[string]bool, stack map[string]bool) bool {
	visited[node] = true
	stack[node] = true

	// check all the resources that depends on this node
	for _, neighbor := range graph.Edges[node] {
		if !visited[neighbor] {
			graph.hasCycle(neighbor, visited, stack)
			return true
		} else if stack[neighbor] {
			// if neighbor is in the current recursion stack ,we found a back-edge (cycle)
			return true
		}
	}
	stack[node] = false
	return false
}
