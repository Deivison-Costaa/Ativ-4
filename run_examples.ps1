$ErrorActionPreference = 'Stop'

$here = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $here

Write-Host 'Compilando o compilador Fun...'
go build -o cmd.exe .

$examples = @(
  'tests\abs.fun',
  'tests\noargs.fun',
  'tests\chain.fun',
  'tests\shadow.fun',
  'tests\fact.fun'
)

foreach ($file in $examples) {
  Write-Host "== $file =="
  & .\cmd.exe $file
}

Write-Host 'Arquivos .s gerados. Para montar e executar os binarios, use Linux/WSL com bash run_tests.sh.'
