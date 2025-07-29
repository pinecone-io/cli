# pcdev.ps1 - PowerShell version for Windows compatibility
# Development script to run the Pinecone CLI from dist folder

param(
    [Parameter(ValueFromRemainingArguments=$true)]
    [string[]]$Arguments
)

# Colors for output
$Red = "Red"
$Green = "Green"
$Yellow = "Yellow"

# Function to print colored output
function Write-Status {
    param([string]$Message)
    Write-Host "[pcdev] $Message" -ForegroundColor $Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[pcdev] $Message" -ForegroundColor $Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[pcdev] $Message" -ForegroundColor $Red
}

# Get script directory
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$DistDir = Join-Path $ScriptDir "dist"

# Check if dist directory exists
if (-not (Test-Path $DistDir)) {
    Write-Error "dist directory not found. Please run 'goreleaser build --single-target --snapshot --clean' first."
    exit 1
}

# Detect OS and architecture
$OS = "windows"
$Arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }

Write-Status "Detected OS: $OS, Architecture: $Arch"

# Determine the correct binary path
$BinaryPath = ""

# On Windows, look for .exe files
$PossiblePaths = @(
    (Join-Path $DistDir "pc_windows_amd64\pc.exe"),
    (Join-Path $DistDir "pc_windows_x86_64\pc.exe"),
    (Join-Path $DistDir "pc_windows_386\pc.exe")
)

foreach ($path in $PossiblePaths) {
    if (Test-Path $path) {
        $BinaryPath = $path
        Write-Status "Using binary: $path"
        break
    }
}

if (-not $BinaryPath) {
    Write-Error "No suitable Windows binary found"
    Write-Status "Available binaries:"
    Get-ChildItem -Path $DistDir -Recurse -Name "*.exe" | ForEach-Object { Write-Host "  $_" }
    exit 1
}

# Check if binary exists
if (-not (Test-Path $BinaryPath)) {
    Write-Error "Binary not found: $BinaryPath"
    exit 1
}

Write-Status "Running: $BinaryPath $($Arguments -join ' ')"

# Execute the binary with all passed arguments
& $BinaryPath @Arguments 