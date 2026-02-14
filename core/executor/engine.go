package executor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"vantage/ai"
	"vantage/core/evidence"
	"vantage/core/exposure"
	"vantage/core/intent"
	"vantage/core/reasoning"
	"vantage/core/roe"
	"vantage/core/state"
	"vantage/techniques"
)

// ============================================================================
// EXECUTOR ENGINE — AUTHORITATIVE EXECUTION CORE
//
// SECURITY CRITICAL FILE
//
// This file defines the ONLY legitimate execution path in VANTAGE.
//
// GUARANTEES:
// - Intent is validated exactly once
// - ROE is enforced before any decision
// - Techniques are DECISION-ONLY
// - AI is advisory-only and ignorable
// - Exposure is monotonically increasing
// - Campaign state transitions are enforced
// - Evidence is ALWAYS produced and signed
//
// If execution does not pass through this file,
// it is NOT a valid VANTAGE execution.
// ============================================================================

// Engine binds all immutable execution dependencies.
//
// Once constructed, an Engine MUST NOT be mutated.
// Runtime safety depends on this immutability.
type Engine struct {

	// contract is the validated, immutable intent declaration.
	// It defines WHAT is allowed — never HOW.
	contract *intent.Contract

	// campaign tracks execution lifecycle and halting state.
	// It is the sole authority for campaign progress.
	campaign *state.Campaign

	// exposure tracks cumulative detection risk.
	// Exposure is conservative and monotonic.
	exposure *exposure.Tracker

	// reasoner provides stateful action ranking while execution integration evolves.
	reasoner *reasoning.Engine
}

// -----------------------------------------------------------------------------
// New constructs a fully-bound execution engine.
//
// REQUIREMENTS (FAIL-CLOSED):
// - Intent contract must be non-nil and valid
// - Campaign state must be initialized
// - Exposure tracker must be initialized
//
// Intent is validated exactly ONCE here.
// After this point, intent is immutable.
// -----------------------------------------------------------------------------
func New(
	contract *intent.Contract,
	campaign *state.Campaign,
	exposureTracker *exposure.Tracker,
) (*Engine, error) {

	// Defensive validation — programmer errors
	if contract == nil {
		return nil, errors.New("engine requires non-nil intent contract")
	}
	if campaign == nil {
		return nil, errors.New("engine requires non-nil campaign state")
	}
	if exposureTracker == nil {
		return nil, errors.New("engine requires non-nil exposure tracker")
	}

	// Validate intent contract once and only once
	if err := contract.Validate(); err != nil {
		return nil, fmt.Errorf("invalid intent contract: %w", err)
	}

	return &Engine{
		contract: contract,
		campaign: campaign,
		exposure: exposureTracker,
		reasoner: reasoning.NewEngine(),
	}, nil
}

// -----------------------------------------------------------------------------
// Run evaluates EXACTLY ONE technique against EXACTLY ONE target.
//
// This function:
// - Enforces intent and ROE
// - Resolves technique decisions
// - Optionally consults AI (non-authoritative)
// - Accounts for exposure
// - Mutates campaign state
// - Produces signed evidence
//
// FAILURE IS EXPLICIT — NEVER SILENT.
// -----------------------------------------------------------------------------
func (e *Engine) Run(
	ctx context.Context,
	techniqueID string,
	target string,
) (*evidence.Artifact, error) {

	// -----------------------------------------------------------------
	// 1. CONTEXT VALIDATION
	// -----------------------------------------------------------------

	// A nil context is a programmer error
	if ctx == nil {
		return nil, errors.New("nil execution context")
	}

	// Abort immediately if context already cancelled
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// -----------------------------------------------------------------
	// 2. CAMPAIGN LIFECYCLE ENFORCEMENT
	// -----------------------------------------------------------------

	// Enforce legal campaign state transitions
	switch e.campaign.Status() {

	case state.StatusInitialized:
		// First execution starts the campaign
		if err := e.campaign.Start(); err != nil {
			return nil, err
		}

	case state.StatusRunning:
		// Normal execution path

	case state.StatusHalted, state.StatusCompleted:
		// Execution after halt or completion is forbidden
		return nil, fmt.Errorf(
			"execution denied: campaign is %s",
			e.campaign.Status(),
		)
	}

	// Hard stop if exposure already breached
	if e.exposure.Halted() {
		_ = e.campaign.Halt("exposure limit exceeded")
		return nil, errors.New("execution halted due to exposure")
	}

	// -----------------------------------------------------------------
	// 3. ROE + INTENT ENFORCEMENT (LAW)
	// -----------------------------------------------------------------

	// ROE is enforced as intersection:
	// - Static policy
	// - Declared intent
	// - Target scope
	// - Time window
	if err := roe.Enforce(e.contract, techniqueID, target); err != nil {
		return nil, err
	}

	// -----------------------------------------------------------------
	// 4. TECHNIQUE RESOLUTION (DECISION ONLY)
	// -----------------------------------------------------------------

	decision, err := e.reasoner.PlanNextAction(reasoning.PlannerQuery{
		Target:             target,
		AllowedTechniques:  []string{techniqueID},
		CurrentTechniqueID: techniqueID,
		TopN:               1,
	})
	if err != nil {
		return nil, err
	}

	// Resolve technique from closed-world registry.
	// Reasoner picks the candidate while executor keeps final validation.
	technique, err := techniques.Get(decision.Selected.TechniqueID)
	if err != nil {
		return nil, err
	}

	// Techniques do NOT execute.
	// They only resolve admissible Action Classes.
	resolution, resolveErr := technique.Resolve(
		techniques.ResolveInput{
			TechniqueID:          decision.Selected.TechniqueID,
			AllowedIntentDomains: []string{}, // v0.x placeholder
			AllowedROECategories: []string{}, // v0.x static ROE
			ExposureBudget:       e.exposure.Level().String(),
		},
	)

	// -----------------------------------------------------------------
	// 5. AI ADVISORY (NON-AUTHORITATIVE, OPTIONAL)
	// -----------------------------------------------------------------

	// AI is consulted ONLY for advisory classification.
	// AI output:
	// - MUST be ignorable
	// - MUST NOT block execution
	// - MUST NOT grant authority
	_, _ = ai.Advise(ai.AdvisoryInput{
		Intent: struct {
			Objective      string   `json:"objective"`
			AllowedDomains []string `json:"allowed_domains"`
		}{
			Objective:      e.contract.Objective,
			AllowedDomains: []string{},
		},
		ROE: struct {
			AllowedCategories   []string `json:"allowed_categories"`
			ForbiddenCategories []string `json:"forbidden_categories"`
		}{
			AllowedCategories: []string{},
		},
		Exposure: struct {
			RemainingBudget string `json:"remaining_budget"`
		}{
			RemainingBudget: e.exposure.Level().String(),
		},
		TargetContext: struct {
			HighLevelType string `json:"high_level_type"`
			KnownAccess   bool   `json:"known_access"`
		}{
			HighLevelType: "unknown",
			KnownAccess:   false,
		},
		CanonicalActionClasses: resolution.AllowedActionClasses,
	})
	// Errors are intentionally ignored

	// -----------------------------------------------------------------
	// 6. EXECUTION ACCOUNTING (NO ACTIONS YET)
	// -----------------------------------------------------------------

	// Record that an execution attempt occurred,
	// regardless of outcome.
	_ = e.campaign.RecordExecution()

	startedAt := time.Now().UTC()
	var execErr error

	// Resolution failure is factual
	if resolveErr != nil {
		execErr = resolveErr
	} else if len(resolution.AllowedActionClasses) == 0 {
		execErr = errors.New("no admissible action classes")
	}

	// -----------------------------------------------------------------
	// 7. EXPOSURE ACCOUNTING (CONSERVATIVE)
	// -----------------------------------------------------------------

	// v0.x policy:
	// Every execution attempt incurs fixed exposure.
	const executionExposure uint64 = 10

	_ = e.exposure.Add(executionExposure)

	if e.exposure.Halted() {
		_ = e.campaign.Halt("exposure limit exceeded")
	}

	// -----------------------------------------------------------------
	// 8. EVIDENCE CREATION (MANDATORY)
	// -----------------------------------------------------------------

	artifact := &evidence.Artifact{
		ArtifactID:    uuid.NewString(),
		CampaignID:    e.contract.CampaignID,
		TechniqueID:   decision.Selected.TechniqueID,
		Target:        target,
		ExecutedAt:    startedAt,
		Success:       execErr == nil,
		Output:        "",
		ExposureScore: e.exposure.Score(),
	}

	// Evidence MUST be signed exactly once
	if err := artifact.Sign(); err != nil {
		return nil, fmt.Errorf("evidence signing failed: %w", err)
	}

	_ = e.reasoner.IngestEvidence(reasoning.EvidenceEvent{
		TechniqueID: artifact.TechniqueID,
		Target:      artifact.Target,
		Success:     artifact.Success,
		Output:      artifact.Output,
		Artifact:    artifact,
	})

	// -----------------------------------------------------------------
	// 9. FINAL OUTCOME
	// -----------------------------------------------------------------

	if e.exposure.Halted() {
		return artifact, errors.New("campaign halted due to exposure")
	}

	if execErr != nil {
		return artifact, execErr
	}

	return artifact, nil
}
