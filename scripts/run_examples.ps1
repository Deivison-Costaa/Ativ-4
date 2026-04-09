$ErrorActionPreference = 'Stop'

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$rootDir = Split-Path -Parent $scriptDir
$positiveDir = Join-Path $rootDir 'testdata\fun\positive'
$tmpDir = Join-Path $env:TEMP 'funcc-examples'
$binDir = Join-Path $tmpDir 'bin'
$outDir = Join-Path $tmpDir 'asm'
$bin = Join-Path $binDir 'funcc.exe'

New-Item -ItemType Directory -Force -Path $binDir | Out-Null
New-Item -ItemType Directory -Force -Path $outDir | Out-Null

Push-Location $rootDir
try {
    go build -o $bin .\cmd\funcc

    Get-ChildItem $positiveDir -Filter *.fun | ForEach-Object {
        $output = Join-Path $outDir ($_.BaseName + '.s')
        & $bin -o $output $_.FullName
        Write-Host ("Gerado: {0}" -f $output)
    }

    Write-Host ""
    Write-Host ("Arquivos gerados em: {0}" -f $outDir)
    Write-Host 'Para montar e executar, prefira Linux/WSL e rode scripts/run_tests.sh a partir da raiz do projeto.'
}
finally {
    Pop-Location
}
