package domain

import "testing"

func TestCycleDetection(t *testing.T) {
	// Scenario: A depends on B, B depends on A (Deadlock)
	g := NewDependencyGraph()

	resA := Resource{ID: "A", DependsOn: []string{"B"}}
	resB := Resource{ID: "B", DependsOn: []string{"A"}}

	// Pass by reference (&resA)
	if err := g.AddResource(&resA); err != nil {
		t.Fatalf("Failed to add resource A: %v", err)
	}
	if err := g.AddResource(&resB); err != nil {
		t.Fatalf("Failed to add resource B: %v", err)
	}

	// This should succeed if logic is correct
	if err := g.BuildEdges(); err != nil {
		t.Fatalf("BuildEdges failed unexpectedly: %v", err)
	}

	// Now check for cycles
	err := g.DetectCycles()
	if err == nil {
		t.Fatal("Expected cycle error, but got nil! Logic is broken.")
	}
}

func TestValidGraph(t *testing.T) {
	// Scenario: VLAN10 -> PortA (Valid)
	g := NewDependencyGraph()

	resOne := Resource{ID: "Vlan-10"}
	resTwo := Resource{ID: "Port-A", DependsOn: []string{"Vlan-10"}}

	// Strict Error Checking
	if err := g.AddResource(&resOne); err != nil {
		t.Fatal(err)
	}
	if err := g.AddResource(&resTwo); err != nil {
		t.Fatal(err)
	}

	if err := g.BuildEdges(); err != nil {
		t.Fatalf("BuildEdges failed: %v", err)
	}

	if err := g.DetectCycles(); err != nil {
		t.Fatalf("Expected valid graph, but got error: %v", err)
	}
}
