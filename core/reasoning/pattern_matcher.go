package reasoning

// MatchPatterns checks whether all graph patterns are satisfied by existing nodes and edges.
// A pattern is satisfied when each required node type and edge type exists at least once.
func MatchPatterns(graph *Graph, patterns []GraphPattern) bool {
	if graph == nil {
		return false
	}
	for _, pattern := range patterns {
		for _, nodeType := range pattern.RequiredNodeTypes {
			if len(graph.NodesByType(nodeType)) == 0 {
				return false
			}
		}
		for _, edgeType := range pattern.RequiredEdges {
			if !graph.HasEdgeType(edgeType) {
				return false
			}
		}
	}
	return true
}
