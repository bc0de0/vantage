# Vantage Research-Grade Campaign Simulator Notes

## Beam search rationale
Vantage now expands attack and campaign candidates with deterministic beam search. At each depth, only top-K candidates are retained (default K=25). This keeps planning tractable while preserving high-value branches.

## Goal bias philosophy
Campaign and path scoring include objective proximity, prioritizing branches that either directly produce the requested objective or reduce estimated distance to it.

## Memory model semantics
State now tracks previous actions, accumulated exposure knowledge, and failed attempts. Repeated failures reduce confidence for repeated actions; repeated reconnaissance increases unlock potential.

## AI overlay architecture
CLI commands serialize top campaigns into structured JSON and emit advisory explanations (realism, top-3 comparison, and defensive implications). Core deterministic reasoning remains unchanged and works without AI mode enabled.

## Scalability strategy
A synthetic 1000-technique stress suite validates bounded runtime and controlled branching under beam search and memory-aware scoring.
