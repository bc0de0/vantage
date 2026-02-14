package main

import (
	"encoding/json"
	"fmt"
	"os"

	"vantage/core/reasoning"
)

type campaignPayload struct {
	Objective string               `json:"objective"`
	Campaigns []reasoning.Campaign `json:"campaigns"`
}

func aiOverlayEnabled() bool {
	return os.Getenv("VANTAGE_AI_MODE") == "advisory_only"
}

func renderCampaignExplanation(objective reasoning.NodeType, campaigns []reasoning.Campaign) string {
	if len(campaigns) == 0 {
		return "No campaigns available."
	}
	if !aiOverlayEnabled() {
		return "AI disabled; deterministic summary only."
	}
	payload := campaignPayload{Objective: string(objective), Campaigns: campaigns}
	b, _ := json.Marshal(payload)
	return fmt.Sprintf("AI advisory payload: %s\nTop campaign realism: grounded in cumulative confidence and low risk.\nTop-3 comparison: first has best score, others offer alternate risk tradeoffs.\nDefensive implications: prioritize detections on early reconnaissance and privilege transitions.", string(b))
}
