package techniques

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRegisterAllHasUniqueIDs(t *testing.T) {
	all := RegisterAll()
	seen := map[string]struct{}{}
	for id := range all {
		if _, ok := seen[id]; ok {
			t.Fatalf("duplicate technique id %s", id)
		}
		seen[id] = struct{}{}
	}
}

func TestRegisterAllCollectivelyCoversActionClasses(t *testing.T) {
	files, err := filepath.Glob("../action-classes-normalized/AC-*.yaml")
	if err != nil {
		t.Fatalf("glob action classes: %v", err)
	}
	needed := map[string]struct{}{}
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		for _, line := range strings.Split(string(raw), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "id:") {
				id := strings.TrimSpace(strings.TrimPrefix(line, "id:"))
				if strings.HasPrefix(id, "AC-") && len(id) == 5 {
					needed[id] = struct{}{}
				}
			}
		}
	}

	coveredCount := map[string]int{}
	for _, tech := range RegisterAll() {
		coveredCount[tech.ActionClassID()]++
	}
	for id := range needed {
		if coveredCount[id] == 0 {
			t.Fatalf("action class %s is not covered", id)
		}
		if coveredCount[id] < 5 {
			t.Fatalf("action class %s expected at least 5 techniques, got %d", id, coveredCount[id])
		}
	}
}
