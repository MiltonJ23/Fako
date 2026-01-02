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

func (graph *DependencyGraph) AddResource(r *Resource) error {

	// let's first of all check if the ID is emtpy
	if r.ID == "" {
		return errors.New("resource id is empty")
	}
	// let's check if the resource exists already in the resource graph
	_, exist := graph.Nodes[r.ID]
	if exist {
		return errors.New("duplicate resource ID found : " + r.ID)
	}
	// reaching here means the resource is not duplicated
	// let's store the pointer
	graph.Nodes[r.ID] = r

	return nil
}

// BuildEdges connects the nodes based on the `depends_on` fields
func (graph *DependencyGraph) BuildEdges() error {
	for id, resource := range graph.Nodes {
		for _, depID := range resource.DependsOn {
			// we first of all check if the dependency actually exists
			_, exist := graph.Nodes[depID]
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
		if !visited[NodeID] { // we don't revisit nodes we visited in previous iterations
			if graph.hasCycle(NodeID, visited, recursionStack) {
				return fmt.Errorf("cycle detected invloving resource %s", NodeID)
			}
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
			if graph.hasCycle(neighbor, visited, stack) {
				return true
			}

		} else if stack[neighbor] {
			// if neighbor is in the current recursion stack ,we found a back-edge (cycle)
			return true
		}
	}
	stack[node] = false
	return false
}

// TopologicalSort return the list of resource in the order they must be created
func (graph *DependencyGraph) TopologicalSort() ([]*Resource, error) {
	var sorted []*Resource

	visited := make(map[string]bool)
	tempMark := make(map[string]bool)

	var visit func(string) error

	visit = func(NodeId string) error {
		if tempMark[NodeId] {
			return fmt.Errorf("cycle detected invloving resource %s", NodeId)
		}
		if visited[NodeId] {
			return nil
		}
		tempMark[NodeId] = true

		// let's visit the parent dependencies first
		res := graph.Nodes[NodeId]
		for _, depID := range res.DependsOn {
			visitingError := visit(depID)
			if visitingError != nil {
				return visitingError
			}
		}
		tempMark[NodeId] = false
		visited[NodeId] = true
		sorted = append(sorted, graph.Nodes[NodeId])
		return nil
	}

	for NodeId := range graph.Nodes {
		if !visited[NodeId] {
			visitingError := visit(NodeId)
			if visitingError != nil {
				return nil, visitingError
			}
		}
	}
	return sorted, nil
}

// ReverseTopologicalSort is to return the list of resources in the order they must be destroyed
func (graph *DependencyGraph) ReverseTopologicalSort() ([]*Resource, error) {
	// We are first of all going to get the list of creating through topological sort
	CreationOrder, TopologicalSortError := graph.TopologicalSort()
	if TopologicalSortError != nil {
		return nil, fmt.Errorf("-> an error occured when trying to obtain the execution order list : %s", TopologicalSortError)
	}

	// we are then going to reverse that list
	sizeExecutionList := len(CreationOrder)
	reversedExecutionList := make([]*Resource, sizeExecutionList)
	for i, resource := range CreationOrder {
		reversedExecutionList[n-1-i] = resource
	}
	return reversedExecutionList, nil
}
