$ErrorActionPreference = 'Stop'

$here = Split-Path -Parent $MyInvocation.MyCommand.Path

Set-Location $here

Write-Host '== tests\input.ec1 ==' 
& go run . .\tests\input.ec1

Write-Host '== tests\ws.ec1 ==' 
& go run . .\tests\ws.ec1

Write-Host '== tests\input2.ec1 ==' 
& go run . .\tests\input2.ec1

Write-Host '== tests\input3.ec1 ==' 
& go run . .\tests\input3.ec1

Write-Host '== tests\err.ec1 (expected lexical error) ==' 
try {
  & go run . .\tests\err.ec1
  exit 1
} catch {
  # go run returns non-zero, so PowerShell throws.
  Write-Host $_.Exception.Message
}

Write-Host '== tests\synerr.ec1 (expected syntax error) ==' 
try {
  & go run . .\tests\synerr.ec1
  exit 1
} catch {
  Write-Host $_.Exception.Message
}

Write-Host '== tests\err1.ec1 (expected runtime error) ==' 
try {
  & go run . .\tests\err1.ec1
  exit 1
} catch {
  Write-Host $_.Exception.Message
}
