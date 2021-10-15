foreach ($item in Get-ChildItem cmd | Where-Object { $_.PsIsContainer -eq $True }) {
  $name = $item.Name  
  Write-Host "building $name"
  go build -o ".\bin\$name.exe" ".\cmd\$name\main.go"

  mkdir -force "C:\bin\$name"
  cp ".\bin\$name.exe" "C:\bin\$name\$name.exe"
  cp *.dll "C:\bin\$name\"
}