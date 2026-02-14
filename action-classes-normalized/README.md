# VANTAGE — Canonical Action Classes

This directory contains the authoritative, normalized set of Action Classes
used by VANTAGE.

Rules:
- Exactly one file per Action Class ID (AC-01 → AC-15)
- Filenames are canonical and stable (`AC-XX-Descriptive-Title.yaml`)
- File `id` values must exactly match filename prefixes (`AC-XX`)
- Techniques may only reference Action Class IDs from this set
- Execution implementations must map 1:1 to these definitions

This directory is the only runtime source for action-class loading.
