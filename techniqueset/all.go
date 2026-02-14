package techniqueset

import "vantage/techniques"

// ForActionClass returns all registered techniques bound to the provided action class ID.
func ForActionClass(actionClassID string) []techniques.Technique {
	all := techniques.RegisterAll()
	out := make([]techniques.Technique, 0)
	for _, t := range all {
		if t.ActionClassID() == actionClassID {
			out = append(out, t)
		}
	}
	return out
}
