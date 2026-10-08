param([Parameter(Mandatory)][ValidateSet('buttonbox', 'handbrake')][string]$Board)

$boards = @{
    buttonbox = @{ Sketch = 'arduino.ino'; Port = 'COM3' }
    handbrake = @{ Sketch = 'arduino2.ino'; Port = 'COM4' }
}
$b = $boards[$Board]

# arduino-cli needs <name>/<name>.ino; stage a copy so the sketches can stay at repo root
$dir = "$PSScriptRoot\build\$Board"
New-Item -ItemType Directory -Force $dir | Out-Null
Copy-Item "$PSScriptRoot\$($b.Sketch)" "$dir\$Board.ino" -Force

# the feeder holds the COM ports open; stop it for the upload and restart it afterwards
$feeder = Get-Process button-box-vjoy-feeder -ErrorAction SilentlyContinue
$feeder | Stop-Process -Force

arduino-cli compile -u --fqbn arduino:megaavr:nona4809 -p $b.Port $dir
$code = $LASTEXITCODE

if ($feeder) { Start-Process "$PSScriptRoot\bin\button-box-vjoy-feeder.exe" }
exit $code
