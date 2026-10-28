<#
.SYNOPSIS
  Builds and launches the lem-in 3D visualizer on Windows.

.EXAMPLE
  .\run-3d.ps1                              # visualizes input.txt
  .\run-3d.ps1 -Map test\example01.txt      # visualizes another map
  .\run-3d.ps1 -Map test\example06.txt -Speed 3
  .\run-3d.ps1 -Check                       # solve + verify only, no window
#>
param(
    [string]$Map = "input.txt",
    [double]$Speed = 0.8,
    [double]$Scale = 1.6,
    [switch]$Check
)

$ErrorActionPreference = "Stop"
$root = $PSScriptRoot

if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    throw "Go is not installed or not on PATH (https://go.dev/dl/)."
}
foreach ($dll in "raylib.dll", "glfw3.dll") {
    if (-not (Test-Path (Join-Path $root $dll))) {
        throw "$dll is missing from $root. The visualizer loads it at start-up (see README > Troubleshooting)."
    }
}

# raylib-go loads raylib.dll at run time (purego), so no C compiler is needed.
$env:CGO_ENABLED = "0"

Push-Location $root
try {
    go build -o lem-in-3d.exe ./3D
    if ($LASTEXITCODE -ne 0) { throw "go build failed" }

    $flags = @("-input", $Map, "-speed", $Speed, "-scale", $Scale)
    if ($Check) { $flags += "-check" }
    # Run from the repository root so Windows finds raylib.dll and glfw3.dll next to the executable.
    & (Join-Path $root "lem-in-3d.exe") @flags
    exit $LASTEXITCODE
}
finally {
    Pop-Location
}
