package main

import (
	"errors"
	"fmt"
	"strings"
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

func parseObjectiveNodeType(raw string) (reasoning.NodeType, error) {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case string(reasoning.NodeTypeDataExposure):
		return reasoning.NodeTypeDataExposure, nil
	case string(reasoning.NodeTypePrivEsc):
		return reasoning.NodeTypePrivEsc, nil
	case string(reasoning.NodeTypeLateralReachability):
		return reasoning.NodeTypeLateralReachability, nil
	default:
		return "", fmt.Errorf("unsupported objective %q", raw)
	}
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
	Short: "Explain planned campaigns and defensive implications",
	RunE: func(cmd *cobra.Command, args []string) error {
		objectiveFlag, _ := cmd.Flags().GetString("objective")
		objective, err := parseObjectiveNodeType(objectiveFlag)
		if err != nil {
			return err
		}
		reasoner := reasoning.NewEngine(nil)
		reasoner.Graph().UpsertNode(&reasoning.Node{ID: "explain-seed", Type: reasoning.NodeTypeEvidence, Label: "seed"})
		campaigns, err := reasoner.PlanCampaign(objective, reasoning.DefaultCampaignOptions())
		if err != nil {
			return err
		}
		fmt.Println(renderCampaignExplanation(objective, campaigns))
		return nil
	},
}

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare top 3 campaigns",
	RunE: func(cmd *cobra.Command, args []string) error {
		objectiveFlag, _ := cmd.Flags().GetString("objective")
		objective, err := parseObjectiveNodeType(objectiveFlag)
		if err != nil {
			return err
		}
		reasoner := reasoning.NewEngine(nil)
		reasoner.Graph().UpsertNode(&reasoning.Node{ID: "compare-seed", Type: reasoning.NodeTypeEvidence, Label: "seed"})
		campaigns, err := reasoner.PlanCampaign(objective, reasoning.CampaignOptions{TopN: 3})
		if err != nil {
			return err
		}
		fmt.Println(renderCampaignExplanation(objective, campaigns))
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

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Plan strategic attack campaigns for a requested objective",
	RunE: func(cmd *cobra.Command, args []string) error {
		objectiveFlag, _ := cmd.Flags().GetString("objective")
		maxDepth, _ := cmd.Flags().GetInt("max-depth")
		riskTolerance, _ := cmd.Flags().GetFloat64("risk")
		confidenceThreshold, _ := cmd.Flags().GetFloat64("confidence")
		beamWidth, _ := cmd.Flags().GetInt("beam-width")

		objective, err := parseObjectiveNodeType(objectiveFlag)
		if err != nil {
			return err
		}

		reasoner := reasoning.NewEngine(nil)
		reasoner.Graph().UpsertNode(&reasoning.Node{ID: "plan-seed", Type: reasoning.NodeTypeEvidence, Label: "planner seed"})
		campaigns, err := reasoner.PlanCampaign(objective, reasoning.CampaignOptions{
			MaxDepth:            maxDepth,
			RiskTolerance:       riskTolerance,
			ConfidenceThreshold: confidenceThreshold,
			BeamWidth:           beamWidth,
		})
		if err != nil {
			return err
		}
		if len(campaigns) == 0 {
			fmt.Println("no campaigns found")
			return nil
		}

		limit := 5
		if len(campaigns) < limit {
			limit = len(campaigns)
		}
		fmt.Println(renderCampaignExplanation(objective, campaigns[:limit]))
		for i := 0; i < limit; i++ {
			campaign := campaigns[i]
			stepIDs := make([]string, 0, len(campaign.Steps))
			for _, step := range campaign.Steps {
				stepIDs = append(stepIDs, step.ActionClassID)
			}
			fmt.Printf("%d. score=%.3f objective=%s attained=%t risk=%.3f confidence=%.3f steps=%s\n", i+1, campaign.Score, campaign.Objective, campaign.Objective == objective, campaign.Risk, campaign.Confidence, strings.Join(stepIDs, " -> "))
		}
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

	explainCmd.Flags().String("objective", "", "Objective node type")
	_ = explainCmd.MarkFlagRequired("objective")

	compareCmd.Flags().String("objective", "", "Objective node type")
	_ = compareCmd.MarkFlagRequired("objective")

	simulateCmd.Flags().StringSlice("technique", nil, "Technique IDs (repeatable)")
	simulateCmd.Flags().String("target", "", "Target identifier")
	_ = simulateCmd.MarkFlagRequired("technique")
	_ = simulateCmd.MarkFlagRequired("target")

	planCmd.Flags().String("objective", "", "Objective node type (DATA_EXPOSURE, PRIV_ESC, LATERAL_REACHABILITY)")
	planCmd.Flags().Int("max-depth", reasoning.DefaultCampaignOptions().MaxDepth, "Maximum campaign depth")
	planCmd.Flags().Float64("risk", reasoning.DefaultCampaignOptions().RiskTolerance, "Maximum cumulative risk tolerance")
	planCmd.Flags().Float64("confidence", reasoning.DefaultCampaignOptions().ConfidenceThreshold, "Minimum average confidence threshold")
	planCmd.Flags().Int("beam-width", reasoning.DefaultCampaignOptions().BeamWidth, "Beam width per depth")
	_ = planCmd.MarkFlagRequired("objective")

	rootCmd.AddCommand(runCmd, loopCmd, graphCmd, explainCmd, simulateCmd, planCmd, compareCmd)
}
