$ErrorActionPreference = "Stop"

function Test-DockerReady {
  docker info *> $null
  return $LASTEXITCODE -eq 0
}

function Wait-DockerReady {
  param([int]$Seconds = 120)

  $deadline = (Get-Date).AddSeconds($Seconds)
  while ((Get-Date) -lt $deadline) {
    if (Test-DockerReady) {
      return $true
    }
    Start-Sleep -Seconds 2
  }

  return $false
}

function Start-DockerDesktop {
  $paths = @(
    "$Env:ProgramFiles\Docker\Docker\Docker Desktop.exe",
    "${Env:ProgramFiles(x86)}\Docker\Docker\Docker Desktop.exe",
    "$Env:LocalAppData\Docker\Docker Desktop.exe"
  )

  foreach ($path in $paths) {
    if ($path -and (Test-Path $path)) {
      Write-Host "Starting Docker Desktop..."
      Start-Process -FilePath $path | Out-Null
      return $true
    }
  }

  return $false
}

if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
  Write-Host @"
Docker is not installed or not available on PATH.

Install Docker Desktop, open it once, then rerun:
  .\sandbox.ps1

This repository does not require Go, Node, npm, SQLite, or cron on the host.
"@
  exit 127
}

docker compose version *> $null
if ($LASTEXITCODE -ne 0) {
  Write-Host @"
Docker is installed, but Docker Compose is not available.

Install or update Docker Desktop, then rerun:
  .\sandbox.ps1
"@
  exit 127
}

if (-not (Test-DockerReady)) {
  if (Start-DockerDesktop) {
    if (-not (Wait-DockerReady -Seconds 120)) {
      Write-Host @"
Docker Desktop started, but Docker did not become ready in time.

Open Docker Desktop, wait until it says it is running, then rerun:
  .\sandbox.ps1
"@
      exit 1
    }
  } else {
    Write-Host @"
Docker is installed but the daemon is not responding.

Start Docker Desktop, then rerun:
  .\sandbox.ps1
"@
    exit 1
  }
}

docker compose up --build -d
if ($LASTEXITCODE -ne 0) {
  exit $LASTEXITCODE
}

Write-Host @"

Sandbox is starting.

Open:
  http://localhost:3000/resources

Useful commands:
  docker compose logs -f
  curl -X POST http://localhost:8080/scrape
  docker compose down
"@
