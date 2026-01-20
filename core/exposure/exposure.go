package exposure

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// -----------------------------------------------------------------------------
// Exposure Accounting — DETECTION RISK GOVERNANCE
//
// Exposure represents the accumulated likelihood that an adversary action
// has been detected or could reasonably be detected.
//
// Exposure is NOT:
// - stealth scoring
// - evasion modeling
// - alert simulation
//
// Exposure IS:
// - conservative risk accounting
// - execution governance
// - a HALTING signal
//
// DESIGN PRINCIPLES:
//
// 1. CONSERVATIVE BY DEFAULT
//    When in doubt, assume higher exposure.
//
// 2. MONOTONIC INCREASE
//    Exposure can only increase, never decrease.
//
// 3. FAIL-CLOSED HALTING
//    Breach of limits triggers mandatory halt.
//
// 4. DETERMINISTIC
//    Same inputs always produce same exposure.
// -----------------------------------------------------------------------------

// Level represents the qualitative exposure state.
type Level int

const (
	// LevelLow indicates minimal observable activity.
	LevelLow Level = iota

	// LevelMedium indicates noticeable activity
	// that could plausibly trigger monitoring.
	LevelMedium

	// LevelHigh indicates activity very likely
	// to be detected by defensive systems.
	LevelHigh

	// LevelCritical indicates exposure beyond acceptable limits.
	// Execution MUST halt.
	LevelCritical
)

// String returns a human-readable exposure level.
func (l Level) String() string {
	switch l {
	case LevelLow:
		return "low"
	case LevelMedium:
		return "medium"
	case LevelHigh:
		return "high"
	case LevelCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// Tracker maintains cumulative exposure for a campaign.
type Tracker struct {

	// maxScore defines the maximum allowable exposure score.
	// Crossing this value mandates a halt.
	maxScore uint64

	// score tracks cumulative exposure.
	score uint64

	// lastUpdated records when exposure last changed.
	lastUpdated time.Time

	// halted indicates whether exposure has breached limits.
	halted bool

	// mu protects all mutable fields.
	mu sync.RWMutex
}

// New creates a new exposure tracker.
//
// PARAMETERS:
//   - maxScore:
//     Maximum permitted exposure score.
//     Must be > 0.
func New(maxScore uint64) (*Tracker, error) {

	if maxScore == 0 {
		return nil, errors.New("exposure tracker requires maxScore > 0")
	}

	return &Tracker{
		maxScore: maxScore,
		score:    0,
	}, nil
}

// Add increments exposure by the provided delta.
//
// RULES:
// - delta must be > 0
// - exposure is monotonic
// - once halted, further updates are rejected
func (t *Tracker) Add(delta uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if delta == 0 {
		return errors.New("exposure delta must be > 0")
	}

	if t.halted {
		return errors.New("exposure already exceeded; execution halted")
	}

	t.score += delta
	t.lastUpdated = time.Now().UTC()

	if t.score >= t.maxScore {
		t.halted = true
	}

	return nil
}

// Score returns the current exposure score.
func (t *Tracker) Score() uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.score
}

// Level returns the qualitative exposure level.
//
// Thresholds are intentionally coarse and conservative.
func (t *Tracker) Level() Level {
	t.mu.RLock()
	defer t.mu.RUnlock()

	switch {
	case t.score >= t.maxScore:
		return LevelCritical
	case t.score >= t.maxScore*3/4:
		return LevelHigh
	case t.score >= t.maxScore/2:
		return LevelMedium
	default:
		return LevelLow
	}
}

// Halted indicates whether exposure limits have been breached.
func (t *Tracker) Halted() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.halted
}

// LastUpdated returns when exposure was last modified.
func (t *Tracker) LastUpdated() time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.lastUpdated
}

// Snapshot returns a stable view of exposure state.
//
// Used for:
// - reporting
// - evidence generation
// - audit output
func (t *Tracker) Snapshot() Snapshot {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return Snapshot{
		Score:       t.score,
		MaxScore:    t.maxScore,
		Level:       t.Level(),
		Halted:      t.halted,
		LastUpdated: t.lastUpdated,
	}
}

// Snapshot represents an immutable exposure state view.
type Snapshot struct {
	Score       uint64
	MaxScore    uint64
	Level       Level
	Halted      bool
	LastUpdated time.Time
}

// String returns a human-readable summary.
func (s Snapshot) String() string {
	return fmt.Sprintf(
		"exposure=%d/%d level=%s halted=%v",
		s.Score,
		s.MaxScore,
		s.Level.String(),
		s.Halted,
	)
}
