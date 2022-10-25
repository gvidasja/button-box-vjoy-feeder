$currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
$isAdmin = $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if ($isAdmin) {
  foreach ($item in Get-ChildItem cmd | Where-Object { $_.PsIsContainer -eq $True }) {
    $name = $item.Name  
    Write-Host "building $name"
    go build -o ".\bin\$name.exe" ".\cmd\$name\main.go"

    Write-Host "stopping..."
    Invoke-Expression -ErrorAction Continue ".\bin\$name.exe stop"
    Write-Host "stopped"

    Write-Host "uninstalling..."
    Invoke-Expression -ErrorAction Continue ".\bin\$name.exe uninstall"
    Write-Host "uninstalled"
    
    Write-Host "copying..."
    Copy-Item *.dll ".\bin"
    Write-Host "copied"
  
    Write-Host "installing..."
    Invoke-Expression -ErrorAction Continue ".\bin\$name.exe install"
    Write-Host "installed"

    Write-Host "starting..."
    Invoke-Expression -ErrorAction Continue ".\bin\$name.exe start"
    Write-Host "started"

    Read-Host "Press any key to continue..."
  }
}
else {
  $myPath = Get-Location
  Start-Process -FilePath powershell -verb runas -ArgumentList "-NoExit Set-Location $myPath ; $myPath\install.ps1"
}