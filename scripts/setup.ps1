# Contributing Project Setup Script

$script = @"
# KubeOps Development Setup Script for Windows

Write-Host "KubeOps Development Setup" -ForegroundColor Green
Write-Host ""

# Check prerequisites
Write-Host "Checking prerequisites..." -ForegroundColor Yellow

$requiredCommands = @("go", "node", "docker", "kubectl", "helm", "kind")
$missing = @()

foreach ($cmd in $requiredCommands) {
    $exists = Get-Command $cmd -ErrorAction SilentlyContinue
    if ($exists) {
        Write-Host "  ✓ $cmd" -ForegroundColor Green
    } else {
        Write-Host "  ✗ $cmd (missing)" -ForegroundColor Red
        $missing += $cmd
    }
}

if ($missing.Count -gt 0) {
    Write-Host ""
    Write-Host "Missing tools: $($missing -join ', ')" -ForegroundColor Red
    Write-Host "Please install missing tools and run again." -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "All prerequisites installed!" -ForegroundColor Green
Write-Host ""

# Install Go dependencies
Write-Host "Installing Go dependencies..." -ForegroundColor Yellow
Set-Location backend
go mod download
Set-Location ..

# Install Node dependencies
Write-Host "Installing Node dependencies..." -ForegroundColor Yellow
Set-Location frontend
npm install
Set-Location ..

Write-Host ""
Write-Host "Setup complete! ✓" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "1. Start databases: docker-compose -f deploy/docker-compose-dev.yaml up -d"
Write-Host "2. Run backend services (see README)"
Write-Host "3. Run frontend: cd frontend && npm run dev"
Write-Host ""
"@

$script | Out-File -FilePath setup.ps1 -Encoding UTF8
