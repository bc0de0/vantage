package reasoning

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vantage/core/state"
)

// OperationPhase aliases lifecycle phases used by the shared campaign state.
type OperationPhase = state.OperationPhase

// ActionClass defines a canonical adversarial action class that can produce graph changes.
type ActionClass struct {
	ID              string
	Name            string
	Phase           OperationPhase
	Preconditions   []GraphPattern
	ProducesNodes   []NodeType
	ProducesEdges   []EdgeType
	RiskWeight      float64
	ImpactWeight    float64
	ConfidenceBoost float64
}

// GraphPattern defines structural graph preconditions for an action class.
type GraphPattern struct {
	RequiredNodeTypes []NodeType
	RequiredEdges     []EdgeType
}

type actionClassYAML struct {
	ID            string
	Name          string
	IntentDomains []string
	Preconditions []string
}

// LoadActionClassesFromDir loads all YAML action class definitions from a directory.
func LoadActionClassesFromDir(dir string) ([]ActionClass, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	classes := make([]ActionClass, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			continue
		}
		if strings.HasPrefix(name, "README") || strings.HasPrefix(name, "_") {
			continue
		}
		ac, err := loadActionClassFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		classes = append(classes, ac)
	}
	return classes, nil
}

func loadActionClassFile(path string) (ActionClass, error) {
	raw, err := parseActionClassYAML(path)
	if err != nil {
		return ActionClass{}, err
	}
	if raw.ID == "" {
		return ActionClass{}, fmt.Errorf("action class %s missing id", path)
	}

	phase := inferPhase(raw.IntentDomains)
	patterns := make([]GraphPattern, 0, len(raw.Preconditions))
	for _, pre := range raw.Preconditions {
		if pattern, ok := preconditionPattern(pre); ok {
			patterns = append(patterns, pattern)
		}
	}

	return ActionClass{
		ID:              raw.ID,
		Name:            raw.Name,
		Phase:           phase,
		Preconditions:   patterns,
		ProducesNodes:   []NodeType{NodeTypeEvidence, NodeTypeHypothesis},
		ProducesEdges:   []EdgeType{EdgeTypeSupports},
		RiskWeight:      0.4,
		ImpactWeight:    0.6,
		ConfidenceBoost: 0.1,
	}, nil
}

func parseActionClassYAML(path string) (actionClassYAML, error) {
	f, err := os.Open(path)
	if err != nil {
		return actionClassYAML{}, err
	}
	defer f.Close()

	out := actionClassYAML{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "id":
			out.ID = trimScalar(value)
		case "name":
			out.Name = trimScalar(value)
		case "intent_domains":
			out.IntentDomains = parseInlineList(value)
		case "preconditions":
			out.Preconditions = parseInlineList(value)
		}
	}
	if err := scanner.Err(); err != nil {
		return actionClassYAML{}, err
	}
	return out, nil
}

func trimScalar(v string) string {
	return strings.Trim(strings.TrimSpace(v), "\"'")
}

func parseInlineList(v string) []string {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "[")
	v = strings.TrimSuffix(v, "]")
	if strings.TrimSpace(v) == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		t := trimScalar(p)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

func inferPhase(domains []string) OperationPhase {
	for _, domain := range domains {
		switch strings.ToLower(domain) {
		case "discovery", "enumeration":
			return state.PhaseRecon
		case "access":
			return state.PhaseInitialAccess
		case "validation":
			return state.PhaseLateralMovement
		case "impact":
			return state.PhaseObjective
		}
	}
	return state.PhaseRecon
}

func preconditionPattern(precondition string) (GraphPattern, bool) {
	switch strings.ToLower(precondition) {
	case "network_reachability":
		return GraphPattern{RequiredNodeTypes: []NodeType{NodeTypeEvidence}}, true
	case "credential_material_present":
		return GraphPattern{RequiredNodeTypes: []NodeType{NodeTypeEvidence, NodeTypeTechnique}}, true
	case "access_established":
		return GraphPattern{RequiredNodeTypes: []NodeType{NodeTypeEvidence, NodeTypeHypothesis}}, true
	case "execution_environment":
		return GraphPattern{RequiredEdges: []EdgeType{EdgeTypeEnables}}, true
	case "user_interaction":
		return GraphPattern{RequiredEdges: []EdgeType{EdgeTypeSupports}}, true
	default:
		return GraphPattern{}, false
	}
}
