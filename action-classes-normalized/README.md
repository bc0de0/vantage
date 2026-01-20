# VANTAGE — Canonical Action Classes

This directory contains the authoritative, normalized set of Action Classes
used by VANTAGE.

Rules:
- Exactly one file per Action Class ID (AC-01 → AC-15)
- Filenames are canonical and stable
- Techniques may only reference Action Class IDs from this set
- Execution implementations must map 1:1 to these definitions

The original `action-classes/` directory is preserved for historical context
and must not be used by the engine.
