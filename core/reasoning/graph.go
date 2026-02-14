package reasoning

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

// NodeType identifies the semantic role of a graph node.
type NodeType string

const (
	// NodeTypeEvidence represents observed execution evidence.
	NodeTypeEvidence NodeType = "evidence"
	// NodeTypeHypothesis represents generated hypotheses.
	NodeTypeHypothesis NodeType = "hypothesis"
	// NodeTypeAttackPath represents potential attack path steps.
	NodeTypeAttackPath NodeType = "attack_path"
	// NodeTypeTechnique represents a technique option.
	NodeTypeTechnique NodeType = "technique"
)

// Node stores a fact in the operational reasoning graph.
type Node struct {
	ID        string
	Type      NodeType
	Label     string
	CreatedAt time.Time
	Metadata  map[string]string
}

// EdgeType identifies how two graph nodes are related.
type EdgeType string

const (
	// EdgeTypeSupports links evidence to hypotheses or paths it supports.
	EdgeTypeSupports EdgeType = "supports"
	// EdgeTypeEnables links facts to reachable next actions.
	EdgeTypeEnables EdgeType = "enables"
	// EdgeTypeRefines links a node to a more specific node.
	EdgeTypeRefines EdgeType = "refines"
)

// Edge connects two nodes in the reasoning graph.
type Edge struct {
	From      string
	To        string
	Type      EdgeType
	Weight    float64
	CreatedAt time.Time
}

// Graph is an in-memory operational graph of evidence and hypotheses.
type Graph struct {
	mu    sync.RWMutex
	nodes map[string]*Node
	edges []*Edge
}

// NewGraph constructs an empty operational graph.
func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[string]*Node),
		edges: make([]*Edge, 0),
	}
}

// UpsertNode inserts or updates a node by ID.
func (g *Graph) UpsertNode(node *Node) {
	if node == nil || node.ID == "" {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if node.CreatedAt.IsZero() {
		node.CreatedAt = time.Now().UTC()
	}
	if node.Metadata == nil {
		node.Metadata = map[string]string{}
	}
	g.nodes[node.ID] = node
}

// AddEdge appends an edge if both endpoint nodes exist.
func (g *Graph) AddEdge(edge *Edge) error {
	if edge == nil || edge.From == "" || edge.To == "" {
		return fmt.Errorf("invalid edge")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.nodes[edge.From]; !ok {
		return fmt.Errorf("from node not found: %s", edge.From)
	}
	if _, ok := g.nodes[edge.To]; !ok {
		return fmt.Errorf("to node not found: %s", edge.To)
	}
	if edge.CreatedAt.IsZero() {
		edge.CreatedAt = time.Now().UTC()
	}
	g.edges = append(g.edges, edge)
	return nil
}

// Node returns a copy-safe pointer to a node by ID.
func (g *Graph) Node(id string) (*Node, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	n, ok := g.nodes[id]
	return n, ok
}

// NodesByType returns all nodes with a specific type.
func (g *Graph) NodesByType(nodeType NodeType) []*Node {
	g.mu.RLock()
	defer g.mu.RUnlock()
	out := make([]*Node, 0)
	for _, n := range g.nodes {
		if n.Type == nodeType {
			out = append(out, n)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// EdgesFrom returns all edges originating at a node.
func (g *Graph) EdgesFrom(nodeID string) []*Edge {
	g.mu.RLock()
	defer g.mu.RUnlock()
	out := make([]*Edge, 0)
	for _, e := range g.edges {
		if e.From == nodeID {
			out = append(out, e)
		}
	}
	return out
}

// HasEdgeType returns true when at least one edge of the requested type exists.
func (g *Graph) HasEdgeType(edgeType EdgeType) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, e := range g.edges {
		if e.Type == edgeType {
			return true
		}
	}
	return false
}

// ToDOT renders the graph as Graphviz DOT text.
func (g *Graph) ToDOT() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	ids := make([]string, 0, len(g.nodes))
	for id := range g.nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	var b strings.Builder
	b.WriteString("digraph reasoning {\n")
	for _, id := range ids {
		n := g.nodes[id]
		b.WriteString(fmt.Sprintf("  \"%s\" [label=\"%s\\n(%s)\"];\n", n.ID, escapeDOT(n.Label), n.Type))
	}
	for _, e := range g.edges {
		b.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\" [label=\"%s:%.2f\"];\n", e.From, e.To, e.Type, e.Weight))
	}
	b.WriteString("}\n")
	return b.String()
}

func escapeDOT(in string) string {
	return strings.ReplaceAll(in, "\"", "\\\"")
}
