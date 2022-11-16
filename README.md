# button-box-vjoy-feeder

vJoy feeder app for a custom button-box

## Installation

1. Build
    ```powershell
    go build -O ".path\to\exe" ".\cmd\button-box-vjoy-feeder\main.go"
    ```
1. Create a Windows Service
    ```powershell
    New-Service -Name "button-box-vjoy-feeder" -BinaryPathName ".\path\to\exe"
    ```
1. Set service user as your own, to write logs