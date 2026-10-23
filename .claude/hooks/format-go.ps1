# PostToolUse hook: format a .go file after Edit/Write (goimports + gofumpt).
# Reads the hook JSON from stdin; never blocks the agent (always exit 0).
# ASCII only: may run under Windows PowerShell 5.1.

$ErrorActionPreference = 'SilentlyContinue'
$payload = [Console]::In.ReadToEnd() | ConvertFrom-Json
$file = $payload.tool_input.file_path
if (-not $file -or -not $file.EndsWith('.go') -or -not (Test-Path $file)) { exit 0 }

# Generated code is left as is.
if ((Get-Content $file -TotalCount 1) -match '^// Code generated .* DO NOT EDIT\.$') { exit 0 }

$root = if ($env:CLAUDE_PROJECT_DIR) { $env:CLAUDE_PROJECT_DIR } else { Split-Path $PSScriptRoot -Parent | Split-Path -Parent }
Push-Location $root
try {
    go tool goimports -w $file 2>&1 | Write-Error
    go tool gofumpt -w $file 2>&1 | Write-Error
} finally {
    Pop-Location
}
exit 0
