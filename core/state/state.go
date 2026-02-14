package state

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// -----------------------------------------------------------------------------
// Campaign State — EXECUTION GOVERNANCE CORE
//
// This module tracks the lifecycle of a single campaign execution.
//
// Campaign state exists to:
// - Enforce halting conditions
// - Track execution progress
// - Prevent uncontrolled or runaway operations
//
// Campaign state is:
// - Authoritative
// - Thread-safe
// - Immutable in intent, mutable only in execution facts
//
// DESIGN PRINCIPLES:
//
// 1. SINGLE CAMPAIGN, SINGLE STATE
//    One campaign == one state object.
//
// 2. FAIL-CLOSED HALTING
//    Once halted, execution must not resume.
//
// 3. NO SIDE EFFECTS
//    State does not perform execution, logging, or evidence creation.
//
// 4. AUDITABLE TRANSITIONS
//    All state changes are explicit and reviewable.
// -----------------------------------------------------------------------------

// Status represents the lifecycle phase of a campaign.
//
// The ordering is intentional and meaningful.
// Transitions MUST move forward only.
type Status int

const (
	// StatusInitialized indicates the campaign state
	// has been created but no execution has occurred.
	StatusInitialized Status = iota

	// StatusRunning indicates at least one technique
	// has been executed or attempted.
	StatusRunning

	// StatusHalted indicates execution has been forcibly stopped
	// due to policy, error, or operator action.
	StatusHalted

	// StatusCompleted indicates all intended execution
	// has finished without violation.
	StatusCompleted
)

// String returns a human-readable status value.
// Used for reporting and diagnostics.
func (s Status) String() string {
	switch s {
	case StatusInitialized:
		return "initialized"
	case StatusRunning:
		return "running"
	case StatusHalted:
		return "halted"
	case StatusCompleted:
		return "completed"
	default:
		return "unknown"
	}
}

// Campaign represents the authoritative execution state
// for a single campaign.
type Campaign struct {

	// campaignID is immutable and binds state to intent.
	campaignID string

	// status tracks the lifecycle phase.
	status Status

	// startedAt records when execution began.
	startedAt time.Time

	// finishedAt records when execution ended.
	finishedAt time.Time

	// executions counts how many techniques were attempted.
	executions uint64

	// memory for multi-cycle adaptation.
	previousActions   []string
	exposureKnowledge map[string]float64
	failedAttempts    map[string]int

	// mu protects all mutable fields.
	mu sync.RWMutex
}

// State is an alias used by components that reason over campaign lifecycle state.
type State = Campaign

// New creates a new campaign state instance.
//
// This function MUST be called exactly once per campaign.
// The returned state starts in StatusInitialized.
func New(campaignID string) (*Campaign, error) {

	if campaignID == "" {
		return nil, errors.New("campaign state requires non-empty campaign ID")
	}

	return &Campaign{
		campaignID:        campaignID,
		status:            StatusInitialized,
		previousActions:   make([]string, 0),
		exposureKnowledge: make(map[string]float64),
		failedAttempts:    make(map[string]int),
	}, nil
}

// Start transitions the campaign into the running state.
//
// This MUST be called immediately before the first execution.
func (c *Campaign) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.status != StatusInitialized {
		return fmt.Errorf(
			"cannot start campaign from state %s",
			c.status.String(),
		)
	}

	c.status = StatusRunning
	c.startedAt = time.Now().UTC()

	return nil
}

// RecordExecution increments the execution counter.
//
// This MUST be called once per technique execution attempt,
// regardless of success or failure.
func (c *Campaign) RecordExecution() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.status != StatusRunning {
		return fmt.Errorf(
			"cannot record execution while campaign is %s",
			c.status.String(),
		)
	}

	c.executions++
	return nil
}

// Halt forcibly stops the campaign.
//
// Once halted, the campaign CANNOT be resumed.
func (c *Campaign) Halt(reason string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.status == StatusCompleted {
		return errors.New("cannot halt a completed campaign")
	}

	if c.status == StatusHalted {
		return nil // idempotent
	}

	c.status = StatusHalted
	c.finishedAt = time.Now().UTC()

	return nil
}

// Complete marks the campaign as successfully finished.
//
// This MUST be called only after all execution is done.
func (c *Campaign) Complete() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.status != StatusRunning {
		return fmt.Errorf(
			"cannot complete campaign from state %s",
			c.status.String(),
		)
	}

	c.status = StatusCompleted
	c.finishedAt = time.Now().UTC()

	return nil
}

// Status returns the current campaign status.
func (c *Campaign) Status() Status {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.status
}

// Executions returns the number of execution attempts.
func (c *Campaign) Executions() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.executions
}

// StartedAt returns the campaign start time.
func (c *Campaign) StartedAt() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.startedAt
}

// FinishedAt returns the campaign end time.
// Zero value indicates campaign still running.
func (c *Campaign) FinishedAt() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.finishedAt
}

// CampaignID returns the immutable campaign identifier.
func (c *Campaign) CampaignID() string {
	return c.campaignID
}

// RecordActionMemory tracks action outcomes across cycles.
func (c *Campaign) RecordActionMemory(actionID string, success bool, recon bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if actionID == "" {
		return
	}
	c.previousActions = append(c.previousActions, actionID)
	if !success {
		c.failedAttempts[actionID]++
	}
	if recon {
		c.exposureKnowledge[actionID] += 0.1
	}
}

func (c *Campaign) PreviousActions() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]string, len(c.previousActions))
	copy(out, c.previousActions)
	return out
}

func (c *Campaign) FailedAttempts(actionID string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.failedAttempts[actionID]
}

func (c *Campaign) ExposureKnowledge() map[string]float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[string]float64, len(c.exposureKnowledge))
	for k, v := range c.exposureKnowledge {
		out[k] = v
	}
	return out
}
