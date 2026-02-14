package knowledge

import (
	"errors"
	"fmt"
	"sync"
)

// NodeType identifies first-class entities in the operational knowledge graph.
type NodeType string

const (
	NodeOrganization   NodeType = "Organization"
	NodeUser           NodeType = "User"
	NodeAsset          NodeType = "Asset"
	NodeCredential     NodeType = "Credential"
	NodeService        NodeType = "Service"
	NodeVulnerability  NodeType = "Vulnerability"
	NodeInfrastructure NodeType = "Infrastructure"
	NodePattern        NodeType = "Pattern"
	NodeBehavior       NodeType = "Behavior"
)

// EdgeType represents relationship semantics used during reasoning.
type EdgeType string

const (
	EdgeOwns            EdgeType = "OWNS"
	EdgeExposes         EdgeType = "EXPOSES"
	EdgeReuses          EdgeType = "REUSES"
	EdgeTrusts          EdgeType = "TRUSTS"
	EdgeCommunicates    EdgeType = "COMMUNICATES_WITH"
	EdgeAuthenticatesTo EdgeType = "AUTHENTICATES_TO"
	EdgeLikelyLinked    EdgeType = "LIKELY_LINKED"
)

type Node struct {
	ID         string
	Type       NodeType
	Properties map[string]string
}

type Edge struct {
	From       string
	To         string
	Type       EdgeType
	Properties map[string]string
}

// Graph is an in-memory, concurrency-safe knowledge graph.
type Graph struct {
	mu    sync.RWMutex
	nodes map[string]Node
	edges []Edge
}

func NewGraph() *Graph {
	return &Graph{nodes: make(map[string]Node), edges: make([]Edge, 0)}
}

func (g *Graph) AddNode(node Node) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if node.ID == "" {
		return errors.New("node ID is required")
	}
	if _, exists := g.nodes[node.ID]; exists {
		return fmt.Errorf("node already exists: %s", node.ID)
	}
	g.nodes[node.ID] = node
	return nil
}

func (g *Graph) AddEdge(edge Edge) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, ok := g.nodes[edge.From]; !ok {
		return fmt.Errorf("unknown source node: %s", edge.From)
	}
	if _, ok := g.nodes[edge.To]; !ok {
		return fmt.Errorf("unknown destination node: %s", edge.To)
	}

	g.edges = append(g.edges, edge)
	return nil
}

func (g *Graph) Snapshot() ([]Node, []Edge) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	nodes := make([]Node, 0, len(g.nodes))
	for _, node := range g.nodes {
		nodes = append(nodes, node)
	}
	edges := append([]Edge(nil), g.edges...)
	return nodes, edges
}
