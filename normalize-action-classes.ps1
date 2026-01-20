# ============================================================================
# VANTAGE — ACTION CLASS NORMALIZATION SCRIPT (APPLY MODE)
# CANONICAL | FAIL-FAST | AUDITABLE
# ============================================================================

$SourceDir = "action-classes"
$TargetDir = "action-classes-normalized"
$DryRun = $false   # APPLY MODE — CHANGES WILL BE MATERIALIZED

# ---------------------------------------------------------------------------
# Canonical Action Class Map (AUTHORITATIVE)
# ---------------------------------------------------------------------------

$CanonicalMap = @{
    "AC-01" = "Passive-Observation"
    "AC-02" = "Active-Surface-Discovery"
    "AC-03" = "Reachability-Validation"
    "AC-04" = "Service-Identification"
    "AC-05" = "Protocol-Metadata"
    "AC-06" = "Version-Enumeration"
    "AC-07" = "Auth-Surface-Analysis"
    "AC-08" = "Credential-Validation"
    "AC-09" = "Access-Establishment"
    "AC-10" = "Privilege-Assessment"
    "AC-11" = "Lateral-Reachability"
    "AC-12" = "Execution-Capability"
    "AC-13" = "Data-Exposure"
    "AC-14" = "Impact-Feasibility"
    "AC-15" = "External-Execution"
}

# ---------------------------------------------------------------------------
# Validate Source Directory
# ---------------------------------------------------------------------------

if (-not (Test-Path $SourceDir)) {
    Write-Error "Source directory '$SourceDir' not found"
    exit 1
}

# ---------------------------------------------------------------------------
# Collect Action Class YAML Files (Exclude Schema)
# ---------------------------------------------------------------------------

$Files = Get-ChildItem -Path $SourceDir -Recurse -Include *.yaml, *.yml |
         Where-Object { $_.Name -notmatch '^_schema\.ya?ml$' }

if ($Files.Count -eq 0) {
    Write-Error "No Action Class YAML files found under '$SourceDir'"
    exit 1
}

Write-Host "Found $($Files.Count) Action Class files"

# ---------------------------------------------------------------------------
# Group Files by AC ID
# ---------------------------------------------------------------------------

$Grouped = @{}

foreach ($File in $Files) {
    if ($File.Name -match '(AC-\d{2})') {
        $ACID = $Matches[1]
    } else {
        Write-Error "❌ Cannot determine AC ID from filename: $($File.FullName)"
        exit 1
    }

    if (-not $CanonicalMap.ContainsKey($ACID)) {
        Write-Error "❌ Unknown Action Class ID detected: $ACID"
        exit 1
    }

    if (-not $Grouped.ContainsKey($ACID)) {
        $Grouped[$ACID] = @()
    }

    $Grouped[$ACID] += $File
}

# ---------------------------------------------------------------------------
# Report Duplicates (EXPECTED, DOCUMENTED)
# ---------------------------------------------------------------------------

foreach ($ACID in $Grouped.Keys) {
    if ($Grouped[$ACID].Count -gt 1) {
        Write-Warning "⚠ Multiple definitions detected for ${ACID}:"
        $Grouped[$ACID] | ForEach-Object {
            Write-Warning "  - $($_.FullName)"
        }
    }
}

# ---------------------------------------------------------------------------
# Create Target Directory (Fresh)
# ---------------------------------------------------------------------------

if (Test-Path $TargetDir) {
    Write-Error "Target directory '$TargetDir' already exists — aborting to prevent overwrite"
    exit 1
}

New-Item -ItemType Directory -Path $TargetDir | Out-Null

# ---------------------------------------------------------------------------
# Materialize Canonical Action Classes
# ---------------------------------------------------------------------------

foreach ($ACID in $CanonicalMap.Keys | Sort-Object) {

    if (-not $Grouped.ContainsKey($ACID)) {
        Write-Error "❌ Missing definition for required Action Class: $ACID"
        exit 1
    }

    $CanonicalName = $CanonicalMap[$ACID]
    $TargetName = "$ACID-$CanonicalName.yaml"

    # Deterministic rule: first file wins (duplicates already logged)
    $SourceFile = $Grouped[$ACID][0]

    Write-Host "✔ Normalizing $ACID → $TargetName"

    Copy-Item `
        -Path $SourceFile.FullName `
        -Destination (Join-Path $TargetDir $TargetName) `
        -Force
}

# ---------------------------------------------------------------------------
# Completion Summary
# ---------------------------------------------------------------------------

Write-Host ""
Write-Host "=== ACTION CLASS NORMALIZATION COMPLETE ==="
Write-Host "Canonical Action Classes: $($CanonicalMap.Count)"
Write-Host "Output Directory        : $TargetDir"
Write-Host "Status                  : APPLIED"
