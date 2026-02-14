package main

import (
	"errors"
	"fmt"
	"time"

	"vantage/core/executor"
	"vantage/core/exposure"
	"vantage/core/intent"
	"vantage/core/reasoning"
	"vantage/core/state"

	"github.com/spf13/cobra"
)

type runtime struct {
	reasoner *reasoning.Engine
	state    *state.State
}

func buildRuntime(campaignID, target string, techniques []string) (*runtime, error) {
	if campaignID == "" || target == "" {
		return nil, errors.New("campaign and target are required")
	}
	if len(techniques) == 0 {
		return nil, errors.New("at least one technique is required")
	}

	contract := &intent.Contract{
		CampaignID:        campaignID,
		Objective:         "CLI-triggered execution",
		AllowedTechniques: techniques,
		Targets:           []string{target},
		NotBefore:         time.Now().UTC().Add(-1 * time.Minute),
		NotAfter:          time.Now().UTC().Add(10 * time.Minute),
	}

	campaign, err := state.New(contract.CampaignID)
	if err != nil {
		return nil, err
	}
	exposureTracker, err := exposure.New(100)
	if err != nil {
		return nil, err
	}
	execEngine, err := executor.New(contract, campaign, exposureTracker)
	if err != nil {
		return nil, err
	}

	reasoner := reasoning.NewEngine(nil)
	reasoner.ConfigureCycle(reasoning.CycleConfig{
		Target:            target,
		AllowedTechniques: techniques,
		Executor:          execEngine,
		Timeout:           30 * time.Second,
	})

	return &runtime{reasoner: reasoner, state: campaign}, nil
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run one reasoning cycle and execute the selected technique",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetString("campaign")
		target, _ := cmd.Flags().GetString("target")
		techniques, _ := cmd.Flags().GetStringSlice("technique")

		rt, err := buildRuntime(campaignID, target, techniques)
		if err != nil {
			return err
		}
		decision, err := rt.reasoner.RunCycle(rt.state)
		if decision != nil {
			fmt.Printf("[+] selected=%s score=%.2f\n", decision.Selected.TechniqueID, decision.Selected.Score)
		}
		return err
	},
}

var loopCmd = &cobra.Command{
	Use:   "loop",
	Short: "Run multiple deterministic reasoning cycles",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetString("campaign")
		target, _ := cmd.Flags().GetString("target")
		techniques, _ := cmd.Flags().GetStringSlice("technique")
		cycles, _ := cmd.Flags().GetInt("cycles")

		rt, err := buildRuntime(campaignID, target, techniques)
		if err != nil {
			return err
		}

		for i := 0; i < cycles; i++ {
			decision, runErr := rt.reasoner.RunCycle(rt.state)
			if decision != nil {
				fmt.Printf("[%d] selected=%s score=%.2f\n", i+1, decision.Selected.TechniqueID, decision.Selected.Score)
			}
			if runErr != nil {
				return runErr
			}
		}
		return nil
	},
}

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Print the reasoning graph as DOT",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetString("campaign")
		target, _ := cmd.Flags().GetString("target")
		techniques, _ := cmd.Flags().GetStringSlice("technique")

		rt, err := buildRuntime(campaignID, target, techniques)
		if err != nil {
			return err
		}
		if _, err := rt.reasoner.RunCycle(rt.state); err != nil {
			return err
		}
		fmt.Println(rt.reasoner.DOT())
		return nil
	},
}

var explainCmd = &cobra.Command{
	Use:   "explain",
	Short: "Explain ranked candidates for the next cycle",
	RunE: func(cmd *cobra.Command, args []string) error {
		target, _ := cmd.Flags().GetString("target")
		techniques, _ := cmd.Flags().GetStringSlice("technique")

		reasoner := reasoning.NewEngine(nil)
		decision, err := reasoner.PlanNextAction(reasoning.PlannerQuery{Target: target, AllowedTechniques: techniques})
		if err != nil {
			return err
		}
		for i, ranked := range decision.Ranked {
			fmt.Printf("%d. %s score=%.2f %s\n", i+1, ranked.TechniqueID, ranked.Score, ranked.Reason)
		}
		return nil
	},
}

var simulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "Simulate reasoning without executor side effects",
	RunE: func(cmd *cobra.Command, args []string) error {
		target, _ := cmd.Flags().GetString("target")
		techniques, _ := cmd.Flags().GetStringSlice("technique")

		reasoner := reasoning.NewEngine(nil)
		decision, err := reasoner.PlanNextAction(reasoning.PlannerQuery{Target: target, AllowedTechniques: techniques, TopN: 3})
		if err != nil {
			return err
		}
		fmt.Printf("selected=%s score=%.2f\n", decision.Selected.TechniqueID, decision.Selected.Score)
		return nil
	},
}

func init() {
	runCmd.Flags().StringSlice("technique", nil, "Technique IDs (repeatable)")
	runCmd.Flags().String("target", "", "Target identifier")
	runCmd.Flags().String("campaign", "", "Campaign identifier")
	_ = runCmd.MarkFlagRequired("technique")
	_ = runCmd.MarkFlagRequired("target")
	_ = runCmd.MarkFlagRequired("campaign")

	loopCmd.Flags().StringSlice("technique", nil, "Technique IDs (repeatable)")
	loopCmd.Flags().String("target", "", "Target identifier")
	loopCmd.Flags().String("campaign", "", "Campaign identifier")
	loopCmd.Flags().Int("cycles", 3, "Number of cycles")
	_ = loopCmd.MarkFlagRequired("technique")
	_ = loopCmd.MarkFlagRequired("target")
	_ = loopCmd.MarkFlagRequired("campaign")

	graphCmd.Flags().StringSlice("technique", nil, "Technique IDs (repeatable)")
	graphCmd.Flags().String("target", "", "Target identifier")
	graphCmd.Flags().String("campaign", "", "Campaign identifier")
	_ = graphCmd.MarkFlagRequired("technique")
	_ = graphCmd.MarkFlagRequired("target")
	_ = graphCmd.MarkFlagRequired("campaign")

	explainCmd.Flags().StringSlice("technique", nil, "Technique IDs (repeatable)")
	explainCmd.Flags().String("target", "", "Target identifier")
	_ = explainCmd.MarkFlagRequired("technique")
	_ = explainCmd.MarkFlagRequired("target")

	simulateCmd.Flags().StringSlice("technique", nil, "Technique IDs (repeatable)")
	simulateCmd.Flags().String("target", "", "Target identifier")
	_ = simulateCmd.MarkFlagRequired("technique")
	_ = simulateCmd.MarkFlagRequired("target")

	rootCmd.AddCommand(runCmd, loopCmd, graphCmd, explainCmd, simulateCmd)
}
