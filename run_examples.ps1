$ErrorActionPreference = 'Stop'

$here = Split-Path -Parent $MyInvocation.MyCommand.Path

Set-Location $here

@'
(33 + (912 * 11))
'@ | Set-Content -NoNewline -Encoding ascii input.ec1

@'
(  3	+
(4+5) )
'@ | Set-Content -NoNewline -Encoding ascii ws.ec1

@'
(3 + a)
'@ | Set-Content -NoNewline -Encoding ascii err.ec1

Write-Host '== input.ec1 ==' 
& go run . input.ec1

Write-Host '== ws.ec1 ==' 
& go run . ws.ec1

Write-Host '== err.ec1 (expected error) ==' 
try {
  & go run . err.ec1
  exit 1
} catch {
  # go run returns non-zero, so PowerShell throws.
  Write-Host $_.Exception.Message
}
