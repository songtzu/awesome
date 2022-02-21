set GOOS=linux
set GOARCH=amd64
set CurrentDir=%cd%
cd ../
set MainPackDir=%cd%
cd ../
cd ../
set GOPATH=%cd%
echo %GOPATH%
go build -o %cd%\bin\game %MainPackDir%\main\main.go
cd %CurrentDir%