$env:PRODUCTION = 'true'
wails3 build
Remove-Item ENV:PRODUCTION
cp vJoyInterface.dll bin
