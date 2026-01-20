package executor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"vantage/core/evidence"
	"vantage/core/exposure"
	"vantage/core/intent"
	"vantage/core/roe"
	"vantage/core/state"
	"vantage/techniques"
)

// -----------------------------------------------------------------------------
// EXECUTOR ENGINE — AUTHORITATIVE EXECUTION CORE
//
// THIS FILE IS SECURITY-CRITICAL.
//
// This is the ONLY place in VANTAGE where:
// - Techniques are executed
// - Campaign state mutates
// - Exposure accumulates
// - Evidence artifacts are created and signed
//
// If execution does not pass through this file,
// it is NOT a valid VANTAGE execution.
//
// DESIGN DOCTRINE:
//
// 1. SINGLE CHOKE POINT
//    All execution must pass through Engine.Run.
//
// 2. POLICY FIRST
//    Intent + ROE are enforced BEFORE execution.
//
// 3. TECHNIQUES ARE UNTRUSTED
//    Techniques may fail, misbehave, or panic.
//    The engine protects the system from them.
//
// 4. EVIDENCE IS MANDATORY
//    Every execution attempt produces signed evidence.
// -----------------------------------------------------------------------------

// Engine binds together the immutable components required
// to safely execute adversary techniques.
type Engine struct {

	// contract is the validated intent declaration.
	// It is immutable for the lifetime of the engine.
	contract *intent.Contract

	// campaign tracks execution lifecycle and halting state.
	state *state.Campaign

	// exposure tracks cumulative detection risk.
	exposure *exposure.Tracker
}

// New constructs a fully bound executor engine.
//
// REQUIREMENTS:
// - Intent contract MUST be valid
// - Campaign state MUST be initialized
// - Exposure tracker MUST be configured
//
// Once created, the engine is immutable.
func New(
	contract *intent.Contract,
	campaign *state.Campaign,
	exposureTracker *exposure.Tracker,
) (*Engine, error) {

	// -----------------------------
	// Defensive Validation
	// -----------------------------

	if contract == nil {
		return nil, errors.New("executor requires non-nil intent contract")
	}
	if campaign == nil {
		return nil, errors.New("executor requires non-nil campaign state")
	}
	if exposureTracker == nil {
		return nil, errors.New("executor requires non-nil exposure tracker")
	}

	// Intent is validated exactly once here.
	// Execution is impossible without valid intent.
	if err := contract.Validate(); err != nil {
		return nil, fmt.Errorf("invalid intent contract: %w", err)
	}

	return &Engine{
		contract: contract,
		state:    campaign,
		exposure: exposureTracker,
	}, nil
}

// Run executes EXACTLY ONE technique against EXACTLY ONE target.
//
// This function:
// - Enforces intent and ROE
// - Records campaign state
// - Accumulates exposure
// - Produces signed evidence
//
// FAILURE IS FACTUAL, NOT SILENT.
func (e *Engine) Run(
	ctx context.Context,
	techniqueID string,
	target string,
) (*evidence.Artifact, error) {

	// -----------------------------------------------------------------
	// 1. Context Validation
	// -----------------------------------------------------------------

	if ctx == nil {
		return nil, errors.New("executor invoked with nil context")
	}

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("execution aborted before start: %w", ctx.Err())
	default:
	}

	// -----------------------------------------------------------------
	// 2. Campaign Lifecycle Enforcement
	// -----------------------------------------------------------------

	switch e.state.Status() {
	case state.StatusInitialized:
		// First execution starts the campaign
		if err := e.state.Start(); err != nil {
			return nil, err
		}
	case state.StatusRunning:
		// Allowed
	case state.StatusHalted, state.StatusCompleted:
		return nil, fmt.Errorf(
			"execution denied: campaign is %s",
			e.state.Status().String(),
		)
	}

	// Hard stop if exposure already breached
	if e.exposure.Halted() {
		_ = e.state.Halt("exposure limit exceeded")
		return nil, errors.New("execution halted due to exposure")
	}

	// -----------------------------------------------------------------
	// 3. ROE + Intent Enforcement (LAW)
	// -----------------------------------------------------------------

	if err := roe.Enforce(e.contract, techniqueID, target); err != nil {
		return nil, err
	}

	// -----------------------------------------------------------------
	// 4. Technique Resolution (Closed World)
	// -----------------------------------------------------------------

	technique, err := techniques.Get(techniqueID)
	if err != nil {
		return nil, err
	}

	// -----------------------------------------------------------------
	// 5. Execution Accounting
	// -----------------------------------------------------------------

	// Record that an execution attempt is occurring,
	// regardless of outcome.
	_ = e.state.RecordExecution()

	// Bound execution lifetime.
	execCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	startedAt := time.Now().UTC()
	var execErr error

	// -----------------------------------------------------------------
	// 6. Execute Untrusted Technique
	// -----------------------------------------------------------------

	execErr = technique.Execute(execCtx, target)

	// -----------------------------------------------------------------
	// 7. Exposure Accounting (Conservative)
	// -----------------------------------------------------------------

	// v0.x policy:
	// Each execution attempt incurs fixed exposure.
	// This will become technique-specific later.
	const executionExposure uint64 = 10

	_ = e.exposure.Add(executionExposure)

	if e.exposure.Halted() {
		_ = e.state.Halt("exposure limit exceeded")
	}

	// -----------------------------------------------------------------
	// 8. Evidence Creation (MANDATORY)
	// -----------------------------------------------------------------

	artifact := &evidence.Artifact{
		ArtifactID:    uuid.NewString(),
		CampaignID:    e.contract.CampaignID,
		TechniqueID:   techniqueID,
		Target:        target,
		ExecutedAt:    startedAt,
		Success:       execErr == nil,
		Output:        "", // v0.x: no stdout capture
		ExposureScore: e.exposure.Score(),
	}

	// Evidence MUST be signed exactly once.
	if err := artifact.Sign(); err != nil {
		return nil, fmt.Errorf("evidence signing failed: %w", err)
	}

	// -----------------------------------------------------------------
	// 9. Final Outcome
	// -----------------------------------------------------------------

	if e.exposure.Halted() {
		return artifact, errors.New("campaign halted due to exposure")
	}

	if execErr != nil {
		return artifact, execErr
	}

	return artifact, nil
}
