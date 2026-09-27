# Starts both projects' dev servers in one Windows Terminal window, each in
# its own pane: api (Go, port 8080) and web (Vite, default port 5173).
#
# Usage: from anywhere, run  pwsh C:\k\kfamily\dev.ps1
# (or cd here first and run .\dev.ps1)

$root = $PSScriptRoot

wt -w -1 `
  -d "$root\api" pwsh -NoExit -Command "go run ./cmd/api" `; `
  split-pane -d "$root\web" pwsh -NoExit -Command "pnpm dev"
