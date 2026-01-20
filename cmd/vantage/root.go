package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// -----------------------------------------------------------------------------
// ROOT COMMAND — CLI ENTRY POINT
//
// The CLI is a THIN orchestration layer.
//
// It MUST NOT:
// - execute techniques directly
// - implement policy
// - mutate state
//
// It MAY:
// - load intent
// - construct required objects
// - invoke the executor
// -----------------------------------------------------------------------------

var rootCmd = &cobra.Command{
	Use:   "vantage",
	Short: "VANTAGE — doctrine-driven adversary execution platform",
	Long: `
VANTAGE is a policy-enforced adversary execution system.

Execution is permitted only when:
- intent is explicitly declared
- ROE permits the action
- exposure limits are respected
- evidence is produced

This CLI is an orchestration layer only.
`,
	SilenceUsage:  true, // Do not print usage on execution errors
	SilenceErrors: true, // Errors are returned explicitly
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
