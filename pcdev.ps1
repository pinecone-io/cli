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
$CacheFile = Join-Path $ScriptDir ".pcdev_cache"

# Function to check if rebuild is needed
function Test-RebuildNeeded {
    # If dist doesn't exist, definitely need to build
    if (-not (Test-Path $DistDir)) {
        return $true
    }
    
    # If cache file doesn't exist, need to build
    if (-not (Test-Path $CacheFile)) {
        return $true
    }
    
    # Check if any Go files have been modified since last build
    $LastBuildTime = Get-Content $CacheFile -ErrorAction SilentlyContinue
    if (-not $LastBuildTime) {
        $LastBuildTime = 0
    }
    
    # Find the most recent modification time of Go files
    $GoFiles = Get-ChildItem -Path $ScriptDir -Filter "*.go" -Recurse -ErrorAction SilentlyContinue
    if ($GoFiles) {
        $LatestFileTime = ($GoFiles | ForEach-Object { $_.LastWriteTime.ToFileTime() } | Measure-Object -Maximum).Maximum
    } else {
        $LatestFileTime = 0
    }
    
    # Compare timestamps
    if ($LatestFileTime -gt $LastBuildTime) {
        return $true
    }
    
    return $false
}

# Function to update cache
function Update-Cache {
    [DateTime]::Now.ToFileTime() | Out-File -FilePath $CacheFile -Encoding UTF8
}

# Check if rebuild is needed
if (Test-RebuildNeeded) {
    Write-Warning "Source files have changed or dist directory missing. Building..."
    goreleaser build --single-target --snapshot --clean
    Update-Cache
    Write-Status "Build completed successfully!"
}

# Check if dist directory exists (final check)
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