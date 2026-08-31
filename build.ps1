windres -i appicon.rc -o appicon_windows.syso --target=pe-x86-64
go build -trimpath -buildvcs=false -ldflags="-w -s -H windowsgui" -o bin/button-box-vjoy-feeder.exe .
Copy-Item vJoyInterface.dll bin -Force
