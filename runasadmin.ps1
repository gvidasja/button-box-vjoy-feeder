$currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
$isAdmin = $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if ($isAdmin) {
  go run .\cmd\button-box-vjoy-feeder\main.go $args[0]
}
else {
  $myPath = Get-Location
  Start-Process -FilePath powershell -verb runas -ArgumentList "-NoExit Set-Location $myPath ; $myPath\runasadmin.ps1 $args"
}