package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"vantage/core/executor"
	"vantage/core/exposure"
	"vantage/core/intent"
	"vantage/core/state"

	"github.com/spf13/cobra"
)

// -----------------------------------------------------------------------------
// EXECUTE COMMAND — CONTROLLED EXECUTION ENTRY
//
// This command:
// - Loads intent
// - Constructs state + exposure
// - Instantiates the executor
// - Executes exactly ONE technique
//
// This command CANNOT:
// - bypass ROE
// - bypass intent
// - bypass exposure
// -----------------------------------------------------------------------------

var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Execute techniques under declared intent and ROE",
}

// runCmd executes a single technique once.
// Loops and batch execution are intentionally excluded in v0.x.
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a single registered technique against a single target",
	RunE: func(cmd *cobra.Command, args []string) error {

		// -----------------------------------------------------------------
		// 1. CLI Argument Retrieval
		// -----------------------------------------------------------------

		techniqueID, _ := cmd.Flags().GetString("technique")
		target, _ := cmd.Flags().GetString("target")
		campaignID, _ := cmd.Flags().GetString("campaign")

		if techniqueID == "" || target == "" || campaignID == "" {
			return errors.New("technique, target, and campaign are required")
		}

		// -----------------------------------------------------------------
		// 2. Intent Construction (v0.x: CLI-supplied)
		// -----------------------------------------------------------------

		// NOTE:
		// In v0.x, intent is constructed from CLI arguments.
		// In v1.x, this will be loaded from a signed intent file.
		contract := &intent.Contract{
			CampaignID:        campaignID,
			Objective:         "CLI-triggered execution",
			AllowedTechniques: []string{techniqueID},
			Targets:           []string{target},
			NotBefore:         time.Now().UTC().Add(-1 * time.Minute),
			NotAfter:          time.Now().UTC().Add(10 * time.Minute),
		}

		// -----------------------------------------------------------------
		// 3. Campaign State Initialization
		// -----------------------------------------------------------------

		campaign, err := state.New(contract.CampaignID)
		if err != nil {
			return err
		}

		// -----------------------------------------------------------------
		// 4. Exposure Tracker Initialization
		// -----------------------------------------------------------------

		// v0.x policy: fixed maximum exposure
		exposureTracker, err := exposure.New(100)
		if err != nil {
			return err
		}

		// -----------------------------------------------------------------
		// 5. Executor Construction (BOUND, IMMUTABLE)
		// -----------------------------------------------------------------

		engine, err := executor.New(contract, campaign, exposureTracker)
		if err != nil {
			return err
		}

		// -----------------------------------------------------------------
		// 6. Bounded Execution Context
		// -----------------------------------------------------------------

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		artifact, err := engine.Run(ctx, techniqueID, target)

		// -----------------------------------------------------------------
		// 7. Outcome Reporting (CLI ONLY)
		// -----------------------------------------------------------------

		if artifact != nil {
			fmt.Printf(
				"[+] Evidence created | id=%s success=%v exposure=%d\n",
				artifact.ArtifactID,
				artifact.Success,
				artifact.ExposureScore,
			)
		}

		return err
	},
}

func init() {

	// -----------------------------
	// Flags
	// -----------------------------

	runCmd.Flags().String("technique", "", "Technique ID (e.g., T1595)")
	runCmd.Flags().String("target", "", "Target identifier")
	runCmd.Flags().String("campaign", "", "Campaign identifier")

	_ = runCmd.MarkFlagRequired("technique")
	_ = runCmd.MarkFlagRequired("target")
	_ = runCmd.MarkFlagRequired("campaign")

	// -----------------------------
	// Command Wiring
	// -----------------------------

	executeCmd.AddCommand(runCmd)
	rootCmd.AddCommand(executeCmd)
}
