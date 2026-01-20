package techniqueset

// ============================================================================
// TECHNIQUE SET — CLOSED WORLD DECLARATION
//
// Importing this package registers ALL supported techniques
// via init() side effects.
//
// RULES:
// - No logic
// - No exports
// - No imports except blank technique imports
//
// This package MUST ONLY be imported by binaries (cmd/*).
// ============================================================================

import (
	_ "vantage/techniques/t1595"
)
